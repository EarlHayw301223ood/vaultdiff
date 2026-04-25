package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// VersionComparison holds the result of comparing two specific versions of a secret.
type VersionComparison struct {
	Path     string
	VersionA int
	VersionB int
	DataA    map[string]string
	DataB    map[string]string
}

// CompareVersions fetches two versions of the same secret path and returns
// the raw data maps for downstream diffing.
func CompareVersions(client *api.Client, path string, versionA, versionB int) (*VersionComparison, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("path must not be empty")
	}
	if versionA <= 0 {
		return nil, fmt.Errorf("versionA must be a positive integer, got %d", versionA)
	}
	if versionB <= 0 {
		return nil, fmt.Errorf("versionB must be a positive integer, got %d", versionB)
	}
	if versionA == versionB {
		return nil, fmt.Errorf("versionA and versionB must differ")
	}

	dataA, err := fetchVersionData(client, path, versionA)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d: %w", versionA, err)
	}

	dataB, err := fetchVersionData(client, path, versionB)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d: %w", versionB, err)
	}

	return &VersionComparison{
		Path:     path,
		VersionA: versionA,
		VersionB: versionB,
		DataA:    dataA,
		DataB:    dataB,
	}, nil
}

func fetchVersionData(client *api.Client, path string, version int) (map[string]string, error) {
	mount, subPath, _ := strings.Cut(strings.TrimPrefix(path, "/"), "/")
	kvPath := fmt.Sprintf("%s/data/%s", mount, subPath)

	secret, err := client.Logical().ReadWithData(kvPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, err
	}
	return extractStringMap(secret), nil
}
