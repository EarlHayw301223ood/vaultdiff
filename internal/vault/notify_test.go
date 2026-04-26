package vault

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNotifyMetaPath_Format(t *testing.T) {
	got := notifyMetaPath("secret/myapp/db")
	if got != "secret/metadata/_notify/myapp/db" {
		t.Errorf("unexpected path: %s", got)
	}
}

func TestNotifyMetaPath_TrimsSlashes(t *testing.T) {
	got := notifyMetaPath("/secret/myapp/")
	if got != "secret/metadata/_notify/myapp" {
		t.Errorf("unexpected path: %s", got)
	}
}

func TestNotifyMetaPath_SingleSegment(t *testing.T) {
	got := notifyMetaPath("secret")
	if !strings.Contains(got, "_notify") {
		t.Errorf("expected _notify in path, got: %s", got)
	}
}

func TestEventMatches_Wildcard(t *testing.T) {
	if !eventMatches("write", []string{"*"}) {
		t.Error("wildcard should match any event")
	}
}

func TestEventMatches_EmptyEvents(t *testing.T) {
	if !eventMatches("delete", []string{}) {
		t.Error("empty events list should match all")
	}
}

func TestEventMatches_CaseInsensitive(t *testing.T) {
	if !eventMatches("WRITE", []string{"write"}) {
		t.Error("event match should be case-insensitive")
	}
}

func TestEventMatches_NoMatch(t *testing.T) {
	if eventMatches("delete", []string{"write", "read"}) {
		t.Error("delete should not match write/read")
	}
}

func TestBuildNotifyBody_DefaultJSON(t *testing.T) {
	event := NotifyEvent{Path: "secret/db", Version: 3, Operation: "write", By: "alice"}
	body := buildNotifyBody(event, NotifyConfig{})
	if !strings.Contains(body, "secret/db") {
		t.Error("body should contain path")
	}
	if !strings.Contains(body, "alice") {
		t.Error("body should contain actor")
	}
}

func TestBuildNotifyBody_Template(t *testing.T) {
	event := NotifyEvent{Path: "secret/db", Version: 2, Operation: "read", By: "bob"}
	cfg := NotifyConfig{Template: "op={{operation}} path={{path}} by={{by}}"}
	body := buildNotifyBody(event, cfg)
	if body != "op=read path=secret/db by=bob" {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestNotify_DispatchesMatchingChannels(t *testing.T) {
	received := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	d := NewDispatcher(nil)
	event := NotifyEvent{Path: "secret/app", Version: 1, Operation: "write", Triggered: time.Now().UTC()}
	configs := []NotifyConfig{
		{Channel: ChannelCustom, Target: ts.URL, Events: []string{"write"}},
		{Channel: ChannelCustom, Target: ts.URL, Events: []string{"delete"}},
	}
	results := Notify(d, event, configs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Success {
		t.Errorf("expected success, got err: %s", results[0].Err)
	}
	if received != 1 {
		t.Errorf("expected 1 HTTP call, got %d", received)
	}
}
