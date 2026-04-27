package vault

import (
	"testing"
)

func TestHashSecretData_Deterministic(t *testing.T) {
	data := map[string]string{"foo": "bar", "baz": "qux"}
	hash1, _ := hashSecretData(data)
	hash2, _ := hashSecretData(data)
	if hash1 != hash2 {
		t.Errorf("expected deterministic hash, got %q and %q", hash1, hash2)
	}
}

func TestHashSecretData_OrderIndependent(t *testing.T) {
	a := map[string]string{"alpha": "1", "beta": "2", "gamma": "3"}
	b := map[string]string{"gamma": "3", "alpha": "1", "beta": "2"}
	hashA, _ := hashSecretData(a)
	hashB, _ := hashSecretData(b)
	if hashA != hashB {
		t.Errorf("expected order-independent hash, got %q and %q", hashA, hashB)
	}
}

func TestHashSecretData_DifferentData(t *testing.T) {
	a := map[string]string{"key": "valueA"}
	b := map[string]string{"key": "valueB"}
	hashA, _ := hashSecretData(a)
	hashB, _ := hashSecretData(b)
	if hashA == hashB {
		t.Error("expected different hashes for different data")
	}
}

func TestHashSecretData_KeyCount(t *testing.T) {
	data := map[string]string{"a": "1", "b": "2", "c": "3"}
	_, count := hashSecretData(data)
	if count != 3 {
		t.Errorf("expected key count 3, got %d", count)
	}
}

func TestHashSecretData_EmptyMap(t *testing.T) {
	hash, count := hashSecretData(map[string]string{})
	if hash == "" {
		t.Error("expected non-empty hash for empty map")
	}
	if count != 0 {
		t.Errorf("expected key count 0, got %d", count)
	}
}

func TestComputeFingerprint_EmptyPath(t *testing.T) {
	_, err := ComputeFingerprint(nil, "", 1)
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestComputeFingerprint_ZeroVersion(t *testing.T) {
	_, err := ComputeFingerprint(nil, "secret/data/myapp", 0)
	if err == nil {
		t.Error("expected error for zero version")
	}
}

func TestComputeFingerprint_NegativeVersion(t *testing.T) {
	_, err := ComputeFingerprint(nil, "secret/data/myapp", -1)
	if err == nil {
		t.Error("expected error for negative version")
	}
}

func TestFingerprintResult_Fields(t *testing.T) {
	r := &FingerprintResult{
		Path:        "secret/myapp",
		Version:     3,
		Fingerprint: "abc123",
		KeyCount:    5,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("unexpected Path: %q", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if r.KeyCount != 5 {
		t.Errorf("unexpected KeyCount: %d", r.KeyCount)
	}
}
