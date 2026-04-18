package vault

import (
	"testing"
)

func TestLintSummary_Empty(t *testing.T) {
	summary := LintSummary(nil)
	if len(summary) != 0 {
		t.Errorf("expected empty summary, got %v", summary)
	}
}

func TestLintSummary_CountsByRule(t *testing.T) {
	results := []LintResult{
		{
			Path: "a",
			Issues: []LintIssue{
				{Rule: RuleNoEmptyValues},
				{Rule: RuleNoUppercaseKeys},
			},
		},
		{
			Path: "b",
			Issues: []LintIssue{
				{Rule: RuleNoEmptyValues},
			},
		},
	}
	summary := LintSummary(results)
	if summary[RuleNoEmptyValues] != 2 {
		t.Errorf("expected 2 for no-empty-values, got %d", summary[RuleNoEmptyValues])
	}
	if summary[RuleNoUppercaseKeys] != 1 {
		t.Errorf("expected 1 for no-uppercase-keys, got %d", summary[RuleNoUppercaseKeys])
	}
}

func TestLintMount_EmptyMount(t *testing.T) {
	_, err := LintMount(nil, nil, "", DefaultRules())
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}
