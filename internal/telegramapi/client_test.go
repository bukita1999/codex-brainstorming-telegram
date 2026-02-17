package telegramapi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	var gotChatID string
	var gotText string

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/bottoken123/sendMessage" {
				t.Fatalf("path = %q", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/x-www-form-urlencoded") {
				t.Fatalf("Content-Type = %q", ct)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			vals, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatalf("ParseQuery() error = %v", err)
			}
			gotChatID = vals.Get("chat_id")
			gotText = vals.Get("text")

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"message_id":42}}`)),
			}, nil
		}),
	}

	client := NewClient("https://api.telegram.test", "token123", httpClient)
	messageID, err := client.SendMessage(context.Background(), "777", "hello")
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	if gotChatID != "777" || gotText != "hello" {
		t.Fatalf("form chat_id=%q text=%q", gotChatID, gotText)
	}
	if messageID != 42 {
		t.Fatalf("messageID = %d, want 42", messageID)
	}
}

func TestGetUpdates(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/bottoken123/getUpdates" {
				t.Fatalf("path = %q", r.URL.Path)
			}

			q := r.URL.Query()
			if q.Get("offset") != "5" {
				t.Fatalf("offset = %q, want 5", q.Get("offset"))
			}
			if q.Get("timeout") != "12" {
				t.Fatalf("timeout = %q, want 12", q.Get("timeout"))
			}

			return &http.Response{
				StatusCode: 200,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":[{"update_id":6,"message":{"text":"123456","chat":{"id":777}}}]}`)),
			}, nil
		}),
	}

	client := NewClient("https://api.telegram.test", "token123", httpClient)
	updates, err := client.GetUpdates(context.Background(), 5, 12)
	if err != nil {
		t.Fatalf("GetUpdates() error = %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("len(updates) = %d, want 1", len(updates))
	}
	if updates[0].UpdateID != 6 {
		t.Fatalf("updateID = %d, want 6", updates[0].UpdateID)
	}
	if updates[0].Message.Chat.ID != 777 {
		t.Fatalf("chatID = %d, want 777", updates[0].Message.Chat.ID)
	}
}
