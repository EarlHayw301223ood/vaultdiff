package vault

import (
	"testing"
)

func TestSignMetaPath_Format(t *testing.T) {
	got := signMetaPath("secret/myapp/db")
	want := "vaultdiff/signatures/secret/myapp/db"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestSignMetaPath_TrimsSlashes(t *testing.T) {
	got := signMetaPath("/secret/myapp/")
	want := "vaultdiff/signatures/secret/myapp"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestComputeSignature_Deterministic(t *testing.T) {
	data := map[string]string{"foo": "bar", "baz": "qux"}
	s1 := computeSignature(data, "mykey")
	s2 := computeSignature(data, "mykey")
	if s1 != s2 {
		t.Errorf("expected deterministic signature, got %q and %q", s1, s2)
	}
}

func TestComputeSignature_DifferentKeys(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	s1 := computeSignature(data, "key1")
	s2 := computeSignature(data, "key2")
	if s1 == s2 {
		t.Error("expected different signatures for different hmac keys")
	}
}

func TestComputeSignature_DifferentData(t *testing.T) {
	s1 := computeSignature(map[string]string{"a": "1"}, "key")
	s2 := computeSignature(map[string]string{"a": "2"}, "key")
	if s1 == s2 {
		t.Error("expected different signatures for different data")
	}
}

func TestSignSecret_EmptyPath(t *testing.T) {
	_, err := SignSecret(nil, "", 1, "key")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSignSecret_InvalidVersion(t *testing.T) {
	_, err := SignSecret(nil, "secret/app", 0, "key")
	if err == nil || err.Error() != "version must be >= 1" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSignSecret_EmptyHMACKey(t *testing.T) {
	_, err := SignSecret(nil, "secret/app", 1, "")
	if err == nil || err.Error() != "hmac key must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifySecret_EmptyPath(t *testing.T) {
	_, err := VerifySecret(nil, "", 1, "key")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifySecret_InvalidVersion(t *testing.T) {
	_, err := VerifySecret(nil, "secret/app", 0, "key")
	if err == nil || err.Error() != "version must be >= 1" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSignResult_Fields(t *testing.T) {
	r := &SignResult{Path: "secret/app", Version: 3, Signature: "abc123"}
	if r.Path != "secret/app" || r.Version != 3 || r.Signature != "abc123" {
		t.Errorf("unexpected fields: %+v", r)
	}
}
