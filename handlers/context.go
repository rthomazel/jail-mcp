package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
)

var (
	skipFstypes  = []string{"proc", "sysfs", "tmpfs", "devpts", "cgroup2", "cgroup", "mqueue", "overlay"}
	skipPrefixes = []string{"/proc", "/sys", "/dev", "/run", "/etc"}
)

func (h *Handler) HandleContext(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gather := func(cmd string) string {
		r := runCommand(ctx, h.cfg, cmd, "/")
		return strings.TrimSpace(r.stdout)
	}

	mounts, err := mountedPaths()
	if err != nil {
		slog.Error("failed to read mounts", "err", err)
		mounts = []string{}
	}

	info := map[string]any{
		"os":                 gather("cat /etc/os-release | grep PRETTY_NAME | cut -d= -f2 | tr -d '\"'"),
		"arch":               gather("uname -m"),
		"shell_exec_timeout": h.cfg.Timeout.String(),
		"projects":           strings.Join(mounts, "\n"),
		"disk":               gather("df -h / | awk 'NR==2{print $4\" free of \"$2}'"),
		"tools": map[string]string{
			"bash":    gather("bash --version | head -1 | cut -d' ' -f4"),
			"git":     gather("git --version | cut -d' ' -f3"),
			"go":      gather("go version | cut -d' ' -f3"),
			"python3": gather("python3 --version | cut -d' ' -f2"),
			"node":    gather("node --version"),
			"make":    gather("make --version | head -1 | cut -d' ' -f3"),
			"jq":      gather("jq --version"),
		},
	}

	b, err := json.Marshal(info)
	if err != nil {
		return mcp.NewToolResultError("failed to encode context"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}

func mountedPaths() ([]string, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var candidates []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}

		mountpoint, fstype := fields[1], fields[2]

		if mountpoint == "/" || lo.Contains(skipFstypes, fstype) {
			continue
		}

		if lo.SomeBy(skipPrefixes, func(p string) bool {
			return mountpoint == p || strings.HasPrefix(mountpoint, p+"/")
		}) {
			continue
		}

		candidates = append(candidates, mountpoint)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// sort by length for parent-first deduplication
	sort.Slice(candidates, func(i, j int) bool {
		return len(candidates[i]) < len(candidates[j])
	})

	kept := lo.Reduce(candidates, func(acc []string, candidate string, _ int) []string {
		if lo.SomeBy(acc, func(k string) bool {
			return strings.HasPrefix(candidate, k+"/")
		}) {
			return acc
		}
		return append(acc, candidate)
	}, []string{})

	sort.Strings(kept)
	return kept, nil
}
