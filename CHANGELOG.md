# CHANGELOG

## [0.2.2](https://github.com/rthomazel/jail-mcp/pull/5) feat: hidden paths & context improvements

### feat

- [`8861fe9`](https://github.com/rthomazel/jail-mcp/commit/8861fe9) **(docker)** hidden paths — volume mounts can shadow files or directories inside a project to hide them from the agent. mount `/dev/null` over a file or an empty host directory over a subdirectory; mount order in the compose file determines precedence.
- [`2acf59e`](https://github.com/rthomazel/jail-mcp/commit/2acf59e) **(run)** add publish script — `run` gains a `publish` command for tagging and pushing new builds.
- [`e651fc1`](https://github.com/rthomazel/jail-mcp/commit/e651fc1) **(handlers/context)** mise shims block, path — `context` now reports the mise shims directory and includes it in the returned `path`, so agents can verify tool availability without manual inspection.

### fix

- [`89528dd`](https://github.com/rthomazel/jail-mcp/commit/89528dd) **(mise)** add shims to path at server start — mise shims are injected into `PATH` when the server starts, making agent-installed tools immediately available to all commands.

### refactor

- [`b58eeec`](https://github.com/rthomazel/jail-mcp/commit/b58eeec) **(handlers)** builder WriteString — handler output construction switched to `WriteString` for consistency.
- [`c32f73c`](https://github.com/rthomazel/jail-mcp/commit/c32f73c) **(handlers)** plain text output — handlers emit plain text instead of structured markup, simplifying agent parsing.

### build

- [`e8d40bc`](https://github.com/rthomazel/jail-mcp/commit/e8d40bc) improve version string — version now embeds commit hash and build timestamp for clearer traceability.

### docs

- [`7e709ad`](https://github.com/rthomazel/jail-mcp/commit/7e709ad) update docs — README updated with hidden mounts documentation and known bugs section covering tool discovery failure after server update.
- [`571cae0`](https://github.com/rthomazel/jail-mcp/commit/571cae0) update ideas.

### misc

- [`0c1faf9`](https://github.com/rthomazel/jail-mcp/commit/0c1faf9) change module name.
- [`cb7a321`](https://github.com/rthomazel/jail-mcp/commit/cb7a321) **(context)** remove go and node from response — trimmed unused fields to reduce noise.

## [0.2.1](https://github.com/rthomazel/jail-mcp/pull/4) feat: http mode & runtime improvements

### features

- **http mode** — new optional HTTP/SSE transport alongside stdio. `bin/jailmcphttp` helper script and `docker-compose-http-sample.yml` added for running the server over HTTP.
- **jujutsu in container** — `jj` binary installed in the runtime image, enabling agents to use jj commands inside the container.
- **setup script support** — `setup` tool now discovers and sources a `setup.sh` (or equivalent at `setup`, `bin/setup`, `script/setup`, `scripts/setup`, `scripts/setup.sh`) before running manifest install commands.
- **version in context** — `context` tool now returns the server build version.

### fixes

- **dockerfile** — fixed jj archive extraction (member path is `./jj`, requires `--strip-components=1`).
- **setup race** — fixed a race condition between concurrent setup jobs sharing output buffers.
- **setup tag & command** — corrected `go install tool` invocation and run script image tagging.

### improvements

- **tool descriptions** — updated `exec_sync` and `exec_background` descriptions to better guide agents on when to use each.
- **CI** — added `release.yml` workflow to push the image to `ghcr.io` on tag push.
- **compose sample** — renamed `docker-compose.sample.yml` → `docker-compose-sample.yml` for consistency; updated to reference the locally built image.
- **docs** — added `doc/architecture.md`, `doc/config.md`; updated `CLAUDE.md` to be agent-directives-only.

## [0.2.0](https://github.com/rthomazel/jail-mcp/pull/3) feat(handlers): setup

### features

- **`setup`** — new tool that installs dependencies for given project paths in parallel. detects supported manifests (`.tool-versions`, `go.mod`, `yarn.lock`, `package.json`, `requirements.txt`, `pyproject.toml`, `Gemfile`, `Cargo.toml`, `mix.exs`) and runs the appropriate install commands.

### improvements

- **Go tool directive** — migrated from `tools.go` pattern to Go 1.24 `tool` directive in `go.mod`. tools are now declared with `go get -tool` and run via `go tool`. `tools.go` removed.
- **`go install tool`** — setup appends `go install tool` to the `go mod download` step. safe when no tools are declared (exits 0 with a warning).

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
- **CI** — `pr.yml` workflow runs build, test, lint, and go mod tidy check on every PR to main.

### fixes

- logs go to stderr only; removed log file path from context output.
- fixed dockerfile GOPATH for alpine builder stage (`/go/bin` not `/root/go/bin`).

### docs

- `doc/ideas.md` — planned features: setup tool, language version management, concurrent context, per-command timeout, command stats.
- `doc/tools.md` — documents the `tools.go` pattern; projects bring their own tool versions, container picks them up at build time.
- `changelog.md` — this file.
