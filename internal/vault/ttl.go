package vault

import (
	"errors"
	"fmt"
	"time"
)

// TTLRecord holds the TTL metadata stored for a secret version.
type TTLRecord struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	ExpiresAt time.Time `json:"expires_at"`
	SetAt     time.Time `json:"set_at"`
	SetBy     string    `json:"set_by"`
}

// IsExpired reports whether the TTL has passed relative to now.
func (r TTLRecord) IsExpired() bool {
	if r.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().UTC().After(r.ExpiresAt)
}

// RemainingTTL returns the duration until expiry, or zero if already expired.
func (r TTLRecord) RemainingTTL() time.Duration {
	if r.ExpiresAt.IsZero() || r.IsExpired() {
		return 0
	}
	return time.Until(r.ExpiresAt).Truncate(time.Second)
}

func ttlMetaPath(path string) string {
	return fmt.Sprintf("vaultdiff/meta/ttl/%s", normalisePath(path))
}

// SetTTL writes a TTL record for the given path and version.
func SetTTL(lc LogicalClient, path string, version int, ttl time.Duration, setBy string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if version <= 0 {
		return errors.New("version must be a positive integer")
	}
	if ttl <= 0 {
		return errors.New("ttl must be a positive duration")
	}

	now := time.Now().UTC()
	record := TTLRecord{
		Path:      path,
		Version:   version,
		ExpiresAt: now.Add(ttl),
		SetAt:     now,
		SetBy:     setBy,
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"path":       record.Path,
			"version":    record.Version,
			"expires_at": record.ExpiresAt.Format(time.RFC3339),
			"set_at":     record.SetAt.Format(time.RFC3339),
			"set_by":     record.SetBy,
		},
	}

	_, err := lc.Write(ttlMetaPath(path), data)
	return err
}

// GetTTL reads the TTL record for the given path.
func GetTTL(lc LogicalClient, path string) (*TTLRecord, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	secret, err := lc.Read(ttlMetaPath(path))
	if err != nil {
		return nil, fmt.Errorf("read ttl: %w", err)
	}
	if secret == nil {
		return nil, nil
	}

	m := extractStringMap(secret)

	parseTime := func(key string) time.Time {
		if v, ok := m[key]; ok {
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				return t.UTC()
			}
		}
		return time.Time{}
	}

	var ver int
	if v, ok := secret.Data["data"]; ok {
		if dm, ok := v.(map[string]interface{}); ok {
			if n, ok := dm["version"].(float64); ok {
				ver = int(n)
			}
		}
	}

	return &TTLRecord{
		Path:      m["path"],
		Version:   ver,
		ExpiresAt: parseTime("expires_at"),
		SetAt:     parseTime("set_at"),
		SetBy:     m["set_by"],
	}, nil
}
