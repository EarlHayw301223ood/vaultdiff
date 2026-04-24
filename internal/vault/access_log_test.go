package vault

import (
	"testing"
	"time"
)

func TestAccessLogMetaPath_Format(t *testing.T) {
	got := accessLogMetaPath("secret/myapp/db")
	want := "vaultdiff/access-log/secret/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAccessLogMetaPath_TrimsSlashes(t *testing.T) {
	got := accessLogMetaPath("/secret/myapp/")
	want := "vaultdiff/access-log/secret/myapp"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAppendAccessLog_EmptyPath(t *testing.T) {
	c := newStubClient()
	err := AppendAccessLog(c, "", AccessEntry{Operation: "read"})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAppendAccessLog_EmptyOperation(t *testing.T) {
	c := newStubClient()
	err := AppendAccessLog(c, "secret/app", AccessEntry{})
	if err == nil {
		t.Fatal("expected error for empty operation")
	}
}

func TestAppendAccessLog_SetsTimestampIfZero(t *testing.T) {
	c := newStubClient()
	before := time.Now().UTC()
	err := AppendAccessLog(c, "secret/app", AccessEntry{
		Path:      "secret/app",
		Version:   1,
		Operation: "read",
		Actor:     "alice",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now().UTC()

	key := accessLogMetaPath("secret/app")
	data := c.written[key]
	ts, ok := data["timestamp"].(string)
	if !ok {
		t.Fatal("timestamp not written")
	}
	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Fatalf("bad timestamp format: %v", err)
	}
	if parsed.Before(before) || parsed.After(after) {
		t.Errorf("timestamp %v out of expected range", parsed)
	}
}

func TestGetAccessLog_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := GetAccessLog(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetAccessLog_NilSecret(t *testing.T) {
	c := newStubClient()
	entry, err := GetAccessLog(c, "secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for missing path")
	}
}

func TestAccessLogEntry_Fields(t *testing.T) {
	e := AccessEntry{
		Path:      "secret/app",
		Version:   3,
		Operation: "promote",
		Actor:     "bob",
		Note:      "promoted to prod",
		Timestamp: time.Now().UTC(),
	}
	if e.Path != "secret/app" {
		t.Errorf("unexpected Path: %s", e.Path)
	}
	if e.Version != 3 {
		t.Errorf("unexpected Version: %d", e.Version)
	}
	if e.Operation != "promote" {
		t.Errorf("unexpected Operation: %s", e.Operation)
	}
}
