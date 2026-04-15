package vault

import (
	"fmt"
	"sort"
)

// DriftResult holds the comparison of a single secret path across two environments.
type DriftResult struct {
	Path     string
	OnlyInA  bool
	OnlyInB  bool
	InSync   bool
	Diverged bool
}

// DriftReport summarizes drift across all paths in a mount.
type DriftReport struct {
	MountA  string
	MountB  string
	Results []DriftResult
}

// Summary returns counts of each drift category.
func (r *DriftReport) Summary() map[string]int {
	summary := map[string]int{
		"in_sync":   0,
		"diverged":  0,
		"only_in_a": 0,
		"only_in_b": 0,
	}
	for _, res := range r.Results {
		switch {
		case res.InSync:
			summary["in_sync"]++
		case res.Diverged:
			summary["diverged"]++
		case res.OnlyInA:
			summary["only_in_a"]++
		case res.OnlyInB:
			summary["only_in_b"]++
		}
	}
	return summary
}

// DetectDrift compares two snapshots and returns a DriftReport.
func DetectDrift(mountA, mountB string, snapA, snapB *Snapshot) (*DriftReport, error) {
	if snapA == nil || snapB == nil {
		return nil, fmt.Errorf("drift: snapshots must not be nil")
	}

	pathSet := map[string]struct{}{}
	for _, p := range snapA.Paths() {
		pathSet[p] = struct{}{}
	}
	for _, p := range snapB.Paths() {
		pathSet[p] = struct{}{}
	}

	paths := make([]string, 0, len(pathSet))
	for p := range pathSet {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	results := make([]DriftResult, 0, len(paths))
	for _, path := range paths {
		secrA, hasA := snapA.Get(path)
		secrB, hasB := snapB.Get(path)

		result := DriftResult{Path: path}
		switch {
		case hasA && !hasB:
			result.OnlyInA = true
		case !hasA && hasB:
			result.OnlyInB = true
		default:
			if mapsEqual(secrA, secrB) {
				result.InSync = true
			} else {
				result.Diverged = true
			}
		}
		results = append(results, result)
	}

	return &DriftReport{MountA: mountA, MountB: mountB, Results: results}, nil
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
