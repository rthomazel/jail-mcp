# jail-mcp — agent context

Call the `context` tool first. It returns mounted project paths, available tool versions, and timeout.

Read `doc/` for architecture and project documentation before making changes — it's faster than reading source.

## file access

All file reads and writes go through `exec_sync` shell commands.

```bash
# read
cat /projects/foo/bar.go

# write (new file or full rewrite)
cat > /path/file << 'HEREDOC'
...
HEREDOC

# edit
# read first, then rewrite with cat > or use sed -i
```

## running commands

```bash
golangci-lint run ./...
godotenv -f .env go run main.go
go build -o bin/server main.go
# Linux: docker is in PATH; macOS: use full path
docker compose -f docker-compose.yml build --build-arg VERSION="$(git rev-parse --short HEAD)"
```

## guidelines

- Run `go mod tidy` after any `go.mod` or dependency changes
- Run the formatter (`gofumpt`) as the last step after code changes
- Do not document obvious things
- Be minimalistic: give the right answer, avoid guessing or workarounds; if blocked, say so explicitly
- Avoid single-letter variable names unless scope is very small (receivers and loop vars are fine)
- Avoid multi-line `if` conditions with `samber/lo` functions
- When refactoring, minimize renames unless asked
- Add tests only when asked; focus on code that is complex or prone to bugs
- Write functions in call order — entry point first, then what it calls
- Do not start background jobs on your own; wait to be asked
- This is a jujutsu repo — do not make commits
