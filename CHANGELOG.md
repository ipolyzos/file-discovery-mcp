# Changelog

All notable changes to this project are documented in this file.

## [0.1.0] - 2026-03-29

### Added

- Go module and `cmd/fd-mcp` stdio MCP server using [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk).
- `get_server_info` tool (read-only metadata; no filesystem or shell access).
- GitHub Actions CI (build, test, `gofmt`, `go vet`).
- Project documentation and Apache 2.0 license.

### Known limitations

- Filesystem search tools are not implemented; `filesystem_search_wip` remains true in tool output.

[0.1.0]: https://github.com/ipolyzos/fd-mcp/releases/tag/v0.1.0
