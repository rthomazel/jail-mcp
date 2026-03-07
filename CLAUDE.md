# jail-mcp — agent context

Call `context` tool first. It returns mounted project paths, available tools, and the log file location.

## Code layout

```
main.go                  server wiring, tool registration
internal/config.go       Config struct, env var loading, defaults
internal/handler.go      HandleContext, HandleExec, runCommand (all in one file)
internal/logger.go       slog.SetDefault, tees to file + stderr, returns io.Closer
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
