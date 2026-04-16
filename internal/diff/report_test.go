package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewReport_SummaryCounts(t *testing.T) {
	changes := []Change{
		{Key: "A", Type: Added, NewValue: "1"},
		{Key: "B", Type: Removed, OldValue: "2"},
		{Key: "C", Type: Modified, OldValue: "3", NewValue: "4"},
		{Key: "D", Type: Added, NewValue: "5"},
	}
	r := NewReport("dev", "prod", "secret/app", "secret/app", changes)

	if r.Summary.Added != 2 {
		t.Errorf("expected 2 added, got %d", r.Summary.Added)
	}
	if r.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", r.Summary.Removed)
	}
	if r.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", r.Summary.Modified)
	}
	if r.Summary.Total != 4 {
		t.Errorf("expected total 4, got %d", r.Summary.Total)
	}
}

func TestNewReport_MetadataStored(t *testing.T) {
	r := NewReport("staging", "prod", "secret/db", "secret/db", nil)
	if r.SourceEnv != "staging" || r.TargetEnv != "prod" {
		t.Errorf("env metadata not stored correctly")
	}
	if r.SourcePath != "secret/db" {
		t.Errorf("source path not stored")
	}
}

func TestReport_WriteJSON(t *testing.T) {
	changes := []Change{
		{Key: "TOKEN", Type: Added, NewValue: "abc"},
	}
	r := NewReport("dev", "prod", "secret/app", "secret/app", changes)

	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded.Changes) != 1 {
		t.Errorf("expected 1 change in JSON, got %d", len(decoded.Changes))
	}
}

func TestReport_WriteText_ContainsEnvAndPath(t *testing.T) {
	r := NewReport("dev", "prod", "secret/app", "secret/app", []Change{})

	var buf bytes.Buffer
	if err := r.WriteText(&buf, false); err != nil {
		t.Fatalf("WriteText failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dev") || !strings.Contains(out, "prod") {
		t.Errorf("expected env names in text output")
	}
	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected path in text output")
	}
}

func TestReport_WriteText_SummaryLine(t *testing.T) {
	changes := []Change{
		{Key: "X", Type: Removed, OldValue: "v"},
	}
	r := NewReport("a", "b", "p", "p", changes)

	var buf bytes.Buffer
	_ = r.WriteText(&buf, false)
	if !strings.Contains(buf.String(), "removed: 1") {
		t.Errorf("expected removed count in summary line")
	}
}

func TestNewReport_EmptyChanges(t *testing.T) {
	r := NewReport("dev", "prod", "secret/app", "secret/app", nil)

	if r.Summary.Total != 0 {
		t.Errorf("expected total 0 for nil changes, got %d", r.Summary.Total)
	}
	if r.Changes == nil {
		t.Errorf("expected Changes to be non-nil slice, got nil")
	}
}
