package vault

import (
	"errors"
	"fmt"
	"strings"
)

const aliasMetaMount = "secret/meta/aliases"

func aliasMetaPath(alias string) string {
	return fmt.Sprintf("%s/%s", aliasMetaMount, strings.Trim(alias, "/"))
}

// AliasEntry maps a short alias name to a vault path and optional version.
type AliasEntry struct {
	Alias   string `json:"alias"`
	Path    string `json:"path"`
	Version int    `json:"version,omitempty"`
}

// SetAlias stores an alias entry pointing to the given vault path.
func SetAlias(client LogicalClient, alias, path string, version int) error {
	if alias == "" {
		return errors.New("alias name must not be empty")
	}
	if path == "" {
		return errors.New("target path must not be empty")
	}
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"alias":   alias,
			"path":    path,
			"version": version,
		},
	}
	_, err := client.Write(kvV2WritePath(aliasMetaPath(alias)), data)
	return err
}

// GetAlias retrieves an alias entry by name.
func GetAlias(client LogicalClient, alias string) (*AliasEntry, error) {
	if alias == "" {
		return nil, errors.New("alias name must not be empty")
	}
	secret, err := client.Read(kvV2ReadPath(aliasMetaPath(alias)))
	if err != nil {
		return nil, fmt.Errorf("read alias: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("alias %q not found", alias)
	}
	data := extractStringMap(secret)
	entry := &AliasEntry{
		Alias: alias,
		Path:  data["path"],
	}
	if v, ok := secret.Data["data"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			if ver, ok := m["version"]; ok {
				if f, ok := ver.(float64); ok {
					entry.Version = int(f)
				}
			}
		}
	}
	return entry, nil
}
