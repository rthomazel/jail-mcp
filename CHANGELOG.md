# CHANGELOG

## [0.1.0](https://github.com/rthomazel/jail-mcp/pull/2) feat(handlers): add exec background and status

### features

- **`exec_background`** — run long-running commands without blocking. returns a `job_id` immediately. background jobs have a separate timeout (`JAIL_MCP_BACKGROUND_TIMEOUT`, default 5m).
- **`exec_status`** — poll a background job for state, stdout, stderr, exit code, and duration.
- **`exec_sync`** — renamed from `shell_exec`.
- **context** — returns `shell_exec_timeout` so agents know the sync timeout upfront.

### improvements

- **dynamic mounts** — `context` no longer hardcodes `/projects`. reads `/proc/mounts` and reports all user-mounted volumes regardless of mount location.
- **Go tools** — `godotenv` and `gofumpt` pinned in `tools.go` and installed from the module graph at docker build time. tool versions stay in sync with the project.
- **Go runtime** — container installs Go 1.25 from upstream instead of the stale apt package.

### fixes

- logs go to stderr only; removed log file path from context output.
- fixed dockerfile GOPATH for alpine builder stage (`/go/bin` not `/root/go/bin`).

### docs

- `doc/ideas.md` — planned features: setup tool, language version management, concurrent context, per-command timeout, command stats.
- `doc/tools.md` — documents the `tools.go` pattern; projects bring their own tool versions, container picks them up at build time.
- `changelog.md` — this file.
