package handlers

import (
	"context"
	"fmt"
	"strconv"
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

	b := strings.Builder{}

	b.WriteString("<metadata>\n")
	b.WriteString("job_id: " + j.id + "\n")
	b.WriteString("done: " + strconv.FormatBool(j.done) + "\n")
	if j.done {
		b.WriteString("exit: " + strconv.Itoa(j.exitCode) + "\n")
	}
	b.WriteString("duration: " + time.Since(j.started).Round(time.Millisecond).String() + "\n")
	if j.err != "" {
		b.WriteString("error: " + j.err + "\n")
	}
	b.WriteString("</metadata>\n")

	stdout := strings.TrimRight(j.stdout.String(), "\n")
	stderr := strings.TrimRight(j.stderr.String(), "\n")

	b.WriteString("\n<stdout>\n" + stdout + "\n</stdout>\n")
	b.WriteString("\n<stderr>\n" + stderr + "\n</stderr>\n")

	return mcp.NewToolResultText(b.String()), nil
}
