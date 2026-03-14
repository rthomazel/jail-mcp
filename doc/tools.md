# tools

Dev tools (formatters, linters, etc.) should be versioned alongside the project that needs them, not installed globally in the container.

## Go

Use a `tools.go` file with a `tools` build tag and blank imports — this pins versions in `go.mod` via the `tool` directive (Go 1.24+). The Dockerfile installs them in the builder stage and copies binaries into the runtime image.

### adding a Go tool

1. Add a blank import to `tools.go`
2. `go get <import-path>` + `go mod tidy`
3. Add `go install <import-path>` in the Dockerfile builder stage
4. Add `COPY --from=builder /root/go/bin/<tool> /usr/local/bin/<tool>` in the runtime stage
5. Rebuild the image

## Other languages

Same pattern: each project owns its tool versions and the container picks them up at build or setup time.

- JavaScript: `package.json` devDependencies
- Python: `Pipfile` or `pyproject.toml`
- etc.
