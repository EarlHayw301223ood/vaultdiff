package vault

import (
	"testing"
)

func TestRename_EmptySourcePath(t *testing.T) {
	c := newRollbackTestClient()
	_, err := Rename(c, "", "dest/secret")
	if err == nil || err.Error() != "rename: source path must not be empty" {
		t.Fatalf("expected empty source error, got %v", err)
	}
}

func TestRename_EmptyDestPath(t *testing.T) {
	c := newRollbackTestClient()
	_, err := Rename(c, "src/secret", "")
	if err == nil || err.Error() != "rename: destination path must not be empty" {
		t.Fatalf("expected empty dest error, got %v", err)
	}
}

func TestRename_SameSourceAndDest(t *testing.T) {
	c := newRollbackTestClient()
	_, err := Rename(c, "src/secret", "src/secret")
	if err == nil {
		t.Fatal("expected error for same source and dest")
	}
}

func TestRenameResult_Fields(t *testing.T) {
	r := &RenameResult{
		SourcePath: "old/key",
		DestPath:   "new/key",
		Version:    3,
	}
	if r.SourcePath != "old/key" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestPath != "new/key" {
		t.Errorf("unexpected DestPath: %s", r.DestPath)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
}
