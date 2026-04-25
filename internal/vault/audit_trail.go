package vault

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

// AuditTrailEntry represents a single recorded event in the audit trail.
type AuditTrailEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Actor     string    `json:"actor"`
	Version   int       `json:"version"`
	Note      string    `json:"note,omitempty"`
}

// AuditTrail holds all entries for a given secret path.
type AuditTrail struct {
	Path    string            `json:"path"`
	Entries []AuditTrailEntry `json:"entries"`
}

func auditTrailMetaPath(path string) string {
	path = strings.Trim(path, "/")
	segments := strings.SplitN(path, "/", 2)
	if len(segments) < 2 {
		return fmt.Sprintf("%s/metadata/_audit_trail/%s", path, path)
	}
	return fmt.Sprintf("%s/metadata/_audit_trail/%s", segments[0], segments[1])
}

// AppendAuditTrail records a new audit trail entry for the given secret path.
func AppendAuditTrail(client *api.Client, path, operation, actor string, version int, note string) error {
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}
	if operation == "" {
		return fmt.Errorf("operation must not be empty")
	}
	if version < 1 {
		return fmt.Errorf("version must be >= 1")
	}

	metaPath := auditTrailMetaPath(path)

	existing, err := GetAuditTrail(client, path)
	if err != nil {
		existing = &AuditTrail{Path: path, Entries: []AuditTrailEntry{}}
	}

	entry := AuditTrailEntry{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		Actor:     actor,
		Version:   version,
		Note:      note,
	}
	existing.Entries = append(existing.Entries, entry)

	data := map[string]interface{}{
		"path":    existing.Path,
		"entries": existing.Entries,
	}
	_, err = client.Logical().Write(metaPath, map[string]interface{}{"data": data})
	return err
}

// GetAuditTrail retrieves the full audit trail for a secret path.
func GetAuditTrail(client *api.Client, path string) (*AuditTrail, error) {
	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}

	metaPath := auditTrailMetaPath(path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read audit trail: %w", err)
	}
	if secret == nil {
		return &AuditTrail{Path: path, Entries: []AuditTrailEntry{}}, nil
	}

	trail := &AuditTrail{Path: path}
	data := extractStringMap(secret)

	if raw, ok := secret.Data["data"]; ok {
		if m, ok := raw.(map[string]interface{}); ok {
			data = m
		}
	}
	_ = data

	if rawEntries, ok := data["entries"]; ok {
		if entries, ok := rawEntries.([]interface{}); ok {
			for _, e := range entries {
				if em, ok := e.(map[string]interface{}); ok {
					var entry AuditTrailEntry
					if v, ok := em["operation"].(string); ok {
						entry.Operation = v
					}
					if v, ok := em["actor"].(string); ok {
						entry.Actor = v
					}
					if v, ok := em["note"].(string); ok {
						entry.Note = v
					}
					trail.Entries = append(trail.Entries, entry)
				}
			}
		}
	}
	return trail, nil
}
