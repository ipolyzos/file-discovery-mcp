# fd-mcp

A **stdio** [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server that exposes **read-only** filesystem discovery (similar to `fd` / `find`) with **path scoping**, **bounded results**, and **structured** tool responses—without arbitrary shell execution.

## Status

Early bootstrap: Python package layout, CI workflow, and placeholder CLI. The stdio MCP server and search tools will land on `develop` via focused PRs.

## Requirements

- Python 3.11 or newer
- Optional: [`fd`](https://github.com/sharkdp/fd) on `PATH` for faster searches when implemented

## Development

```bash
python3 -m venv .venv
source .venv/bin/activate   # or .venv\Scripts\activate on Windows
pip install -e ".[dev]"
ruff check src tests
pytest
```

Install as a console script after editable install: `fd-mcp` (currently exits with code 2 until the server is wired up).

## Goals

- Safe defaults: scoped roots, hard caps on matches and traversal depth where applicable.
- Portable behavior: prefer `fd` when available, fall back to a constrained `find`-style search.
- Clear errors and machine-readable output for tools and tests.

## License

Copyright 2026 Ioannis Polyzos. Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
