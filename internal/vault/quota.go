package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// QuotaRecord holds metadata about a path-level write quota.
type QuotaRecord struct {
	Path      string    `json:"path"`
	MaxWrites int       `json:"max_writes"`
	WindowSec int       `json:"window_sec"`
	Writes    int       `json:"writes"`
	ResetAt   time.Time `json:"reset_at"`
}

// Exceeded reports whether the write count has reached the maximum.
func (q QuotaRecord) Exceeded() bool {
	return q.Writes >= q.MaxWrites
}

// Remaining returns how many writes are still allowed in the current window.
func (q QuotaRecord) Remaining() int {
	r := q.MaxWrites - q.Writes
	if r < 0 {
		return 0
	}
	return r
}

func quotaMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/__quota", path)
}

// SetQuota stores a quota policy for the given secret path.
func SetQuota(c LogicalClient, path string, maxWrites, windowSec int) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if maxWrites <= 0 {
		return errors.New("max_writes must be greater than zero")
	}
	if windowSec <= 0 {
		return errors.New("window_sec must be greater than zero")
	}

	record := map[string]interface{}{
		"path":       path,
		"max_writes": maxWrites,
		"window_sec": windowSec,
		"writes":     0,
		"reset_at":   time.Now().UTC().Add(time.Duration(windowSec) * time.Second).Format(time.RFC3339),
	}

	_, err := c.Write(quotaMetaPath(path), map[string]interface{}{"data": record})
	return err
}

// GetQuota retrieves the quota record for the given secret path.
func GetQuota(c LogicalClient, path string) (*QuotaRecord, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	secret, err := c.Read(quotaMetaPath(path))
	if err != nil {
		return nil, fmt.Errorf("read quota: %w", err)
	}
	if secret == nil {
		return nil, nil
	}

	data := extractStringMap(secret)
	record := &QuotaRecord{
		Path:      stringVal(data, "path"),
		MaxWrites: intVal(data, "max_writes"),
		WindowSec: intVal(data, "window_sec"),
		Writes:    intVal(data, "writes"),
	}

	if raw, ok := data["reset_at"]; ok {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			record.ResetAt = t
		}
	}

	return record, nil
}

func stringVal(m map[string]string, key string) string { return m[key] }

func intVal(m map[string]string, key string) int {
	var n int
	fmt.Sscanf(m[key], "%d", &n)
	return n
}
