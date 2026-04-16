package vault

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

type stubVersionClient struct {
	calls    int
	responses []*api.Secret
}

func (s *stubVersionClient) Read(path string) (*api.Secret, error) {
	idx := s.calls
	if idx >= len(s.responses) {
		idx = len(s.responses) - 1
	}
	s.calls++
	return s.responses[idx], nil
}

func (s *stubVersionClient) List(path string) (*api.Secret, error) {
	return nil, nil
}

func TestNewWatcher_ReturnsWatcher(t *testing.T) {
	w := NewWatcher(&stubVersionClient{}, "secret/data/foo", time.Second)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestWatch_EmitsEventOnVersionChange(t *testing.T) {
	v1 := makeVersionsSecret(map[string]int{"1": 1})
	v2 := makeVersionsSecret(map[string]int{"1": 1, "2": 2})

	client := &stubVersionClient{
		responses: []*api.Secret{v1, v1, v2},
	}

	w := NewWatcher(client, "secret/data/foo", 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case evt := <-ch:
		if evt.NewVersion != 2 {
			t.Errorf("expected new version 2, got %d", evt.NewVersion)
		}
		if evt.OldVersion != 1 {
			t.Errorf("expected old version 1, got %d", evt.OldVersion)
		}
		if evt.Path != "secret/data/foo" {
			t.Errorf("unexpected path: %s", evt.Path)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatch_NoEventWhenVersionUnchanged(t *testing.T) {
	v1 := makeVersionsSecret(map[string]int{"1": 1})
	client := &stubVersionClient{
		responses: []*api.Secret{v1, v1, v1},
	}

	w := NewWatcher(client, "secret/data/foo", 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case evt := <-ch:
		t.Errorf("unexpected event received: %+v", evt)
	case <-ctx.Done():
		// expected: no event emitted when version does not change
	}
}

func TestWatch_ClosesChannelOnCancel(t *testing.T) {
	v1 := makeVersionsSecret(map[string]int{"1": 1})
	client := &stubVersionClient{responses: []*api.Secret{v1}}

	w := NewWatcher(client, "secret/data/foo", 50*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cancel()

	timer := time.NewTimer(200 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return // channel closed as expected
			}
		case <-timer.C:
			t.Fatal("channel not closed after cancel")
		}
	}
}
