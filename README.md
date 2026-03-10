# jail-mcp

MCP server providing shell access to clients, jailed in a container.

> **Running outside Docker is dangerous.**
> The server runs as root in a container that dies at session end.

| tool            | use case                          |
| --------------- | --------------------------------- |
| context         | project and environment discovery |
| exec sync       | run foreground commands           |
| exec background | run background jobs               |
| status          | pool job status                   |
| setup           | install project dependencies      |

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
The example configuration shows how to add read-only paths, for things you don't want to risk, like .git.
See environment section in docker-compose.yml to add global values to the container.

_Linux:_ consider using [rootless docker](https://docs.docker.com/engine/security/rootless)

**2. Pull image**

```bash
docker pull ghcr.io/rthomazel/jail-mcp:latest
```

**3. Wire up clients**

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

**4. Discovery and setup**

To discover projects, the agent can call the context tool.
It's expected that the language will be versioned using a `.tool-versions` file or similar for each project.
The container has only bash and python3, for basic scripting.
[Mise](https://mise.jdx.dev) is provided for language version management.
Have the agent run the setup tool in the project directory at the start of the session.
Setup will install the language and project dependencies.
We support half a dozen languages, including go, javascript (npm or yarn classic, autodetected) and python.
Check definitions on `handlers/setup.go`.

**4.1. Setup script**

For further project bootstraping, the setup tool will look for a `setup.sh` bash script and execute it.
There are a few locations we expect to find this file besides the project root.
Check definitions on `handlers/setup.go`.

## Logs

Logs are written in plain text to stderr.

## Dev

Check run script.
