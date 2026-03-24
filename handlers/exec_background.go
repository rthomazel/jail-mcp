package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) HandleExecBackground(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	commands, ok := parseStringSlice(req.Params.Arguments["commands"])
	if !ok || len(commands) == 0 {
		return mcp.NewToolResultError("missing required parameter: commands"), nil
	}

	cwd, _ := req.Params.Arguments["cwd"].(string)
	if cwd == "" {
		cwd = "/"
	}

	multi := len(commands) > 1
	b := strings.Builder{}

	openTag(&b, "metadata")

	for i, cmd := range commands {
		j := h.startJob(cmd, cwd)
		if multi {
			fmt.Fprintf(&b, "command_%d: %s\n", i, cmd)
			fmt.Fprintf(&b, "job_id_%d: %s\n", i, j.id)
		} else {
			b.WriteString("job_id: " + j.id + "\n")
		}
	}

	closeTag(&b, "metadata")

	return mcp.NewToolResultText(b.String()), nil
}
