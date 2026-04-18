package vault

import (
	"testing"
)

// TestSearchOptions_BothPatternsEmpty verifies that a search with both
// key and value patterns empty returns no results (nothing to match).
func TestSearchOptions_BothPatternsEmpty(t *testing.T) {
	opts := SearchOptions{KeyPattern: "", ValuePattern: ""}
	// matchesPattern returns false for empty pattern
	if matchesPattern("anykey", opts.KeyPattern, false) {
		t.Fatal("empty key pattern should not match")
	}
	if matchesPattern("anyvalue", opts.ValuePattern, false) {
		t.Fatal("empty value pattern should not match")
	}
}

func TestSearchOptions_KeyPatternOnly(t *testing.T) {
	opts := SearchOptions{KeyPattern: "pass", ValuePattern: ""}
	if !matchesPattern("password", opts.KeyPattern, false) {
		t.Fatal("expected key pattern to match")
	}
	if matchesPattern("somevalue", opts.ValuePattern, false) {
		t.Fatal("empty value pattern should not match")
	}
}

func TestSearchOptions_ValuePatternOnly(t *testing.T) {
	opts := SearchOptions{KeyPattern: "", ValuePattern: "prod"}
	if matchesPattern("mykey", opts.KeyPattern, false) {
		t.Fatal("empty key pattern should not match")
	}
	if !matchesPattern("production", opts.ValuePattern, false) {
		t.Fatal("expected value pattern to match")
	}
}

func TestSearchResult_EmptyMatchedKeys(t *testing.T) {
	r := SearchResult{Path: "x", MatchedKeys: nil}
	if len(r.MatchedKeys) != 0 {
		t.Fatal("expected empty matched keys")
	}
}
