package main

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/ipolyzos/file-discovery-mcp/internal/mcpapp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Regression: mcp.Server.Run must report context.Canceled when ctx is canceled,
// so main can exit 0 on SIGINT (see runExitCode).
func TestServerRun_returnsCanceledWhenContextAlreadyCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	rDown, wDown := io.Pipe()
	rUp, wUp := io.Pipe()
	t.Cleanup(func() {
		_ = rDown.Close()
		_ = wDown.Close()
		_ = rUp.Close()
		_ = wUp.Close()
	})

	tr := &mcp.IOTransport{Reader: rUp, Writer: wDown}
	srv, err := mcpapp.BuildServer()
	if err != nil {
		t.Fatal(err)
	}
	err = srv.Run(ctx, tr)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run: got %v; want Canceled", err)
	}
}
