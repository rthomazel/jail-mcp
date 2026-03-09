# tools

Dev tools (formatters, runners, etc.) should be versioned alongside the project that needs them, not installed globally in the container.

For Go projects, use a `tools.go` with a `tools` build tag and blank imports — this pins versions in `go.mod`. The Dockerfile installs them in the builder stage and copies the binaries into the runtime image.

Other languages have equivalents: `package.json` devDependencies, `Pipfile`, etc. The pattern is the same: each project owns its tool versions, and the container picks them up at build time.

## adding a Go tool

1. Add a blank import to `tools.go`
2. `go get <import-path>` + `go mod tidy`
3. Add `go install <import-path>` in the Dockerfile builder stage
4. Add `COPY --from=builder /root/go/bin/<tool> /usr/local/bin/<tool>` in the runtime stage
5. Rebuild
