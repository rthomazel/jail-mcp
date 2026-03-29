package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rthomazel/jail-mcp/internal/pathsnapshot"
	"github.com/samber/lo"
)

var (
	skipFSTypes  = []string{"proc", "sysfs", "tmpfs", "devpts", "cgroup2", "cgroup", "mqueue", "overlay"}
	skipPrefixes = []string{"/proc", "/sys", "/dev", "/run", "/etc", "/mise", "/root"}
)

const miseShimsDir = "/mise/shims"

var preInstalled = map[string]string{
	"bash":    "bash --version | head -1 | cut -d' ' -f4",
	"git":     "git --version | cut -d' ' -f3",
	"jujutsu": "jj version",
	"mise":    "mise v 2>/dev/null | tail -1",
	"python3": "python3 --version | cut -d' ' -f2",
	"pip3":    "pip3 --version 2>/dev/null | cut -d' ' -f2",
	"rg":      "rg --version | head -1 | cut -d' ' -f2",
	"make":    "make --version | head -1 | cut -d' ' -f3",
	"jq":      "jq --version",
	"curl":    "curl --version | head -1",
}

func (h *Handler) HandleContext(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gather := func(cmd string) string {
		r := runCommand(ctx, h.cfg, cmd, "/")
		return strings.TrimSpace(r.Stdout)
	}

	osName := gather("cat /etc/os-release | grep PRETTY_NAME | cut -d= -f2 | tr -d '\"'")
	arch := gather("uname -m")
	disk := gather("df -h / | awk 'NR==2{print $4\" free of \"$2}'")
	path := os.Getenv("PATH")

	var mounts []string
	file, err := os.Open("/proc/mounts")
	if err != nil {
		slog.Error("failed to read mounts", "err", err)
	} else {
		defer func() { _ = file.Close() }()
		mounts, err = parseMounts(file)
		if err != nil {
			slog.Error("failed to parse mounts", "err", err)
		}
	}

	versions := make(map[string]string, len(preInstalled))
	for name, cmd := range preInstalled {
		v := gather(cmd)
		if v == "" {
			v = "-"
		}
		versions[name] = v
	}

	miseShims := discoverMiseShims()
	detected := pathsnapshot.Diff()

	return mcp.NewToolResultText(formatPlainTextContext(osName, arch, disk, path, h.cfg.Timeout.String(), h.version, mounts, versions, miseShims, detected)), nil
}

// discoverMiseShims returns executable filenames in miseShimsDir, sorted, skipping
// non-executable files and wrapper scripts (.cmd, .js).
func discoverMiseShims() []string {
	entries, err := os.ReadDir(miseShimsDir)
	if err != nil {
		return nil
	}

	var shims []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.Contains(name, ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.Mode()&0o111 == 0 {
			continue
		}
		shims = append(shims, filepath.Base(name))
	}

	sort.Strings(shims)
	return shims
}

func formatPlainTextContext(osName, arch, disk, path, timeout, version string, projects []string, versions map[string]string, miseShims []string, detected []pathsnapshot.Entry) string {
	b := strings.Builder{}

	b.WriteString("<metadata>\n")
	b.WriteString("os: " + osName + "\n")
	b.WriteString("arch: " + arch + "\n")
	b.WriteString("disk: " + disk + "\n")
	b.WriteString("path: " + path + "\n")
	b.WriteString("shell_exec_timeout: " + timeout + "\n")
	b.WriteString("version: " + version + "\n")

	b.WriteString("projects:\n")
	for _, p := range projects {
		b.WriteString("  " + p + "\n")
	}

	b.WriteString("installed:\n")
	maxLen := 0

	for name := range preInstalled {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for name := range preInstalled {
		b.WriteString("  " + fmt.Sprintf("%-*s", maxLen+1, name+":") + " " + versions[name] + "\n")
	}

	if len(miseShims) > 0 {
		b.WriteString("mise shims:\n")
		for _, s := range miseShims {
			b.WriteString("  " + s + "\n")
		}
	}

	if len(detected) > 0 {
		b.WriteString("auto-detected in path:\n")
		for _, e := range detected {
			b.WriteString("  " + e.Name + " " + e.Path + "\n")
		}
	}

	b.WriteString("</metadata>\n")

	return b.String()
}

func parseMounts(r io.Reader) ([]string, error) {
	var candidates []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}

		mountpoint, fstype := fields[1], fields[2]

		if mountpoint == "/" || lo.Contains(skipFSTypes, fstype) {
			continue
		}

		isSkipped := lo.SomeBy(skipPrefixes, func(p string) bool {
			return mountpoint == p || strings.HasPrefix(mountpoint, p+"/")
		})
		if isSkipped {
			continue
		}

		candidates = append(candidates, mountpoint)
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, scanErr
	}

	sort.Slice(candidates, func(i, j int) bool {
		return len(candidates[i]) < len(candidates[j])
	})

	kept := lo.Reduce(candidates, func(acc []string, candidate string, _ int) []string {
		isChild := lo.SomeBy(acc, func(k string) bool {
			return strings.HasPrefix(candidate, k+"/")
		})
		if isChild {
			return acc
		}
		return append(acc, candidate)
	}, []string{})

	sort.Strings(kept)
	return kept, nil
}
