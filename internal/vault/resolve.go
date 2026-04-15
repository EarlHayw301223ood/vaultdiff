package vault

import (
	"fmt"
	"strconv"
	"strings"
)

// VersionRef represents a resolved version reference for a secret path.
type VersionRef struct {
	Path    string
	Version int
	IsLatest bool
}

// ResolveVersionRef parses a path that may contain an optional version suffix.
// Accepted formats:
//   - "secret/myapp"         → latest version
//   - "secret/myapp@3"       → explicit version 3
//   - "secret/myapp@latest"  → latest version (explicit)
func ResolveVersionRef(raw string) (VersionRef, error) {
	if raw == "" {
		return VersionRef{}, fmt.Errorf("path must not be empty")
	}

	parts := strings.SplitN(raw, "@", 2)
	path := parts[0]

	if path == "" {
		return VersionRef{}, fmt.Errorf("path segment before '@' must not be empty")
	}

	if len(parts) == 1 {
		return VersionRef{Path: path, IsLatest: true}, nil
	}

	versionStr := parts[1]
	if versionStr == "latest" || versionStr == "" {
		return VersionRef{Path: path, IsLatest: true}, nil
	}

	v, err := strconv.Atoi(versionStr)
	if err != nil || v < 1 {
		return VersionRef{}, fmt.Errorf("invalid version %q: must be a positive integer or 'latest'", versionStr)
	}

	return VersionRef{Path: path, Version: v, IsLatest: false}, nil
}

// String returns the canonical string representation of a VersionRef.
func (r VersionRef) String() string {
	if r.IsLatest {
		return r.Path + "@latest"
	}
	return fmt.Sprintf("%s@%d", r.Path, r.Version)
}
