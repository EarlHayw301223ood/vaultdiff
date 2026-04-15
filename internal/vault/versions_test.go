package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

// mockLogical is a minimal stub satisfying the logical interface used by Client.
type mockLogical struct {
	readFn func(path string) (*vaultapi.Secret, error)
}

func (m *mockLogical) ReadWithContext(_ context.Context, path string) (*vaultapi.Secret, error) {
	return m.readFn(path)
}

func makeVersionsSecret(versions map[string]interface{}) *vaultapi.Secret {
	return &vaultapi.Secret{
		Data: map[string]interface{}{
			"versions": versions,
		},
	}
}

func TestListVersions_ReturnsVersionsSorted(t *testing.T) {
	raw := map[string]interface{}{
		"3": map[string]interface{}{"created_time": "2024-03-01T00:00:00Z", "deletion_time": "", "destroyed": false},
		"1": map[string]interface{}{"created_time": "2024-01-01T00:00:00Z", "deletion_time": "", "destroyed": false},
		"2": map[string]interface{}{"created_time": "2024-02-01T00:00:00Z", "deletion_time": "", "destroyed": false},
	}

	c := &Client{logical: &mockLogical{readFn: func(_ string) (*vaultapi.Secret, error) {
		return makeVersionsSecret(raw), nil
	}}}

	versions, err := c.ListVersions(context.Background(), "secret", "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(versions))
	}
	for i, want := range []int{1, 2, 3} {
		if versions[i].Version != want {
			t.Errorf("versions[%d].Version = %d, want %d", i, versions[i].Version, want)
		}
	}
}

func TestListVersions_NilSecret(t *testing.T) {
	c := &Client{logical: &mockLogical{readFn: func(_ string) (*vaultapi.Secret, error) {
		return nil, nil
	}}}

	_, err := c.ListVersions(context.Background(), "secret", "missing/path")
	if err == nil {
		t.Fatal("expected error for nil secret, got nil")
	}
}

func TestLatestVersion_ReturnsHighest(t *testing.T) {
	raw := map[string]interface{}{
		"1": map[string]interface{}{"created_time": "2024-01-01T00:00:00Z", "deletion_time": "", "destroyed": false},
		"5": map[string]interface{}{"created_time": "2024-05-01T00:00:00Z", "deletion_time": "", "destroyed": false},
		"3": map[string]interface{}{"created_time": "2024-03-01T00:00:00Z", "deletion_time": "", "destroyed": false},
	}

	c := &Client{logical: &mockLogical{readFn: func(_ string) (*vaultapi.Secret, error) {
		return makeVersionsSecret(raw), nil
	}}}

	latest, err := c.LatestVersion(context.Background(), "secret", "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != 5 {
		t.Errorf("expected latest version 5, got %d", latest)
	}
}
