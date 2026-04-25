package vault

import (
	"testing"
)

func TestCompareVersions_EmptyPath(t *testing.T) {
	_, err := CompareVersions(nil, "", 1, 2)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCompareVersions_ZeroVersionA(t *testing.T) {
	_, err := CompareVersions(nil, "secret/myapp", 0, 2)
	if err == nil {
		t.Fatal("expected error for zero versionA")
	}
}

func TestCompareVersions_ZeroVersionB(t *testing.T) {
	_, err := CompareVersions(nil, "secret/myapp", 1, 0)
	if err == nil {
		t.Fatal("expected error for zero versionB")
	}
}

func TestCompareVersions_NegativeVersion(t *testing.T) {
	_, err := CompareVersions(nil, "secret/myapp", -1, 2)
	if err == nil {
		t.Fatal("expected error for negative version")
	}
}

func TestCompareVersions_SameVersions(t *testing.T) {
	_, err := CompareVersions(nil, "secret/myapp", 3, 3)
	if err == nil {
		t.Fatal("expected error when versionA == versionB")
	}
	if err.Error() != "versionA and versionB must differ" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestVersionComparison_Fields(t *testing.T) {
	cmp := &VersionComparison{
		Path:     "secret/myapp",
		VersionA: 1,
		VersionB: 3,
		DataA:    map[string]string{"key": "old"},
		DataB:    map[string]string{"key": "new"},
	}
	if cmp.Path != "secret/myapp" {
		t.Errorf("unexpected path: %s", cmp.Path)
	}
	if cmp.DataA["key"] != "old" {
		t.Errorf("expected DataA key=old, got %s", cmp.DataA["key"])
	}
	if cmp.DataB["key"] != "new" {
		t.Errorf("expected DataB key=new, got %s", cmp.DataB["key"])
	}
}

func TestCompareVersions_WhitespacePath(t *testing.T) {
	_, err := CompareVersions(nil, "   ", 1, 2)
	if err == nil {
		t.Fatal("expected error for whitespace-only path")
	}
}
