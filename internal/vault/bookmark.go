package vault

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

func bookmarkMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/bookmarks", path)
}

// SetBookmark stores a named bookmark pointing to a specific version of a secret path.
func SetBookmark(client *api.Client, path, name string, version int) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if name == "" {
		return errors.New("bookmark name must not be empty")
	}
	if version < 1 {
		return errors.New("version must be >= 1")
	}

	metaPath := bookmarkMetaPath(path)
	existing, _ := client.Logical().Read(metaPath)

	data := map[string]interface{}{}
	if existing != nil && existing.Data != nil {
		for k, v := range existing.Data {
			data[k] = v
		}
	}
	data[name] = version

	_, err := client.Logical().Write(metaPath, data)
	return err
}

// GetBookmark retrieves the version number associated with a named bookmark.
func GetBookmark(client *api.Client, path, name string) (int, error) {
	if path == "" {
		return 0, errors.New("path must not be empty")
	}
	if name == "" {
		return 0, errors.New("bookmark name must not be empty")
	}

	metaPath := bookmarkMetaPath(path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return 0, err
	}
	if secret == nil || secret.Data == nil {
		return 0, fmt.Errorf("no bookmarks found for path %q", path)
	}

	val, ok := secret.Data[name]
	if !ok {
		return 0, fmt.Errorf("bookmark %q not found for path %q", name, path)
	}

	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("unexpected type for bookmark version: %T", val)
	}
}

// DeleteBookmark removes a named bookmark from a secret path.
func DeleteBookmark(client *api.Client, path, name string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if name == "" {
		return errors.New("bookmark name must not be empty")
	}

	metaPath := bookmarkMetaPath(path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return err
	}
	if secret == nil || secret.Data == nil {
		return fmt.Errorf("no bookmarks found for path %q", path)
	}

	if _, ok := secret.Data[name]; !ok {
		return fmt.Errorf("bookmark %q not found for path %q", name, path)
	}

	delete(secret.Data, name)
	_, err = client.Logical().Write(metaPath, secret.Data)
	return err
}
