package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
}

func TestRunPrintsOnlyStatusToTerminal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TELEGRAM_BOT_TOKEN=token\nTELEGRAM_CHAT_ID=123\nTELEGRAM_REPLY_TIMEOUT=1m\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	orig := runSession
	runSession = func(ctx context.Context, _ sessionAPI, chatID string, timeout time.Duration) (sessionResult, error) {
		if chatID != "123" {
			t.Fatalf("chatID = %q, want 123", chatID)
		}
		if timeout != time.Minute {
			t.Fatalf("timeout = %s, want 1m", timeout)
		}
		return sessionResult{}, nil
	}
	t.Cleanup(func() {
		runSession = orig
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(context.Background(), &stdout, &stderr, []string{"--env", envPath})
	if exitCode != 0 {
		t.Fatalf("run() exitCode = %d, stderr = %s", exitCode, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "程序正在运行中") {
		t.Fatalf("stdout = %q, want running status", out)
	}
	if strings.Contains(out, "请回复 1/2/3") {
		t.Fatalf("stdout should not include brainstorming questions: %q", out)
	}
}
