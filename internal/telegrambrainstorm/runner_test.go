package telegrambrainstorm

import (
	"context"
	"errors"
	"testing"
	"time"

	"codex-brainstorming-telegram/internal/telegramapi"
)

type fakeAPI struct {
	polls    [][]telegramapi.Update
	pollIdx  int
	sentText []string
}

func (f *fakeAPI) SendMessage(_ context.Context, _ string, text string) (int64, error) {
	f.sentText = append(f.sentText, text)
	return int64(len(f.sentText)), nil
}

func (f *fakeAPI) GetUpdates(_ context.Context, _ int64, _ int) ([]telegramapi.Update, error) {
	if f.pollIdx >= len(f.polls) {
		return nil, nil
	}
	out := f.polls[f.pollIdx]
	f.pollIdx++
	return out, nil
}

func TestRunPromptSuccess(t *testing.T) {
	t.Parallel()

	u1 := telegramapi.Update{UpdateID: 2}
	u1.Message.Chat.ID = 1001
	u1.Message.Text = "B"

	api := &fakeAPI{
		polls: [][]telegramapi.Update{
			nil,
			{u1},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prompt := "请选择方案：\nA) 低风险\nB) 平衡\nC) 激进\n请回复 A/B/C。"
	result, err := RunPrompt(ctx, api, "1001", prompt, 2*time.Second)
	if err != nil {
		t.Fatalf("RunPrompt() error = %v", err)
	}

	if got, want := result.RawReply, "B"; got != want {
		t.Fatalf("result.RawReply = %q, want %q", got, want)
	}
	if got, want := result.NormalizedReply, "B"; got != want {
		t.Fatalf("result.NormalizedReply = %q, want %q", got, want)
	}
	if len(api.sentText) != 1 {
		t.Fatalf("sent count = %d, want 1", len(api.sentText))
	}
	if got := api.sentText[0]; got != prompt {
		t.Fatalf("sent prompt = %q, want %q", got, prompt)
	}
}

func TestRunPromptTimeout(t *testing.T) {
	t.Parallel()

	api := &fakeAPI{polls: [][]telegramapi.Update{nil, nil, nil}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := RunPrompt(ctx, api, "1001", "A/B/C?", 20*time.Millisecond)
	if err == nil {
		t.Fatal("RunPrompt() error = nil, want timeout error")
	}
	if !errors.Is(err, ErrSessionTimeout) {
		t.Fatalf("RunPrompt() error = %v, want %v", err, ErrSessionTimeout)
	}
}
