package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// NamespaceInfo holds metadata about a Vault namespace.
type NamespaceInfo struct {
	Path     string
	FullPath string
	Children []string
}

// ListNamespaces returns the child namespaces under the given namespace path.
// Pass an empty string to list from the root namespace.
func ListNamespaces(client *api.Client, namespacePath string) ([]NamespaceInfo, error) {
	if client == nil {
		return nil, fmt.Errorf("client must not be nil")
	}

	listPath := "sys/namespaces"
	if namespacePath != "" {
		ns := strings.Trim(namespacePath, "/")
		listPath = fmt.Sprintf("%s/sys/namespaces", ns)
	}

	secret, err := client.Logical().List(listPath)
	if err != nil {
		return nil, fmt.Errorf("listing namespaces at %q: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return []NamespaceInfo{}, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return []NamespaceInfo{}, nil
	}

	result := make([]NamespaceInfo, 0, len(keys))
	for _, k := range keys {
		name, ok := k.(string)
		if !ok {
			continue
		}
		name = strings.TrimSuffix(name, "/")
		full := name
		if namespacePath != "" {
			full = strings.Trim(namespacePath, "/") + "/" + name
		}
		result = append(result, NamespaceInfo{
			Path:     name,
			FullPath: full,
		})
	}
	return result, nil
}

// NamespacePath normalises a namespace path by trimming slashes and
// lower-casing, consistent with how Vault stores namespace identifiers.
func NamespacePath(raw string) string {
	return strings.ToLower(strings.Trim(raw, "/"))
}
