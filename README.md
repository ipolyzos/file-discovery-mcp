# fd-mcp

## Purpose

**fd-mcp** is a stdio [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server that will expose **read-only**, **scoped** filesystem discovery in the spirit of `fd` / `find`, with structured results and **no arbitrary shell execution**.

**Stack:** Go 1.25+ (see `go.mod`), [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) for MCP.

## Current behavior

- The process speaks MCP over **stdin/stdout** (`fd-mcp` with no arguments).
- Interrupt (**SIGINT**) stops the server cleanly.
- Tool **`get_server_info`**: returns version and capability flags only (**no filesystem access**). Structured JSON includes `read_only: true` and `filesystem_search_wip: true` until search tools exist.
- **Filesystem search tools** are not implemented yet; the server does not run `fd`, `find`, or a shell.

## Non-goals (for now)

- Arbitrary commands, shell integration, or write access to the filesystem
- Broad indexing, background daemons, or cloud sync

## Build locally

```bash
go build -o fd-mcp ./cmd/fd-mcp
```

Run (for MCP clients that launch the binary; do not attach a TTY that interferes with stdio):

```bash
./fd-mcp
```

## Test locally

```bash
go test ./...
go vet ./...
test -z "$(gofmt -l .)"
```

## Development workflow

1. Branch from **`develop`** (`feature/*`, `chore/*`, etc.).
2. Run `gofmt`, `go vet`, and `go test ./...` before pushing.
3. Open a PR into **`develop`**; CI runs build, tests, `gofmt`, and `go vet`.

## License

Copyright 2026 Ioannis Polyzos. Licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE).

## Roadmap (high level)

1. Define safety model (roots, limits, error contract) and MCP tool surface for search.
2. Implement read-only discovery (prefer `fd` or in-process walk; **never** pass user input to a shell).
3. Releases via `main` with version tags when the first usable build exists.

Remote: [github.com/ipolyzos/fd-mcp](https://github.com/ipolyzos/fd-mcp).
