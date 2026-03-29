// Command fd-mcp is the stdio MCP server entrypoint (implementation pending).
package main

import (
	"fmt"
	"os"

	"github.com/ipolyzos/fd-mcp/internal/meta"
)

func main() {
	fmt.Fprintf(os.Stderr, "fd-mcp %s: server wiring not yet implemented\n", meta.Version())
	os.Exit(2)
}
