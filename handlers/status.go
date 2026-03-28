package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) HandleStatus(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ids, ok := parseStringSlice(req.Params.Arguments["job_ids"])
	if !ok || len(ids) == 0 {
		return mcp.NewToolResultError("missing required parameter: job_ids"), nil
	}

	multi := len(ids) > 1
	var b xmlBuilder

	for i, id := range ids {
		if multi {
			b.openTag("command", "index", strconv.Itoa(i))
		}

		formatJobStatus(&b, h, id, multi)

		if multi {
			b.closeTag("command", true)
		}
	}

	return mcp.NewToolResultText(b.String()), nil
}

func formatJobStatus(b *xmlBuilder, h *Handler, id string, includeCommand bool) {
	h.mu.RLock()
	j, exists := h.jobs[id]
	h.mu.RUnlock()

	b.openTag("metadata")
	if !exists {
		b.WriteString("job_id: " + id + "\n")
		b.WriteString("error: job not found\n")
		b.closeTag("metadata", false)
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	if includeCommand {
		b.WriteString("command: " + j.cmd + "\n")
	}

	b.WriteString("job_id: " + j.id + "\n")
	b.WriteString("done: " + strconv.FormatBool(j.done) + "\n")

	if j.done {
		b.WriteString("exit: " + strconv.Itoa(j.exitCode) + "\n")
	}

	b.WriteString("duration: " + time.Since(j.started).Round(time.Millisecond).String() + "\n")

	if j.err != "" {
		b.WriteString("error: " + j.err + "\n")
	}

	b.closeTag("metadata", true)
	b.tag("stdout", j.stdout.String(), true)
	b.tag("stderr", j.stderr.String(), false)
}
