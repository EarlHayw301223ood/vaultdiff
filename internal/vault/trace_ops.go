package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// AppendTrace appends a new trace entry to the trace log stored in Vault metadata.
func AppendTrace(client *Client, path, operation, user, note string, version int) error {
	if path == "" {
		return errors.New("trace: path must not be empty")
	}
	if operation == "" {
		return errors.New("trace: operation must not be empty")
	}

	existing, _ := GetTrace(client, path)
	if existing == nil {
		existing = &TraceLog{Path: path}
	}

	entry := TraceEntry{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		Path:      path,
		Version:   version,
		User:      user,
		Note:      note,
	}
	existing.Entries = append(existing.Entries, entry)

	raw, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("trace: marshal: %w", err)
	}

	metaPath := traceMetaPath(path)
	_, err = client.Logical().Write(metaPath, map[string]interface{}{
		"data": string(raw),
	})
	if err != nil {
		return fmt.Errorf("trace: write: %w", err)
	}
	return nil
}

// GetTrace retrieves the trace log for the given secret path.
func GetTrace(client *Client, path string) (*TraceLog, error) {
	if path == "" {
		return nil, errors.New("trace: path must not be empty")
	}

	metaPath := traceMetaPath(path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("trace: read: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return &TraceLog{Path: path}, nil
	}

	raw, ok := secret.Data["data"].(string)
	if !ok || raw == "" {
		return &TraceLog{Path: path}, nil
	}

	var log TraceLog
	if err := json.Unmarshal([]byte(raw), &log); err != nil {
		return nil, fmt.Errorf("trace: unmarshal: %w", err)
	}
	return &log, nil
}
