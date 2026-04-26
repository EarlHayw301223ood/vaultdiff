package vault

import (
	"testing"
)

func TestImmutableMetaPath_Format(t *testing.T) {
	got := immutableMetaPath("secret/myapp/db")
	want := "vaultdiff/meta/secret/myapp/db/immutable"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestImmutableMetaPath_TrimsSlashes(t *testing.T) {
	got := immutableMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/immutable"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetImmutable_EmptyPath(t *testing.T) {
	c := newStubClient()
	err := SetImmutable(c, "", "alice", "freeze for audit")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetImmutable_EmptySetBy(t *testing.T) {
	c := newStubClient()
	err := SetImmutable(c, "secret/app", "", "reason")
	if err == nil {
		t.Fatal("expected error for empty setBy")
	}
}

func TestGetImmutable_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := GetImmutable(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetImmutable_NotSet_ReturnsDisabled(t *testing.T) {
	c := newStubClient()
	rec, err := GetImmutable(c, "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Enabled {
		t.Error("expected Enabled=false when no record exists")
	}
}

func TestClearImmutable_EmptyPath(t *testing.T) {
	c := newStubClient()
	err := ClearImmutable(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestImmutableRecord_Fields(t *testing.T) {
	rec := ImmutableRecord{
		Enabled: true,
		SetBy:   "alice",
		Reason:  "compliance freeze",
	}
	if !rec.Enabled {
		t.Error("expected Enabled=true")
	}
	if rec.SetBy != "alice" {
		t.Errorf("got SetBy=%q, want %q", rec.SetBy, "alice")
	}
	if rec.Reason != "compliance freeze" {
		t.Errorf("got Reason=%q, want %q", rec.Reason, "compliance freeze")
	}
}
