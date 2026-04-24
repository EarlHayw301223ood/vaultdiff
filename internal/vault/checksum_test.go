package vault

import (
	"testing"
)

func TestChecksumMetaPath_Format(t *testing.T) {
	got := checksumMetaPath("secret/myapp/db")
	want := "vaultdiff/meta/secret/myapp/db/checksum"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestChecksumMetaPath_TrimsSlashes(t *testing.T) {
	got := checksumMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/checksum"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestComputeChecksum_Deterministic(t *testing.T) {
	data := map[string]string{"foo": "bar", "baz": "qux"}
	a := ComputeChecksum(data)
	b := ComputeChecksum(data)
	if a != b {
		t.Errorf("checksum not deterministic: %q vs %q", a, b)
	}
}

func TestComputeChecksum_OrderIndependent(t *testing.T) {
	a := ComputeChecksum(map[string]string{"x": "1", "y": "2"})
	b := ComputeChecksum(map[string]string{"y": "2", "x": "1"})
	if a != b {
		t.Errorf("checksum differs by insertion order: %q vs %q", a, b)
	}
}

func TestComputeChecksum_DifferentData(t *testing.T) {
	a := ComputeChecksum(map[string]string{"key": "value1"})
	b := ComputeChecksum(map[string]string{"key": "value2"})
	if a == b {
		t.Error("expected different checksums for different data")
	}
}

func TestComputeChecksum_EmptyMap(t *testing.T) {
	got := ComputeChecksum(map[string]string{})
	if got == "" {
		t.Error("expected non-empty checksum for empty map")
	}
}

func TestSaveChecksum_EmptyPath(t *testing.T) {
	err := SaveChecksum(nil, "", 1, "abc")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestSaveChecksum_InvalidVersion(t *testing.T) {
	err := SaveChecksum(nil, "secret/app", 0, "abc")
	if err == nil {
		t.Error("expected error for version <= 0")
	}
}

func TestSaveChecksum_EmptyChecksum(t *testing.T) {
	err := SaveChecksum(nil, "secret/app", 1, "")
	if err == nil {
		t.Error("expected error for empty checksum")
	}
}

func TestGetChecksum_EmptyPath(t *testing.T) {
	_, err := GetChecksum(nil, "", 1)
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestGetChecksum_InvalidVersion(t *testing.T) {
	_, err := GetChecksum(nil, "secret/app", -1)
	if err == nil {
		t.Error("expected error for version <= 0")
	}
}
