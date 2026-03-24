package handlers

import (
	"fmt"
	"io"
)

// openTag writes an opening XML tag to w. Optional attrs are key-value pairs:
// openTag(w, "command", "index", "0") writes <command index="0">.
func openTag(w io.Writer, tag string, attrs ...string) {
	_, _ = fmt.Fprintf(w, "<%s", tag)
	for i := 0; i+1 < len(attrs); i += 2 {
		_, _ = fmt.Fprintf(w, " %s=%q", attrs[i], attrs[i+1])
	}
	_, _ = fmt.Fprintf(w, ">\n")
}

// closeTag writes a closing XML tag to w.
func closeTag(w io.Writer, tag string) {
	_, _ = fmt.Fprintf(w, "</%s>\n", tag)
}

// parseStringSlice coerces a []any (as returned by mcp-go for array params)
// into a []string. Returns false if v is not a slice or contains non-string elements.
func parseStringSlice(v any) ([]string, bool) {
	raw, ok := v.([]any)
	if !ok {
		return nil, false
	}

	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}

	return out, true
}
