package vault

import (
	"testing"
)

func TestRestore_EmptySourcePath(t *testing.T) {
	c := &Client{}
	_, err := Restore(c, "", "dest/secret", 1)
	if err == nil || err.Error() != "source path must not be empty" {
		t.Fatalf("expected source path error, got %v", err)
	}
}

func TestRestore_EmptyDestPath(t *testing.T) {
	c := &Client{}
	_, err := Restore(c, "src/secret", "", 1)
	if err == nil || err.Error() != "destination path must not be empty" {
		t.Fatalf("expected dest path error, got %v", err)
	}
}

func TestRestore_InvalidVersion(t *testing.T) {
	c := &Client{}
	_, err := Restore(c, "src/secret", "dest/secret", 0)
	if err == nil {
		t.Fatal("expected version error, got nil")
	}
}

func TestRestoreResult_Fields(t *testing.T) {
	r := &RestoreResult{
		SourcePath: "src/secret",
		DestPath:   "dest/secret",
		Version:    3,
		RestoredTo: "dest/secret",
	}
	if r.SourcePath != "src/secret" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if r.RestoredTo != r.DestPath {
		t.Errorf("RestoredTo should match DestPath")
	}
}
