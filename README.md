# jail-mcp

MCP server providing shell access to clients, jailed in a container.

## Setup

**1. Build**

```bash
go mod tidy
docker compose build
```

**2. Configure volume mounts**

Edit `docker-compose.yml` and replace the placeholder paths:

```yaml
volumes:
  - /Users/you/myproject:/workspace
  - /Users/you/myproject/.git:/workspace/.git:ro
  - /Users/you/.jail-mcp-logs:/var/log/jail-mcp
```

The `.git` mount is read-only — the AI can read it but cannot modify or delete it.

To expose additional directories add more volume entries and update `JAIL_MCP_DIRS`:

```yaml
environment:
  JAIL_MCP_DIRS: /workspace:/data
volumes:
  - /Users/you/myproject:/workspace
  - /Users/you/somedata:/data
```

**3. Wire up Claude Desktop**

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "jail-mcp": {
      "command": "docker",
      "args": ["compose", "-f", "/Users/you/Desktop/jail-mcp/docker-compose.yml", "run", "--rm", "-i", "jail-mcp"]
    }
  }
}
```

Restart Claude Desktop.

## Configuration

| Variable           | Required | Default                      |
|--------------------|----------|------------------------------|
| `JAIL_MCP_DIRS`    | yes      | —                            |
| `JAIL_MCP_TIMEOUT` | no       | `30s`                        |
| `JAIL_MCP_LOG`     | no       | `/var/log/jail-mcp/jail.log` |

`JAIL_MCP_DIRS` is a colon-separated list of dirs the AI can use as cwd, e.g. `/workspace:/data`.

## Logs

Logs are written in plain text to `JAIL_MCP_LOG` and teed to stderr.

```
time=2026-03-05T14:32:01Z level=INFO msg="exec start" cmd="go build ./..." cwd=/workspace
time=2026-03-05T14:32:03Z level=INFO msg="exec done" cmd="go build ./..." exit_code=0 duration=1.82s
```

```bash
tail -f ~/.jail-mcp-logs/jail.log
```
