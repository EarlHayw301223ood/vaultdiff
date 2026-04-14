package vault

import (
	"errors"
	"os"
)

// Config holds the configuration needed to connect to a Vault instance.
type Config struct {
	// Address is the Vault server URL, e.g. https://vault.example.com:8200
	Address string
	// Token is the Vault authentication token.
	Token string
	// MountPath is the KV v2 secrets engine mount path (default: "secret").
	MountPath string
}

// ConfigFromEnv builds a Config by reading standard Vault environment variables.
// VAULT_ADDR, VAULT_TOKEN, and optionally VAULT_MOUNT_PATH.
func ConfigFromEnv() (Config, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		addr = "http://127.0.0.1:8200"
	}

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return Config{}, errors.New("VAULT_TOKEN environment variable is required")
	}

	mount := os.Getenv("VAULT_MOUNT_PATH")
	if mount == "" {
		mount = "secret"
	}

	return Config{
		Address:   addr,
		Token:     token,
		MountPath: mount,
	}, nil
}

// Validate checks that the Config has the required fields set.
func (c Config) Validate() error {
	if c.Address == "" {
		return errors.New("vault address must not be empty")
	}
	if c.Token == "" {
		return errors.New("vault token must not be empty")
	}
	if c.MountPath == "" {
		return errors.New("vault mount path must not be empty")
	}
	return nil
}
