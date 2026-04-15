package vault

import (
	"context"
	"fmt"
	"strings"
)

// ScanResult holds the path and version count for a secret found during a scan.
type ScanResult struct {
	Path     string
	Versions int
	Mount    string
}

// ScanMount walks all secrets under the given mount and returns metadata about
// each secret path, including how many versions exist.
func ScanMount(ctx context.Context, c *Client, mount string) ([]ScanResult, error) {
	if mount == "" {
		return nil, fmt.Errorf("mount must not be empty")
	}

	mount = strings.TrimSuffix(mount, "/")

	paths, err := ListTree(ctx, c, mount)
	if err != nil {
		return nil, fmt.Errorf("scan: list tree: %w", err)
	}

	results := make([]ScanResult, 0, len(paths))
	for _, p := range paths {
		versions, err := ListVersions(ctx, c, p)
		if err != nil {
			// Non-fatal: record zero versions and continue.
			results = append(results, ScanResult{Path: p, Mount: mount, Versions: 0})
			continue
		}
		results = append(results, ScanResult{
			Path:     p,
			Mount:    mount,
			Versions: len(versions),
		})
	}

	return results, nil
}

// FilterScanResults returns only results whose version count satisfies the
// provided predicate.
func FilterScanResults(results []ScanResult, pred func(ScanResult) bool) []ScanResult {
	out := make([]ScanResult, 0)
	for _, r := range results {
		if pred(r) {
			out = append(out, r)
		}
	}
	return out
}
