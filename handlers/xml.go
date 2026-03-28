package handlers

import (
	"fmt"
	"strings"
)

// xmlBuilder wraps strings.Builder with XML tag helpers.
type xmlBuilder struct{ strings.Builder }

// openTag writes an opening XML tag. Optional attrs are key-value pairs:
// b.openTag("command", "index", "0") writes <command index="0">.
func (b *xmlBuilder) openTag(tag string, attrs ...string) {
	_, _ = fmt.Fprintf(&b.Builder, "<%s", tag)

	for i := 0; i+1 < len(attrs); i += 2 {
		_, _ = fmt.Fprintf(&b.Builder, " %s=%q", attrs[i], attrs[i+1])
	}

	b.WriteString(">\n")
}

// closeTag writes a closing XML tag with optional trailing newline.
func (b *xmlBuilder) closeTag(tag string, newline bool) {
	_, _ = fmt.Fprintf(&b.Builder, "</%s>\n", tag)

	if newline {
		b.WriteString("\n")
	}
}

// tag writes <name attrs...>\ncontents\n</name>\n — with optional trailing newline.
func (b *xmlBuilder) tag(name, contents string, newline bool, attrs ...string) {
	b.openTag(name, attrs...)

	contents = strings.TrimRight(contents, "\n")

	b.WriteString(contents)
	b.WriteString("\n")
	b.closeTag(name, newline)
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
