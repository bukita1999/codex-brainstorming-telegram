package virtualcodex

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRespondKnownPrompt(t *testing.T) {
	t.Parallel()

	engine := NewEngine(Config{ProcessingDelay: 0})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	got, err := engine.Respond(ctx, "brainstorm login flow")
	if err != nil {
		t.Fatalf("Respond() error = %v", err)
	}

	want := "VirtualCodex: 好的，我们先收敛需求。请回复 1/2/3：1) 目标用户 2) 关键约束 3) 成功标准"
	if got != want {
		t.Fatalf("Respond() = %q, want %q", got, want)
	}
}

func TestRespondUnknownPromptFallsBackToEcho(t *testing.T) {
	t.Parallel()

	engine := NewEngine(Config{ProcessingDelay: 0})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	got, err := engine.Respond(ctx, "hello from cli")
	if err != nil {
		t.Fatalf("Respond() error = %v", err)
	}

	want := "VirtualCodex: 收到输入 -> hello from cli"
	if got != want {
		t.Fatalf("Respond() = %q, want %q", got, want)
	}
}

func TestRespondEmptyPromptReturnsError(t *testing.T) {
	t.Parallel()

	engine := NewEngine(Config{ProcessingDelay: 0})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := engine.Respond(ctx, "   ")
	if !errors.Is(err, ErrEmptyInput) {
		t.Fatalf("Respond() error = %v, want %v", err, ErrEmptyInput)
	}
}

func TestRespondTimeout(t *testing.T) {
	t.Parallel()

	engine := NewEngine(Config{ProcessingDelay: 80 * time.Millisecond})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := engine.Respond(ctx, "brainstorm anything")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Respond() error = %v, want %v", err, context.DeadlineExceeded)
	}
}
