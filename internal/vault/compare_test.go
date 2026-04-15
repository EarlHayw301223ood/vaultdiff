package vault

import (
	"testing"
)

func TestSameMount_SamePrefix(t *testing.T) {
	if !SameMount("secret/app/prod", "secret/app/staging") {
		t.Error("expected same mount for paths sharing 'secret' prefix")
	}
}

func TestSameMount_DifferentPrefix(t *testing.T) {
	if SameMount("kv/app/prod", "secret/app/prod") {
		t.Error("expected different mount for 'kv' vs 'secret'")
	}
}

func TestSameMount_NoSlash(t *testing.T) {
	if !SameMount("secret", "secret") {
		t.Error("expected same mount when paths have no slash")
	}
}

func TestSameMount_EmptyPaths(t *testing.T) {
	if !SameMount("", "") {
		t.Error("expected same mount for two empty paths")
	}
}

func TestMountPrefix_ReturnsFirstSegment(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"secret/foo/bar", "secret"},
		{"kv/prod/db", "kv"},
		{"nosegment", "nosegment"},
		{"", ""},
	}
	for _, tc := range cases {
		got := mountPrefix(tc.path)
		if got != tc.want {
			t.Errorf("mountPrefix(%q) = %q; want %q", tc.path, got, tc.want)
		}
	}
}

func TestVersionPair_Fields(t *testing.T) {
	vp := VersionPair{
		PathA:    "secret/prod/db",
		VersionA: 3,
		PathB:    "secret/staging/db",
		VersionB: 5,
	}
	if vp.PathA != "secret/prod/db" || vp.VersionA != 3 {
		t.Errorf("unexpected VersionPair fields: %+v", vp)
	}
	if vp.PathB != "secret/staging/db" || vp.VersionB != 5 {
		t.Errorf("unexpected VersionPair fields: %+v", vp)
	}
}
