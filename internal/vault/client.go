package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with convenience methods.
type Client struct {
	vc      *vaultapi.Client
	MountPath string
}

// NewClient creates a new Vault client from the provided config.
func NewClient(cfg Config) (*Client, error) {
	vaultCfg := vaultapi.DefaultConfig()
	vaultCfg.Address = cfg.Address

	vc, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	vc.SetToken(cfg.Token)

	return &Client{
		vc:        vc,
		MountPath: cfg.MountPath,
	}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// Pass version 0 to read the latest version.
func (c *Client) ReadSecretVersion(path string, version int) (map[string]interface{}, error) {
	var versionParam map[string][]string
	if version > 0 {
		versionParam = map[string][]string{
			"version": {fmt.Sprintf("%d", version)},
		}
	}

	secret, err := c.vc.KVv2(c.MountPath).GetVersion(nil, path, version)
	if err != nil {
		_ = versionParam
		return nil, fmt.Errorf("reading secret %q version %d: %w", path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q version %d not found", path, version)
	}

	return secret.Data, nil
}

// ListSecretVersions returns metadata for all versions of a KV v2 secret.
func (c *Client) ListSecretVersions(path string) (*vaultapi.KVMetadata, error) {
	meta, err := c.vc.KVv2(c.MountPath).GetMetadata(nil, path)
	if err != nil {
		return nil, fmt.Errorf("listing versions for %q: %w", path, err)
	}
	return meta, nil
}
