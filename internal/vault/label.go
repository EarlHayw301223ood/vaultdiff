package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func labelMetaPath(path, labelName string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%s/__meta__/labels/%s", path, labelName)
}

// LabelEntry holds the value and metadata for a label attached to a secret version.
type LabelEntry struct {
	Path    string    `json:"path"`
	Name    string    `json:"name"`
	Value   string    `json:"value"`
	Version int       `json:"version"`
	SetAt   time.Time `json:"set_at"`
}

// SetLabel attaches a named label with a string value to a specific version of a secret.
func SetLabel(client *Client, path, name, value string, version int) (*LabelEntry, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if name == "" {
		return nil, errors.New("label name must not be empty")
	}
	if version <= 0 {
		return nil, errors.New("version must be a positive integer")
	}

	entry := &LabelEntry{
		Path:    path,
		Name:    name,
		Value:   value,
		Version: version,
		SetAt:   time.Now().UTC(),
	}

	data := map[string]interface{}{
		"value":   entry.Value,
		"version": entry.Version,
		"set_at":  entry.SetAt.Format(time.RFC3339),
	}

	writePath := kvV2WritePath(labelMetaPath(path, name))
	_, err := client.Logical().Write(writePath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("writing label: %w", err)
	}
	return entry, nil
}

// GetLabel retrieves a named label for a secret path.
func GetLabel(client *Client, path, name string) (*LabelEntry, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if name == "" {
		return nil, errors.New("label name must not be empty")
	}

	readPath := kvV2ReadPath(labelMetaPath(path, name))
	secret, err := client.Logical().Read(readPath)
	if err != nil {
		return nil, fmt.Errorf("reading label: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("label %q not found for path %q", name, path)
	}

	data := extractStringMap(secret)
	version := 0
	if v, ok := secret.Data["version"]; ok {
		if vf, ok := v.(float64); ok {
			version = int(vf)
		}
	}

	var setAt time.Time
	if s, ok := data["set_at"]; ok {
		setAt, _ = time.Parse(time.RFC3339, s)
	}

	return &LabelEntry{
		Path:    path,
		Name:    name,
		Value:   data["value"],
		Version: version,
		SetAt:   setAt,
	}, nil
}
