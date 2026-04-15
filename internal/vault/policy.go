package vault

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// PolicyAccess represents the access level for a given path.
type PolicyAccess struct {
	Path         string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

// PolicyChecker checks whether the authenticated token has access to a set of paths.
type PolicyChecker struct {
	client LogicalClient
}

// NewPolicyChecker creates a new PolicyChecker using the provided logical client.
func NewPolicyChecker(client LogicalClient) *PolicyChecker {
	return &PolicyChecker{client: client}
}

// CheckPaths queries Vault's sys/capabilities-self endpoint for each path
// and returns a slice of PolicyAccess results sorted by path.
func (p *PolicyChecker) CheckPaths(ctx context.Context, paths []string) ([]PolicyAccess, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	results := make([]PolicyAccess, 0, len(paths))

	for _, path := range paths {
		caps, err := p.fetchCapabilities(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("checking capabilities for %q: %w", path, err)
		}
		results = append(results, PolicyAccess{
			Path:         path,
			Capabilities: caps,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	return results, nil
}

// HasCapability returns true if the given PolicyAccess includes the specified capability.
func HasCapability(access PolicyAccess, capability string) bool {
	for _, c := range access.Capabilities {
		if strings.EqualFold(c, capability) {
			return true
		}
	}
	return false
}

func (p *PolicyChecker) fetchCapabilities(ctx context.Context, path string) ([]string, error) {
	secret, err := p.client.WriteWithContext(ctx, "sys/capabilities-self", map[string]interface{}{
		"paths": []string{path},
	})
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return []string{"deny"}, nil
	}

	raw, ok := secret.Data[path]
	if !ok {
		return []string{"deny"}, nil
	}

	iface, ok := raw.([]interface{})
	if !ok {
		return []string{"deny"}, nil
	}

	caps := make([]string, 0, len(iface))
	for _, v := range iface {
		if s, ok := v.(string); ok {
			caps = append(caps, s)
		}
	}
	return caps, nil
}
