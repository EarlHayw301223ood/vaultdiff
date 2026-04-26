package vault

import (
	"testing"
)

func TestFreezeMetaPath_Format(t *testing.T) {
	got := freezeMetaPath("secret/myapp/db")
	want := "vaultdiff/freeze/secret/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFreezeMetaPath_TrimsSlashes(t *testing.T) {
	got := freezeMetaPath("/secret/myapp/")
	want := "vaultdiff/freeze/secret/myapp"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFreezeSecret_EmptyPath(t *testing.T) {
	client := newRollbackTestClient(nil, nil)
	_, err := FreezeSecret(client, "", 1, "safety", "alice")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestFreezeSecret_ZeroVersion(t *testing.T) {
	client := newRollbackTestClient(nil, nil)
	_, err := FreezeSecret(client, "secret/app", 0, "safety", "alice")
	if err == nil {
		t.Fatal("expected error for zero version")
	}
}

func TestFreezeSecret_EmptyReason(t *testing.T) {
	client := newRollbackTestClient(nil, nil)
	_, err := FreezeSecret(client, "secret/app", 1, "", "alice")
	if err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestFreezeRecord_Fields(t *testing.T) {
	rec := FreezeRecord{
		Frozen:   true,
		Version:  3,
		Reason:   "audit",
		FrozenBy: "bob",
	}
	if !rec.Frozen {
		t.Error("expected Frozen to be true")
	}
	if rec.Version != 3 {
		t.Errorf("expected Version 3, got %d", rec.Version)
	}
	if rec.Reason != "audit" {
		t.Errorf("expected Reason 'audit', got %q", rec.Reason)
	}
	if rec.FrozenBy != "bob" {
		t.Errorf("expected FrozenBy 'bob', got %q", rec.FrozenBy)
	}
}

func TestGetFreeze_EmptyPath(t *testing.T) {
	client := newRollbackTestClient(nil, nil)
	_, err := GetFreeze(client, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetFreeze_NilSecret_ReturnsFalse(t *testing.T) {
	client := newRollbackTestClient(nil, nil)
	rec, err := GetFreeze(client, "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Frozen {
		t.Error("expected Frozen to be false when no record exists")
	}
}

func TestUnfreezeSecret_EmptyPath(t *testing.T) {
	client := newRollbackTestClient(nil, nil)
	err := UnfreezeSecret(client, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}
