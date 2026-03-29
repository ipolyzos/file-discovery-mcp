// Package findfiles performs read-only filename glob discovery under a directory.
// It prefers the fd binary when available; otherwise it uses a filepath walk (no shell).
package findfiles

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Params holds inputs for a single search (root must already be absolute and allowed).
type Params struct {
	RootAbs    string
	Pattern    string
	MaxResults int
	MaxDepth   *int
}

// Result is the outcome of a search (machine-oriented; use for JSON tools).
type Result struct {
	Matches   []string `json:"matches"`
	Backend   string   `json:"backend"`
	Truncated bool     `json:"truncated"`
}

// Run executes the search with context cancellation forwarded to fd.
func Run(ctx context.Context, p Params) (Result, error) {
	if p.MaxResults < 1 {
		return Result{}, fmt.Errorf("findfiles: maxResults must be >= 1")
	}
	if strings.TrimSpace(p.Pattern) == "" {
		return Result{}, fmt.Errorf("findfiles: pattern is required")
	}
	root := filepath.Clean(p.RootAbs)
	if fdPath, err := exec.LookPath("fd"); err == nil {
		return runFD(ctx, fdPath, root, p.Pattern, p.MaxResults, p.MaxDepth)
	}
	return runWalk(root, p.Pattern, p.MaxResults, p.MaxDepth)
}

func runFD(ctx context.Context, fdExe, root, pattern string, maxResults int, maxDepth *int) (Result, error) {
	// argv only — never invoke a shell.
	args := []string{
		"--glob", pattern,
		"--max-results", strconv.Itoa(maxResults),
	}
	if maxDepth != nil {
		args = append(args, "--max-depth", strconv.Itoa(*maxDepth))
	}
	args = append(args, root)
	cmd := exec.CommandContext(ctx, fdExe, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("findfiles: fd: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	var lines []string
	s := bufio.NewScanner(&stdout)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" {
			lines = append(lines, filepath.Clean(line))
		}
	}
	if err := s.Err(); err != nil {
		return Result{}, fmt.Errorf("findfiles: read fd output: %w", err)
	}
	truncated := len(lines) >= maxResults
	return Result{Matches: lines, Backend: "fd", Truncated: truncated}, nil
}

func runWalk(root, pattern string, maxResults int, maxDepth *int) (Result, error) {
	var matches []string
	truncated := false

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if maxDepth != nil {
			sep := string(filepath.Separator)
			dc := strings.Count(rel, sep)
			if d.IsDir() && rel != "." && dc >= *maxDepth {
				return filepath.SkipDir
			}
			if !d.IsDir() && dc > *maxDepth {
				return nil
			}
		}
		if d.IsDir() {
			return nil
		}
		ok, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return fmt.Errorf("findfiles: bad pattern: %w", err)
		}
		if !ok {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		matches = append(matches, filepath.Clean(abs))
		if len(matches) >= maxResults {
			truncated = true
			return fs.SkipAll
		}
		return nil
	})
	if err != nil {
		return Result{}, err
	}
	return Result{Matches: matches, Backend: "walk", Truncated: truncated}, nil
}
