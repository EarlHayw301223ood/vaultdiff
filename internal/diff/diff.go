package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType describes the kind of change for a secret key.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
	ChangeUnchanged ChangeType = "unchanged"
)

// KeyDiff represents the diff result for a single key.
type KeyDiff struct {
	Key      string
	Change   ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret versions.
type Result struct {
	Path    string
	FromVersion int
	ToVersion   int
	Changes []KeyDiff
}

// HasChanges returns true if any keys were added, removed, or modified.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Change != ChangeUnchanged {
			return true
		}
	}
	return false
}

// Compare computes the diff between two maps of secret key/value pairs.
func Compare(path string, fromVersion int, from map[string]string, toVersion int, to map[string]string) *Result {
	result := &Result{
		Path:        path,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
	}

	keys := unionKeys(from, to)
	for _, key := range keys {
		oldVal, inFrom := from[key]
		newVal, inTo := to[key]

		var change ChangeType
		switch {
		case inFrom && !inTo:
			change = ChangeRemoved
		case !inFrom && inTo:
			change = ChangeAdded
		case oldVal != newVal:
			change = ChangeModified
		default:
			change = ChangeUnchanged
		}

		result.Changes = append(result.Changes, KeyDiff{
			Key:      key,
			Change:   change,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}
	return result
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Summary returns a human-readable one-line summary of the diff result.
func (r *Result) Summary() string {
	var added, removed, modified int
	for _, c := range r.Changes {
		switch c.Change {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeModified:
			modified++
		}
	}
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("+%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d removed", removed))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("~%d modified", modified))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}
