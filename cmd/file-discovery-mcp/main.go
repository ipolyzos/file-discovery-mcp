// Command file-discovery-mcp is the file-discovery-mcp stdio MCP server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/ipolyzos/file-discovery-mcp/internal/mcpapp"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	srv, err := mcpapp.BuildServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "file-discovery-mcp: %v\n", err)
		os.Exit(1)
	}
	err = srv.Run(ctx, &mcp.StdioTransport{})
	if code := runExitCode(err); code != 0 {
		fmt.Fprintf(os.Stderr, "file-discovery-mcp: %v\n", err)
		os.Exit(code)
	}
}
