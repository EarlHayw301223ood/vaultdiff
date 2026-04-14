package diff

import (
	"strings"
	"testing"
)

func TestRender_NoChanges(t *testing.T) {
	result := &Result{
		Path:        "secret/app",
		FromVersion: 1,
		ToVersion:   2,
		Changes: []KeyDiff{
			{Key: "KEY", Change: ChangeUnchanged, OldValue: "val", NewValue: "val"},
		},
	}

	var sb strings.Builder
	Render(&sb, result, FormatOptions{})
	out := sb.String()

	if !strings.Contains(out, "(no changes)") {
		t.Errorf("expected '(no changes)' in output, got:\n%s", out)
	}
}

func TestRender_ShowsAddedKey(t *testing.T) {
	result := &Result{
		Path:        "secret/app",
		FromVersion: 1,
		ToVersion:   2,
		Changes: []KeyDiff{
			{Key: "NEW_KEY", Change: ChangeAdded, NewValue: "newval"},
		},
	}

	var sb strings.Builder
	Render(&sb, result, FormatOptions{})
	out := sb.String()

	if !strings.Contains(out, "+ NEW_KEY = newval") {
		t.Errorf("expected added key in output, got:\n%s", out)
	}
}

func TestRender_MasksValues(t *testing.T) {
	result := &Result{
		Path:        "secret/db",
		FromVersion: 1,
		ToVersion:   2,
		Changes: []KeyDiff{
			{Key: "DB_PASS", Change: ChangeModified, OldValue: "old_pass", NewValue: "new_pass"},
		},
	}

	var sb strings.Builder
	Render(&sb, result, FormatOptions{MaskValues: true})
	out := sb.String()

	if strings.Contains(out, "old_pass") || strings.Contains(out, "new_pass") {
		t.Errorf("expected values to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "********") {
		t.Errorf("expected masked value '********' in output, got:\n%s", out)
	}
}

func TestRender_SummaryLine(t *testing.T) {
	result := Compare("secret/app", 1,
		map[string]string{"A": "1"},
		2,
		map[string]string{"A": "2"},
	)

	var sb strings.Builder
	Render(&sb, result, FormatOptions{})
	out := sb.String()

	if !strings.Contains(out, "~1 modified") {
		t.Errorf("expected summary in output, got:\n%s", out)
	}
}
