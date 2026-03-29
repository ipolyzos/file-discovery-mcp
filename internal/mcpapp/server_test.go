package mcpapp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetServerInfoReturnsVersionAndReadOnly(t *testing.T) {
	ctx := context.Background()
	srv := BuildServer()

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "fd-mcp-test", Version: "0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "get_server_info"})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool returned error: %#v", res)
	}

	if res.StructuredContent == nil {
		t.Fatal("expected StructuredContent from typed tool output")
	}

	raw, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var info ServerInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if info.Version == "" {
		t.Fatal("empty version")
	}
	if !info.ReadOnly {
		t.Fatal("expected ReadOnly true")
	}
	if !info.FilesystemSearchWIP {
		t.Fatal("expected FilesystemSearchWIP true until search ships")
	}
}
