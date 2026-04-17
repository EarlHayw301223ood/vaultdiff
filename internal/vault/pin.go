package vault

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

func pinMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/pin", path)
}

// PinResult holds the result of a pin operation.
type PinResult struct {
	Path    string
	Version int
	Pinned  bool
}

// SetPin marks a specific version of a secret as pinned.
func SetPin(client *api.Client, path string, version int) (*PinResult, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if version <= 0 {
		return nil, errors.New("version must be a positive integer")
	}

	metaPath := pinMetaPath(path)
	data := map[string]interface{}{
		"version": version,
		"pinned":  true,
	}
	_, err := client.Logical().Write(metaPath, data)
	if err != nil {
		return nil, fmt.Errorf("set pin: %w", err)
	}
	return &PinResult{Path: path, Version: version, Pinned: true}, nil
}

// GetPin retrieves the pinned version for a secret path, if any.
func GetPin(client *api.Client, path string) (*PinResult, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	metaPath := pinMetaPath(path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("get pin: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return &PinResult{Path: path, Pinned: false}, nil
	}

	v, _ := secret.Data["version"].(float64)
	pinned, _ := secret.Data["pinned"].(bool)
	return &PinResult{Path: path, Version: int(v), Pinned: pinned}, nil
}

// ClearPin removes the pin from a secret path.
func ClearPin(client *api.Client, path string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	metaPath := pinMetaPath(path)
	_, err := client.Logical().Delete(metaPath)
	if err != nil {
		return fmt.Errorf("clear pin: %w", err)
	}
	return nil
}
