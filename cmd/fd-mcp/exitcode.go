package main

import (
	"context"
	"errors"
)

// runExitCode maps errors from [mcp.Server.Run] to a process exit status.
// Context cancellation (including SIGINT via [signal.NotifyContext]) is treated as success.
func runExitCode(err error) int {
	if err == nil {
		return 0
	}
	if errors.Is(err, context.Canceled) {
		return 0
	}
	return 1
}
