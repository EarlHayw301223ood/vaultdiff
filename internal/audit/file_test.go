package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/diff"
)

func TestNewFileLogger_CreatesDirectory(t *testing.T) {
	dir := t.TempDir() + "/audit-logs"

	fl, err := NewFileLogger(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer fl.Close()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected audit log directory to be created")
	}
}

func TestFileLogger_WritesReadableEntries(t *testing.T) {
	dir := t.TempDir()

	fl, err := NewFileLogger(dir)
	if err != nil {
		t.Fatalf("NewFileLogger: %v", err)
	}
	defer fl.Close()

	changes := []diff.Change{
		{Key: "SECRET", Type: diff.Removed, OldValue: "x"},
	}

	if err := fl.Record("prod", "secret/svc", 5, 6, changes); err != nil {
		t.Fatalf("Record: %v", err)
	}

	path := fl.LogPath()
	if path == "" {
		t.Fatal("expected non-empty log path")
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one log line")
	}

	var entry Entry
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if entry.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", entry.Summary.Removed)
	}
}

func TestFileLogger_LogPathContainsDate(t *testing.T) {
	dir := t.TempDir()

	fl, err := NewFileLogger(dir)
	if err != nil {
		t.Fatalf("NewFileLogger: %v", err)
	}
	defer fl.Close()

	path := fl.LogPath()
	if !strings.Contains(path, "vaultdiff-audit-") {
		t.Errorf("expected log path to contain date prefix, got: %s", path)
	}
	if !strings.HasSuffix(path, ".jsonl") {
		t.Errorf("expected .jsonl extension, got: %s", path)
	}
}
