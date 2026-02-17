package telegramtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"codex-brainstorming-telegram/internal/telegramapi"
)

type fakeAPI struct {
	sendChatID string
	sendText   string
	polls      [][]telegramapi.Update
	pollIndex  int
}

func (f *fakeAPI) SendMessage(_ context.Context, chatID string, text string) (int64, error) {
	f.sendChatID = chatID
	f.sendText = text
	return 1, nil
}

func (f *fakeAPI) GetUpdates(_ context.Context, _ int64, _ int) ([]telegramapi.Update, error) {
	if f.pollIndex >= len(f.polls) {
		return nil, nil
	}
	out := f.polls[f.pollIndex]
	f.pollIndex++
	return out, nil
}

func TestRunChallengeSuccess(t *testing.T) {
	t.Parallel()

	update := telegramapi.Update{UpdateID: 2}
	update.Message.Chat.ID = 123
	update.Message.Text = "654321"

	api := &fakeAPI{
		polls: [][]telegramapi.Update{
			nil,
			{update},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := RunChallenge(ctx, api, "123", "654321", 3*time.Second)
	if err != nil {
		t.Fatalf("RunChallenge() error = %v", err)
	}

	if api.sendChatID != "123" {
		t.Fatalf("send chat_id = %q, want %q", api.sendChatID, "123")
	}
	if api.sendText != `这是一个测试，请回复 "[654321]"` {
		t.Fatalf("send text = %q", api.sendText)
	}
}

func TestRunChallengeTimeout(t *testing.T) {
	t.Parallel()

	api := &fakeAPI{polls: [][]telegramapi.Update{nil, nil, nil, nil}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := RunChallenge(ctx, api, "123", "654321", 20*time.Millisecond)
	if err == nil {
		t.Fatal("RunChallenge() error = nil, want timeout error")
	}
	if !errors.Is(err, ErrChallengeTimeout) {
		t.Fatalf("RunChallenge() error = %v, want %v", err, ErrChallengeTimeout)
	}
}
