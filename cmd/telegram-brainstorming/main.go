package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"codex-brainstorming-telegram/internal/config"
	"codex-brainstorming-telegram/internal/telegramapi"
	"codex-brainstorming-telegram/internal/telegrambrainstorm"
)

type sessionAPI interface {
	SendMessage(ctx context.Context, chatID string, text string) (int64, error)
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegramapi.Update, error)
}

type sessionResult = telegrambrainstorm.SessionResult

var runSession = func(ctx context.Context, api sessionAPI, chatID string, timeout time.Duration) (sessionResult, error) {
	return telegrambrainstorm.RunSession(ctx, api, chatID, timeout)
}

func main() {
	os.Exit(run(context.Background(), os.Stdout, os.Stderr, os.Args[1:]))
}

func run(parent context.Context, stdout io.Writer, stderr io.Writer, args []string) int {
	fs := flag.NewFlagSet("telegram-brainstorming", flag.ContinueOnError)
	fs.SetOutput(stderr)

	envPath := fs.String("env", ".env", "path to .env file")
	apiBase := fs.String("api-base", "https://api.telegram.org", "telegram API base URL")
	overrideTimeout := fs.Duration("session-timeout", 0, "override TELEGRAM_REPLY_TIMEOUT")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *overrideTimeout < 0 {
		fmt.Fprintln(stderr, "session-timeout must be >= 0")
		return 2
	}

	cfg, err := config.LoadTelegramConfig(*envPath)
	if err != nil {
		fmt.Fprintf(stderr, "load config failed: %v\n", err)
		return 2
	}
	if *overrideTimeout > 0 {
		cfg.ReplyTimeout = *overrideTimeout
	}

	httpClient, err := buildHTTPClient(cfg.ProxyURL)
	if err != nil {
		fmt.Fprintf(stderr, "proxy config error: %v\n", err)
		return 2
	}

	apiClient := telegramapi.NewClient(*apiBase, cfg.BotToken, httpClient)

	fmt.Fprintln(stdout, "程序正在运行中，请前往 Telegram 完成 brainstorming 对话。")
	fmt.Fprintln(stdout, "终端仅显示运行状态，不显示提问内容。")

	ctx, cancel := context.WithTimeout(parent, cfg.ReplyTimeout+30*time.Second)
	defer cancel()

	if _, err := runSession(ctx, apiClient, cfg.ChatID, cfg.ReplyTimeout); err != nil {
		if errors.Is(err, telegrambrainstorm.ErrSessionTimeout) {
			fmt.Fprintln(stderr, "会话超时：未在规定时间内完成 Telegram 对话")
			return 1
		}
		fmt.Fprintf(stderr, "会话失败：%v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, "会话完成：结果已发送到 Telegram。")
	return 0
}

func buildHTTPClient(proxyURL string) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment

	if proxyURL != "" {
		u, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid TELEGRAM_PROXY_URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(u)
	}

	return &http.Client{
		Transport: transport,
		Timeout:   35 * time.Second,
	}, nil
}
