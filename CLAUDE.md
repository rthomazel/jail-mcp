# jail-mcp — agent context

Call `context` tool first. It returns mounted volume paths, available tools, and the log file location.

## Code layout

```
main.go                  server wiring, tool registration
internal/config.go       Config struct (Timeout), env var loading, defaults
internal/handler.go      HandleContext, HandleExec, runCommand (all in one file)
```

## Design decisions worth preserving

- No command filtering — container is the security boundary, not the server
- `bash -c` so the AI gets pipes, redirects, `&&`, subshells
- `slog.SetDefault` at startup — no logger passed around anywhere
- Config from env vars only, no flags, no config files
- `internal/` for everything except `main.go` so `go run main.go` works
- `image: jail-mcp` in both compose files so builds and runs share the same image tag

## Running commands (agents can't execute the run script directly)

```bash
golangci-lint run ./...
/Users/user/go/bin/godotenv -f .env go run main.go
go build -o bin/server main.go
# on linux, docker will be in PATH
/Applications/Docker.app/Contents/Resources/bin/docker compose -f docker-compose.yml build --build-arg VERSION="$(git rev-parse --short HEAD)"
```

## Module

`github.com/tcodes0/jail-mcp` — Go 1.25 — single dep: `github.com/mark3labs/mcp-go v0.18.0`

## Other Guidelines

go: run go mod tidy after making changes to go.mod and dependencies.
do not document obvious things
be more minimalistic: being helpful is good but we need to right answer, avoid guessing or crazy workarounds, if you are blocked, be explicit.
avoid single letter vars if their scope is not small; go: receivers, loop vars are an exception.
go: avoid multi line if conditions with samber/lo functions.
when we refactor, minimize renames unless asked for.
add tests when asked for; look for code that is complex or prone to change/ bugs; if tests never break they add no value.
run formatter as last step after making code changes.