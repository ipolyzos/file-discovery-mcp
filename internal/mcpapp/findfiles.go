package mcpapp

import (
	"context"
	"strings"

	"github.com/ipolyzos/fd-mcp/internal/allowlist"
	"github.com/ipolyzos/fd-mcp/internal/findfiles"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindFilesInput is the JSON input for the find_files tool.
type FindFilesInput struct {
	Pattern    string `json:"pattern" jsonschema:"Filename glob (e.g. *.go). Passed to fd --glob or filepath.Match."`
	Root       string `json:"root" jsonschema:"Directory to search; must fall under FD_MCP_ROOTS (default: process working directory)."`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"Max paths to return; capped by FD_MCP_ABSOLUTE_MAX. 0 uses FD_MCP_MAX_RESULTS default."`
	MaxDepth   *int   `json:"max_depth,omitempty" jsonschema:"Max directory depth under root. Omitted means unlimited."`
}

// FindFilesOutput is the structured JSON response for find_files.
type FindFilesOutput struct {
	OK         bool     `json:"ok"`
	Matches    []string `json:"matches,omitempty"`
	Truncated  bool     `json:"truncated,omitempty"`
	Backend    string   `json:"backend,omitempty"`
	Root       string   `json:"root,omitempty"`
	Pattern    string   `json:"pattern,omitempty"`
	MaxResults int      `json:"max_results,omitempty"`
	Code       string   `json:"code,omitempty"`
	Error      string   `json:"error,omitempty"`
}

func addFindFilesTool(s *mcp.Server, roots []string) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "find_files",
		Description: "List files under an allowed root matching a filename glob. Read-only; uses fd if available, else filepath walk. No shell.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in FindFilesInput) (*mcp.CallToolResult, FindFilesOutput, error) {
		pattern := strings.TrimSpace(in.Pattern)
		if pattern == "" {
			return nil, errOut("INVALID_PATTERN", "pattern is required"), nil
		}
		rootRaw := strings.TrimSpace(in.Root)
		if rootRaw == "" {
			return nil, errOut("INVALID_ROOT", "root is required"), nil
		}
		rootAbs, err := allowlist.ResolveSearchRoot(rootRaw)
		if err != nil {
			return nil, errOut("INVALID_ROOT", err.Error()), nil
		}
		if err := allowlist.MustBeUnderRoots(rootAbs, roots); err != nil {
			return nil, errOut("PATH_NOT_ALLOWED", err.Error()), nil
		}
		max := clampMaxResults(in.MaxResults)
		if in.MaxDepth != nil && *in.MaxDepth < 0 {
			return nil, errOut("INVALID_DEPTH", "max_depth must be >= 0"), nil
		}
		res, err := findfiles.Run(ctx, findfiles.Params{
			RootAbs:    rootAbs,
			Pattern:    pattern,
			MaxResults: max,
			MaxDepth:   in.MaxDepth,
		})
		if err != nil {
			return nil, errOut("SEARCH_FAILED", err.Error()), nil
		}
		return nil, FindFilesOutput{
			OK:         true,
			Matches:    res.Matches,
			Truncated:  res.Truncated,
			Backend:    res.Backend,
			Root:       rootAbs,
			Pattern:    pattern,
			MaxResults: max,
		}, nil
	})
}

func errOut(code, msg string) FindFilesOutput {
	return FindFilesOutput{OK: false, Code: code, Error: msg}
}
