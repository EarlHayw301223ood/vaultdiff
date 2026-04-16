package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestAliasMetaPath_Format(t *testing.T) {
	got := aliasMetaPath("prod-db")
	want := "secret/meta/aliases/prod-db"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestAliasMetaPath_TrimsSlashes(t *testing.T) {
	got := aliasMetaPath("/prod-db/")
	want := "secret/meta/aliases/prod-db"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestSetAlias_EmptyAlias(t *testing.T) {
	c := &mockLogical{}
	if err := SetAlias(c, "", "secret/data/foo", 0); err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestSetAlias_EmptyPath(t *testing.T) {
	c := &mockLogical{}
	if err := SetAlias(c, "myalias", "", 0); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetAlias_Writes(t *testing.T) {
	c := &mockLogical{}
	if err := SetAlias(c, "dev-db", "secret/data/dev/db", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.writes) != 1 {
		t.Fatalf("expected 1 write, got %d", len(c.writes))
	}
}

func TestGetAlias_EmptyAlias(t *testing.T) {
	c := &mockLogical{}
	if _, err := GetAlias(c, ""); err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestGetAlias_NotFound(t *testing.T) {
	c := &mockLogical{readSecret: nil}
	_, err := GetAlias(c, "missing")
	if err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestGetAlias_ReturnsEntry(t *testing.T) {
	c := &mockLogical{
		readSecret: &api.Secret{
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"alias":   "prod-db",
					"path":    "secret/data/prod/db",
					"version": float64(3),
				},
			},
		},
	}
	entry, err := GetAlias(c, "prod-db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Path != "secret/data/prod/db" {
		t.Errorf("unexpected path: %s", entry.Path)
	}
	if entry.Version != 3 {
		t.Errorf("unexpected version: %d", entry.Version)
	}
}
