package vault

import (
	"context"
	"fmt"
	"path"
)

// SecretVersion represents a single version of a Vault KV secret.
type SecretVersion struct {
	Path    string
	Version int
	Data    map[string]string
	Metadata SecretMetadata
}

// SecretMetadata holds version metadata returned by Vault.
type SecretMetadata struct {
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
	Version      int
}

// GetSecretVersion retrieves a specific version of a KV v2 secret.
// If version is 0, the latest version is returned.
func (c *Client) GetSecretVersion(ctx context.Context, mountPath, secretPath string, version int) (*SecretVersion, error) {
	kvPath := path.Join(mountPath, "data", secretPath)

	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := c.logical.ReadWithDataWithContext(ctx, kvPath, params)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q (version %d): %w", secretPath, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q not found", secretPath)
	}

	data, err := extractStringMap(secret.Data, "data")
	if err != nil {
		return nil, fmt.Errorf("extracting data from secret %q: %w", secretPath, err)
	}

	meta, err := extractMetadata(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("extracting metadata from secret %q: %w", secretPath, err)
	}

	return &SecretVersion{
		Path:     secretPath,
		Version:  meta.Version,
		Data:     data,
		Metadata: meta,
	}, nil
}

func extractStringMap(raw map[string]interface{}, key string) (map[string]string, error) {
	v, ok := raw[key]
	if !ok {
		return nil, fmt.Errorf("key %q not present", key)
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("key %q is not a map", key)
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		out[k] = fmt.Sprintf("%v", val)
	}
	return out, nil
}

func extractMetadata(raw map[string]interface{}) (SecretMetadata, error) {
	v, ok := raw["metadata"]
	if !ok {
		return SecretMetadata{}, fmt.Errorf("metadata key not present")
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return SecretMetadata{}, fmt.Errorf("metadata is not a map")
	}
	meta := SecretMetadata{}
	if ct, ok := m["created_time"].(string); ok {
		meta.CreatedTime = ct
	}
	if dt, ok := m["deletion_time"].(string); ok {
		meta.DeletionTime = dt
	}
	if d, ok := m["destroyed"].(bool); ok {
		meta.Destroyed = d
	}
	if ver, ok := m["version"].(float64); ok {
		meta.Version = int(ver)
	}
	return meta, nil
}
