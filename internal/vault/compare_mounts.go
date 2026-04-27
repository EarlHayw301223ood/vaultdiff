package vault

import (
	"fmt"
	"strings"
)

// MountCompareResult holds the outcome of comparing two KV mounts.
type MountCompareResult struct {
	MountA    string
	MountB    string
	OnlyInA   []string
	OnlyInB   []string
	Diverged  []string
	InSync    []string
}

// Summary returns a human-readable one-liner for the comparison.
func (r *MountCompareResult) Summary() string {
	return fmt.Sprintf(
		"only-in-%s=%d only-in-%s=%d diverged=%d in-sync=%d",
		r.MountA, len(r.OnlyInA),
		r.MountB, len(r.OnlyInB),
		len(r.Diverged),
		len(r.InSync),
	)
}

// CompareMounts walks two KV v2 mounts and returns a MountCompareResult.
// prefix filters paths to those that start with the given string.
func CompareMounts(c Logical, mountA, mountB, prefix string) (*MountCompareResult, error) {
	if mountA == "" {
		return nil, fmt.Errorf("mountA must not be empty")
	}
	if mountB == "" {
		return nil, fmt.Errorf("mountB must not be empty")
	}

	snapshotA, err := SnapshotMount(c, mountA)
	if err != nil {
		return nil, fmt.Errorf("snapshot %s: %w", mountA, err)
	}
	snapshotB, err := SnapshotMount(c, mountB)
	if err != nil {
		return nil, fmt.Errorf("snapshot %s: %w", mountB, err)
	}

	result := &MountCompareResult{MountA: mountA, MountB: mountB}

	pathsA := filterPaths(snapshotA.Paths(), prefix)
	pathsB := filterPaths(snapshotB.Paths(), prefix)

	setB := make(map[string]bool, len(pathsB))
	for _, p := range pathsB {
		setB[p] = true
	}
	setA := make(map[string]bool, len(pathsA))
	for _, p := range pathsA {
		setA[p] = true
	}

	for _, p := range pathsA {
		if !setB[p] {
			result.OnlyInA = append(result.OnlyInA, p)
			continue
		}
		dataA, _ := snapshotA.Get(p)
		dataB, _ := snapshotB.Get(p)
		if mapsEqual(dataA, dataB) {
			result.InSync = append(result.InSync, p)
		} else {
			result.Diverged = append(result.Diverged, p)
		}
	}
	for _, p := range pathsB {
		if !setA[p] {
			result.OnlyInB = append(result.OnlyInB, p)
		}
	}
	return result, nil
}

func filterPaths(paths []string, prefix string) []string {
	if prefix == "" {
		return paths
	}
	var out []string
	for _, p := range paths {
		if strings.HasPrefix(p, prefix) {
			out = append(out, p)
		}
	}
	return out
}
