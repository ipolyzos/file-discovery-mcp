package main

import (
	"context"
	"errors"
	"testing"
)

func TestRunExitCode_nilMeansSuccess(t *testing.T) {
	if got := runExitCode(nil); got != 0 {
		t.Fatalf("runExitCode(nil) = %d; want 0", got)
	}
}

func TestRunExitCode_contextCanceledMeansSuccess(t *testing.T) {
	if got := runExitCode(context.Canceled); got != 0 {
		t.Fatalf("runExitCode(Canceled) = %d; want 0", got)
	}
}

func TestRunExitCode_wrappedCancelMeansSuccess(t *testing.T) {
	wrapped := errors.Join(errors.New("x"), context.Canceled)
	if got := runExitCode(wrapped); got != 0 {
		t.Fatalf("runExitCode(wrapped Canceled) = %d; want 0", got)
	}
}

func TestRunExitCode_otherErrorMeansFailure(t *testing.T) {
	if got := runExitCode(errors.New("boom")); got != 1 {
		t.Fatalf("runExitCode(boom) = %d; want 1", got)
	}
}
