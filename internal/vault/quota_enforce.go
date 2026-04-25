package vault

import (
	"errors"
	"fmt"
	"time"
)

// ErrQuotaExceeded is returned when a write would exceed the configured quota.
var ErrQuotaExceeded = errors.New("write quota exceeded for path")

// RecordWrite increments the write counter for a path and returns an error
// if the quota has been exceeded. If the window has expired the counter resets.
func RecordWrite(c LogicalClient, path string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}

	record, err := GetQuota(c, path)
	if err != nil {
		return fmt.Errorf("fetch quota: %w", err)
	}
	if record == nil {
		// No quota configured — writes are unrestricted.
		return nil
	}

	now := time.Now().UTC()

	// Reset the window if it has expired.
	if now.After(record.ResetAt) {
		record.Writes = 0
		record.ResetAt = now.Add(time.Duration(record.WindowSec) * time.Second)
	}

	if record.Exceeded() {
		return fmt.Errorf("%w: %s (resets at %s)",
			ErrQuotaExceeded, path, record.ResetAt.Format(time.RFC3339))
	}

	record.Writes++

	update := map[string]interface{}{
		"path":       record.Path,
		"max_writes": record.MaxWrites,
		"window_sec": record.WindowSec,
		"writes":     record.Writes,
		"reset_at":   record.ResetAt.Format(time.RFC3339),
	}
	_, err = c.Write(quotaMetaPath(path), map[string]interface{}{"data": update})
	if err != nil {
		return fmt.Errorf("persist quota: %w", err)
	}

	return nil
}

// ResetQuota zeroes the write counter for a path without changing the policy.
func ResetQuota(c LogicalClient, path string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}

	record, err := GetQuota(c, path)
	if err != nil {
		return fmt.Errorf("fetch quota: %w", err)
	}
	if record == nil {
		return nil
	}

	record.Writes = 0
	record.ResetAt = time.Now().UTC().Add(time.Duration(record.WindowSec) * time.Second)

	update := map[string]interface{}{
		"path":       record.Path,
		"max_writes": record.MaxWrites,
		"window_sec": record.WindowSec,
		"writes":     0,
		"reset_at":   record.ResetAt.Format(time.RFC3339),
	}
	_, err = c.Write(quotaMetaPath(path), map[string]interface{}{"data": update})
	return err
}
