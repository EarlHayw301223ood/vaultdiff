package vault

import (
	"testing"
)

func TestBookmarkMetaPath_Format(t *testing.T) {
	got := bookmarkMetaPath("secret/myapp/config")
	want := "vaultdiff/meta/secret/myapp/config/bookmarks"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBookmarkMetaPath_TrimsSlashes(t *testing.T) {
	got := bookmarkMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/bookmarks"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetBookmark_EmptyPath(t *testing.T) {
	err := SetBookmark(nil, "", "release", 1)
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestSetBookmark_EmptyName(t *testing.T) {
	err := SetBookmark(nil, "secret/app", "", 1)
	if err == nil || err.Error() != "bookmark name must not be empty" {
		t.Errorf("expected name error, got %v", err)
	}
}

func TestSetBookmark_InvalidVersion(t *testing.T) {
	err := SetBookmark(nil, "secret/app", "release", 0)
	if err == nil || err.Error() != "version must be >= 1" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestGetBookmark_EmptyPath(t *testing.T) {
	_, err := GetBookmark(nil, "", "release")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestGetBookmark_EmptyName(t *testing.T) {
	_, err := GetBookmark(nil, "secret/app", "")
	if err == nil || err.Error() != "bookmark name must not be empty" {
		t.Errorf("expected name error, got %v", err)
	}
}

func TestDeleteBookmark_EmptyPath(t *testing.T) {
	err := DeleteBookmark(nil, "", "release")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestDeleteBookmark_EmptyName(t *testing.T) {
	err := DeleteBookmark(nil, "secret/app", "")
	if err == nil || err.Error() != "bookmark name must not be empty" {
		t.Errorf("expected name error, got %v", err)
	}
}
