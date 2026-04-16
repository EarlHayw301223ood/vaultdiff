package vault

import (
	"testing"
)

func TestAnnotationMetaPath_Format(t *testing.T) {
	got := annotationMetaPath("secret/myapp", 3)
	want := ".annotations/secret/myapp/v3"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAnnotationMetaPath_SingleSegment(t *testing.T) {
	got := annotationMetaPath("myapp", 1)
	want := ".annotations/myapp/v1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetAnnotation_EmptyPath(t *testing.T) {
	_, err := SetAnnotation(nil, "", 1, "note", "alice")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestSetAnnotation_InvalidVersion(t *testing.T) {
	_, err := SetAnnotation(nil, "secret/app", 0, "note", "alice")
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestSetAnnotation_EmptyNote(t *testing.T) {
	_, err := SetAnnotation(nil, "secret/app", 1, "", "alice")
	if err == nil || err.Error() != "note must not be empty" {
		t.Errorf("expected note error, got %v", err)
	}
}

func TestGetAnnotation_EmptyPath(t *testing.T) {
	_, err := GetAnnotation(nil, "", 1)
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestGetAnnotation_InvalidVersion(t *testing.T) {
	_, err := GetAnnotation(nil, "secret/app", -1)
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestAnnotation_Fields(t *testing.T) {
	a := &Annotation{
		Path:    "secret/app",
		Version: 2,
		Note:    "bumped db password",
		Author:  "bob",
	}
	if a.Path != "secret/app" {
		t.Errorf("unexpected path: %s", a.Path)
	}
	if a.Version != 2 {
		t.Errorf("unexpected version: %d", a.Version)
	}
	if a.Note != "bumped db password" {
		t.Errorf("unexpected note: %s", a.Note)
	}
}
