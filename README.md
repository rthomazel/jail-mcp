# jail-mcp

MCP server providing shell access to clients, jailed in a container.

> **Running outside Docker is dangerous.**
The server runs as root in a container that dies at session end.

## Setup

**1. Build**

```bash
go mod tidy
./run docker-build
```

**2. Configure your mounts**

`docker-compose.yml` is a sample file — do not edit it.
Copy it to `docker-compose.user.yml` and edit that instead:

```bash
cp docker-compose.yml docker-compose.user.yml
```

Update the volume paths to point to your real projects.
Paths bind-mounted as volumes _can be modified in your machine_ which is what you want for the agent to work for you.
The example configurations shows how to add read-only paths, for things you don't want to risk, like .git.

**2.1. Configuration**

See environment section in docker-compose.yml.

**3. Wire up clients**

For Claude desktop, add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "jail-mcp": {
      "command": "/Applications/Docker.app/Contents/Resources/bin/docker",
      "args": [
        "compose",
        "-f",
        "/Users/you/jail-mcp/docker-compose.user.yml",
        "run",
        "--rm",
        "-i",
        "jail-mcp"
      ]
    }
  }
}
```

Restart client.

## Logs

Logs are written in plain text to `JAIL_MCP_LOG_FILE` and to stderr.

## Dev

Check run script.