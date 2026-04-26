package vault

import (
	"fmt"
	"strings"
	"time"
)

// NotifyChannel represents a supported notification channel type.
type NotifyChannel string

const (
	ChannelSlack  NotifyChannel = "slack"
	ChannelEmail  NotifyChannel = "email"
	ChannelCustom NotifyChannel = "custom"
)

// NotifyConfig holds configuration for a single notification rule.
type NotifyConfig struct {
	Channel  NotifyChannel     `json:"channel"`
	Target   string            `json:"target"`
	Events   []string          `json:"events"`
	Headers  map[string]string `json:"headers,omitempty"`
	Template string            `json:"template,omitempty"`
}

// NotifyEvent is emitted when a secret changes.
type NotifyEvent struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	Operation string    `json:"operation"`
	Triggered time.Time `json:"triggered"`
	By        string    `json:"by,omitempty"`
}

// NotifyResult records the outcome of dispatching a notification.
type NotifyResult struct {
	Channel  NotifyChannel `json:"channel"`
	Target   string        `json:"target"`
	Success  bool          `json:"success"`
	Err      string        `json:"err,omitempty"`
	SentAt   time.Time     `json:"sent_at"`
}

func notifyMetaPath(path string) string {
	path = strings.Trim(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return fmt.Sprintf("%s/metadata/_notify/%s", parts[0], parts[0])
	}
	return fmt.Sprintf("%s/metadata/_notify/%s", parts[0], parts[1])
}

// Notify dispatches a NotifyEvent to all configured channels using the Dispatcher.
func Notify(d *Dispatcher, event NotifyEvent, configs []NotifyConfig) []NotifyResult {
	results := make([]NotifyResult, 0, len(configs))
	for _, cfg := range configs {
		if !eventMatches(event.Operation, cfg.Events) {
			continue
		}
		body := buildNotifyBody(event, cfg)
		err := d.Send(cfg.Target, body, cfg.Headers)
		res := NotifyResult{
			Channel: cfg.Channel,
			Target:  cfg.Target,
			Success: err == nil,
			SentAt:  time.Now().UTC(),
		}
		if err != nil {
			res.Err = err.Error()
		}
		results = append(results, res)
	}
	return results
}

func eventMatches(op string, events []string) bool {
	if len(events) == 0 {
		return true
	}
	op = strings.ToLower(op)
	for _, e := range events {
		if strings.ToLower(e) == op || e == "*" {
			return true
		}
	}
	return false
}

func buildNotifyBody(event NotifyEvent, cfg NotifyConfig) string {
	if cfg.Template != "" {
		body := cfg.Template
		body = strings.ReplaceAll(body, "{{path}}", event.Path)
		body = strings.ReplaceAll(body, "{{version}}", fmt.Sprintf("%d", event.Version))
		body = strings.ReplaceAll(body, "{{operation}}", event.Operation)
		body = strings.ReplaceAll(body, "{{by}}", event.By)
		return body
	}
	return fmt.Sprintf(`{"path":%q,"version":%d,"operation":%q,"by":%q}`,
		event.Path, event.Version, event.Operation, event.By)
}
