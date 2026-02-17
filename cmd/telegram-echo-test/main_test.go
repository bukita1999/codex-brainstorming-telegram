package main

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunShowsEnvCreationHintWhenEnvMissing(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	missing := filepath.Join(t.TempDir(), ".env")
	exitCode := run(context.Background(), &stdout, &stderr, []string{"--env", missing})
	if exitCode != 2 {
		t.Fatalf("run() exitCode = %d, want 2", exitCode)
	}
	if !strings.Contains(stderr.String(), ".env.example") {
		t.Fatalf("stderr = %q, want mention .env.example", stderr.String())
	}
	if !strings.Contains(stderr.String(), "创建对应的 .env 文件") {
		t.Fatalf("stderr = %q, want creation hint", stderr.String())
	}
}
