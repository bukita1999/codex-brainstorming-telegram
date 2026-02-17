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

type Question struct {
	Key     string
	Title   string
	Options []string
}

type Answer struct {
	QuestionKey   string
	QuestionTitle string
	RawInput      string
	AnswerText    string
}

type SessionResult struct {
	Answers []Answer
}

func RunSession(ctx context.Context, api sessionAPI, chatID string, sessionTimeout time.Duration) (SessionResult, error) {
	if strings.TrimSpace(chatID) == "" {
		return SessionResult{}, errors.New("chatID is required")
	}
	if sessionTimeout <= 0 {
		return SessionResult{}, errors.New("session timeout must be greater than 0")
	}

	offset, err := latestUpdateOffset(ctx, api)
	if err != nil {
		return SessionResult{}, fmt.Errorf("read latest update offset: %w", err)
	}

	questions := defaultQuestions()
	if _, err := api.SendMessage(ctx, chatID, introMessage()); err != nil {
		return SessionResult{}, fmt.Errorf("send intro: %w", err)
	}
	if _, err := api.SendMessage(ctx, chatID, buildQuestionMessage(0, len(questions), questions[0])); err != nil {
		return SessionResult{}, fmt.Errorf("send first question: %w", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, sessionTimeout)
	defer cancel()
	deadline := time.Now().Add(sessionTimeout)

	current := 0
	answers := make([]Answer, 0, len(questions))

	for current < len(questions) {
		if err := waitCtx.Err(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return SessionResult{}, ErrSessionTimeout
			}
			return SessionResult{}, err
		}

		pollTimeout := computePollTimeout(time.Until(deadline))
		updates, err := api.GetUpdates(waitCtx, offset, pollTimeout)
		if err != nil {
			if waitCtx.Err() != nil {
				continue
			}
			return SessionResult{}, fmt.Errorf("poll updates: %w", err)
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

			q := questions[current]
			answers = append(answers, Answer{
				QuestionKey:   q.Key,
				QuestionTitle: q.Title,
				RawInput:      raw,
				AnswerText:    normalizeAnswer(raw, q.Options),
			})

			current++
			if current >= len(questions) {
				summary := buildSummaryMessage(answers)
				if _, err := api.SendMessage(waitCtx, chatID, summary); err != nil {
					return SessionResult{}, fmt.Errorf("send summary: %w", err)
				}
				return SessionResult{Answers: answers}, nil
			}

			next := buildQuestionMessage(current, len(questions), questions[current])
			if _, err := api.SendMessage(waitCtx, chatID, next); err != nil {
				return SessionResult{}, fmt.Errorf("send next question: %w", err)
			}
		}
	}

	return SessionResult{Answers: answers}, nil
}

func defaultQuestions() []Question {
	return []Question{
		{
			Key:   "target_user",
			Title: "目标用户是谁？",
			Options: []string{
				"新用户优先",
				"现有活跃用户",
				"内部运营或管理用户",
			},
		},
		{
			Key:   "constraints",
			Title: "当前最关键的约束是什么？",
			Options: []string{
				"上线时间优先",
				"性能和稳定性优先",
				"开发成本优先",
			},
		},
		{
			Key:   "success_criteria",
			Title: "你希望以什么标准判定成功？",
			Options: []string{
				"用户反馈明显变好",
				"业务数据明显提升",
				"有可验证的指标达成",
			},
		},
	}
}

func introMessage() string {
	return "我们开始 brainstorming。接下来我会一次只问一个问题。"
}

func buildQuestionMessage(index int, total int, q Question) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[Brainstorming %d/%d]\n", index+1, total))
	b.WriteString(q.Title)
	b.WriteString("\n")
	for i, opt := range q.Options {
		b.WriteString(fmt.Sprintf("%d) %s\n", i+1, opt))
	}
	b.WriteString("请回复 1/2/3，或直接输入你的说明。")
	return b.String()
}

func buildSummaryMessage(answers []Answer) string {
	var b strings.Builder
	b.WriteString("已完成本轮 brainstorming 收敛，摘要如下：\n")
	for i, answer := range answers {
		b.WriteString(fmt.Sprintf("%d) %s: %s\n", i+1, answer.QuestionTitle, answer.AnswerText))
	}
	b.WriteString("如需继续细化，请继续回复你的补充。")
	return b.String()
}

func normalizeAnswer(raw string, options []string) string {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) == 1 {
		idx := int(trimmed[0] - '1')
		if idx >= 0 && idx < len(options) {
			return options[idx]
		}
	}
	return trimmed
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
