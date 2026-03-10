# jail-mcp

MCP server providing shell access to clients, jailed in a container.

> **Running outside Docker is dangerous.**
> The server runs as root in a container that dies at session end.

## Setup

**1. Configure container**

Do not edit `docker-compose.sample.yml`.
Copy it to `docker-compose.yml` and edit that instead:

```bash
cp docker-compose.sample.yml docker-compose.yml
```

Update the volume paths to point to your real work.
The server discovers them dynamically.
Paths bind-mounted as volumes _can be modified in your machine_ which is what you want for the agent to work for you.
The example configurations shows how to add read-only paths, for things you don't want to risk, like .git.

_Linux:_ consider using [rootless docker](https://docs.docker.com/engine/security/rootless)

**1.1. Language versions**

The container provides [mise](https://mise.jdx.dev) for language version management.
Add a `.tool-versions` file or similar to each project.
Have the agent run `mise install` in the project directory.
`python3` and `pip` are available by default for agent's scripting needs.

**3. Build**

```bash
go mod tidy
./run docker-build
```

**3.1. Configuration**

See environment section in docker-compose.yml.
.env does not affect the containerized application, just local development.

**4. Wire up clients**

For Claude desktop, add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

_Linux:_ `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "jail-mcp": {
      // Linux: just "docker"
      "command": "/Applications/Docker.app/Contents/Resources/bin/docker",
      "args": [
        "compose",
        "-f",
        // Linux: /home/you
        "/Users/you/jail-mcp/docker-compose.yml",
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

Logs are written in plain text to stderr.

## Dev

Check run script.
