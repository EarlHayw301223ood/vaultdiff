package vault

import (
	"context"
	"fmt"
	"sort"
)

// VersionMeta holds metadata about a specific secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns all available version metadata for a KV v2 secret path.
func (c *Client) ListVersions(ctx context.Context, mount, secretPath string) ([]VersionMeta, error) {
	path := fmt.Sprintf("%s/metadata/%s", mount, secretPath)

	secret, err := c.logical.ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %q: %w", secretPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", secretPath)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata for %q", secretPath)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for %q", secretPath)
	}

	var versions []VersionMeta
	for vKey, vVal := range versionsMap {
		var vNum int
		fmt.Sscanf(vKey, "%d", &vNum)

		vMap, ok := vVal.(map[string]interface{})
		if !ok {
			continue
		}

		meta := VersionMeta{Version: vNum}
		if ct, ok := vMap["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := vMap["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if destroyed, ok := vMap["destroyed"].(bool); ok {
			meta.Destroyed = destroyed
		}
		versions = append(versions, meta)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version < versions[j].Version
	})

	return versions, nil
}

// LatestVersion returns the highest available version number for a secret path.
func (c *Client) LatestVersion(ctx context.Context, mount, secretPath string) (int, error) {
	versions, err := c.ListVersions(ctx, mount, secretPath)
	if err != nil {
		return 0, err
	}
	if len(versions) == 0 {
		return 0, fmt.Errorf("no versions found for %q", secretPath)
	}
	return versions[len(versions)-1].Version, nil
}
