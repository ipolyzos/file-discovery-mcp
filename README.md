# fd-mcp

A **stdio** [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server that exposes **read-only** filesystem discovery (similar to `fd` / `find`) with **path scoping**, **bounded results**, and **structured** tool responses—without arbitrary shell execution.

## Status

Early bootstrap: repository metadata and branching only. Implementation, CI, and releases will land on `develop` via focused PRs.

## Goals

- Safe defaults: scoped roots, hard caps on matches and traversal depth where applicable.
- Portable behavior: prefer `fd` when available, fall back to a constrained `find`-style search.
- Clear errors and machine-readable output for tools and tests.

## License

Copyright 2026 Ioannis Polyzos. Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
