package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tcodes0/jail-mcp/internal"
)

type commandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration string
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

	return mcp.NewToolResultText(formatPlainText(result)), nil
}

func formatPlainText(r *commandResult) string {
	b := strings.Builder{}

	b.WriteString("<metadata>\n")
	b.WriteString("exit: " + strconv.Itoa(r.ExitCode) + "\n")
	b.WriteString("duration: " + r.Duration + "\n")
	b.WriteString("</metadata>\n")

	b.WriteString("\n<stdout>\n" + r.Stdout + "\n</stdout>\n")
	b.WriteString("\n<stderr>\n" + r.Stderr + "\n</stderr>\n")

	return b.String()
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
				Duration: duration.Round(1_000_000).String(),
				ExitCode: -1,
				err:      fmt.Sprintf("could not start process: %v", err),
			}
		}
	}

	slog.Info("exec done", "cmd", command, "exit_code", exitCode, "duration", duration.Round(time.Millisecond))

	return &commandResult{
		Stdout:   strings.TrimRight(stdout.String(), "\n"),
		Stderr:   strings.TrimRight(stderr.String(), "\n"),
		ExitCode: exitCode,
		Duration: duration.Round(1_000_000).String(),
	}
}
