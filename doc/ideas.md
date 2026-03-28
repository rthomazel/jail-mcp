# ideas

## concurrent context

`context` tool runs subprocesses serially.
Could run them with goroutines and be meaningfully faster.

## per-command timeout

Timeout is global via `JAIL_MCP_TIMEOUT`.
Letting `exec_sync` accept an optional `timeout` param would be useful for known slow commands.

## sqlite db with command stats

Server would tokenize commands with weights, base command has higher weight, then flags.
Normalize input.
Expose historic command stats to allow planning when to use exec sync or background.

# what not to add

- filesystem MCP tool — redundant, shell already does cat/ls/cp/find
- command allowlists — defeats the purpose, Docker is the boundary

## indented xml output

`xmlBuilder` has no depth tracking. A `depth int` field incremented by `openTag` /
decremented by `closeTag` would let all write methods prepend indentation.
Metadata fields written directly via `WriteString` would need a `b.line(s)` helper
to respect depth — a wider refactor touching all handlers.

## path snapshot registration file

Setup scripts could write a `.jail-mcp-extras` file in the project root —
one `name: /path/to/binary` pair per line. `context` reads all such files under
known project roots and surfaces them alongside the `auto-detected in path:` block.
Explicit opt-in, works for non-PATH installs, but requires setup scripts to be
authored with the convention.
