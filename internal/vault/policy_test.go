package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestHasCapability_Match(t *testing.T) {
	access := PolicyAccess{
		Path:         "secret/data/foo",
		Capabilities: []string{"read", "list"},
	}
	if !HasCapability(access, "read") {
		t.Error("expected HasCapability to return true for 'read'")
	}
}

func TestHasCapability_NoMatch(t *testing.T) {
	access := PolicyAccess{
		Path:         "secret/data/foo",
		Capabilities: []string{"list"},
	}
	if HasCapability(access, "write") {
		t.Error("expected HasCapability to return false for 'write'")
	}
}

func TestHasCapability_CaseInsensitive(t *testing.T) {
	access := PolicyAccess{
		Path:         "secret/data/foo",
		Capabilities: []string{"READ"},
	}
	if !HasCapability(access, "read") {
		t.Error("expected HasCapability to be case-insensitive")
	}
}

func TestCheckPaths_EmptyInput(t *testing.T) {
	checker := NewPolicyChecker(&mockLogical{})
	results, err := checker.CheckPaths(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty input, got %v", results)
	}
}

func TestCheckPaths_ReturnsSorted(t *testing.T) {
	logical := &mockLogical{
		writeFunc: func(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
			paths := data["paths"].([]string)
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					paths[0]: []interface{}{"read"},
				},
			}, nil
		},
	}

	checker := NewPolicyChecker(logical)
	paths := []string{"secret/data/z", "secret/data/a", "secret/data/m"}
	results, err := checker.CheckPaths(context.Background(), paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Path != "secret/data/a" {
		t.Errorf("expected first result to be 'secret/data/a', got %q", results[0].Path)
	}
	if results[2].Path != "secret/data/z" {
		t.Errorf("expected last result to be 'secret/data/z', got %q", results[2].Path)
	}
}

func TestCheckPaths_NilSecretReturnsDeny(t *testing.T) {
	logical := &mockLogical{
		writeFunc: func(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
			return nil, nil
		},
	}

	checker := NewPolicyChecker(logical)
	results, err := checker.CheckPaths(context.Background(), []string{"secret/data/foo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if !HasCapability(results[0], "deny") {
		t.Error("expected 'deny' capability for nil secret response")
	}
}
