package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
)

type Handler struct {
	cfg *Config
}

func NewHandler(cfg *Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) HandleContext(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gather := func(cmd string) string {
		r := runCommand(ctx, h.cfg, cmd, "/")
		return strings.TrimSpace(r.Stdout)
	}

	mounts, err := mountedPaths()
	if err != nil {
		slog.Error("failed to read mounts", "err", err)
		mounts = []string{}
	}

	info := map[string]any{
		"os":       gather("cat /etc/os-release | grep PRETTY_NAME | cut -d= -f2 | tr -d '\"'"),
		"arch":     gather("uname -m"),
		"projects": strings.Join(mounts, "\n"),
		"disk":     gather("df -h / | awk 'NR==2{print $4\" free of \"$2}'"),
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

var (
	skipFstypes  = []string{"proc", "sysfs", "tmpfs", "devpts", "cgroup2", "cgroup", "mqueue", "overlay"}
	skipPrefixes = []string{"/proc", "/sys", "/dev", "/run", "/etc"}
)

func mountedPaths() ([]string, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer f.Close()

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

		shouldSkip := lo.SomeBy(skipPrefixes, func(p string) bool {
			return mountpoint == p || strings.HasPrefix(mountpoint, p+"/")
		})

		if shouldSkip {
			continue
		}

		candidates = append(candidates, mountpoint)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// sort by length for deduplication
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

func (h *Handler) HandleExec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	command, ok := args["command"].(string)
	if !ok || command == "" {
		return mcp.NewToolResultError("missing required parameter: command"), nil
	}

	cwd, _ := args["cwd"].(string)
	if cwd == "" {
		cwd = "/"
	}

	result := runCommand(ctx, h.cfg, command, cwd)

	if result.Error != "" {
		return mcp.NewToolResultError(result.Error), nil
	}

	resp := map[string]any{
		"stdout":    result.Stdout,
		"stderr":    result.Stderr,
		"exit_code": result.ExitCode,
		"duration":  result.Duration.Round(1_000_000).String(),
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}

type result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	Error    string
}

func runCommand(ctx context.Context, cfg *Config, command, cwd string) *result {
	start := time.Now()
	slog.Info("exec start", "cmd", command, "cwd", cwd)

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = cwd

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			slog.Error("exec failed to start", "cmd", command, "err", err)
			return &result{
				Duration: duration,
				ExitCode: -1,
				Error:    fmt.Sprintf("could not start process: %v", err),
			}
		}
	}

	slog.Info("exec done", "cmd", command, "exit_code", exitCode, "duration", duration.Round(time.Millisecond))

	return &result{
		Stdout:   strings.TrimRight(stdout.String(), "\n"),
		Stderr:   strings.TrimRight(stderr.String(), "\n"),
		ExitCode: exitCode,
		Duration: duration,
	}
}
