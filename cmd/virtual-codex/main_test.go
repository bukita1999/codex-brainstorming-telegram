package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunPrintsVirtualCodexResponse(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(context.Background(), &stdout, &stderr, []string{"--input", "brainstorm login flow", "--delay", "0s"})
	if exitCode != 0 {
		t.Fatalf("run() exitCode = %d, stderr = %s", exitCode, stderr.String())
	}

	got := strings.TrimSpace(stdout.String())
	want := "VirtualCodex: 好的，我们先收敛需求。请回复 1/2/3：1) 目标用户 2) 关键约束 3) 成功标准"
	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunMissingInputReturnsUsageError(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(context.Background(), &stdout, &stderr, nil)
	if exitCode != 2 {
		t.Fatalf("run() exitCode = %d, want 2", exitCode)
	}

	if !strings.Contains(stderr.String(), "input is required") {
		t.Fatalf("stderr = %q, want contains %q", stderr.String(), "input is required")
	}
}

func TestRunRespectsTimeout(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(context.Background(), &stdout, &stderr, []string{"--input", "hello", "--delay", "100ms", "--timeout", "10ms"})
	if exitCode != 1 {
		t.Fatalf("run() exitCode = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr.String(), "deadline exceeded") {
		t.Fatalf("stderr = %q, want timeout error", stderr.String())
	}
}
