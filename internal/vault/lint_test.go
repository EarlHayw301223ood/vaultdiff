package vault

import (
	"testing"
)

func TestLintSecret_NoIssues(t *testing.T) {
	data := map[string]string{"db_host": "localhost", "db_port": "5432"}
	result := LintSecret("secret/app", data, DefaultRules())
	if !result.Clean() {
		t.Fatalf("expected clean result, got %d issues", len(result.Issues))
	}
}

func TestLintSecret_EmptyValue(t *testing.T) {
	data := map[string]string{"api_key": ""}
	result := LintSecret("secret/app", data, []LintRule{RuleNoEmptyValues})
	if result.Clean() {
		t.Fatal("expected issue for empty value")
	}
	if result.Issues[0].Rule != RuleNoEmptyValues {
		t.Errorf("unexpected rule: %s", result.Issues[0].Rule)
	}
}

func TestLintSecret_UppercaseKey(t *testing.T) {
	data := map[string]string{"API_KEY": "abc123"}
	result := LintSecret("secret/app", data, []LintRule{RuleNoUppercaseKeys})
	if result.Clean() {
		t.Fatal("expected issue for uppercase key")
	}
	if result.Issues[0].Key != "API_KEY" {
		t.Errorf("unexpected key: %s", result.Issues[0].Key)
	}
}

func TestLintSecret_SpaceInKey(t *testing.T) {
	data := map[string]string{"my key": "val"}
	result := LintSecret("secret/app", data, []LintRule{RuleNoSpacesInKeys})
	if result.Clean() {
		t.Fatal("expected issue for space in key")
	}
	if result.Issues[0].Rule != RuleNoSpacesInKeys {
		t.Errorf("unexpected rule: %s", result.Issues[0].Rule)
	}
}

func TestLintSecret_MultipleIssues(t *testing.T) {
	data := map[string]string{"BAD KEY": ""}
	result := LintSecret("secret/app", data, DefaultRules())
	if len(result.Issues) < 2 {
		t.Fatalf("expected at least 2 issues, got %d", len(result.Issues))
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{Path: "secret/app", Key: "FOO", Rule: RuleNoUppercaseKeys, Detail: "key contains uppercase characters"}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestDefaultRules_NotEmpty(t *testing.T) {
	if len(DefaultRules()) == 0 {
		t.Error("expected non-empty default rules")
	}
}
