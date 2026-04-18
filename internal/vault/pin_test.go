package vault

import (
	"testing"
)

func TestPinMetaPath_Format(t *testing.T) {
	got := pinMetaPath("secret/myapp/db")
	want := "vaultdiff/meta/secret/myapp/db/pin"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPinMetaPath_TrimsSlashes(t *testing.T) {
	got := pinMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/pin"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetPin_EmptyPath(t *testing.T) {
	_, err := SetPin(nil, "", 1)
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestSetPin_InvalidVersion(t *testing.T) {
	_, err := SetPin(nil, "secret/app", 0)
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected invalid version error, got %v", err)
	}
}

func TestSetPin_NegativeVersion(t *testing.T) {
	_, err := SetPin(nil, "secret/app", -5)
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected invalid version error for negative version, got %v", err)
	}
}

func TestGetPin_EmptyPath(t *testing.T) {
	_, err := GetPin(nil, "")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestClearPin_EmptyPath(t *testing.T) {
	err := ClearPin(nil, "")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestPinResult_Fields(t *testing.T) {
	r := &PinResult{Path: "secret/app", Version: 3, Pinned: true}
	if r.Path != "secret/app" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if !r.Pinned {
		t.Error("expected Pinned to be true")
	}
}
