package vault

import (
	"testing"
)

func TestTagMetaPath_Standard(t *testing.T) {
	got := tagMetaPath("secret", "myapp/config", "stable")
	want := "secret/metadata/tags/myapp/config/stable"
	if got != want {
		t.Errorf("tagMetaPath = %q, want %q", got, want)
	}
}

func TestSetTag_EmptyPath(t *testing.T) {
	c := &Client{logical: &mockLogical{}}
	err := SetTag(c, "secret", "", "stable", 1)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetTag_EmptyTagName(t *testing.T) {
	c := &Client{logical: &mockLogical{}}
	err := SetTag(c, "secret", "myapp/config", "", 1)
	if err == nil {
		t.Fatal("expected error for empty tag name")
	}
}

func TestSetTag_InvalidVersion(t *testing.T) {
	c := &Client{logical: &mockLogical{}}
	err := SetTag(c, "secret", "myapp/config", "stable", 0)
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestGetTag_EmptyPath(t *testing.T) {
	c := &Client{logical: &mockLogical{}}
	_, err := GetTag(c, "secret", "", "stable")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetTag_EmptyTagName(t *testing.T) {
	c := &Client{logical: &mockLogical{}}
	_, err := GetTag(c, "secret", "myapp/config", "")
	if err == nil {
		t.Fatal("expected error for empty tag name")
	}
}

func TestGetTag_NotFound(t *testing.T) {
	ml := &mockLogical{readSecret: nil}
	c := &Client{logical: ml}
	_, err := GetTag(c, "secret", "myapp/config", "stable")
	if err == nil {
		t.Fatal("expected error when tag secret is nil")
	}
}

func TestGetTag_ReturnsTag(t *testing.T) {
	ml := &mockLogical{
		readSecret: makeVersionsSecret(map[string]interface{}{
			"version": float64(3),
			"path":    "myapp/config",
		}),
	}
	c := &Client{logical: ml}
	tag, err := GetTag(c, "secret", "myapp/config", "stable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Version != 3 {
		t.Errorf("tag.Version = %d, want 3", tag.Version)
	}
	if tag.Name != "stable" {
		t.Errorf("tag.Name = %q, want %q", tag.Name, "stable")
	}
	if tag.Path != "myapp/config" {
		t.Errorf("tag.Path = %q, want %q", tag.Path, "myapp/config")
	}
}
