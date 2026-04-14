package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	from := map[string]string{"KEY": "value", "OTHER": "same"}
	to := map[string]string{"KEY": "value", "OTHER": "same"}

	result := Compare("secret/app", 1, from, 2, to)

	if result.HasChanges() {
		t.Errorf("expected no changes, got: %s", result.Summary())
	}
	if len(result.Changes) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Changes))
	}
}

func TestCompare_Added(t *testing.T) {
	from := map[string]string{"KEY": "value"}
	to := map[string]string{"KEY": "value", "NEW_KEY": "new"}

	result := Compare("secret/app", 1, from, 2, to)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	if result.Summary() != "+1 added" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestCompare_Removed(t *testing.T) {
	from := map[string]string{"KEY": "value", "OLD_KEY": "old"}
	to := map[string]string{"KEY": "value"}

	result := Compare("secret/app", 2, from, 3, to)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	if result.Summary() != "-1 removed" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestCompare_Modified(t *testing.T) {
	from := map[string]string{"DB_PASS": "old_pass"}
	to := map[string]string{"DB_PASS": "new_pass"}

	result := Compare("secret/db", 1, from, 2, to)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "DB_PASS" && c.Change == ChangeModified {
			found = true
			if c.OldValue != "old_pass" || c.NewValue != "new_pass" {
				t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
			}
		}
	}
	if !found {
		t.Error("expected DB_PASS to be modified")
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	from := map[string]string{"A": "1", "B": "2", "C": "3"}
	to := map[string]string{"A": "1", "B": "changed", "D": "4"}

	result := Compare("secret/mixed", 1, from, 2, to)

	if result.Summary() != "+1 added, -1 removed, ~1 modified" {
		t.Errorf("unexpected summary: %q", result.Summary())
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := Compare("secret/empty", 0, map[string]string{}, 1, map[string]string{})
	if result.HasChanges() {
		t.Error("expected no changes for two empty maps")
	}
}
