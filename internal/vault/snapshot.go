package vault

import (
	"context"
	"fmt"
	"sort"
)

// SecretSnapshot holds all secrets read from a mount at a point in time.
type SecretSnapshot struct {
	Mount   string
	Secrets map[string]map[string]string // path -> key/value pairs
}

// SnapshotMount reads every secret under the given mount and returns a
// SecretSnapshot containing the latest data for each path.
func SnapshotMount(ctx context.Context, c *Client, mount string) (*SecretSnapshot, error) {
	paths, err := ListTree(ctx, c, mount)
	if err != nil {
		return nil, fmt.Errorf("snapshot: list tree: %w", err)
	}

	snapshot := &SecretSnapshot{
		Mount:   mount,
		Secrets: make(map[string]map[string]string, len(paths)),
	}

	for _, path := range paths {
		secret, err := FetchAtRef(ctx, c, path, "latest")
		if err != nil {
			// Skip paths we cannot read rather than aborting the whole snapshot.
			continue
		}
		snapshot.Secrets[path] = secret
	}

	return snapshot, nil
}

// Paths returns the sorted list of secret paths in the snapshot.
func (s *SecretSnapshot) Paths() []string {
	paths := make([]string, 0, len(s.Secrets))
	for p := range s.Secrets {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// Get returns the key/value map for a given path, or nil if not present.
func (s *SecretSnapshot) Get(path string) map[string]string {
	return s.Secrets[path]
}
