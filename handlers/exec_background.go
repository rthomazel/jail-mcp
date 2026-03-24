package handlers

import (
	"context"
	"strings"

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

	b := strings.Builder{}
	b.WriteString("<metadata>\n")
	b.WriteString("job_id: " + j.id + "\n")
	b.WriteString("</metadata>\n")

	return mcp.NewToolResultText(b.String()), nil
}
