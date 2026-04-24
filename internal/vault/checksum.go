package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ChecksumResult holds the computed checksum for a secret version.
type ChecksumResult struct {
	Path    string
	Version int
	Checksum string
}

// checksumMetaPath returns the metadata path used to store a checksum.
func checksumMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/checksum", path)
}

// ComputeChecksum returns a deterministic SHA-256 hex digest over the
// key-value pairs of data. Keys are sorted before hashing so that
// insertion order does not affect the result.
func ComputeChecksum(data map[string]string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, data[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// SaveChecksum writes the checksum for the given path and version into
// the Vault metadata namespace so it can be verified later.
func SaveChecksum(c *Client, path string, version int, checksum string) error {
	if path == "" {
		return errors.New("checksum: path must not be empty")
	}
	if version <= 0 {
		return errors.New("checksum: version must be a positive integer")
	}
	if checksum == "" {
		return errors.New("checksum: checksum must not be empty")
	}

	metaPath := checksumMetaPath(path)
	key := fmt.Sprintf("v%d", version)
	_, err := c.Logical().Write(metaPath, map[string]interface{}{
		key: checksum,
	})
	return err
}

// GetChecksum retrieves the stored checksum for path at the given version.
func GetChecksum(c *Client, path string, version int) (string, error) {
	if path == "" {
		return "", errors.New("checksum: path must not be empty")
	}
	if version <= 0 {
		return "", errors.New("checksum: version must be a positive integer")
	}

	metaPath := checksumMetaPath(path)
	secret, err := c.Logical().Read(metaPath)
	if err != nil {
		return "", err
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("checksum: no entry found for %s", path)
	}

	key := fmt.Sprintf("v%d", version)
	v, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("checksum: version %d not found for %s", version, path)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("checksum: unexpected type for stored checksum")
	}
	return s, nil
}
