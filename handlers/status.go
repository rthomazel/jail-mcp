package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) HandleStatus(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := req.Params.Arguments["job_id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("missing required parameter: job_id"), nil
	}

	h.mu.RLock()
	j, exists := h.jobs[id]
	h.mu.RUnlock()

	if !exists {
		return mcp.NewToolResultError(fmt.Sprintf("job %s not found", id)), nil
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	resp := map[string]any{
		"job_id":   j.id,
		"done":     j.done,
		"stdout":   strings.TrimRight(j.stdout.String(), "\n"),
		"stderr":   strings.TrimRight(j.stderr.String(), "\n"),
		"duration": time.Since(j.started).Round(time.Millisecond).String(),
	}
	if j.done {
		resp["exit_code"] = j.exitCode
	}
	if j.err != "" {
		resp["error"] = j.err
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}
