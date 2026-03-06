package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type Handler struct {
	executor *Executor
	cfg      *Config
	log      *Logger
}

func newHandler(executor *Executor, cfg *Config, log *Logger) *Handler {
	return &Handler{executor: executor, cfg: cfg, log: log}
}

func (h *Handler) handleExec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command, err := req.RequireString("command")
	if err != nil {
		return mcp.NewToolResultError("missing required parameter: command"), nil
	}

	// Default cwd to first allowed dir if caller did not specify one.
	cwd := req.GetString("cwd", h.cfg.AllowedDirs[0])

	if !h.isAllowedDir(cwd) {
		return mcp.NewToolResultError(fmt.Sprintf(
			"cwd %q is not under an allowed directory. Allowed: %v",
			cwd, h.cfg.AllowedDirs,
		)), nil
	}

	result := h.executor.Run(ctx, command, cwd)

	// Process could not start at all — surface as a tool-level error.
	if result.Error != "" {
		return mcp.NewToolResultError(result.Error), nil
	}

	// Command ran (even if exit_code != 0) — return full structured result.
	// The caller can inspect exit_code and stderr to decide what happened.
	resp := map[string]any{
		"stdout":    result.Stdout,
		"stderr":    result.Stderr,
		"exit_code": result.ExitCode,
		"duration":  result.Duration.Round(1_000_000).String(), // round to ms
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}

func (h *Handler) handleListDirs(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(strings.Join(h.cfg.AllowedDirs, "\n")), nil
}

// isAllowedDir returns true if cwd is exactly an allowed dir or a subpath of one.
func (h *Handler) isAllowedDir(cwd string) bool {
	for _, allowed := range h.cfg.AllowedDirs {
		if cwd == allowed || strings.HasPrefix(cwd, allowed+"/") {
			return true
		}
	}
	return false
}
