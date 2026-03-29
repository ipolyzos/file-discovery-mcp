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

	srv, err := mcpapp.BuildServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fd-mcp: %v\n", err)
		os.Exit(1)
	}
	err = srv.Run(ctx, &mcp.StdioTransport{})
	if code := runExitCode(err); code != 0 {
		fmt.Fprintf(os.Stderr, "fd-mcp: %v\n", err)
		os.Exit(code)
	}
}
