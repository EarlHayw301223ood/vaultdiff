package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Path       string
	FromVersion int
	ToVersion   int
}

// Rollback restores a KV v2 secret at path to the given target version.
// It reads the data at targetVersion and writes it back as a new version.
func Rollback(ctx context.Context, client *api.Client, path string, targetVersion int) (*RollbackResult, error) {
	if path == "" {
		return nil, fmt.Errorf("rollback: path must not be empty")
	}
	if targetVersion < 1 {
		return nil, fmt.Errorf("rollback: targetVersion must be >= 1, got %d", targetVersion)
	}

	// Fetch the secret at the target version.
	secret, err := FetchAtRef(ctx, client, path, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("rollback: fetch version %d of %q: %w", targetVersion, path, err)
	}

	data, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("rollback: no data field in secret at path %q version %d", path, targetVersion)
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("rollback: unexpected data format at path %q", path)
	}

	// Determine current version before writing.
	versions, err := ListVersions(ctx, client, path)
	var currentVersion int
	if err == nil && len(versions) > 0 {
		currentVersion = versions[len(versions)-1].Version
	}

	// Write the old data back as a new version.
	writePath := kvV2WritePath(path)
	_, err = client.Logical().WriteWithContext(ctx, writePath, map[string]interface{}{
		"data": dataMap,
	})
	if err != nil {
		return nil, fmt.Errorf("rollback: write to %q: %w", path, err)
	}

	return &RollbackResult{
		Path:        path,
		FromVersion: currentVersion,
		ToVersion:   targetVersion,
	}, nil
}
