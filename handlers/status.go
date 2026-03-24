package handlers

import (
	"context"
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

	var b strings.Builder

	fmt.Fprintf(&b, "<metadata>\n")
	fmt.Fprintf(&b, "job_id: %s\n", j.id)
	fmt.Fprintf(&b, "done: %t\n", j.done)
	if j.done {
		fmt.Fprintf(&b, "exit: %d\n", j.exitCode)
	}
	fmt.Fprintf(&b, "duration: %s\n", time.Since(j.started).Round(time.Millisecond))
	if j.err != "" {
		fmt.Fprintf(&b, "error: %s\n", j.err)
	}
	fmt.Fprintf(&b, "</metadata>\n")

	stdout := strings.TrimRight(j.stdout.String(), "\n")
	stderr := strings.TrimRight(j.stderr.String(), "\n")

	fmt.Fprintf(&b, "\n<stdout>\n%s\n</stdout>\n", stdout)
	fmt.Fprintf(&b, "\n<stderr>\n%s\n</stderr>\n", stderr)

	return mcp.NewToolResultText(b.String()), nil
}
