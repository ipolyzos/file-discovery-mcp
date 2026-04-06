package mcpapp

import "testing"

func TestClampMaxResults_usesFILE_DISCOVERY_MCP_env_defaults(t *testing.T) {
	t.Setenv("FILE_DISCOVERY_MCP_MAX_RESULTS", "10")
	t.Setenv("FILE_DISCOVERY_MCP_ABSOLUTE_MAX", "20")
	if got := clampMaxResults(0); got != 10 {
		t.Fatalf("clampMaxResults(0) = %d; want 10", got)
	}
	if got := clampMaxResults(100); got != 20 {
		t.Fatalf("clampMaxResults(100) = %d; want 20", got)
	}
}
