package handlers

import (
	"context"
	"fmt"
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

	var b strings.Builder
	fmt.Fprintf(&b, "<metadata>\n")
	fmt.Fprintf(&b, "job_id: %s\n", j.id)
	fmt.Fprintf(&b, "</metadata>\n")

	return mcp.NewToolResultText(b.String()), nil
}
