package vault

import (
	"errors"
	"fmt"
	"strings"
)

const protectMetaPrefix = "vaultdiff/protect"

func protectMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%s/%s", protectMetaPrefix, path)
}

// SetProtection marks a secret path as protected with an optional reason.
func SetProtection(client LogicalClient, path, reason string) error {
	if path == "" {
		return errors.New("path is required")
	}
	metaPath := protectMetaPath(path)
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"protected": "true",
			"reason":    reason,
		},
	}
	_, err := client.Write(metaPath, data)
	return err
}

// GetProtection retrieves the protection status and reason for a secret path.
func GetProtection(client LogicalClient, path string) (protected bool, reason string, err error) {
	if path == "" {
		return false, "", errors.New("path is required")
	}
	metaPath := protectMetaPath(path)
	secret, err := client.Read(metaPath)
	if err != nil {
		return false, "", err
	}
	if secret == nil {
		return false, "", nil
	}
	data := extractStringMap(secret)
	protected = data["protected"] == "true"
	reason = data["reason"]
	return protected, reason, nil
}

// ClearProtection removes protection from a secret path.
func ClearProtection(client LogicalClient, path string) error {
	if path == "" {
		return errors.New("path is required")
	}
	metaPath := protectMetaPath(path)
	_, err := client.Delete(metaPath)
	return err
}
