package telegrambrainstorm

import (
	"context"
	"errors"
	"strings"
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

func TestRunSessionSuccess(t *testing.T) {
	t.Parallel()

	u1 := telegramapi.Update{UpdateID: 2}
	u1.Message.Chat.ID = 1001
	u1.Message.Text = "1"

	u2 := telegramapi.Update{UpdateID: 3}
	u2.Message.Chat.ID = 1001
	u2.Message.Text = "性能和稳定性优先"

	u3 := telegramapi.Update{UpdateID: 4}
	u3.Message.Chat.ID = 1001
	u3.Message.Text = "3"

	api := &fakeAPI{
		polls: [][]telegramapi.Update{
			nil,
			{u1},
			{u2},
			{u3},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := RunSession(ctx, api, "1001", 2*time.Second)
	if err != nil {
		t.Fatalf("RunSession() error = %v", err)
	}

	if len(result.Answers) != 3 {
		t.Fatalf("len(result.Answers) = %d, want 3", len(result.Answers))
	}
	if got := result.Answers[0].AnswerText; !strings.Contains(got, "新用户") {
		t.Fatalf("first answer = %q, want normalized option text", got)
	}
	if got := result.Answers[2].AnswerText; !strings.Contains(got, "验证") {
		t.Fatalf("third answer = %q, want normalized option text", got)
	}

	if len(api.sentText) < 5 {
		t.Fatalf("sent count = %d, want at least 5", len(api.sentText))
	}
	if !strings.Contains(api.sentText[1], "请回复 1/2/3，或直接输入你的说明。") {
		t.Fatalf("question message missing reply hint: %q", api.sentText[1])
	}
}

func TestRunSessionTimeout(t *testing.T) {
	t.Parallel()

	api := &fakeAPI{polls: [][]telegramapi.Update{nil, nil, nil}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := RunSession(ctx, api, "1001", 20*time.Millisecond)
	if err == nil {
		t.Fatal("RunSession() error = nil, want timeout error")
	}
	if !errors.Is(err, ErrSessionTimeout) {
		t.Fatalf("RunSession() error = %v, want %v", err, ErrSessionTimeout)
	}
}
