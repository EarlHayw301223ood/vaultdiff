package vault

import (
	"testing"
)

func TestSanitizeSecret_EmptyPath(t *testing.T) {
	client := newStubClient(t, nil)
	_, err := SanitizeSecret(client, "", DefaultSanitizeOptions())
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSanitizeSecret_TrimSpace(t *testing.T) {
	data := map[string]interface{}{
		"key": "  hello  ",
		"clean": "world",
	}
	client := newStubClient(t, data)

	opts := SanitizeOptions{TrimSpace: true}
	res, err := SanitizeSecret(client, "secret/data/test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Sanitized["key"] != "hello" {
		t.Errorf("expected trimmed value, got %q", res.Sanitized["key"])
	}
	if res.Sanitized["clean"] != "world" {
		t.Errorf("expected unchanged value, got %q", res.Sanitized["clean"])
	}
	if len(res.ChangedKeys) != 1 || res.ChangedKeys[0] != "key" {
		t.Errorf("expected only 'key' in changed list, got %v", res.ChangedKeys)
	}
}

func TestSanitizeSecret_NormalizeKeys(t *testing.T) {
	data := map[string]interface{}{
		"My Key": "value",
		"already_ok": "fine",
	}
	client := newStubClient(t, data)

	opts := SanitizeOptions{NormalizeKeys: true}
	res, err := SanitizeSecret(client, "secret/data/test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := res.Sanitized["my_key"]; !ok {
		t.Error("expected key 'my_key' after normalization")
	}
	if _, ok := res.Sanitized["already_ok"]; !ok {
		t.Error("expected key 'already_ok' to remain")
	}
}

func TestSanitizeSecret_StripNonPrintable(t *testing.T) {
	data := map[string]interface{}{
		"token": "abc\x00def\x1f",
	}
	client := newStubClient(t, data)

	opts := SanitizeOptions{StripNonPrintable: true}
	res, err := SanitizeSecret(client, "secret/data/test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Sanitized["token"] != "abcdef" {
		t.Errorf("expected non-printable chars stripped, got %q", res.Sanitized["token"])
	}
}

func TestSanitizeSecret_NoChanges(t *testing.T) {
	data := map[string]interface{}{
		"key": "clean_value",
	}
	client := newStubClient(t, data)

	res, err := SanitizeSecret(client, "secret/data/test", DefaultSanitizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.ChangedKeys) != 0 {
		t.Errorf("expected no changed keys, got %v", res.ChangedKeys)
	}
}

func TestDefaultSanitizeOptions(t *testing.T) {
	opts := DefaultSanitizeOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true by default")
	}
	if !opts.StripNonPrintable {
		t.Error("expected StripNonPrintable to be true by default")
	}
	if opts.NormalizeKeys {
		t.Error("expected NormalizeKeys to be false by default")
	}
}
