package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultReplyTimeout = 5 * time.Minute

type TelegramConfig struct {
	BotToken     string
	ChatID       string
	ProxyURL     string
	ReplyTimeout time.Duration
}

func LoadTelegramConfig(path string) (TelegramConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return TelegramConfig{}, fmt.Errorf("%s 不存在，请根据 .env.example 创建对应的 .env 文件", path)
		}
		return TelegramConfig{}, fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		values[key] = val
	}
	if err := scanner.Err(); err != nil {
		return TelegramConfig{}, fmt.Errorf("scan env file: %w", err)
	}

	cfg := TelegramConfig{
		BotToken:     strings.TrimSpace(values["TELEGRAM_BOT_TOKEN"]),
		ChatID:       strings.TrimSpace(values["TELEGRAM_CHAT_ID"]),
		ProxyURL:     strings.TrimSpace(values["TELEGRAM_PROXY_URL"]),
		ReplyTimeout: defaultReplyTimeout,
	}

	if raw := strings.TrimSpace(values["TELEGRAM_REPLY_TIMEOUT"]); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return TelegramConfig{}, fmt.Errorf("parse TELEGRAM_REPLY_TIMEOUT: %w", err)
		}
		if d <= 0 {
			return TelegramConfig{}, errors.New("TELEGRAM_REPLY_TIMEOUT must be greater than 0")
		}
		cfg.ReplyTimeout = d
	}

	if cfg.BotToken == "" {
		return TelegramConfig{}, errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.ChatID == "" {
		return TelegramConfig{}, errors.New("TELEGRAM_CHAT_ID is required")
	}

	return cfg, nil
}
