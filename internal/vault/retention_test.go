package vault

import (
	"testing"
	"time"
)

func TestRetentionMetaPath_Format(t *testing.T) {
	got := retentionMetaPath("secret/myapp/db")
	want := "secret/myapp/db/" + retentionMetaSuffix
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRetentionMetaPath_TrimsSlashes(t *testing.T) {
	got := retentionMetaPath("/secret/myapp/")
	want := "secret/myapp/" + retentionMetaSuffix
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetRetention_EmptyPath(t *testing.T) {
	client := newRollbackTestClient()
	err := SetRetention(client, "", 1, time.Hour, 10, "alice")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestSetRetention_InvalidVersion(t *testing.T) {
	client := newRollbackTestClient()
	err := SetRetention(client, "secret/app", 0, time.Hour, 10, "alice")
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestSetRetention_NegativeMaxAge(t *testing.T) {
	client := newRollbackTestClient()
	err := SetRetention(client, "secret/app", 1, -time.Hour, 10, "alice")
	if err == nil || err.Error() != "max_age must not be negative" {
		t.Errorf("expected max_age error, got %v", err)
	}
}

func TestSetRetention_NegativeMaxVersions(t *testing.T) {
	client := newRollbackTestClient()
	err := SetRetention(client, "secret/app", 1, time.Hour, -1, "alice")
	if err == nil || err.Error() != "max_versions must not be negative" {
		t.Errorf("expected max_versions error, got %v", err)
	}
}

func TestGetRetention_EmptyPath(t *testing.T) {
	client := newRollbackTestClient()
	_, err := GetRetention(client, "")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestRetentionPolicy_Fields(t *testing.T) {
	p := RetentionPolicy{
		Path:        "secret/app",
		Version:     3,
		MaxAge:      24 * time.Hour,
		MaxVersions: 5,
		SetBy:       "bob",
		SetAt:       time.Now().UTC(),
	}
	if p.Path != "secret/app" {
		t.Errorf("unexpected path: %s", p.Path)
	}
	if p.Version != 3 {
		t.Errorf("unexpected version: %d", p.Version)
	}
	if p.MaxVersions != 5 {
		t.Errorf("unexpected max_versions: %d", p.MaxVersions)
	}
	if p.SetBy != "bob" {
		t.Errorf("unexpected set_by: %s", p.SetBy)
	}
}
