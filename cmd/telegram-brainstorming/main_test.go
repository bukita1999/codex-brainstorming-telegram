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

func TestRunRequiresPrompt(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TELEGRAM_BOT_TOKEN=token\nTELEGRAM_CHAT_ID=123\nTELEGRAM_REPLY_TIMEOUT=1m\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(context.Background(), &stdout, &stderr, []string{"--env", envPath})
	if exitCode != 2 {
		t.Fatalf("run() exitCode = %d, want 2", exitCode)
	}
	if !strings.Contains(stderr.String(), "prompt is required") {
		t.Fatalf("stderr = %q, want prompt required error", stderr.String())
	}
}

func TestRunPositionalPromptAndReplyOnlyOutput(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TELEGRAM_BOT_TOKEN=token\nTELEGRAM_CHAT_ID=123\nTELEGRAM_REPLY_TIMEOUT=1m\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	promptText := "请选择方案：\nA) 稳健\nB) 平衡\nC) 激进\n请回复 A/B/C。"
	orig := runPrompt
	runPrompt = func(ctx context.Context, _ promptAPI, chatID string, prompt string, timeout time.Duration) (promptResult, error) {
		if chatID != "123" {
			t.Fatalf("chatID = %q, want 123", chatID)
		}
		if prompt != promptText {
			t.Fatalf("prompt = %q, want %q", prompt, promptText)
		}
		if timeout != time.Minute {
			t.Fatalf("timeout = %s, want 1m", timeout)
		}
		return promptResult{
			RawReply:        "B",
			NormalizedReply: "B",
		}, nil
	}
	t.Cleanup(func() {
		runPrompt = orig
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(context.Background(), &stdout, &stderr, []string{"--env", envPath, promptText})
	if exitCode != 0 {
		t.Fatalf("run() exitCode = %d, stderr = %s", exitCode, stderr.String())
	}

	if got := stdout.String(); got != "B\n" {
		t.Fatalf("stdout = %q, want only reply output", got)
	}

	status := stderr.String()
	if !strings.Contains(status, "程序正在运行中") {
		t.Fatalf("stderr = %q, want running status", status)
	}
	if strings.Contains(status, promptText) {
		t.Fatalf("stderr should not include prompt text: %q", status)
	}
}
