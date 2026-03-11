# architecture

## layout

```
main.go                       server wiring, tool registration, MCP server init
internal/config.go            Config struct, env var loading, defaults
handlers/handler.go           Handler struct, job store, startJob, background job GC
handlers/context.go           HandleContext, parseMounts
handlers/exec_sync.go         HandleExec, runCommand (shared by context)
handlers/exec_background.go   HandleExecBackground
handlers/status.go            HandleStatus
handlers/setup.go             HandleSetup, orderedRules, setupScriptCandidates
```

## request flow

All tools go through `mcp-go` → handler method → JSON response.

`exec_sync` and `context` run commands synchronously via `runCommand`, which wraps `bash -c` with a context timeout.

`exec_background` calls `startJob`, which spawns a goroutine, assigns a random 4-digit ID, and returns immediately. The caller polls with `status`.

`setup` detects the project's package manager by checking for known manifest files in order, builds a compound shell command (`&&`-joined), and launches it as a background job per path. If a `setup.sh` (or equivalent) is found it runs first.

`context` reads `/proc/mounts`, filters noise (proc/sysfs/tmpfs/overlay/etc.), deduplicates child mounts, and collects tool versions in parallel via `runCommand`.

## concurrency

`Handler.jobs` is a `map[string]*job` guarded by `Handler.mu` (RWMutex). Each job has its own `sync.Mutex` protecting its output buffers and done flag. A background goroutine sweeps completed jobs older than 1 hour every 5 minutes.

## configuration

Config is env-var only. See [config.md](config.md).

## design decisions

- No command filtering — the container is the security boundary, not the server
- `bash -c` gives agents pipes, redirects, `&&`, subshells
- `slog.SetDefault` at startup — no logger threaded through the codebase
- `internal/` for everything except `main.go`
- Both compose files use `image: jail-mcp` so builds and runs share the same tag
