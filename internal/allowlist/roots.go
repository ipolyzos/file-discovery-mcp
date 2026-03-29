// Package allowlist enforces search roots from FD_MCP_ROOTS and path containment.
package allowlist

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const envRoots = "FD_MCP_ROOTS"

// LoadRoots returns absolute directory roots allowed for search.
// If FD_MCP_ROOTS is empty, the current working directory is the only root.
func LoadRoots() ([]string, error) {
	raw := strings.TrimSpace(os.Getenv(envRoots))
	if raw == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("allowlist: get working directory: %w", err)
		}
		abs, err := filepath.Abs(wd)
		if err != nil {
			return nil, fmt.Errorf("allowlist: abs cwd: %w", err)
		}
		return []string{filepath.Clean(abs)}, nil
	}
	parts := filepath.SplitList(raw)
	if len(parts) == 0 {
		return LoadRoots()
	}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		abs, err := filepath.Abs(filepath.Clean(p))
		if err != nil {
			return nil, fmt.Errorf("allowlist: root %q: %w", p, err)
		}
		out = append(out, filepath.Clean(abs))
	}
	if len(out) == 0 {
		return LoadRoots()
	}
	return out, nil
}

// ResolveSearchRoot turns user input into a clean absolute path.
func ResolveSearchRoot(userRoot string) (string, error) {
	userRoot = strings.TrimSpace(userRoot)
	if userRoot == "" {
		return "", fmt.Errorf("allowlist: root is required")
	}
	p := filepath.Clean(userRoot)
	if !filepath.IsAbs(p) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("allowlist: get working directory: %w", err)
		}
		p = filepath.Join(wd, p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("allowlist: abs root: %w", err)
	}
	return filepath.Clean(abs), nil
}

// ErrNotUnderRoot is returned when a path escapes configured roots.
type ErrNotUnderRoot struct {
	Path  string
	Roots []string
}

func (e *ErrNotUnderRoot) Error() string {
	return fmt.Sprintf("allowlist: path %q is not under allowed roots", e.Path)
}

// MustBeUnderRoots returns nil if absPath is contained in at least one root.
func MustBeUnderRoots(absPath string, roots []string) error {
	if len(roots) == 0 {
		return fmt.Errorf("allowlist: no roots configured")
	}
	for _, r := range roots {
		if pathWithinRoot(r, absPath) {
			return nil
		}
	}
	return &ErrNotUnderRoot{Path: absPath, Roots: append([]string(nil), roots...)}
}

func pathWithinRoot(root, target string) bool {
	root = filepath.Clean(root)
	target = filepath.Clean(target)
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
