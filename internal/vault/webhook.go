package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// WebhookEvent represents the payload sent to a webhook endpoint.
type WebhookEvent struct {
	Path      string            `json:"path"`
	Operation string            `json:"operation"`
	Version   int               `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// WebhookConfig holds configuration for a single webhook target.
type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Timeout time.Duration     `json:"timeout"`
}

// WebhookResult captures the outcome of a webhook dispatch.
type WebhookResult struct {
	URL        string
	StatusCode int
	Err        error
}

// Dispatcher sends webhook events to configured targets.
type Dispatcher struct {
	client *http.Client
}

// NewDispatcher creates a Dispatcher with a default HTTP client.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches a WebhookEvent to all provided configs concurrently
// and returns one result per config.
func (d *Dispatcher) Send(ctx context.Context, event WebhookEvent, configs []WebhookConfig) []WebhookResult {
	results := make([]WebhookResult, len(configs))
	for i, cfg := range configs {
		results[i] = d.sendOne(ctx, event, cfg)
	}
	return results
}

func (d *Dispatcher) sendOne(ctx context.Context, event WebhookEvent, cfg WebhookConfig) WebhookResult {
	res := WebhookResult{URL: cfg.URL}

	body, err := json.Marshal(event)
	if err != nil {
		res.Err = fmt.Errorf("marshal event: %w", err)
		return res
	}

	if !strings.HasPrefix(cfg.URL, "http") {
		res.Err = fmt.Errorf("invalid webhook URL: %s", cfg.URL)
		return res
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.URL, bytes.NewReader(body))
	if err != nil {
		res.Err = fmt.Errorf("build request: %w", err)
		return res
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		res.Err = fmt.Errorf("send request: %w", err)
		return res
	}
	defer resp.Body.Close()
	res.StatusCode = resp.StatusCode
	if resp.StatusCode >= 300 {
		res.Err = fmt.Errorf("unexpected status %d from %s", resp.StatusCode, cfg.URL)
	}
	return res
}
