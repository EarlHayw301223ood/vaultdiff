package vault

import (
	"testing"
)

func TestMatchesPattern_CaseInsensitive(t *testing.T) {
	if !matchesPattern("MySecret", "secret", false) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestMatchesPattern_CaseSensitive_NoMatch(t *testing.T) {
	if matchesPattern("MySecret", "secret", true) {
		t.Fatal("expected no match with case-sensitive")
	}
}

func TestMatchesPattern_CaseSensitive_Match(t *testing.T) {
	if !matchesPattern("MySecret", "Secret", true) {
		t.Fatal("expected match")
	}
}

func TestMatchesPattern_EmptyPattern(t *testing.T) {
	if matchesPattern("anything", "", false) {
		t.Fatal("empty pattern should not match")
	}
}

func TestKvV2DataPath_Format(t *testing.T) {
	got := kvV2DataPath("secret", "myapp/config")
	want := "secret/data/myapp/config"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestKvV2DataPath_TrimsSlashes(t *testing.T) {
	got := kvV2DataPath("/secret/", "/myapp/config/")
	want := "secret/data/myapp/config"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{Path: "myapp/db", MatchedKeys: []string{"password"}}
	if r.Path != "myapp/db" {
		t.Fatalf("unexpected path: %s", r.Path)
	}
	if len(r.MatchedKeys) != 1 || r.MatchedKeys[0] != "password" {
		t.Fatalf("unexpected matched keys: %v", r.MatchedKeys)
	}
}
