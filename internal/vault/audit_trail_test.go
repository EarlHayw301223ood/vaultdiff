package vault

import (
	"testing"
)

func TestAuditTrailMetaPath_Format(t *testing.T) {
	got := auditTrailMetaPath("secret/myapp/db")
	want := "secret/metadata/_audit_trail/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAuditTrailMetaPath_TrimsSlashes(t *testing.T) {
	got := auditTrailMetaPath("/secret/myapp/db/")
	want := "secret/metadata/_audit_trail/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAuditTrailMetaPath_SingleSegment(t *testing.T) {
	got := auditTrailMetaPath("secret")
	if got == "" {
		t.Error("expected non-empty path")
	}
}

func TestAppendAuditTrail_EmptyPath(t *testing.T) {
	err := AppendAuditTrail(nil, "", "write", "alice", 1, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAppendAuditTrail_EmptyOperation(t *testing.T) {
	err := AppendAuditTrail(nil, "secret/myapp", "", "alice", 1, "")
	if err == nil {
		t.Fatal("expected error for empty operation")
	}
}

func TestAppendAuditTrail_InvalidVersion(t *testing.T) {
	err := AppendAuditTrail(nil, "secret/myapp", "write", "alice", 0, "")
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestGetAuditTrail_EmptyPath(t *testing.T) {
	_, err := GetAuditTrail(nil, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAuditTrailEntry_Fields(t *testing.T) {
	entry := AuditTrailEntry{
		Operation: "write",
		Actor:     "alice",
		Version:   3,
		Note:      "deploy",
	}
	if entry.Operation != "write" {
		t.Errorf("unexpected operation: %s", entry.Operation)
	}
	if entry.Actor != "alice" {
		t.Errorf("unexpected actor: %s", entry.Actor)
	}
	if entry.Version != 3 {
		t.Errorf("unexpected version: %d", entry.Version)
	}
	if entry.Note != "deploy" {
		t.Errorf("unexpected note: %s", entry.Note)
	}
}

func TestAuditTrail_EmptyEntries(t *testing.T) {
	trail := &AuditTrail{Path: "secret/myapp", Entries: []AuditTrailEntry{}}
	if len(trail.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(trail.Entries))
	}
}
