# fd-mcp

A **stdio** [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server that exposes **read-only** filesystem discovery (similar to `fd` / `find`) with **path scoping**, **bounded results**, and **structured** tool responses—without arbitrary shell execution.

## Status

`develop` includes a working stdio server with a single tool, `search_files`, backed by `fd` when available or a pure Python directory walk otherwise.

## Requirements

- Python 3.11 or newer
- Optional: [`fd`](https://github.com/sharkdp/fd) on `PATH` for faster searches

## Running

After install (`pip install .` or editable `pip install -e .`):

```bash
fd-mcp
```

The process speaks MCP over stdin/stdout; configure your MCP client to launch this command.

### Example MCP client config

```json
{
  "mcpServers": {
    "fd-mcp": {
      "command": "fd-mcp",
      "env": {
        "FD_MCP_ROOTS": "/home/you/safe-root"
      }
    }
  }
}
```

On Windows, use `;` as the separator for multiple paths in `FD_MCP_ROOTS` (same as `PATH`).

## Configuration (environment)

| Variable | Meaning |
|----------|---------|
| `FD_MCP_ROOTS` | Allowed directories (OS path separator). Default: current working directory at startup. |
| `FD_MCP_MAX_RESULTS` | Default cap per search (default `1000`). |
| `FD_MCP_ABSOLUTE_MAX` | Hard ceiling for requested `max_results` (default `50000`). |
| `FD_MCP_MAX_DEPTH` | Optional default max directory depth for searches (omit for unlimited). |

Every `search_files` call must use a `root` path that resolves inside one of the allowed roots.

## Tool: `search_files`

Parameters:

- `pattern` (string): glob for the filename (e.g. `*.py`).
- `root` (string): directory to search under; must stay within `FD_MCP_ROOTS`.
- `max_results` (optional int): capped by `FD_MCP_ABSOLUTE_MAX`.
- `max_depth` (optional int): maximum directory depth under `root` (`0` is the root only).

Structured result (JSON), for example:

```json
{
  "ok": true,
  "matches": ["/path/to/file.py"],
  "truncated": false,
  "backend": "fd",
  "root": "/path/to/root",
  "pattern": "*.py",
  "max_results": 1000
}
```

On errors (invalid pattern, path outside roots, I/O failure), `ok` is `false` and `code` / `error` describe the problem.

## Development

```bash
python3 -m venv .venv
source .venv/bin/activate   # or .venv\Scripts\activate on Windows
pip install -e ".[dev]"
ruff check src tests
pytest
```

Console script after editable install: `fd-mcp`.

## Goals

- Safe defaults: scoped roots, hard caps on matches and traversal depth where applicable.
- Portable behavior: prefer `fd` when available, fall back to a safe in-process walk (no shell).
- Clear errors and machine-readable output for tools and tests.

## License

Copyright 2026 Ioannis Polyzos. Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
