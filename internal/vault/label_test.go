package vault

import (
	"testing"
)

func TestLabelMetaPath_Format(t *testing.T) {
	got := labelMetaPath("secret/myapp/db", "env")
	want := "secret/myapp/db/__meta__/labels/env"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLabelMetaPath_TrimsSlashes(t *testing.T) {
	got := labelMetaPath("/secret/myapp/", "tier")
	want := "secret/myapp/__meta__/labels/tier"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetLabel_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := SetLabel(client, "", "env", "production", 1)
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestSetLabel_EmptyName(t *testing.T) {
	client := &Client{}
	_, err := SetLabel(client, "secret/myapp", "", "production", 1)
	if err == nil || err.Error() != "label name must not be empty" {
		t.Errorf("expected name error, got %v", err)
	}
}

func TestSetLabel_InvalidVersion(t *testing.T) {
	client := &Client{}
	_, err := SetLabel(client, "secret/myapp", "env", "production", 0)
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestSetLabel_NegativeVersion(t *testing.T) {
	client := &Client{}
	_, err := SetLabel(client, "secret/myapp", "env", "production", -3)
	if err == nil {
		t.Error("expected error for negative version")
	}
}

func TestGetLabel_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := GetLabel(client, "", "env")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestGetLabel_EmptyName(t *testing.T) {
	client := &Client{}
	_, err := GetLabel(client, "secret/myapp", "")
	if err == nil || err.Error() != "label name must not be empty" {
		t.Errorf("expected name error, got %v", err)
	}
}

func TestLabelEntry_Fields(t *testing.T) {
	entry := &LabelEntry{
		Path:    "secret/myapp/db",
		Name:    "env",
		Value:   "production",
		Version: 3,
	}
	if entry.Path != "secret/myapp/db" {
		t.Errorf("unexpected Path: %s", entry.Path)
	}
	if entry.Name != "env" {
		t.Errorf("unexpected Name: %s", entry.Name)
	}
	if entry.Value != "production" {
		t.Errorf("unexpected Value: %s", entry.Value)
	}
	if entry.Version != 3 {
		t.Errorf("unexpected Version: %d", entry.Version)
	}
}
