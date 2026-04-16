package vault

import (
	"testing"
)

func TestCopy_EmptySourcePath(t *testing.T) {
	c := &Client{}
	_, err := Copy(c, "", "dest/secret", "latest")
	if err == nil || err.Error() != "copy: source path must not be empty" {
		t.Fatalf("expected empty source error, got %v", err)
	}
}

func TestCopy_EmptyDestPath(t *testing.T) {
	c := &Client{}
	_, err := Copy(c, "src/secret", "", "latest")
	if err == nil || err.Error() != "copy: destination path must not be empty" {
		t.Fatalf("expected empty dest error, got %v", err)
	}
}

func TestCopy_SameSourceAndDest(t *testing.T) {
	c := &Client{}
	_, err := Copy(c, "secret/foo", "secret/foo", "latest")
	if err == nil {
		t.Fatal("expected error when source == dest")
	}
}

func TestCopyResult_Fields(t *testing.T) {
	r := &CopyResult{
		SourcePath: "src/a",
		DestPath:   "dst/b",
		Version:    3,
		Keys:       5,
	}
	if r.SourcePath != "src/a" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if r.Keys != 5 {
		t.Errorf("unexpected Keys: %d", r.Keys)
	}
}
