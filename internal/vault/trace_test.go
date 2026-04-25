package vault

import (
	"testing"
	"time"
)

func TestTraceMetaPath_Format(t *testing.T) {
	got := traceMetaPath("secret/myapp/db")
	want := "vaultdiff-meta/secret/myapp/db/__trace__"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTraceMetaPath_TrimsSlashes(t *testing.T) {
	got := traceMetaPath("/secret/myapp/")
	want := "vaultdiff-meta/secret/myapp/__trace__"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTraceLog_Latest_Empty(t *testing.T) {
	log := &TraceLog{}
	if log.Latest() != nil {
		t.Error("expected nil for empty log")
	}
}

func TestTraceLog_Latest_ReturnsLast(t *testing.T) {
	log := &TraceLog{
		Entries: []TraceEntry{
			{Operation: "read", Timestamp: time.Now().Add(-time.Minute)},
			{Operation: "write", Timestamp: time.Now()},
		},
	}
	got := log.Latest()
	if got == nil || got.Operation != "write" {
		t.Errorf("expected write, got %v", got)
	}
}

func TestTraceLog_FilterByOperation(t *testing.T) {
	log := &TraceLog{
		Entries: []TraceEntry{
			{Operation: "read"},
			{Operation: "write"},
			{Operation: "READ"},
		},
	}
	results := log.FilterByOperation("read")
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestTraceLog_FilterByOperation_NoMatch(t *testing.T) {
	log := &TraceLog{
		Entries: []TraceEntry{
			{Operation: "read"},
		},
	}
	results := log.FilterByOperation("delete")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestAppendTrace_EmptyPath(t *testing.T) {
	client := &Client{}
	err := AppendTrace(client, "", "read", "user", "", 0)
	if err == nil || err.Error() != "trace: path must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAppendTrace_EmptyOperation(t *testing.T) {
	client := &Client{}
	err := AppendTrace(client, "secret/app", "", "user", "", 0)
	if err == nil || err.Error() != "trace: operation must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetTrace_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := GetTrace(client, "")
	if err == nil || err.Error() != "trace: path must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}
