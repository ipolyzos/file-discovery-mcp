// Package mcpapp configures the fd-mcp Model Context Protocol server (tools only; stdio in cmd).
package mcpapp

import (
	"context"

	"github.com/ipolyzos/fd-mcp/internal/meta"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerInfo is the structured result of the get_server_info tool.
type ServerInfo struct {
	Version             string `json:"version" jsonschema:"Project version string"`
	Description         string `json:"description" jsonschema:"What this server provides"`
	ReadOnly            bool   `json:"read_only" jsonschema:"True: no writes or shell execution"`
	FilesystemSearchWIP bool   `json:"filesystem_search_wip" jsonschema:"True until search tools are implemented"`
}

// BuildServer returns an MCP server with read-only tools registered.
func BuildServer() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "fd-mcp",
		Version: meta.Version(),
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_server_info",
		Description: "Return fd-mcp version and capability flags. Does not access the filesystem.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ServerInfo, error) {
		return nil, ServerInfo{
			Version:             meta.Version(),
			Description:         "fd-mcp: read-only filesystem discovery MCP server",
			ReadOnly:            true,
			FilesystemSearchWIP: true,
		}, nil
	})

	return s
}
