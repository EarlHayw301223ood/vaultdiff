package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/vaultdiff/internal/diff"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	changes := []diff.Change{
		{Key: "DB_PASS", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		{Key: "API_KEY", Type: diff.Added, NewValue: "abc"},
	}

	if err := logger.Record("staging", "secret/app", 2, 3, changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(line, "{") {
		t.Fatalf("expected JSON output, got: %s", line)
	}

	var entry Entry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if entry.Environment != "staging" {
		t.Errorf("expected environment 'staging', got %q", entry.Environment)
	}
	if entry.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", entry.Path)
	}
	if entry.FromVersion != 2 || entry.ToVersion != 3 {
		t.Errorf("unexpected versions: from=%d to=%d", entry.FromVersion, entry.ToVersion)
	}
	if entry.Summary.Modified != 1 || entry.Summary.Added != 1 || entry.Summary.Total != 2 {
		t.Errorf("unexpected summary: %+v", entry.Summary)
	}
}

func TestRecord_TimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	_ = logger.Record("prod", "secret/db", 1, 2, nil)

	var entry Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)

	if entry.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", entry.Timestamp.Location())
	}
}

func TestBuildSummary_Counts(t *testing.T) {
	changes := []diff.Change{
		{Type: diff.Added},
		{Type: diff.Added},
		{Type: diff.Removed},
		{Type: diff.Modified},
		{Type: diff.Modified},
		{Type: diff.Modified},
	}
	s := buildSummary(changes)
	if s.Added != 2 || s.Removed != 1 || s.Modified != 3 || s.Total != 6 {
		t.Errorf("unexpected summary counts: %+v", s)
	}
}

func TestNewLogger_DefaultsToStdout(t *testing.T) {
	l := NewLogger(nil)
	if l.out == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
