# fd-mcp

## Purpose

**fd-mcp** will be a stdio [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server that exposes **read-only**, **scoped** filesystem discovery in the spirit of `fd` / `find`, with structured results and **no arbitrary shell execution**.

This repository is in **bootstrap** phase: documentation and licensing only; there is no runtime implementation here yet.

## Non-goals (for now)

- Shipping executable MCP server code or packaging
- Arbitrary commands, shell integration, or write access to the filesystem
- Broad indexing, background daemons, or cloud sync

## License

Copyright 2026 Ioannis Polyzos. Licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE).

## Roadmap (placeholder)

1. Define safety model (roots, limits, error contract) and protocol surface.
2. Add implementation (stdio MCP, tests, CI) on feature branches from `develop`.
3. Releases via `main` with version tags when the first usable build exists.

Remote: [github.com/ipolyzos/fd-mcp](https://github.com/ipolyzos/fd-mcp).
