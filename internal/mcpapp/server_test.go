package mcpapp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipolyzos/fd-mcp/internal/allowlist"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetServerInfoReturnsVersionAndReadOnly(t *testing.T) {
	ctx := context.Background()
	srv, err := BuildServer()
	if err != nil {
		t.Fatal(err)
	}

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "file-discovery-mcp-test", Version: "0"}, nil)
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
	if info.FilesystemSearchWIP {
		t.Fatal("expected FilesystemSearchWIP false when find_files exists")
	}
	if !strings.Contains(info.Description, "file-discovery-mcp") {
		t.Fatalf("Description %q should mention file-discovery-mcp", info.Description)
	}
}

func TestFindFilesTool_findsInAllowedRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.go"), []byte("p"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(allowlist.EnvRoots, root)
	t.Setenv("PATH", t.TempDir()) // force walk backend for stability in CI

	srv, err := BuildServer()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "file-discovery-mcp-test", Version: "0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "find_files",
		Arguments: map[string]any{
			"pattern": "*.go",
			"root":    root,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("tool error: %#v", res)
	}
	raw, _ := json.Marshal(res.StructuredContent)
	var out FindFilesOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	if !out.OK || len(out.Matches) != 1 {
		t.Fatalf("out=%+v", out)
	}
	if filepath.Base(out.Matches[0]) != "x.go" {
		t.Fatal(out.Matches)
	}
}

func TestFindFilesTool_rejectsPathOutsideRoots(t *testing.T) {
	allowed := t.TempDir()
	outside := t.TempDir()
	t.Setenv(allowlist.EnvRoots, allowed)

	srv, err := BuildServer()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "file-discovery-mcp-test", Version: "0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "find_files",
		Arguments: map[string]any{
			"pattern": "*.go",
			"root":    outside,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected protocol error: %#v", res)
	}
	raw, _ := json.Marshal(res.StructuredContent)
	var out FindFilesOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	if out.OK || out.Code != "PATH_NOT_ALLOWED" {
		t.Fatalf("out=%+v", out)
	}
}
