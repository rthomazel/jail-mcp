package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tcodes0/jail-mcp/internal"
)

type commandResult struct {
	stdout   string
	stderr   string
	exitCode int
	duration time.Duration
	err      string
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

	if result.err != "" {
		return mcp.NewToolResultError(result.err), nil
	}

	b, err := json.Marshal(map[string]any{
		"stdout":    result.stdout,
		"stderr":    result.stderr,
		"exit_code": result.exitCode,
		"duration":  result.duration.Round(1_000_000).String(),
	})
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}

func runCommand(ctx context.Context, cfg *internal.Config, command, cwd string) *commandResult {
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
			return &commandResult{
				duration: duration,
				exitCode: -1,
				err:      fmt.Sprintf("could not start process: %v", err),
			}
		}
	}

	slog.Info("exec done", "cmd", command, "exit_code", exitCode, "duration", duration.Round(time.Millisecond))

	return &commandResult{
		stdout:   strings.TrimRight(stdout.String(), "\n"),
		stderr:   strings.TrimRight(stderr.String(), "\n"),
		exitCode: exitCode,
		duration: duration,
	}
}
