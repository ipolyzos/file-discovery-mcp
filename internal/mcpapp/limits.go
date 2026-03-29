package mcpapp

import (
	"os"
	"strconv"
	"strings"
)

const (
	defaultMaxResults     = 500
	absoluteMaxResultsCap = 50000
)

func clampMaxResults(requested int) int {
	if requested <= 0 {
		requested = envIntDefault("FD_MCP_MAX_RESULTS", defaultMaxResults)
	}
	if requested < 1 {
		requested = 1
	}
	capLimit := envIntDefault("FD_MCP_ABSOLUTE_MAX", absoluteMaxResultsCap)
	if capLimit < 1 {
		capLimit = absoluteMaxResultsCap
	}
	if requested > capLimit {
		return capLimit
	}
	return requested
}

func envIntDefault(key string, defaultVal int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return defaultVal
	}
	return n
}
