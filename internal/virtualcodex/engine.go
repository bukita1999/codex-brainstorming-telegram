package virtualcodex

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrEmptyInput = errors.New("input is empty")

type Config struct {
	ProcessingDelay time.Duration
}

type Engine struct {
	processingDelay time.Duration
}

func NewEngine(cfg Config) *Engine {
	delay := cfg.ProcessingDelay
	if delay < 0 {
		delay = 0
	}
	return &Engine{processingDelay: delay}
}

func (e *Engine) Respond(ctx context.Context, input string) (string, error) {
	normalized := strings.TrimSpace(input)
	if normalized == "" {
		return "", ErrEmptyInput
	}

	if err := waitOrTimeout(ctx, e.processingDelay); err != nil {
		return "", err
	}

	lower := strings.ToLower(normalized)
	if strings.Contains(lower, "brainstorm") || strings.Contains(lower, "login") || strings.Contains(lower, "头脑") {
		return "VirtualCodex: 好的，我们先收敛需求。请回复 1/2/3：1) 目标用户 2) 关键约束 3) 成功标准", nil
	}

	return fmt.Sprintf("VirtualCodex: 收到输入 -> %s", normalized), nil
}

func waitOrTimeout(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
