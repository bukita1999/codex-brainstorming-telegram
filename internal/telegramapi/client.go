package telegramapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	botToken   string
	httpClient *http.Client
}

type Update struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type sendMessageResult struct {
	MessageID int64 `json:"message_id"`
}

func NewClient(baseURL string, botToken string, httpClient *http.Client) *Client {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		trimmed = "https://api.telegram.org"
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    trimmed,
		botToken:   strings.TrimSpace(botToken),
		httpClient: httpClient,
	}
}

func (c *Client) SendMessage(ctx context.Context, chatID string, text string) (int64, error) {
	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", text)

	respBody, err := c.postForm(ctx, "sendMessage", form)
	if err != nil {
		return 0, err
	}

	var apiResp struct {
		OK          bool              `json:"ok"`
		Description string            `json:"description"`
		Result      sendMessageResult `json:"result"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return 0, fmt.Errorf("decode sendMessage response: %w", err)
	}
	if !apiResp.OK {
		return 0, fmt.Errorf("sendMessage failed: %s", apiResp.Description)
	}

	return apiResp.Result.MessageID, nil
}

func (c *Client) GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]Update, error) {
	q := url.Values{}
	q.Set("offset", fmt.Sprintf("%d", offset))
	if timeoutSec > 0 {
		q.Set("timeout", fmt.Sprintf("%d", timeoutSec))
	}

	respBody, err := c.get(ctx, "getUpdates", q)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		OK          bool     `json:"ok"`
		Description string   `json:"description"`
		Result      []Update `json:"result"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("decode getUpdates response: %w", err)
	}
	if !apiResp.OK {
		return nil, fmt.Errorf("getUpdates failed: %s", apiResp.Description)
	}

	return apiResp.Result, nil
}

func (c *Client) get(ctx context.Context, method string, q url.Values) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.botToken, method)
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request %s: %w", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request %s: status %d", method, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", method, err)
	}

	return body, nil
}

func (c *Client) postForm(ctx context.Context, method string, form url.Values) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/bot%s/%s", c.baseURL, c.botToken, method)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request %s: %w", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request %s: status %d", method, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", method, err)
	}

	return body, nil
}
