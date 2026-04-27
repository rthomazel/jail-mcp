# jail-mcp

MCP server providing shell access to clients, jailed in a container.

> **Running outside Docker is dangerous.**
> The server runs as root in a container.

| tool            | use case                          |
| --------------- | --------------------------------- |
| context         | project and environment discovery |
| exec sync       | run foreground commands           |
| exec background | run background jobs               |
| status          | pool job status                   |
| setup           | install project dependencies      |

## Setup

### Overview

- 1 write your compose file with projects as volumes
- 2 pull image
- 3 configure clients

**1. Configure container**

Two sample compose files are provided depending on your client:

| file                             | mode  | use case                     |
| -------------------------------- | ----- | ---------------------------- |
| `docker-compose-sample.yml`      | stdio | Claude Desktop, CLI clients  |
| `docker-compose-http-sample.yml` | HTTP  | Open WebUI, HTTP MCP clients |

Do not edit the sample files. Copy the one you need and edit that instead:

```bash
# stdio mode (Claude Desktop)
cp docker-compose-sample.yml docker-compose.yml

# HTTP mode (Open WebUI)
cp docker-compose-http-sample.yml docker-compose-http.yml
```

Update the volume paths to point to your real work.
The server discovers them dynamically, `/projects` is a suggestion.
Paths bind-mounted as volumes _can be modified in your machine_ which is what you want for the agent to work for you.
See environment section in docker-compose.yml to add global values to the container.

_Linux:_ consider using [rootless docker](https://docs.docker.com/engine/security/rootless)

#### Read-only paths

The example configuration shows how to add read-only paths, i.e `.git`.

```yaml
# :ro adds a path as read-only, must come after the parent path
- /Users/you/helloworld/.git:/projects/helloworld/.git:ro
```

#### Hidden mounts

Sensitive files or directories inside a mounted project can be hidden from the agent using Docker volume mounts — no server changes needed.

Docker applies mounts in declaration order. A second mount over a subpath of an already-mounted project shadows it before the container process starts. The container has no `CAP_SYS_ADMIN` so runtime mounts are not possible; this must be done in the compose file.

**Hide a file** — mount `/dev/null` over it:

```yaml
volumes:
  - /Users/you/myproject:/projects/myproject
  - /dev/null:/projects/myproject/.env
```

**Hide a directory** — mount an empty host directory over it:

```yaml
volumes:
  - /Users/you/myproject:/projects/myproject
  - /tmp/jail-hidden:/projects/myproject/secrets
```

The empty dir must exist on the host (`mkdir -p /tmp/jail-hidden`). Mount order matters — the hide entry must come after the parent project mount, same rule as `:ro` overlays.

**2. Pull image**

```bash
docker pull ghcr.io/rthomazel/jail-mcp:latest
```

**3. Wire up clients**

### Claude Desktop (stdio)

Spawns a fresh container per session via `docker compose run`.
The container is ephemeral — `--rm` removes it after each session. Only named volumes (`/mise`, `/root`) persist.
To install ad-hoc tools that survive across sessions, install to `$HOME/bin` (`/root/bin`), which is on the `jail-mcp-root` volume.

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

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
        "/Users/you/your-compose-file/docker-compose.yml",
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

### Open WebUI / HTTP clients

Runs a persistent container exposing an HTTP MCP endpoint on port 8001.

```bash
docker compose -f docker-compose-http.yml up -d
```

Then add `http://localhost:8001` as a tool server in your client.

The HTTP transport is configured via `JAIL_MCP_TRANSPORT` in the container environment — `mcpo` for OpenAI-compatible REST (Open WebUI) or `mcp-proxy` for native MCP/SSE (LibreChat, Claude Desktop). See `docker-compose-http-sample.yml` for an example.

#### Known (client) Bugs

When updating the MCP server to a new build, Claude desktop may show errors or fail to discover tools.
This has been observed to happen when changing permission settings as well.
This can be fixed by renaming the server in the configuration above (e.g. `jail-mcp` → `1_jail-mcp`), which forces the client to treat it as a new server and re-register the tools.
Renaming the first letter seems to be important.

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
