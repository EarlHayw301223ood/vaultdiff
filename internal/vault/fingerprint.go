package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
)

// FingerprintResult holds the computed fingerprint for a secret version.
type FingerprintResult struct {
	Path        string
	Version     int
	Fingerprint string
	KeyCount    int
}

fingerprintMetaPath := func(path, name string) string {
	return fmt.Sprintf("%s/__meta__/fingerprint/%s", strings.Trim(path, "/"), name)
}

// ComputeFingerprint returns a stable SHA-256 digest of the key-value pairs
// at the given path and version. The digest is order-independent.
func ComputeFingerprint(client *api.Client, path string, version int) (*FingerprintResult, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("path must not be empty")
	}
	if version < 1 {
		return nil, fmt.Errorf("version must be >= 1, got %d", version)
	}

	data, err := fetchVersionData(client, path, version)
	if err != nil {
		return nil, fmt.Errorf("fetch version data: %w", err)
	}

	fingerprint, keyCount := hashSecretData(data)
	return &FingerprintResult{
		Path:        path,
		Version:     version,
		Fingerprint: fingerprint,
		KeyCount:    keyCount,
	}, nil
}

// hashSecretData produces a deterministic hex digest from a map of secret data.
func hashSecretData(data map[string]string) (string, int) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, data[k])
	}
	return hex.EncodeToString(h.Sum(nil)), len(keys)
}
