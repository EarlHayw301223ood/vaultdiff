package vault

import (
	"context"
	"fmt"
	"strings"
)

// TreeNode represents a single path entry in a Vault KV mount.
type TreeNode struct {
	Path     string
	IsPrefix bool // true if this is a directory-like prefix (ends with /)
}

// ListTree recursively lists all secret paths under the given prefix in a KV v2 mount.
// It returns a flat slice of leaf paths (non-prefix entries).
func ListTree(ctx context.Context, client *Client, mount, prefix string) ([]string, error) {
	if mount == "" {
		return nil, fmt.Errorf("mount must not be empty")
	}

	var results []string
	err := walkTree(ctx, client, mount, prefix, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func walkTree(ctx context.Context, client *Client, mount, prefix string, out *[]string) error {
	listPath := fmt.Sprintf("%s/metadata/%s", mount, prefix)

	secret, err := client.Logical().ListWithContext(ctx, listPath)
	if err != nil {
		return fmt.Errorf("listing %q: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil
	}

	keys, ok := secret.Data["keys"]
	if !ok {
		return nil
	}

	rawKeys, ok := keys.([]interface{})
	if !ok {
		return fmt.Errorf("unexpected keys type at %q", listPath)
	}

	for _, k := range rawKeys {
		key, ok := k.(string)
		if !ok {
			continue
		}
		full := strings.TrimSuffix(prefix, "/") + "/" + key
		if strings.HasSuffix(key, "/") {
			if err := walkTree(ctx, client, mount, full, out); err != nil {
				return err
			}
		} else {
			*out = append(*out, full)
		}
	}
	return nil
}
