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
