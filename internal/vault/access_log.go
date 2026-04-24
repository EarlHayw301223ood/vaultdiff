package vault

import (
	"fmt"
	"strings"
	"time"
)

// AccessEntry records a single access event for a secret path.
type AccessEntry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	Operation string    `json:"operation"`
	Actor     string    `json:"actor"`
	Note      string    `json:"note,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// accessLogMetaPath returns the metadata path used to store the access log.
func accessLogMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/access-log/%s", path)
}

// AppendAccessLog writes a new AccessEntry to the access log stored in Vault
// metadata for the given secret path.
func AppendAccessLog(client LogicalClient, path string, entry AccessEntry) error {
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}
	if entry.Operation == "" {
		return fmt.Errorf("operation must not be empty")
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	metaPath := accessLogMetaPath(path)

	data := map[string]interface{}{
		"path":      entry.Path,
		"version":   entry.Version,
		"operation": entry.Operation,
		"actor":     entry.Actor,
		"note":      entry.Note,
		"timestamp": entry.Timestamp.Format(time.RFC3339),
	}

	_, err := client.Write(metaPath, data)
	return err
}

// GetAccessLog reads the most recent access log entry for the given path.
func GetAccessLog(client LogicalClient, path string) (*AccessEntry, error) {
	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}

	metaPath := accessLogMetaPath(path)
	secret, err := client.Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read access log: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	entry := &AccessEntry{}
	if v, ok := secret.Data["path"].(string); ok {
		entry.Path = v
	}
	if v, ok := secret.Data["operation"].(string); ok {
		entry.Operation = v
	}
	if v, ok := secret.Data["actor"].(string); ok {
		entry.Actor = v
	}
	if v, ok := secret.Data["note"].(string); ok {
		entry.Note = v
	}
	if v, ok := secret.Data["timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			entry.Timestamp = t
		}
	}

	return entry, nil
}
