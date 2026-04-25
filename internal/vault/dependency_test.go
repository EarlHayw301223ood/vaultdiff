package vault

import (
	"testing"
)

func TestDependencyMetaPath_Format(t *testing.T) {
	got := dependencyMetaPath("secret/app/db")
	want := "vaultdiff/meta/secret/app/db/dependencies"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDependencyMetaPath_TrimsSlashes(t *testing.T) {
	got := dependencyMetaPath("/secret/app/")
	want := "vaultdiff/meta/secret/app/dependencies"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAddDependency_EmptySourcePath(t *testing.T) {
	client := newStubClient()
	err := AddDependency(client, "", "secret/b", "")
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestAddDependency_EmptyTargetPath(t *testing.T) {
	client := newStubClient()
	err := AddDependency(client, "secret/a", "", "")
	if err == nil {
		t.Fatal("expected error for empty target path")
	}
}

func TestAddDependency_SameSourceAndTarget(t *testing.T) {
	client := newStubClient()
	err := AddDependency(client, "secret/a", "secret/a", "")
	if err == nil {
		t.Fatal("expected error when source equals target")
	}
}

func TestAddDependency_DuplicateTarget(t *testing.T) {
	client := newStubClient()
	if err := AddDependency(client, "secret/a", "secret/b", "db"); err != nil {
		t.Fatalf("first add failed: %v", err)
	}
	err := AddDependency(client, "secret/a", "secret/b", "db")
	if err == nil {
		t.Fatal("expected error for duplicate dependency")
	}
}

func TestGetDependencies_EmptyPath(t *testing.T) {
	client := newStubClient()
	_, err := GetDependencies(client, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetDependencies_NoneRegistered(t *testing.T) {
	client := newStubClient()
	list, err := GetDependencies(client, "secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list.Dependencies) != 0 {
		t.Errorf("expected empty list, got %d entries", len(list.Dependencies))
	}
}

func TestDependencyResult_Fields(t *testing.T) {
	client := newStubClient()
	if err := AddDependency(client, "secret/a", "secret/b", "cache"); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	list, err := GetDependencies(client, "secret/a")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if len(list.Dependencies) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(list.Dependencies))
	}
	d := list.Dependencies[0]
	if d.TargetPath != "secret/b" {
		t.Errorf("TargetPath: got %q, want %q", d.TargetPath, "secret/b")
	}
	if d.Label != "cache" {
		t.Errorf("Label: got %q, want %q", d.Label, "cache")
	}
	if d.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}
