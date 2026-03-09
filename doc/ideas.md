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

## project setup tool

A dedicated `setup` tool that discovers and installs project dependencies for mounted volumes.
Scans each mounted path for known manifests (`go.mod`+`tools.go`, `package.json`) and runs
the appropriate install command as a background job per project. Jobs run in series to benefit
from shared tools across projects. Returns a `{path: job_id}` map. Errors only visible on status poll.

`context` could accept a `run_setup` param (default false) to fire setup jobs and include their
IDs in the response, for convenience without changing the default behavior.

## language version management

The Dockerfile currently pins language versions, which is a project/user concern.
The right answer is mise installed in the container. The Dockerfile provides the version manager,
not the language version.

Projects declare versions in `.tool-versions` at their root:

```
nodejs 22.0.0
go 1.25.0
python 3.12.0
```

The `setup` tool (see above) would run `mise install` in each mounted project directory.
mise reads `.tool-versions` and installs the declared versions.

Users add `.tool-versions` to their projects and use mise locally as well, keeping
local and container environments in sync.

## hidden mounts

overwrite sensitive directories and files with blank mounts

# what not to add

- filesystem MCP tool — redundant, shell already does cat/ls/cp/find
- command allowlists — defeats the purpose, Docker is the boundary
- http/sse transport — only needed if running the server persistently and remotely
