package vault

import (
	"context"
	"fmt"
)

// VersionPair holds two secret versions to be compared.
type VersionPair struct {
	PathA    string
	VersionA int
	PathB    string
	VersionB int
}

// SecretPair holds the resolved secret data for two versions.
type SecretPair struct {
	A        map[string]string
	B        map[string]string
	MetaA    SecretMetadata
	MetaB    SecretMetadata
}

// FetchPair retrieves two secret versions (possibly from different paths)
// and returns them as a SecretPair ready for diffing.
func FetchPair(ctx context.Context, client *Client, pair VersionPair) (*SecretPair, error) {
	secA, err := client.ReadSecretVersion(ctx, pair.PathA, pair.VersionA)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d of %q: %w", pair.VersionA, pair.PathA, err)
	}

	secB, err := client.ReadSecretVersion(ctx, pair.PathB, pair.VersionB)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d of.VersionB, pair.PathB, err)
	}

	return &SecretPair{
		A:     secA.Data,
		B:     secB.Data,
		MetaA: secA.Metadata,
		MetaB: secB.Metadata,
	}, nil
}

// SameMount returns true when both paths share the same KV mount prefix.
func SameMount(pathA, pathB string) bool {
	mountA := mountPrefix(pathA)
	mountB := mountPrefix(pathB)
	return mountA == mountB
}

func mountPrefix(path string) string {
	for i, c := range path {
		if c == '/' {
			return path[:i]
		}
	}
	return path
}
