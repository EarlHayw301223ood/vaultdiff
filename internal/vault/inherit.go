package vault

import (
	"errors"
	"fmt"
	"strings"
)

// InheritResult holds the merged secret data produced by InheritSecret.
type InheritResult struct {
	Path     string
	Data     map[string]string
	Override int // number of keys overridden by child
	Inherited int // number of keys taken from parent
}

// InheritSecret reads a parent secret and a child secret, returning a merged
// map where child values take precedence over parent values. Keys present only
// in the parent are inherited as-is. Keys present only in the child are kept.
// The version parameters accept the same ref syntax as ResolveVersionRef
// ("latest", or an explicit integer string).
func InheritSecret(l LogicalClient, parentPath, childPath, parentRef, childRef string) (InheritResult, error) {
	if strings.TrimSpace(parentPath) == "" {
		return InheritResult{}, errors.New("parent path must not be empty")
	}
	if strings.TrimSpace(childPath) == "" {
		return InheritResult{}, errors.New("child path must not be empty")
	}
	if parentPath == childPath {
		return InheritResult{}, errors.New("parent and child paths must differ")
	}

	parentData, err := fetchStringMap(l, parentPath, parentRef)
	if err != nil {
		return InheritResult{}, fmt.Errorf("reading parent %q: %w", parentPath, err)
	}

	childData, err := fetchStringMap(l, childPath, childRef)
	if err != nil {
		return InheritResult{}, fmt.Errorf("reading child %q: %w", childPath, err)
	}

	merged := make(map[string]string, len(parentData))
	inherited := 0
	overridden := 0

	for k, v := range parentData {
		if cv, ok := childData[k]; ok {
			merged[k] = cv
			overridden++
		} else {
			merged[k] = v
			inherited++
		}
	}
	for k, v := range childData {
		if _, exists := merged[k]; !exists {
			merged[k] = v
		}
	}

	return InheritResult{
		Path:      childPath,
		Data:      merged,
		Override:  overridden,
		Inherited: inherited,
	}, nil
}

// fetchStringMap resolves a version ref and returns the secret's string data.
func fetchStringMap(l LogicalClient, path, ref string) (map[string]string, error) {
	secret, err := FetchAtRef(l, path, ref)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return map[string]string{}, nil
	}
	return extractStringMap(secret), nil
}
