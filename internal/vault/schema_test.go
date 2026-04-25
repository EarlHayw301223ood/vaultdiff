package vault

import (
	"regexp"
	"testing"
)

func TestValidateSecret_NoViolations(t *testing.T) {
	data := map[string]string{"api_key": "abc123", "db_pass": "secret"}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "api_key", Required: true, MinLen: 3},
		},
	}
	vs := ValidateSecret(data, schema)
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(vs), vs)
	}
}

func TestValidateSecret_RequiredKeyMissing(t *testing.T) {
	data := map[string]string{"other": "value"}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "api_key", Required: true},
		},
	}
	vs := ValidateSecret(data, schema)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "api_key" {
		t.Errorf("expected violation for api_key, got %q", vs[0].Key)
	}
}

func TestValidateSecret_ValueTooShort(t *testing.T) {
	data := map[string]string{"token": "ab"}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "token", Required: true, MinLen: 8},
		},
	}
	vs := ValidateSecret(data, schema)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
}

func TestValidateSecret_FormatMismatch(t *testing.T) {
	data := map[string]string{"port": "not-a-number"}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "port", Format: regexp.MustCompile(`^\d+$`)},
		},
	}
	vs := ValidateSecret(data, schema)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
}

func TestValidateSecret_PatternMatchesMultipleKeys(t *testing.T) {
	data := map[string]string{"svc_timeout": "0", "svc_retries": "x", "unrelated": "ok"}
	schema := Schema{
		Rules: []SchemaRule{
			{Pattern: regexp.MustCompile(`^svc_`), MinLen: 1, Format: regexp.MustCompile(`^\d+$`)},
		},
	}
	vs := ValidateSecret(data, schema)
	// both svc_timeout and svc_retries fail format; svc_timeout also has len>=1
	if len(vs) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(vs), vs)
	}
}

func TestViolationsToError_Nil(t *testing.T) {
	if err := ViolationsToError(nil); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestViolationsToError_NonEmpty(t *testing.T) {
	vs := []SchemaViolation{{Key: "a", Message: "bad"}, {Key: "b", Message: "missing"}}
	err := ViolationsToError(vs)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}
