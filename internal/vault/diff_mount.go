package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// MountDiffResult holds the outcome of comparing two mount snapshots.
type MountDiffResult struct {
	// OnlyInA contains paths present in mount A but not mount B.
	OnlyInA []string
	// OnlyInB contains paths present in mount B but not mount A.
	OnlyInB []string
	// Modified contains paths present in both mounts whose secret data differs.
	Modified []string
	// Identical contains paths whose data is byte-for-byte equal across mounts.
	Identical []string
}

// Total returns the number of paths examined across both mounts.
func (r *MountDiffResult) Total() int {
	return len(r.OnlyInA) + len(r.OnlyInB) + len(r.Modified) + len(r.Identical)
}

// HasDifferences returns true when any path differs between the two mounts.
func (r *MountDiffResult) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Modified) > 0
}

// DiffMountOptions controls the behaviour of DiffMounts.
type DiffMountOptions struct {
	// Prefix restricts comparison to paths that begin with the given string.
	Prefix string
}

// DiffMounts compares every secret path under mountA against mountB using the
// provided Vault client. Secrets are fetched at their latest version.
//
// Both mountA and mountB should be KV-v2 mount names, e.g. "secret".
func DiffMounts(client *api.Client, mountA, mountB string, opts DiffMountOptions) (*MountDiffResult, error) {
	if mountA == "" {
		return nil, fmt.Errorf("mountA must not be empty")
	}
	if mountB == "" {
		return nil, fmt.Errorf("mountB must not be empty")
	}

	snapshotA, err := SnapshotMount(client, mountA)
	if err != nil {
		return nil, fmt.Errorf("snapshot mount %q: %w", mountA, err)
	}

	snapshotB, err := SnapshotMount(client, mountB)
	if err != nil {
		return nil, fmt.Errorf("snapshot mount %q: %w", mountB, err)
	}

	return diffSnapshots(snapshotA, snapshotB, opts), nil
}

// diffSnapshots performs the in-memory comparison between two Snapshots.
func diffSnapshots(a, b *Snapshot, opts DiffMountOptions) *MountDiffResult {
	result := &MountDiffResult{}

	pathsA := a.Paths()
	pathsB := b.Paths()

	setB := make(map[string]struct{}, len(pathsB))
	for _, p := range pathsB {
		setB[p] = struct{}{}
	}

	setA := make(map[string]struct{}, len(pathsA))
	for _, p := range pathsA {
		setA[p] = struct{}{}
	}

	for _, path := range pathsA {
		if opts.Prefix != "" && !hasPrefix(path, opts.Prefix) {
			continue
		}
		if _, ok := setB[path]; !ok {
			result.OnlyInA = append(result.OnlyInA, path)
			continue
		}
		dataA, _ := a.Get(path)
		dataB, _ := b.Get(path)
		if mapsEqual(dataA, dataB) {
			result.Identical = append(result.Identical, path)
		} else {
			result.Modified = append(result.Modified, path)
		}
	}

	for _, path := range pathsB {
		if opts.Prefix != "" && !hasPrefix(path, opts.Prefix) {
			continue
		}
		if _, ok := setA[path]; !ok {
			result.OnlyInB = append(result.OnlyInB, path)
		}
	}

	return result
}

// hasPrefix reports whether path starts with prefix, treating prefix as a
// path segment boundary (i.e. prefix "foo" matches "foo/bar" but not "foobar").
func hasPrefix(path, prefix string) bool {
	if len(path) < len(prefix) {
		return false
	}
	if path == prefix {
		return true
	}
	if path[:len(prefix)] == prefix {
		return len(path) > len(prefix) && path[len(prefix)] == '/'
	}
	return false
}
