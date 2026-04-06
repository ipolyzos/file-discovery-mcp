# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

### Changed

- **Breaking:** Renamed the Go module to `github.com/ipolyzos/file-discovery-mcp` and the CLI entrypoint to `cmd/file-discovery-mcp` (binary name **`file-discovery-mcp`**).
- **Breaking:** Environment variables **`FD_MCP_*`** are replaced with **`FILE_DISCOVERY_MCP_*`** (`ROOTS`, `MAX_RESULTS`, `ABSOLUTE_MAX`).
- MCP server `Implementation.name` and user-facing descriptions now use **`file-discovery-mcp`**; stderr messages are prefixed with **`file-discovery-mcp:`**.

## [0.2.1] - 2026-03-30

### Added

- [Contributor Covenant 2.1](CODE_OF_CONDUCT.md) with reporting via GitHub issues; security issues cross-linked to [docs/security.md](docs/security.md). Linked from README, CONTRIBUTING, and docs index.

[0.2.1]: https://github.com/ipolyzos/file-discovery-mcp/releases/tag/v0.2.1

## [0.2.0] - 2026-03-29

### Added

- Documentation set: [client contract](docs/client-contract.md), [user guide](docs/user-guide.md), [using](docs/using.md), [examples](docs/examples.md), [troubleshooting](docs/troubleshooting.md), [security](docs/security.md), [deploying](docs/deploying.md), and [releases](docs/releases.md) (promotion and tagging process).
- README: quick start, security summary, and inline `find_files` example pointers.

[0.2.0]: https://github.com/ipolyzos/file-discovery-mcp/releases/tag/v0.2.0

## [0.1.0] - 2026-03-29

### Added

- Go module `github.com/ipolyzos/fd-mcp` and `cmd/fd-mcp` stdio MCP server using [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk).
- `get_server_info` tool (read-only metadata; no filesystem or shell access).
- GitHub Actions CI (build, test, `gofmt`, `go vet`).
- Project documentation and Apache 2.0 license.

[0.1.0]: https://github.com/ipolyzos/file-discovery-mcp/releases/tag/v0.1.0
