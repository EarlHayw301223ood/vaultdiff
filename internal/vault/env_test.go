package vault

import (
	"testing"
)

func TestToEnvKey_LowercaseToUpper(t *testing.T) {
	if got := toEnvKey("database_url"); got != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", got)
	}
}

func TestToEnvKey_HyphenBecomesUnderscore(t *testing.T) {
	if got := toEnvKey("my-secret-key"); got != "MY_SECRET_KEY" {
		t.Errorf("expected MY_SECRET_KEY, got %s", got)
	}
}

func TestToEnvKey_SpecialCharsBecomesUnderscore(t *testing.T) {
	if got := toEnvKey("app.port"); got != "APP_PORT" {
		t.Errorf("expected APP_PORT, got %s", got)
	}
}

func TestToEnvKey_AlreadyUppercase(t *testing.T) {
	if got := toEnvKey("TOKEN"); got != "TOKEN" {
		t.Errorf("expected TOKEN, got %s", got)
	}
}

func TestShellEscape_SimpleValue(t *testing.T) {
	if got := shellEscape("hello"); got != "'hello'" {
		t.Errorf("expected 'hello', got %s", got)
	}
}

func TestShellEscape_ValueWithSingleQuote(t *testing.T) {
	got := shellEscape("it's")
	expected := "'it'\\''s'"
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestShellEscape_EmptyValue(t *testing.T) {
	if got := shellEscape(""); got != "''" {
		t.Errorf("expected empty quoted string, got %s", got)
	}
}

func TestExportEnv_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := ExportEnv(client, "", "latest")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestExportEnv_WhitespacePath(t *testing.T) {
	client := &Client{}
	_, err := ExportEnv(client, "   ", "latest")
	if err == nil {
		t.Fatal("expected error for whitespace-only path")
	}
}
