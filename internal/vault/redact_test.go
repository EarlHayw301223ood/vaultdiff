package vault

import (
	"testing"
)

func TestRedactSecretData_NoSensitiveKeys(t *testing.T) {
	data := map[string]string{
		"host": "localhost",
		"port": "5432",
	}
	rules := DefaultRedactRules()
	result := RedactSecretData("secret/db", data, rules)

	if result.Redacted["host"] != "localhost" {
		t.Errorf("expected host to be unchanged, got %q", result.Redacted["host"])
	}
	if len(result.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", result.RedactedKeys)
	}
}

func TestRedactSecretData_PasswordKey(t *testing.T) {
	data := map[string]string{
		"username": "admin",
		"password": "s3cr3t",
	}
	rules := DefaultRedactRules()
	result := RedactSecretData("secret/app", data, rules)

	if result.Redacted["password"] != "[redacted:password]" {
		t.Errorf("expected password to be redacted, got %q", result.Redacted["password"])
	}
	if result.Redacted["username"] != "admin" {
		t.Errorf("expected username unchanged, got %q", result.Redacted["username"])
	}
	if len(result.RedactedKeys) != 1 || result.RedactedKeys[0] != "password" {
		t.Errorf("expected redacted keys [password], got %v", result.RedactedKeys)
	}
}

func TestRedactSecretData_TokenKey(t *testing.T) {
	data := map[string]string{
		"api_key": "abc123",
		"endpoint": "https://example.com",
	}
	rules := DefaultRedactRules()
	result := RedactSecretData("secret/svc", data, rules)

	if result.Redacted["api_key"] != "[redacted:token]" {
		t.Errorf("expected api_key to be redacted, got %q", result.Redacted["api_key"])
	}
}

func TestRedactSecretData_PathPreserved(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := RedactSecretData("secret/mypath", data, DefaultRedactRules())
	if result.Path != "secret/mypath" {
		t.Errorf("expected path to be preserved, got %q", result.Path)
	}
}

func TestRedactKey_Match(t *testing.T) {
	rules := DefaultRedactRules()
	if !RedactKey("db_password", rules) {
		t.Error("expected db_password to match redact rules")
	}
}

func TestRedactKey_NoMatch(t *testing.T) {
	rules := DefaultRedactRules()
	if RedactKey("hostname", rules) {
		t.Error("expected hostname to not match redact rules")
	}
}

func TestLabelForKey_ReturnsLabel(t *testing.T) {
	rules := DefaultRedactRules()
	label := LabelForKey("PRIVATE_KEY", rules)
	if label != "[redacted:private_key]" {
		t.Errorf("expected private_key label, got %q", label)
	}
}

func TestLabelForKey_NoMatch(t *testing.T) {
	rules := DefaultRedactRules()
	label := LabelForKey("region", rules)
	if label != "" {
		t.Errorf("expected empty label for non-sensitive key, got %q", label)
	}
}

func TestRedactSecretData_EmptyData(t *testing.T) {
	result := RedactSecretData("secret/empty", map[string]string{}, DefaultRedactRules())
	if len(result.Redacted) != 0 {
		t.Errorf("expected empty redacted map, got %v", result.Redacted)
	}
	if len(result.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", result.RedactedKeys)
	}
}
