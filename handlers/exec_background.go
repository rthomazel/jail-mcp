package handlers

import (
	"context"
	"encoding/json"

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

	j := h.startJob(command, cwd)

	b, err := json.Marshal(map[string]any{"job_id": j.id})
	if err != nil {
		return mcp.NewToolResultError("failed to encode result"), nil
	}

	return mcp.NewToolResultText(string(b)), nil
}
