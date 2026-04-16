package vault

import (
	"errors"
	"fmt"
	"time"
)

// ExpireResult holds the outcome of an expiration check or action.
type ExpireResult struct {
	Path      string
	Version   int
	CreatedAt time.Time
	ExpiredAt time.Time
	Expired   bool
}

// CheckExpiry returns an ExpireResult indicating whether the given secret version
// is older than the provided TTL duration.
func CheckExpiry(client *Client, path string, version int, ttl time.Duration) (*ExpireResult, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if ttl <= 0 {
		return nil, errors.New("ttl must be positive")
	}

	secret, err := FetchAtRef(client, path, fmt.Sprintf("%d", version))
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %s version %d", path, version)
	}

	meta := extractMetadata(secret)
	created := meta.CreatedTime
	expiredAt := created.Add(ttl)

	return &ExpireResult{
		Path:      path,
		Version:   version,
		CreatedAt: created,
		ExpiredAt: expiredAt,
		Expired:   time.Now().UTC().After(expiredAt),
	}, nil
}

// ScanExpired walks all paths under a mount and returns those whose latest
// version is older than ttl.
func ScanExpired(client *Client, mount string, ttl time.Duration) ([]*ExpireResult, error) {
	paths, err := ListTree(client, mount)
	if err != nil {
		return nil, fmt.Errorf("list tree: %w", err)
	}

	var results []*ExpireResult
	for _, p := range paths {
		res, err := CheckExpiry(client, p, 0, ttl)
		if err != nil {
			continue
		}
		if res.Expired {
			results = append(results, res)
		}
	}
	return results, nil
}
