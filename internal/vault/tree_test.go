package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestListTree_EmptyMount(t *testing.T) {
	client := &Client{logical: &mockLogical{}}
	_, err := ListTree(context.Background(), client, "", "")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestListTree_NilSecret(t *testing.T) {
	m := &mockLogical{
		listFn: func(path string) (*vaultapi.Secret, error) {
			return nil, nil
		},
	}
	client := &Client{logical: m}
	paths, err := ListTree(context.Background(), client, "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("expected empty result, got %v", paths)
	}
}

func TestListTree_FlatKeys(t *testing.T) {
	m := &mockLogical{
		listFn: func(path string) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"keys": []interface{}{"alpha", "beta", "gamma"},
				},
			}, nil
		},
	}
	client := &Client{logical: m}
	paths, err := ListTree(context.Background(), client, "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d: %v", len(paths), paths)
	}
}

func TestListTree_NestedPrefix(t *testing.T) {
	callCount := 0
	m := &mockLogical{
		listFn: func(path string) (*vaultapi.Secret, error) {
			callCount++
			if callCount == 1 {
				return &vaultapi.Secret{
					Data: map[string]interface{}{
						"keys": []interface{}{"subdir/", "top-level"},
					},
				}, nil
			}
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"keys": []interface{}{"nested-key"},
				},
			}, nil
		},
	}
	client := &Client{logical: m}
	paths, err := ListTree(context.Background(), client, "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 leaf paths, got %d: %v", len(paths), paths)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 list calls for recursive walk, got %d", callCount)
	}
}
