# fd-mcp

## Purpose

**fd-mcp** is a stdio [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server for **read-only**, **scoped** filesystem discovery (similar to `fd` / `find`), with structured JSON and **no shell**.

**Stack:** Go 1.25+ (see `go.mod`), [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) for MCP.

**Current release:** v0.2.1 — see [CHANGELOG.md](CHANGELOG.md).

## Quick start

1. **Build** the binary: `go build -o fd-mcp ./cmd/fd-mcp` (or use a published release).
2. **Register** `fd-mcp` in your MCP host as a **stdio** server (command = path to the binary; set **cwd** and optional **`FD_MCP_ROOTS`** if you need specific search boundaries).
3. **Use** the assistant’s tools **`get_server_info`** and **`find_files`** as exposed by the host.

Step-by-step orientation, example searches (by name, by extension), what to expect in results, and limitations: **[docs/user-guide.md](docs/user-guide.md)**.

**Examples:** call **`find_files`** with `{"pattern": "*.go", "root": "/path/to/repo"}` for all Go files under that directory; add `"max_depth": 2` and `"max_results": 50` to cap depth and list size. Sample MCP configs, success/failure JSON, and edge cases: **[docs/examples.md](docs/examples.md)**.

## Documentation

Extended guides and outlines live under [**docs/**](docs/). Start at [**docs/index.md**](docs/index.md) for the full map (architecture, development, security, client contract, and more).

## Architecture

fd-mcp is a **stdio MCP server** with a thin MCP layer (`internal/mcpapp`), a **safety layer** (`internal/allowlist`, env limits), and a **read-only discovery layer** (`internal/findfiles`: optional `fd` subprocess with fixed argv, else in-process walk). For the full component model, request/response flow, principles, and non-goals, see [**docs/architecture.md**](docs/architecture.md).

## Development

Clone, build, test, run the stdio server, branching conventions, commit habits, and what **GitHub Actions** runs: see [**docs/development.md**](docs/development.md).

## Contributing

Branch naming, commits, pull requests, code and documentation expectations, and how to report issues: see [**CONTRIBUTING.md**](CONTRIBUTING.md) (full guide: [**docs/contributing.md**](docs/contributing.md)). Community standards: [**CODE_OF_CONDUCT.md**](CODE_OF_CONDUCT.md).

## Client contract

Normative MCP tool names, input/output JSON fields, error codes, safety guarantees, and versioning expectations for integrators: see [**docs/client-contract.md**](docs/client-contract.md).

## Security

Threat model (command execution and path exposure), implemented mitigations, safe operational defaults, residual risks, and **responsible disclosure**: [**docs/security.md**](docs/security.md).

## Implemented functionality

- **Transport:** MCP over **stdin/stdout** (`fd-mcp` with no arguments).
- **Shutdown:** **SIGINT** → clean exit status **0** (see `runExitCode` in `cmd/fd-mcp`).
- **`get_server_info`:** version and flags (no filesystem I/O).
- **`find_files`:** list files under an allowed directory whose **basename** matches a **glob** (e.g. `*.go`).
  - **Backend:** uses **`fd`** on `PATH` when available (`exec.Command`, fixed argv only); otherwise a **`filepath.WalkDir`** implementation. **No** `sh -c`, no user-controlled executables.
  - **Limits:** `max_results` is clamped by **`FD_MCP_ABSOLUTE_MAX`** (default cap **50000**). Default per-request ceiling when omitted comes from **`FD_MCP_MAX_RESULTS`** (default **500**).
  - **Optional:** `max_depth` (non-negative) limits directory depth below `root`.

## Safety constraints

| Mechanism | Behavior |
|-----------|----------|
| **Root allowlist** | Env **`FD_MCP_ROOTS`**: list of directories (OS path-list separator, same as `PATH`). If unset, **current working directory** at process start is the only root. Every `find_files` **`root`** must resolve inside one of these directories (after `filepath.Clean` / abs); otherwise **`PATH_NOT_ALLOWED`**. |
| **Read-only** | Tools do not create, modify, or delete files. |
| **No arbitrary execution** | `fd` is invoked with a fixed argument shape; patterns are **not** passed through a shell. |
| **Bounded output** | Results are capped by **`max_results`** (after clamping). **`truncated: true`** when the cap may have cut off further matches. |

## Current limitations

- Glob is **filename / basename** only (not full-path patterns); semantics follow **`filepath.Match`** for the walk backend and **`fd --glob`** when `fd` is used (minor differences possible).
- Symlink-heavy trees may behave subtly vs. a fully canonicalized policy; containment checks use `filepath` semantics as documented in `internal/allowlist`.
- **MCP resource/prompt** surfaces are not implemented.

## Non-goals (still)

- Write access, arbitrary commands, interactive shell integration.
- Background indexers or network sync.

## Build locally

```bash
go build -o fd-mcp ./cmd/fd-mcp
```

Run (typical MCP client configuration: command = built binary, optional env):

```bash
./fd-mcp
```

Example **restrict roots** (Unix):

```bash
FD_MCP_ROOTS=/safe/project:/tmp/readonly ./fd-mcp
```

## Test locally

```bash
go test ./...
go vet ./...
test -z "$(gofmt -l .)"
```

## Development workflow

1. Branch from **`develop`** (`feature/*`, `bugfix/*`, `chore/*`, etc.).
2. Run `gofmt`, `go vet`, and `go test ./...` before pushing.
3. Open a PR into **`develop`**; CI runs build, tests, `gofmt`, and `go vet`.

## License

Copyright 2026 Ioannis Polyzos. Licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE).

## Roadmap (high level)

1. Tighten symlink / cross-platform edge cases; richer error codes if needed.
2. Optional depth-first streaming or paging for huge directories.
3. Ongoing releases via `main` with version tags.

Remote: [github.com/ipolyzos/fd-mcp](https://github.com/ipolyzos/fd-mcp).
