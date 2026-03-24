# configuration

Config is loaded from environment variables only — no flags, no config files.

| variable                      | default | description                                  |
| ----------------------------- | ------- | -------------------------------------------- |
| `JAIL_MCP_TIMEOUT`            | `15s`   | Timeout for `exec_sync` commands             |
| `JAIL_MCP_BACKGROUND_TIMEOUT` | `5m`    | Timeout for `exec_background` / `setup` jobs |

Values must be valid Go duration strings (e.g. `30s`, `2m`, `1h`).

Set these in the `environment:` section of `docker-compose.yml`.
