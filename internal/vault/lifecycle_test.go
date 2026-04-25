package vault

import (
	"testing"
)

func TestLifecycleMetaPath_Format(t *testing.T) {
	got := lifecycleMetaPath("secret/myapp/db")
	want := "vaultdiff/meta/secret/myapp/db/lifecycle"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLifecycleMetaPath_TrimsSlashes(t *testing.T) {
	got := lifecycleMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/lifecycle"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetLifecycle_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := SetLifecycle(c, "", 1, StageActive, "user", "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetLifecycle_InvalidVersion(t *testing.T) {
	c := newStubClient()
	_, err := SetLifecycle(c, "secret/app", 0, StageActive, "user", "")
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestSetLifecycle_UnknownStage(t *testing.T) {
	c := newStubClient()
	_, err := SetLifecycle(c, "secret/app", 1, LifecycleStage("unknown"), "user", "")
	if err == nil {
		t.Fatal("expected error for unknown stage")
	}
}

func TestSetLifecycle_WritesRecord(t *testing.T) {
	c := newStubClient()
	rec, err := SetLifecycle(c, "secret/app", 2, StageDeprecated, "alice", "rotating")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Stage != StageDeprecated {
		t.Errorf("got stage %q, want %q", rec.Stage, StageDeprecated)
	}
	if rec.Version != 2 {
		t.Errorf("got version %d, want 2", rec.Version)
	}
	if rec.ChangedBy != "alice" {
		t.Errorf("got changed_by %q, want alice", rec.ChangedBy)
	}
	if rec.Reason != "rotating" {
		t.Errorf("got reason %q, want rotating", rec.Reason)
	}
	if rec.ChangedAt.IsZero() {
		t.Error("expected non-zero ChangedAt")
	}
}

func TestGetLifecycle_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := GetLifecycle(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetLifecycle_MissingRecord(t *testing.T) {
	c := newStubClient()
	rec, err := GetLifecycle(c, "secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Error("expected nil record for missing path")
	}
}

func TestLifecycleStages_Constants(t *testing.T) {
	stages := []LifecycleStage{StageActive, StageDeprecated, StageRetired, StageReview}
	if len(stages) != 4 {
		t.Errorf("expected 4 lifecycle stages, got %d", len(stages))
	}
}
