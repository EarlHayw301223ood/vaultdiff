package vault

import (
	"testing"
	"time"
)

func TestCheckExpiry_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := CheckExpiry(client, "", 1, time.Hour)
	if err == nil || err.Error() != "path must not be empty" {
		t.Fatalf("expected empty path error, got %v", err)
	}
}

func TestCheckExpiry_ZeroTTL(t *testing.T) {
	client := &Client{}
	_, err := CheckExpiry(client, "secret/foo", 1, 0)
	if err == nil || err.Error() != "ttl must be positive" {
		t.Fatalf("expected ttl error, got %v", err)
	}
}

func TestCheckExpiry_NegativeTTL(t *testing.T) {
	client := &Client{}
	_, err := CheckExpiry(client, "secret/foo", 1, -time.Minute)
	if err == nil {
		t.Fatal("expected error for negative ttl")
	}
}

func TestExpireResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	ttl := 24 * time.Hour
	res := &ExpireResult{
		Path:      "secret/myapp/db",
		Version:   3,
		CreatedAt: now.Add(-48 * time.Hour),
		ExpiredAt: now.Add(-48 * time.Hour).Add(ttl),
		Expired:   true,
	}

	if res.Path != "secret/myapp/db" {
		t.Errorf("unexpected path: %s", res.Path)
	}
	if res.Version != 3 {
		t.Errorf("unexpected version: %d", res.Version)
	}
	if !res.Expired {
		t.Error("expected expired=true")
	}
	if !res.ExpiredAt.Before(time.Now().UTC()) {
		t.Error("expected expiredAt to be in the past")
	}
}

func TestExpireResult_NotExpired(t *testing.T) {
	now := time.Now().UTC()
	res := &ExpireResult{
		Path:      "secret/fresh",
		Version:   1,
		CreatedAt: now.Add(-1 * time.Hour),
		ExpiredAt: now.Add(23 * time.Hour),
		Expired:   false,
	}
	if res.Expired {
		t.Error("expected expired=false")
	}
}
