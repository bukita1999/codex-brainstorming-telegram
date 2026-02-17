package main

import (
	"context"
	"crypto/rand"
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
	"codex-brainstorming-telegram/internal/telegramtest"
)

func main() {
	os.Exit(run(context.Background(), os.Stdout, os.Stderr, os.Args[1:]))
}

func run(parent context.Context, stdout io.Writer, stderr io.Writer, args []string) int {
	fs := flag.NewFlagSet("telegram-echo-test", flag.ContinueOnError)
	fs.SetOutput(stderr)

	envPath := fs.String("env", ".env", "path to .env file")
	apiBase := fs.String("api-base", "https://api.telegram.org", "telegram API base URL")
	overrideTimeout := fs.Duration("reply-timeout", 0, "override TELEGRAM_REPLY_TIMEOUT")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *overrideTimeout < 0 {
		fmt.Fprintln(stderr, "reply-timeout must be >= 0")
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

	code, err := telegramtest.GenerateCode(rand.Reader)
	if err != nil {
		fmt.Fprintf(stderr, "generate code failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "即将发送验证消息到 chat_id=%s\n", cfg.ChatID)
	fmt.Fprintf(stdout, "挑战码: %s\n", code)
	fmt.Fprintf(stdout, "消息内容: %s\n", telegramtest.BuildChallengeMessage(code))
	fmt.Fprintln(stdout, "请在 Telegram 中回复完全相同的六码。")

	apiClient := telegramapi.NewClient(*apiBase, cfg.BotToken, httpClient)
	ctx, cancel := context.WithTimeout(parent, cfg.ReplyTimeout+30*time.Second)
	defer cancel()

	err = telegramtest.RunChallenge(ctx, apiClient, cfg.ChatID, code, cfg.ReplyTimeout)
	if err != nil {
		if errors.Is(err, telegramtest.ErrChallengeTimeout) {
			fmt.Fprintln(stderr, "测试失败: 等待超时，未收到匹配回复")
			return 1
		}
		fmt.Fprintf(stderr, "测试失败: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, "测试成功: 收到匹配回复，链路未被篡改")
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
