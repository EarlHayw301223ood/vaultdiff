package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// LintMount runs lint rules against all secrets under a mount prefix.
func LintMount(ctx context.Context, client *vaultapi.Client, mount string, rules []LintRule) ([]LintResult, error) {
	if mount == "" {
		return nil, fmt.Errorf("mount must not be empty")
	}

	paths, err := ListTree(ctx, client, mount)
	if err != nil {
		return nil, fmt.Errorf("list tree: %w", err)
	}

	var results []LintResult
	for _, path := range paths {
		secret, err := client.KVv2(mount).Get(ctx, path)
		if err != nil || secret == nil {
			continue
		}
		data := extractStringMap(secret.Data)
		r := LintSecret(path, data, rules)
		if !r.Clean() {
			results = append(results, r)
		}
	}

	return results, nil
}

// LintSummary returns counts of issues per rule.
func LintSummary(results []LintResult) map[LintRule]int {
	counts := make(map[LintRule]int)
	for _, r := range results {
		for _, issue := range r.Issues {
			counts[issue.Rule]++
		}
	}
	return counts
}
