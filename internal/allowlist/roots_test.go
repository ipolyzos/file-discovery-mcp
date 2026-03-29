package allowlist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSearchRoot_rejectsEmpty(t *testing.T) {
	_, err := ResolveSearchRoot("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMustBeUnderRoots_acceptsInside(t *testing.T) {
	root := filepath.Join(t.TempDir(), "r")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	roots := []string{root}
	if err := MustBeUnderRoots(sub, roots); err != nil {
		t.Fatal(err)
	}
}

func TestMustBeUnderRoots_rejectsParentEscape(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "allowed")
	other := filepath.Join(base, "other")
	for _, p := range []string{root, other} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// From within root, resolved path could point to sibling.
	chdir := filepath.Join(root, "d")
	if err := os.MkdirAll(chdir, 0o755); err != nil {
		t.Fatal(err)
	}
	old, _ := os.Getwd()
	if err := os.Chdir(chdir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	absOther, err := filepath.Abs(other)
	if err != nil {
		t.Fatal(err)
	}
	if err := MustBeUnderRoots(absOther, []string{root}); err == nil {
		t.Fatal("expected not under root")
	}
}

func TestMustBeUnderRoots_rootItself(t *testing.T) {
	root := t.TempDir()
	if err := MustBeUnderRoots(root, []string{root}); err != nil {
		t.Fatal(err)
	}
}

func TestLoadRoots_fromEnvSplitList(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()
	aa, err := filepath.Abs(a)
	if err != nil {
		t.Fatal(err)
	}
	bb, err := filepath.Abs(b)
	if err != nil {
		t.Fatal(err)
	}
	sep := string(os.PathListSeparator)
	t.Setenv(envRoots, aa+sep+bb)
	roots, err := LoadRoots()
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 2 {
		t.Fatalf("got %d roots", len(roots))
	}
	if filepath.Clean(roots[0]) != filepath.Clean(aa) || filepath.Clean(roots[1]) != filepath.Clean(bb) {
		t.Fatalf("roots %v want %s and %s", roots, aa, bb)
	}
}
