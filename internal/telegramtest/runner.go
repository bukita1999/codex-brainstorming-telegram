package telegramtest

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"codex-brainstorming-telegram/internal/telegramapi"
)

var ErrChallengeTimeout = errors.New("did not receive matching reply before timeout")

type challengeAPI interface {
	SendMessage(ctx context.Context, chatID string, text string) (int64, error)
	GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]telegramapi.Update, error)
}

func RunChallenge(ctx context.Context, api challengeAPI, chatID string, code string, replyTimeout time.Duration) error {
	if replyTimeout <= 0 {
		return errors.New("reply timeout must be greater than 0")
	}

	latestOffset, err := latestUpdateOffset(ctx, api)
	if err != nil {
		return fmt.Errorf("read latest update offset: %w", err)
	}

	message := BuildChallengeMessage(code)
	if _, err := api.SendMessage(ctx, chatID, message); err != nil {
		return fmt.Errorf("send challenge message: %w", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, replyTimeout)
	defer cancel()
	deadline := time.Now().Add(replyTimeout)

	offset := latestOffset
	for {
		if err := waitCtx.Err(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return ErrChallengeTimeout
			}
			return err
		}

		pollTimeout := computePollTimeout(time.Until(deadline))
		updates, err := api.GetUpdates(waitCtx, offset, pollTimeout)
		if err != nil {
			if waitCtx.Err() != nil {
				continue
			}
			return fmt.Errorf("poll updates: %w", err)
		}

		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}

			if fmt.Sprintf("%d", update.Message.Chat.ID) != chatID {
				continue
			}
			if IsMatchingReply(update.Message.Text, code) {
				return nil
			}
		}
	}
}

func latestUpdateOffset(ctx context.Context, api challengeAPI) (int64, error) {
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
