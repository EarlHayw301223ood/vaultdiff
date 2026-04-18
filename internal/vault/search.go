package vault

import (
	"fmt"
	"strings"
)

// SearchResult holds a matched secret path and the matching keys.
type SearchResult struct {
	Path        string
	MatchedKeys []string
}

// SearchOptions controls how secrets are searched.
type SearchOptions struct {
	KeyPattern   string
	ValuePattern string
	CaseSensitive bool
}

// SearchMount scans all secrets under a mount and returns paths whose
// keys or values match the provided patterns.
func SearchMount(client *Client, mount string, opts SearchOptions) ([]SearchResult, error) {
	paths, err := ListTree(client, mount)
	if err != nil {
		return nil, fmt.Errorf("list tree: %w", err)
	}

	var results []SearchResult
	for _, p := range paths {
		secret, err := client.Logical().Read(kvV2DataPath(mount, p))
		if err != nil || secret == nil {
			continue
		}
		data := extractStringMap(secret)
		var matched []string
		for k, v := range data {
			if matchesPattern(k, opts.KeyPattern, opts.CaseSensitive) ||
				matchesPattern(v, opts.ValuePattern, opts.CaseSensitive) {
				matched = append(matched, k)
			}
		}
		if len(matched) > 0 {
			results = append(results, SearchResult{Path: p, MatchedKeys: matched})
		}
	}
	return results, nil
}

func matchesPattern(s, pattern string, caseSensitive bool) bool {
	if pattern == "" {
		return false
	}
	if !caseSensitive {
		s = strings.ToLower(s)
		pattern = strings.ToLower(pattern)
	}
	return strings.Contains(s, pattern)
}

func kvV2DataPath(mount, path string) string {
	mount = strings.Trim(mount, "/")
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%s/data/%s", mount, path)
}
