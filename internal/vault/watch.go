package vault

import (
	"context"
	"time"
)

// WatchEvent represents a change detected during a watch poll.
type WatchEvent struct {
	Path       string
	OldVersion int
	NewVersion int
	DetectedAt time.Time
}

// Watcher polls a Vault KV path at a given interval and emits events
// when the current version changes.
type Watcher struct {
	client   LogicalClient
	path     string
	interval time.Duration
	last     int
}

// NewWatcher creates a Watcher for the given path and poll interval.
func NewWatcher(client LogicalClient, path string, interval time.Duration) *Watcher {
	return &Watcher{
		client:   client,
		path:     path,
		interval: interval,
	}
}

// Watch starts polling and sends WatchEvents on the returned channel.
// It stops when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) (<-chan WatchEvent, error) {
	current, err := latestVersion(w.client, w.path)
	if err != nil {
		return nil, err
	}
	w.last = current

	ch := make(chan WatchEvent, 8)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				v, err := latestVersion(w.client, w.path)
				if err != nil || v == w.last {
					continue
				}
				ch <- WatchEvent{
					Path:       w.path,
					OldVersion: w.last,
					NewVersion: v,
					DetectedAt: time.Now().UTC(),
				}
				w.last = v
			}
		}
	}()
	return ch, nil
}

// latestVersion returns the current latest version number for a KV v2 path.
func latestVersion(client LogicalClient, path string) (int, error) {
	versions, err := ListVersions(client, path)
	if err != nil {
		return 0, err
	}
	if len(versions) == 0 {
		return 0, nil
	}
	return versions[len(versions)-1].Version, nil
}
