# jail-mcp

An MCP server written in Go that gives an AI full shell access inside a Docker container. No command filtering, no allowlists — isolation is handled entirely by Docker. The container has a generous Ubuntu-based toolchain so the AI can actually get things done.

## How it works

Claude Desktop spawns `docker compose run` as a child process and pipes stdio to it. The Go binary inside the container speaks the MCP protocol over that stdio connection. Two tools are exposed:

- `shell_exec` — run any shell command, returns stdout / stderr / exit_code / duration
- `list_dirs` — list the directories available as working dirs

## Setup

### 1. Build the image

```bash
cd ~/Desktop/jail-mcp
go mod tidy
docker compose build
```

### 2. Configure your volume mounts

Edit `docker-compose.yml`. Replace the placeholder paths with your real ones:

```yaml
volumes:
  - /Users/you/myproject:/workspace        # rw — AI works here
  - /Users/you/myproject/.git:/workspace/.git:ro  # ro — .git is protected
  - /Users/you/.jail-mcp-logs:/var/log/jail-mcp   # logs survive restarts
```

The `.git` read-only mount stacks on top of the writable workspace mount. The AI can read it (`git log`, `git status`) but cannot delete or modify it.

To expose additional directories, add more volume entries and update `JAIL_MCP_DIRS` to include the container-side paths:

```yaml
environment:
  JAIL_MCP_DIRS: /workspace:/data
volumes:
  - /Users/you/myproject:/workspace
  - /Users/you/somedata:/data
```

### 3. Wire up Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "jail-mcp": {
      "command": "docker",
      "args": [
        "compose",
        "-f", "/Users/you/Desktop/jail-mcp/docker-compose.yml",
        "run", "--rm", "-i",
        "jail-mcp"
      ]
    }
  }
}
```

Restart Claude Desktop.

## Configuration

All configuration is via environment variables set in `docker-compose.yml`. No config files, no CLI flags.

| Variable          | Required | Default                        | Description                                              |
|-------------------|----------|--------------------------------|----------------------------------------------------------|
| `JAIL_MCP_DIRS`   | yes      | —                              | Colon-separated list of dirs the AI can use as cwd, e.g. `/workspace:/data` |
| `JAIL_MCP_TIMEOUT`| no       | `30s`                          | Max command execution time. Any Go duration: `60s`, `5m` |
| `JAIL_MCP_LOG`    | no       | `/var/log/jail-mcp/jail.log`   | Log file path inside the container                        |

## Logs

Logs are written in plain text to the file specified by `JAIL_MCP_LOG`. Bind-mounting a local directory to `/var/log/jail-mcp` keeps them after the container exits.

```
time=2026-03-05T14:32:01Z level=INFO msg="exec start" cmd="go build ./..." cwd=/workspace
time=2026-03-05T14:32:03Z level=INFO msg="exec done" cmd="go build ./..." exit_code=0 duration=1.82s
```

To tail live:

```bash
tail -f ~/.jail-mcp-logs/jail.log
```

## Security model

There is none inside the server — that's intentional. The container is the jail. If the AI does something destructive inside the container, that's contained. The only things that can be affected on your host are the directories you explicitly bind-mount.

The `.git` read-only mount is the one meaningful host-side protection: even a `rm -rf /workspace` inside the container cannot touch your git history.
