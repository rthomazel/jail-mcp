package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) HandleExecBackground(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	command, ok := args["command"].(string)
	if !ok || command == "" {
		return mcp.NewToolResultError("missing required parameter: command"), nil
	}

	cwd, _ := args["cwd"].(string)
	if cwd == "" {
		cwd = "/"
	}

	job := &job{
		cmd:     command,
		started: time.Now(),
	}
	h.addJob(job)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), h.cfg.BackgroundTimeout)
		defer cancel()

		slog.Info("bg exec start", "job", job.id, "cmd", command, "cwd", cwd)

		cmd := exec.CommandContext(ctx, "bash", "-c", command)
		cmd.Dir = cwd

		job.mu.Lock()
		cmd.Stdout = &job.stdout
		cmd.Stderr = &job.stderr
		job.mu.Unlock()

		err := cmd.Run()

		job.mu.Lock()
		defer job.mu.Unlock()

		job.done = true

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				job.exitCode = exitErr.ExitCode()
			} else {
				job.err = fmt.Sprintf("could not start process: %v", err)
				job.exitCode = -1
			}
		}

		slog.Info("bg exec done", "job", job.id, "exit_code", job.exitCode, "duration", time.Since(job.started).Round(time.Millisecond))
	}()

	b, err := json.Marshal(map[string]any{"job_id": job.id})
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}
