package vault_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_MOUNT_PATH", "")

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got %q", cfg.Address)
	}
	if cfg.MountPath != "secret" {
		t.Errorf("expected default mount path 'secret', got %q", cfg.MountPath)
	}
	if cfg.Token != "test-token" {
		t.Errorf("expected token 'test-token', got %q", cfg.Token)
	}
}

func TestConfigFromEnv_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")

	_, err := vault.ConfigFromEnv()
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing")
	}
}

func TestConfigFromEnv_CustomValues(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://vault.prod.example.com:8200")
	t.Setenv("VAULT_TOKEN", "s.abc123")
	t.Setenv("VAULT_MOUNT_PATH", "kv")

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Address != "https://vault.prod.example.com:8200" {
		t.Errorf("unexpected address: %q", cfg.Address)
	}
	if cfg.MountPath != "kv" {
		t.Errorf("unexpected mount path: %q", cfg.MountPath)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     vault.Config
		wantErr bool
	}{
		{"valid", vault.Config{Address: "http://localhost:8200", Token: "tok", MountPath: "secret"}, false},
		{"missing address", vault.Config{Token: "tok", MountPath: "secret"}, true},
		{"missing token", vault.Config{Address: "http://localhost:8200", MountPath: "secret"}, true},
		{"missing mount", vault.Config{Address: "http://localhost:8200", Token: "tok"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
