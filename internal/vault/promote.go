package vault

import (
	"context"
	"fmt"
)

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	SourcePath  string
	DestPath    string
	SourceVersion int
	DestVersion   int
	Keys          int
}

// Promote copies the secret at sourcePath (at the given version ref) into
// destPath, writing it as a new version. It returns a PromoteResult
// describing what was written.
func Promote(ctx context.Context, client *Client, sourcePath, destPath, ref string) (*PromoteResult, error) {
	if sourcePath == "" || destPath == "" {
		return nil, fmt.Errorf("promote: source and destination paths must not be empty")
	}

	secret, version, err := FetchAtRef(ctx, client, sourcePath, ref)
	if err != nil {
		return nil, fmt.Errorf("promote: fetch source %q@%s: %w", sourcePath, ref, err)
	}

	data := extractStringMap(secret)
	if len(data) == 0 {
		return nil, fmt.Errorf("promote: source secret %q has no data", sourcePath)
	}

	payload := map[string]interface{}{"data": toInterfaceMap(data)}
	mountedPath := kvV2WritePath(destPath)

	writeSecret, err := client.Logical().WriteWithContext(ctx, mountedPath, payload)
	if err != nil {
		return nil, fmt.Errorf("promote: write dest %q: %w", destPath, err)
	}

	destVersion := 0
	if writeSecret != nil {
		if v, ok := writeSecret.Data["version"]; ok {
			if n, ok := v.(int); ok {
				destVersion = n
			}
		}
	}

	return &PromoteResult{
		SourcePath:    sourcePath,
		DestPath:      destPath,
		SourceVersion: version,
		DestVersion:   destVersion,
		Keys:          len(data),
	}, nil
}

// toInterfaceMap converts map[string]string to map[string]interface{}.
func toInterfaceMap(m map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// kvV2WritePath converts a logical secret path to its KV v2 write path.
func kvV2WritePath(path string) string {
	mount := mountPrefix(path)
	if mount == "" {
		return path
	}
	suffix := path[len(mount)+1:]
	return fmt.Sprintf("%s/data/%s", mount, suffix)
}
