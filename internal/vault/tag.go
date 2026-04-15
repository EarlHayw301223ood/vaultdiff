package vault

import (
	"fmt"
	"strings"
)

// Tag represents a named pointer to a specific secret version.
type Tag struct {
	Name    string
	Path    string
	Version int
}

// tagMetaPath returns the internal Vault path used to store tag metadata.
func tagMetaPath(mount, secretPath, tag string) string {
	return fmt.Sprintf("%s/metadata/tags/%s/%s", mount, secretPath, tag)
}

// SetTag writes a named tag pointing to a specific version of a secret.
// Tags are stored as KV metadata entries under a reserved "tags/" namespace.
func SetTag(client *Client, mount, secretPath, tagName string, version int) error {
	if strings.TrimSpace(secretPath) == "" {
		return fmt.Errorf("secret path must not be empty")
	}
	if strings.TrimSpace(tagName) == "" {
		return fmt.Errorf("tag name must not be empty")
	}
	if version < 1 {
		return fmt.Errorf("version must be >= 1, got %d", version)
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"version": version,
			"path":    secretPath,
		},
	}

	writePath := kvV2WritePath(mount, fmt.Sprintf("tags/%s/%s", secretPath, tagName))
	_, err := client.Logical().Write(writePath, data)
	if err != nil {
		return fmt.Errorf("set tag %q on %s: %w", tagName, secretPath, err)
	}
	return nil
}

// GetTag retrieves the Tag record for the given tag name on a secret path.
func GetTag(client *Client, mount, secretPath, tagName string) (*Tag, error) {
	if strings.TrimSpace(secretPath) == "" {
		return nil, fmt.Errorf("secret path must not be empty")
	}
	if strings.TrimSpace(tagName) == "" {
		return nil, fmt.Errorf("tag name must not be empty")
	}

	readPath := fmt.Sprintf("%s/data/tags/%s/%s", mount, secretPath, tagName)
	secret, err := client.Logical().Read(readPath)
	if err != nil {
		return nil, fmt.Errorf("get tag %q on %s: %w", tagName, secretPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("tag %q not found for path %s", tagName, secretPath)
	}

	data := extractStringMap(secret)
	versionRaw, ok := data["version"]
	if !ok {
		return nil, fmt.Errorf("tag %q missing version field", tagName)
	}

	var version int
	switch v := interface{}(versionRaw).(type) {
	case float64:
		version = int(v)
	case int:
		version = v
	default:
		return nil, fmt.Errorf("tag %q has unexpected version type %T", tagName, versionRaw)
	}

	return &Tag{Name: tagName, Path: secretPath, Version: version}, nil
}
