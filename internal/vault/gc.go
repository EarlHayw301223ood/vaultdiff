package vault

import (
	"fmt"
	"sort"
	"time"
)

// GCOptions controls which versions are retained during garbage collection.
type GCOptions struct {
	// KeepLast is the minimum number of most-recent versions to retain.
	KeepLast int
	// MaxAge is the maximum age of a version to retain. Zero means no age limit.
	MaxAge time.Duration
	// DryRun reports what would be deleted without performing any writes.
	DryRun bool
}

// GCResult describes the outcome of a garbage-collection run on one path.
type GCResult struct {
	Path            string
	DeletedVersions []int
	RetainedVersions []int
	DryRun          bool
}

// GarbageCollect removes old secret versions at path according to opts.
// It always retains at least opts.KeepLast versions regardless of age.
func GarbageCollect(l LogicalClient, path string, opts GCOptions) (*GCResult, error) {
	if path == "" {
		return nil, fmt.Errorf("gc: path must not be empty")
	}
	if opts.KeepLast < 1 {
		opts.KeepLast = 1
	}

	versions, err := ListVersions(l, path)
	if err != nil {
		return nil, fmt.Errorf("gc: list versions %q: %w", path, err)
	}

	// Sort descending so index 0 is the newest.
	sort.Slice(versions, func(i, j int) bool { return versions[i] > versions[j] })

	now := time.Now().UTC()
	var toDelete, toKeep []int

	for i, v := range versions {
		if i < opts.KeepLast {
			toKeep = append(toKeep, v)
			continue
		}
		if opts.MaxAge > 0 {
			meta, err := fetchVersionMeta(l, path, v)
			if err != nil || now.Sub(meta.CreatedTime) <= opts.MaxAge {
				toKeep = append(toKeep, v)
				continue
			}
		}
		toDelete = append(toDelete, v)
	}

	result := &GCResult{
		Path:             path,
		DeletedVersions:  toDelete,
		RetainedVersions: toKeep,
		DryRun:           opts.DryRun,
	}

	if opts.DryRun || len(toDelete) == 0 {
		return result, nil
	}

	if err := destroyVersions(l, path, toDelete); err != nil {
		return nil, fmt.Errorf("gc: destroy versions at %q: %w", path, err)
	}
	return result, nil
}
