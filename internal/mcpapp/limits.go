package mcpapp

import (
	"os"
	"strconv"
	"strings"
)

const (
	defaultMaxResults     = 500
	absoluteMaxResultsCap = 50000
	envMaxResults         = "FILE_DISCOVERY_MCP_MAX_RESULTS"
	envAbsoluteMax        = "FILE_DISCOVERY_MCP_ABSOLUTE_MAX"
)

func clampMaxResults(requested int) int {
	if requested <= 0 {
		requested = envIntDefault(envMaxResults, defaultMaxResults)
	}
	if requested < 1 {
		requested = 1
	}
	capLimit := envIntDefault(envAbsoluteMax, absoluteMaxResultsCap)
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
