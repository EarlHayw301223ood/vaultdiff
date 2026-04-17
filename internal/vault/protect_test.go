package vault

import (
	"testing"
)

func TestProtectMetaPath_Format(t *testing.T) {
	got := protectMetaPath("secrets/myapp/db")
	want := "vaultdiff/protect/secrets/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProtectMetaPath_TrimsSlashes(t *testing.T) {
	got := protectMetaPath("/secrets/myapp/")
	want := "vaultdiff/protect/secrets/myapp"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetProtection_EmptyPath(t *testing.T) {
	c := &mockLogical{}
	err := SetProtection(c, "", "safety")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetProtection_EmptyPath(t *testing.T) {
	c := &mockLogical{}
	_, _, err := GetProtection(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestClearProtection_EmptyPath(t *testing.T) {
	c := &mockLogical{}
	err := ClearProtection(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetProtection_NilSecret(t *testing.T) {
	c := &mockLogical{readSecret: nil}
	protected, reason, err := GetProtection(c, "secrets/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if protected {
		t.Error("expected not protected")
	}
	if reason != "" {
		t.Errorf("expected empty reason, got %q", reason)
	}
}

func TestSetProtection_Writes(t *testing.T) {
	c := &mockLogical{}
	err := SetProtection(c, "secrets/myapp", "do not touch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.writePath != "vaultdiff/protect/secrets/myapp" {
		t.Errorf("unexpected write path: %q", c.writePath)
	}
}
