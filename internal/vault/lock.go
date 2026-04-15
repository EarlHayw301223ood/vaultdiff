package vault

import (
	"context"
	"fmt"
	"time"
)

// LockResult holds information about an acquired advisory lock.
type LockResult struct {
	Path      string
	AcquiredAt time.Time
	Owner     string
}

// LockOptions configures lock acquisition behaviour.
type LockOptions struct {
	// TTL is how long the lock should be held (written as metadata).
	TTL time.Duration
	// Owner identifies the process/user acquiring the lock.
	Owner string
}

// AcquireLock writes an advisory lock entry under the given KV path.
// It returns an error if a lock already exists and has not expired.
func AcquireLock(ctx context.Context, lc LogicalClient, path string, opts LockOptions) (*LockResult, error) {
	if path == "" {
		return nil, fmt.Errorf("lock path must not be empty")
	}

	existing, _ := lc.ReadWithContext(ctx, kvV2ReadPath(path))
	if existing != nil && existing.Data != nil {
		data, _ := existing.Data["data"].(map[string]interface{})
		if expStr, ok := data["expires_at"].(string); ok {
			exp, err := time.Parse(time.RFC3339, expStr)
			if err == nil && time.Now().UTC().Before(exp) {
				return nil, fmt.Errorf("lock at %q is held by %q until %s", path, data["owner"], exp.Format(time.RFC3339))
			}
		}
	}

	now := time.Now().UTC()
	expires := now.Add(opts.TTL)

	_, err := lc.WriteWithContext(ctx, kvV2WritePath(path), map[string]interface{}{
		"data": map[string]interface{}{
			"owner":      opts.Owner,
			"acquired_at": now.Format(time.RFC3339),
			"expires_at":  expires.Format(time.RFC3339),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("writing lock: %w", err)
	}

	return &LockResult{
		Path:       path,
		AcquiredAt: now,
		Owner:      opts.Owner,
	}, nil
}

// ReleaseLock deletes the advisory lock at the given KV path.
func ReleaseLock(ctx context.Context, lc LogicalClient, path string) error {
	if path == "" {
		return fmt.Errorf("lock path must not be empty")
	}
	_, err := lc.WriteWithContext(ctx, "secret/metadata/"+path, map[string]interface{}{
		"versions": []int{},
	})
	// Best-effort delete via metadata; ignore not-found.
	_ = err
	_, err = lc.WriteWithContext(ctx, kvV2WritePath(path), map[string]interface{}{
		"data": map[string]interface{}{},
	})
	return err
}

// kvV2ReadPath converts a logical path to its KV v2 data read path.
func kvV2ReadPath(path string) string {
	return "secret/data/" + path
}
