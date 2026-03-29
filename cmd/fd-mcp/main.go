// Command fd-mcp is the fd-mcp stdio MCP server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/ipolyzos/fd-mcp/internal/mcpapp"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	srv := mcpapp.BuildServer()
	if err := srv.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "fd-mcp: %v\n", err)
		os.Exit(1)
	}
}
