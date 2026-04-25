package vault

import (
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// ShadowResult holds the outcome of a shadow-write operation.
type ShadowResult struct {
	SourcePath string
	ShadowPath string
	Version    int
	DryRun     bool
}

// shadowMetaPath returns the metadata path used to record shadow mappings.
func shadowMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/shadow/%s", path)
}

// Shadow writes a copy of the secret at sourcePath into shadowMount under the
// same relative path, then records the mapping in metadata. When dryRun is
// true the write is skipped and the result is returned with DryRun set.
func Shadow(client *vaultapi.Client, sourcePath, shadowMount string, version int, dryRun bool) (*ShadowResult, error) {
	if sourcePath == "" {
		return nil, fmt.Errorf("shadow: source path must not be empty")
	}
	if shadowMount == "" {
		return nil, fmt.Errorf("shadow: shadow mount must not be empty")
	}
	if version < 1 {
		return nil, fmt.Errorf("shadow: version must be >= 1, got %d", version)
	}

	secret, err := FetchAtRef(client, sourcePath, fmt.Sprintf("%d", version))
	if err != nil {
		return nil, fmt.Errorf("shadow: fetch source: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("shadow: no secret found at %q version %d", sourcePath, version)
	}

	data := extractStringMap(secret)
	rel := strings.TrimPrefix(strings.Trim(sourcePath, "/"), firstSegment(sourcePath)+"/")
	shadowPath := fmt.Sprintf("%s/data/%s", strings.Trim(shadowMount, "/"), rel)

	result := &ShadowResult{
		SourcePath: sourcePath,
		ShadowPath: shadowPath,
		Version:    version,
		DryRun:     dryRun,
	}

	if dryRun {
		return result, nil
	}

	payload := map[string]interface{}{"data": toInterfaceMap(data)}
	_, err = client.Logical().Write(shadowPath, payload)
	if err != nil {
		return nil, fmt.Errorf("shadow: write to shadow mount: %w", err)
	}

	return result, nil
}

// firstSegment returns the first path segment (the mount name).
func firstSegment(path string) string {
	path = strings.Trim(path, "/")
	if idx := strings.Index(path, "/"); idx >= 0 {
		return path[:idx]
	}
	return path
}
