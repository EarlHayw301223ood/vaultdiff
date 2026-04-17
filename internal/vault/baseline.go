package vault

import (
	"errors"
	"fmt"
	"time"
)

// Baseline represents a named snapshot saved for future comparison.
type Baseline struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	Data      map[string]string `json:"data"`
}

const baselineMetaPrefix = "vaultdiff/baselines"

func baselineMetaPath(path, name string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%s/%s/%s", baselineMetaPrefix, path, name)
}

// SaveBaseline stores the current secret version as a named baseline.
func SaveBaseline(client *Client, path, name string, version int) (*Baseline, error) {
	if path == "" {
		return nil, errors.New("path is required")
	}
	if name == "" {
		return nil, errors.New("baseline name is required")
	}
	if version < 1 {
		return nil, errors.New("version must be >= 1")
	}

	secret, err := FetchAtRef(client, path, version)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	if secret == nil {
		return nil, errors.New("secret not found")
	}

	data := extractStringMap(secret)
	bl := &Baseline{
		Name:      name,
		Path:      path,
		Version:   version,
		CreatedAt: time.Now().UTC(),
		Data:      data,
	}

	payload := map[string]interface{}{
		"name":       bl.Name,
		"path":       bl.Path,
		"version":    bl.Version,
		"created_at": bl.CreatedAt.Format(time.RFC3339),
		"data":       bl.Data,
	}

	metaPath := kvV2WritePath(baselineMetaPath(path, name))
	_, err = client.Logical().Write(metaPath, map[string]interface{}{"data": payload})
	if err != nil {
		return nil, fmt.Errorf("write baseline: %w", err)
	}
	return bl, nil
}

// GetBaseline retrieves a previously saved baseline by name.
func GetBaseline(client *Client, path, name string) (*Baseline, error) {
	if path == "" {
		return nil, errors.New("path is required")
	}
	if name == "" {
		return nil, errors.New("baseline name is required")
	}

	readPath := fmt.Sprintf("%s/data/%s", strings.SplitN(path, "/", 2)[0], baselineMetaPath(path, name))
	secret, err := client.Logical().Read(readPath)
	if err != nil {
		return nil, fmt.Errorf("read baseline: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("baseline %q not found for path %q", name, path)
	}

	inner := extractStringMap(secret)
	bl := &Baseline{
		Name:    name,
		Path:    path,
		Data:    inner,
	}
	return bl, nil
}
