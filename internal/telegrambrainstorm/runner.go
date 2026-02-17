package telegrambrainstorm

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"codex-brainstorming-telegram/internal/telegramapi"
)

var ErrSessionTimeout = errors.New("brainstorming session timed out")

type sessionAPI interface {
	SendMessage(ctx context.Context, chatID string, text string) (int64, error)
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegramapi.Update, error)
}

type PromptResult struct {
	RawReply        string
	NormalizedReply string
}

func RunPrompt(ctx context.Context, api sessionAPI, chatID string, prompt string, sessionTimeout time.Duration) (PromptResult, error) {
	chatID = strings.TrimSpace(chatID)
	if chatID == "" {
		return PromptResult{}, errors.New("chatID is required")
	}

	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return PromptResult{}, errors.New("prompt is required")
	}
	if sessionTimeout <= 0 {
		return PromptResult{}, errors.New("session timeout must be greater than 0")
	}

	offset, err := latestUpdateOffset(ctx, api)
	if err != nil {
		return PromptResult{}, fmt.Errorf("read latest update offset: %w", err)
	}

	if _, err := api.SendMessage(ctx, chatID, prompt); err != nil {
		return PromptResult{}, fmt.Errorf("send prompt: %w", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()
	deadline := time.Now().Add(sessionTimeout)

	for {
		if err := waitCtx.Err(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return PromptResult{}, ErrSessionTimeout
			}
			return PromptResult{}, err
		}

		pollTimeout := computePollTimeout(time.Until(deadline))
		updates, err := api.GetUpdates(waitCtx, offset, pollTimeout)
		if err != nil {
			if waitCtx.Err() != nil {
				continue
			}
			return PromptResult{}, fmt.Errorf("poll updates: %w", err)
		}

		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			if fmt.Sprintf("%d", update.Message.Chat.ID) != chatID {
				continue
			}

			raw := strings.TrimSpace(update.Message.Text)
			if raw == "" {
				continue
			}

			return PromptResult{
				RawReply:        raw,
				NormalizedReply: normalizeReply(raw),
			}, nil
		}
	}
}

func normalizeReply(raw string) string {
	return strings.TrimSpace(raw)
}

func latestUpdateOffset(ctx context.Context, api sessionAPI) (int64, error) {
	updates, err := api.GetUpdates(ctx, 0, 0)
	if err != nil {
		return 0, err
	}

	var offset int64
	for _, update := range updates {
		if update.UpdateID >= offset {
			offset = update.UpdateID + 1
		}
	}

	return offset, nil
}

func computePollTimeout(remaining time.Duration) int {
	if remaining <= 0 {
		return 1
	}
	sec := int(math.Ceil(remaining.Seconds()))
	if sec > 20 {
		return 20
	}
	if sec < 1 {
		return 1
	}
	return sec
}
