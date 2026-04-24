package vault

import (
	"testing"
	"time"
)

func TestTTLMetaPath_Format(t *testing.T) {
	got := ttlMetaPath("secret/my-app")
	want := "vaultdiff/meta/ttl/secret/my-app"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTTLMetaPath_TrimsSlashes(t *testing.T) {
	got := ttlMetaPath("/secret/my-app/")
	want := "vaultdiff/meta/ttl/secret/my-app"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetTTL_EmptyPath(t *testing.T) {
	c := newStubClient()
	err := SetTTL(c, "", 1, time.Hour, "user")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetTTL_InvalidVersion(t *testing.T) {
	c := newStubClient()
	err := SetTTL(c, "secret/app", 0, time.Hour, "user")
	if err == nil {
		t.Fatal("expected error for version <= 0")
	}
}

func TestSetTTL_NegativeDuration(t *testing.T) {
	c := newStubClient()
	err := SetTTL(c, "secret/app", 1, -time.Minute, "user")
	if err == nil {
		t.Fatal("expected error for non-positive ttl")
	}
}

func TestGetTTL_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := GetTTL(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetTTL_NilSecret(t *testing.T) {
	c := newStubClient()
	rec, err := GetTTL(c, "secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Errorf("expected nil record for missing path, got %+v", rec)
	}
}

func TestTTLRecord_IsExpired_ZeroTime(t *testing.T) {
	r := TTLRecord{}
	if r.IsExpired() {
		t.Error("zero ExpiresAt should not be expired")
	}
}

func TestTTLRecord_IsExpired_Future(t *testing.T) {
	r := TTLRecord{ExpiresAt: time.Now().UTC().Add(time.Hour)}
	if r.IsExpired() {
		t.Error("future ExpiresAt should not be expired")
	}
}

func TestTTLRecord_IsExpired_Past(t *testing.T) {
	r := TTLRecord{ExpiresAt: time.Now().UTC().Add(-time.Second)}
	if !r.IsExpired() {
		t.Error("past ExpiresAt should be expired")
	}
}

func TestTTLRecord_RemainingTTL_Zero(t *testing.T) {
	r := TTLRecord{ExpiresAt: time.Now().UTC().Add(-time.Minute)}
	if r.RemainingTTL() != 0 {
		t.Error("expired record should return 0 remaining")
	}
}

func TestTTLRecord_RemainingTTL_Positive(t *testing.T) {
	r := TTLRecord{ExpiresAt: time.Now().UTC().Add(2 * time.Hour)}
	rem := r.RemainingTTL()
	if rem <= 0 || rem > 2*time.Hour {
		t.Errorf("unexpected remaining TTL: %v", rem)
	}
}
