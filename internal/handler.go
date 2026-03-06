package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type Handler struct {
	cfg *Config
}

func (h *Handler) HandleExec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	command, ok := args["command"].(string)
	if !ok || command == "" {
		return mcp.NewToolResultError("missing required parameter: command"), nil
	}

	cwd := h.cfg.AllowedDirs[0]
	if v, ok := args["cwd"].(string); ok && v != "" {
		cwd = v
	}

	if !h.isAllowedDir(cwd) {
		slog.Warn("cwd not allowed", "cwd", cwd, "allowed", h.cfg.AllowedDirs)
		return mcp.NewToolResultError(fmt.Sprintf(
			"cwd %q is not under an allowed directory. Allowed: %v",
			cwd, h.cfg.AllowedDirs,
		)), nil
	}

	result := RunCommand(ctx, h.cfg, command, cwd)

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

func (h *Handler) HandleListDirs(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(strings.Join(h.cfg.AllowedDirs, "\n")), nil
}

func (h *Handler) isAllowedDir(cwd string) bool {
	for _, allowed := range h.cfg.AllowedDirs {
		if cwd == allowed || strings.HasPrefix(cwd, allowed+"/") {
			return true
		}
	}
	return false
}

func NewHandler(cfg *Config) *Handler {
	return &Handler{cfg: cfg}
}
