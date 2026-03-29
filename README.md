# fd-mcp

## Purpose

**fd-mcp** will be a stdio [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server that exposes **read-only**, **scoped** filesystem discovery in the spirit of `fd` / `find`, with structured results and **no arbitrary shell execution**.

The codebase is **Go** (small static binary, strong stdlib, good fit for a stdio MCP server). Toolchain: **Go 1.23+** (see `go.mod`).

## Non-goals (for now)

- Arbitrary commands, shell integration, or write access to the filesystem
- Broad indexing, background daemons, or cloud sync

## Build locally

From the repository root:

```bash
go build -o fd-mcp ./cmd/fd-mcp
```

The binary exits with status `2` until the MCP server is wired up.

## Test locally

```bash
go test ./...
go vet ./...
```

Check formatting (same as CI):

```bash
test -z "$(gofmt -l .)"
```

## Development workflow

1. Branch from **`develop`** (`feature/*`, `chore/*`, etc.).
2. Run `gofmt`, `go vet`, and `go test ./...` before pushing.
3. Open a PR into **`develop`**; CI runs build, tests, `gofmt`, and `go vet`.

## License

Copyright 2026 Ioannis Polyzos. Licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE).

## Roadmap (high level)

1. Define safety model (roots, limits, error contract) and MCP tool surface.
2. Implement stdio MCP and filesystem discovery with tests.
3. Releases via `main` with version tags when the first usable build exists.

Remote: [github.com/ipolyzos/fd-mcp](https://github.com/ipolyzos/fd-mcp).
