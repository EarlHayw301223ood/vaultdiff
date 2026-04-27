package vault

import (
	"testing"
)

func TestCompareMounts_EmptyMountA(t *testing.T) {
	c := newStubClient(nil)
	_, err := CompareMounts(c, "", "secret", "")
	if err == nil {
		t.Fatal("expected error for empty mountA")
	}
}

func TestCompareMounts_EmptyMountB(t *testing.T) {
	c := newStubClient(nil)
	_, err := CompareMounts(c, "secret", "", "")
	if err == nil {
		t.Fatal("expected error for empty mountB")
	}
}

func TestMountCompareResult_Summary(t *testing.T) {
	r := &MountCompareResult{
		MountA:   "dev",
		MountB:   "prod",
		OnlyInA:  []string{"a"},
		OnlyInB:  []string{"b", "c"},
		Diverged: []string{"d"},
		InSync:   []string{"e", "f", "g"},
	}
	got := r.Summary()
	if got == "" {
		t.Fatal("expected non-empty summary")
	}
	for _, want := range []string{"dev", "prod", "1", "2", "3"} {
		if !containsStr(got, want) {
			t.Errorf("summary %q missing %q", got, want)
		}
	}
}

func TestFilterPaths_EmptyPrefix(t *testing.T) {
	paths := []string{"a/b", "c/d"}
	got := filterPaths(paths, "")
	if len(got) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(got))
	}
}

func TestFilterPaths_MatchingPrefix(t *testing.T) {
	paths := []string{"app/db", "app/api", "infra/net"}
	got := filterPaths(paths, "app/")
	if len(got) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(got))
	}
}

func TestFilterPaths_NoneMatch(t *testing.T) {
	paths := []string{"a/b", "c/d"}
	got := filterPaths(paths, "z/")
	if len(got) != 0 {
		t.Fatalf("expected 0 paths, got %d", len(got))
	}
}

// containsStr is a small helper to avoid importing strings in tests.
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
