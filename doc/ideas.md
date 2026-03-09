# ideas

## background jobs

Long-running commands (big builds, test suites, slow installs) get killed by the timeout.
A `run_background` tool returning a job ID and a `job_status` tool to poll it would fix this.
Genuine gap — shell can't work around it.

## concurrent context

`context` tool runs subprocesses serially.
Could run them with goroutines and be meaningfully faster.

## per-command timeout

Timeout is global via `JAIL_MCP_TIMEOUT`.
Letting `shell_exec` accept an optional `timeout` param would be useful for known slow commands.

## what not to add

- filesystem MCP tool — redundant, shell already does cat/ls/cp/find
- command allowlists — defeats the purpose, Docker is the boundary
- http/sse transport — only needed if running the server persistently and remotely
