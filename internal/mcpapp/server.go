// Package mcpapp configures the file-discovery-mcp Model Context Protocol server (tools only; stdio in cmd).
package mcpapp

import (
	"context"
	"fmt"

	"github.com/ipolyzos/file-discovery-mcp/internal/allowlist"
	"github.com/ipolyzos/file-discovery-mcp/internal/meta"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerInfo is the structured result of the get_server_info tool.
type ServerInfo struct {
	Version             string `json:"version" jsonschema:"Project version string"`
	Description         string `json:"description" jsonschema:"What this server provides"`
	ReadOnly            bool   `json:"read_only" jsonschema:"True: no writes or shell execution"`
	FilesystemSearchWIP bool   `json:"filesystem_search_wip" jsonschema:"False when find_files is available"`
}

// BuildServer returns an MCP server with read-only tools registered.
func BuildServer() (*mcp.Server, error) {
	roots, err := allowlist.LoadRoots()
	if err != nil {
		return nil, fmt.Errorf("mcpapp: %w", err)
	}

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "file-discovery-mcp",
		Version: meta.Version(),
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_server_info",
		Description: "Return file-discovery-mcp version and capability flags. Does not access the filesystem.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, ServerInfo, error) {
		return nil, ServerInfo{
			Version:             meta.Version(),
			Description:         "file-discovery-mcp: read-only filesystem discovery MCP server",
			ReadOnly:            true,
			FilesystemSearchWIP: false,
		}, nil
	})

	addFindFilesTool(s, roots)

	return s, nil
}
