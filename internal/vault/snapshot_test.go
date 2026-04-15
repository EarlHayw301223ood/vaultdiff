package vault

import (
	"testing"
)

func makeSnapshot(mount string, secrets map[string]map[string]string) *SecretSnapshot {
	return &SecretSnapshot{
		Mount:   mount,
		Secrets: secrets,
	}
}

func TestSnapshot_Paths_Sorted(t *testing.T) {
	s := makeSnapshot("secret", map[string]map[string]string{
		"secret/z": {"k": "v"},
		"secret/a": {"k": "v"},
		"secret/m": {"k": "v"},
	})

	paths := s.Paths()
	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(paths))
	}
	if paths[0] != "secret/a" || paths[1] != "secret/m" || paths[2] != "secret/z" {
		t.Errorf("paths not sorted: %v", paths)
	}
}

func TestSnapshot_Paths_Empty(t *testing.T) {
	s := makeSnapshot("secret", map[string]map[string]string{})
	paths := s.Paths()
	if len(paths) != 0 {
		t.Errorf("expected empty paths, got %v", paths)
	}
}

func TestSnapshot_Get_ExistingPath(t *testing.T) {
	s := makeSnapshot("secret", map[string]map[string]string{
		"secret/foo": {"user": "admin", "pass": "secret"},
	})

	kv := s.Get("secret/foo")
	if kv == nil {
		t.Fatal("expected non-nil map for existing path")
	}
	if kv["user"] != "admin" {
		t.Errorf("expected user=admin, got %s", kv["user"])
	}
}

func TestSnapshot_Get_MissingPath(t *testing.T) {
	s := makeSnapshot("secret", map[string]map[string]string{})
	kv := s.Get("secret/nonexistent")
	if kv != nil {
		t.Errorf("expected nil for missing path, got %v", kv)
	}
}

func TestSnapshot_MountStored(t *testing.T) {
	s := makeSnapshot("kv", map[string]map[string]string{})
	if s.Mount != "kv" {
		t.Errorf("expected mount=kv, got %s", s.Mount)
	}
}
