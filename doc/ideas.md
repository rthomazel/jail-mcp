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

## hidden mounts

Overwrite sensitive directories and files with blank mounts

## enhance docs with agent guidance

Add to docs prompts and tips to get agents to use jail mcp tools more effectively
Add docs on how to obtain project files (probably clone the repo) or package in some way

# what not to add

- filesystem MCP tool — redundant, shell already does cat/ls/cp/find
- command allowlists — defeats the purpose, Docker is the boundary
- http/sse transport — only needed if running the server persistently and remotely
