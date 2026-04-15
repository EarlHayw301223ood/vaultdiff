package vault

import (
	"testing"
)

func TestResolveVersionRef_Latest(t *testing.T) {
	ref, err := ResolveVersionRef("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ref.IsLatest {
		t.Error("expected IsLatest to be true")
	}
	if ref.Path != "secret/myapp" {
		t.Errorf("expected path 'secret/myapp', got %q", ref.Path)
	}
}

func TestResolveVersionRef_ExplicitLatest(t *testing.T) {
	ref, err := ResolveVersionRef("secret/myapp@latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ref.IsLatest {
		t.Error("expected IsLatest to be true")
	}
}

func TestResolveVersionRef_ExplicitVersion(t *testing.T) {
	ref, err := ResolveVersionRef("secret/myapp@5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.IsLatest {
		t.Error("expected IsLatest to be false")
	}
	if ref.Version != 5 {
		t.Errorf("expected version 5, got %d", ref.Version)
	}
	if ref.Path != "secret/myapp" {
		t.Errorf("expected path 'secret/myapp', got %q", ref.Path)
	}
}

func TestResolveVersionRef_InvalidVersion(t *testing.T) {
	cases := []string{
		"secret/myapp@abc",
		"secret/myapp@0",
		"secret/myapp@-1",
	}
	for _, tc := range cases {
		_, err := ResolveVersionRef(tc)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", tc)
		}
	}
}

func TestResolveVersionRef_EmptyPath(t *testing.T) {
	_, err := ResolveVersionRef("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestResolveVersionRef_EmptySegmentBeforeAt(t *testing.T) {
	_, err := ResolveVersionRef("@3")
	if err == nil {
		t.Error("expected error when path segment before '@' is empty")
	}
}

func TestVersionRef_String(t *testing.T) {
	latest := VersionRef{Path: "secret/app", IsLatest: true}
	if latest.String() != "secret/app@latest" {
		t.Errorf("unexpected string: %s", latest.String())
	}

	explicit := VersionRef{Path: "secret/app", Version: 3}
	if explicit.String() != "secret/app@3" {
		t.Errorf("unexpected string: %s", explicit.String())
	}
}
