package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const retentionMetaSuffix = "__vaultdiff_retention__"

// RetentionPolicy defines how long versions of a secret should be kept.
type RetentionPolicy struct {
	Path       string        `json:"path"`
	Version    int           `json:"version"`
	MaxAge     time.Duration `json:"max_age_ns"`
	MaxVersions int          `json:"max_versions"`
	SetBy      string        `json:"set_by"`
	SetAt      time.Time     `json:"set_at"`
}

func retentionMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%s/%s", path, retentionMetaSuffix)
}

// SetRetention writes a retention policy for the given secret path and version.
func SetRetention(client LogicalClient, path string, version int, maxAge time.Duration, maxVersions int, setBy string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if version <= 0 {
		return errors.New("version must be a positive integer")
	}
	if maxAge < 0 {
		return errors.New("max_age must not be negative")
	}
	if maxVersions < 0 {
		return errors.New("max_versions must not be negative")
	}
	if setBy == "" {
		setBy = "unknown"
	}

	policy := RetentionPolicy{
		Path:        path,
		Version:     version,
		MaxAge:      maxAge,
		MaxVersions: maxVersions,
		SetBy:       setBy,
		SetAt:       time.Now().UTC(),
	}

	data := map[string]interface{}{
		"path":         policy.Path,
		"version":      policy.Version,
		"max_age_ns":   int64(policy.MaxAge),
		"max_versions": policy.MaxVersions,
		"set_by":       policy.SetBy,
		"set_at":       policy.SetAt.Format(time.RFC3339),
	}

	metaPath := retentionMetaPath(path)
	_, err := client.Write(kvV2WritePath(metaPath), map[string]interface{}{"data": data})
	return err
}

// GetRetention reads the retention policy for a secret path.
func GetRetention(client LogicalClient, path string) (*RetentionPolicy, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	metaPath := retentionMetaPath(path)
	secret, err := client.Read(kvV2DataPath(metaPath))
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no retention policy found for %q", path)
	}

	m := extractStringMap(secret)

	parsedAt, _ := time.Parse(time.RFC3339, m["set_at"])
	var version int
	fmt.Sscanf(m["version"], "%d", &version)
	var maxVersions int
	fmt.Sscanf(m["max_versions"], "%d", &maxVersions)
	var maxAgeNs int64
	fmt.Sscanf(m["max_age_ns"], "%d", &maxAgeNs)

	return &RetentionPolicy{
		Path:        m["path"],
		Version:     version,
		MaxAge:      time.Duration(maxAgeNs),
		MaxVersions: maxVersions,
		SetBy:       m["set_by"],
		SetAt:       parsedAt,
	}, nil
}
