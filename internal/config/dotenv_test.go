package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadTelegramConfigFromEnvFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "# sample\nTELEGRAM_BOT_TOKEN=abc123\nTELEGRAM_CHAT_ID=987654\nTELEGRAM_PROXY_URL=http://127.0.0.1:7890\nTELEGRAM_REPLY_TIMEOUT=4m\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadTelegramConfig(envPath)
	if err != nil {
		t.Fatalf("LoadTelegramConfig() error = %v", err)
	}

	if cfg.BotToken != "abc123" {
		t.Fatalf("BotToken = %q, want %q", cfg.BotToken, "abc123")
	}
	if cfg.ChatID != "987654" {
		t.Fatalf("ChatID = %q, want %q", cfg.ChatID, "987654")
	}
	if cfg.ProxyURL != "http://127.0.0.1:7890" {
		t.Fatalf("ProxyURL = %q, want %q", cfg.ProxyURL, "http://127.0.0.1:7890")
	}
	if cfg.ReplyTimeout != 4*time.Minute {
		t.Fatalf("ReplyTimeout = %s, want %s", cfg.ReplyTimeout, 4*time.Minute)
	}
}

func TestLoadTelegramConfigMissingEnvFileHasActionableMessage(t *testing.T) {
	t.Parallel()

	missing := filepath.Join(t.TempDir(), ".env")
	_, err := LoadTelegramConfig(missing)
	if err == nil {
		t.Fatal("LoadTelegramConfig() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), ".env.example") {
		t.Fatalf("error = %q, want mention .env.example", err.Error())
	}
	if !strings.Contains(err.Error(), "创建对应的 .env 文件") {
		t.Fatalf("error = %q, want creation hint", err.Error())
	}
}

func TestLoadTelegramConfigRequiresTokenAndChatID(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TELEGRAM_PROXY_URL=http://127.0.0.1:7890\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := LoadTelegramConfig(envPath)
	if err == nil {
		t.Fatal("LoadTelegramConfig() error = nil, want non-nil")
	}
}

func TestLoadTelegramConfigUsesDefaultTimeout(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	content := "TELEGRAM_BOT_TOKEN=abc123\nTELEGRAM_CHAT_ID=10001\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := LoadTelegramConfig(envPath)
	if err != nil {
		t.Fatalf("LoadTelegramConfig() error = %v", err)
	}

	if cfg.ReplyTimeout != 5*time.Minute {
		t.Fatalf("ReplyTimeout = %s, want %s", cfg.ReplyTimeout, 5*time.Minute)
	}
}
