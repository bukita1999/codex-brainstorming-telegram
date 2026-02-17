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
	"strings"
	"time"

	"codex-brainstorming-telegram/internal/config"
	"codex-brainstorming-telegram/internal/telegramapi"
	"codex-brainstorming-telegram/internal/telegrambrainstorm"
)

type promptAPI interface {
	SendMessage(ctx context.Context, chatID string, text string) (int64, error)
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegramapi.Update, error)
}

type promptResult = telegrambrainstorm.PromptResult

var runPrompt = func(ctx context.Context, api promptAPI, chatID string, prompt string, timeout time.Duration) (promptResult, error) {
	return telegrambrainstorm.RunPrompt(ctx, api, chatID, prompt, timeout)
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
	promptFlag := fs.String("prompt", "", "prompt text to send to Telegram")

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

	promptText, err := buildPromptText(*promptFlag, fs.Args())
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	httpClient, err := buildHTTPClient(cfg.ProxyURL)
	if err != nil {
		fmt.Fprintf(stderr, "proxy config error: %v\n", err)
		return 2
	}

	apiClient := telegramapi.NewClient(*apiBase, cfg.BotToken, httpClient)

	fmt.Fprintln(stderr, "程序正在运行中，请前往 Telegram 查看并回复。")
	fmt.Fprintln(stderr, "终端仅显示运行状态，不显示提问内容。")

	ctx, cancel := context.WithTimeout(parent, cfg.ReplyTimeout+30*time.Second)
	defer cancel()

	result, err := runPrompt(ctx, apiClient, cfg.ChatID, promptText, cfg.ReplyTimeout)
	if err != nil {
		if errors.Is(err, telegrambrainstorm.ErrSessionTimeout) {
			fmt.Fprintln(stderr, "会话超时：未在规定时间内完成 Telegram 对话")
			return 1
		}
		fmt.Fprintf(stderr, "会话失败：%v\n", err)
		return 1
	}

	fmt.Fprintln(stderr, "会话完成：已收到 Telegram 回复。")
	fmt.Fprintln(stdout, result.NormalizedReply)
	return 0
}

func buildPromptText(promptFlag string, positional []string) (string, error) {
	fromFlag := strings.TrimSpace(promptFlag)
	fromPositional := strings.TrimSpace(strings.Join(positional, " "))

	switch {
	case fromFlag != "" && fromPositional != "":
		return "", errors.New("use either --prompt or positional prompt, not both")
	case fromFlag != "":
		return unescapePromptText(fromFlag), nil
	case fromPositional != "":
		return unescapePromptText(fromPositional), nil
	default:
		return "", errors.New("prompt is required: pass --prompt \"...\" or provide positional text")
	}
}

func unescapePromptText(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch != '\\' || i+1 >= len(s) {
			b.WriteByte(ch)
			continue
		}

		switch s[i+1] {
		case 'n':
			b.WriteByte('\n')
			i++
		case 'r':
			b.WriteByte('\r')
			i++
		case 't':
			b.WriteByte('\t')
			i++
		case '\\':
			b.WriteByte('\\')
			i++
		default:
			// Keep unknown escape sequences unchanged.
			b.WriteByte('\\')
		}
	}

	return b.String()
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
