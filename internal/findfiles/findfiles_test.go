package findfiles

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestRun_rejectsEmptyPattern(t *testing.T) {
	_, err := Run(context.Background(), Params{RootAbs: t.TempDir(), Pattern: " ", MaxResults: 10})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_walkFindsByGlob(t *testing.T) {
	root := t.TempDir()
	mustFile(t, filepath.Join(root, "a.go"), "x")
	mustMk(t, filepath.Join(root, "nested"))
	mustFile(t, filepath.Join(root, "nested", "b.go"), "x")
	mustFile(t, filepath.Join(root, "skip.txt"), "x")

	// Force walk backend: PATH points at an empty dir so `fd` is not found.
	t.Setenv("PATH", t.TempDir())

	res, err := Run(context.Background(), Params{
		RootAbs:    root,
		Pattern:    "*.go",
		MaxResults: 50,
		MaxDepth:   nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Backend != "walk" {
		t.Fatalf("backend %q", res.Backend)
	}
	if len(res.Matches) != 2 {
		t.Fatalf("matches %#v", res.Matches)
	}
	sort.Strings(res.Matches)
}

func TestRun_walkRespectsMaxDepth(t *testing.T) {
	root := t.TempDir()
	mustFile(t, filepath.Join(root, "r.go"), "x")
	mustMk(t, filepath.Join(root, "one"))
	mustFile(t, filepath.Join(root, "one", "o.go"), "x")
	d2 := filepath.Join(root, "one", "two")
	mustMk(t, d2)
	mustFile(t, filepath.Join(d2, "deep.go"), "x")

	t.Setenv("PATH", t.TempDir())

	max := 1
	res, err := Run(context.Background(), Params{
		RootAbs:    root,
		Pattern:    "*.go",
		MaxResults: 50,
		MaxDepth:   &max,
	})
	if err != nil {
		t.Fatal(err)
	}
	names := make(map[string]struct{})
	for _, p := range res.Matches {
		names[filepath.Base(p)] = struct{}{}
	}
	if _, ok := names["r.go"]; !ok {
		t.Fatal("missing r.go")
	}
	if _, ok := names["o.go"]; !ok {
		t.Fatal("missing o.go")
	}
	if _, ok := names["deep.go"]; ok {
		t.Fatal("deep.go should not match maxDepth=1")
	}
}

func TestRun_truncatesAtCap(t *testing.T) {
	root := t.TempDir()
	for i := range 5 {
		mustFile(t, filepath.Join(root, fmt.Sprintf("f%d.txt", i)), "x")
	}
	t.Setenv("PATH", t.TempDir())

	res, err := Run(context.Background(), Params{
		RootAbs:    root,
		Pattern:    "*.txt",
		MaxResults: 3,
		MaxDepth:   nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Matches) != 3 || !res.Truncated {
		t.Fatalf("matches=%d truncated=%v", len(res.Matches), res.Truncated)
	}
}

func mustMk(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustFile(t *testing.T, p, content string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
