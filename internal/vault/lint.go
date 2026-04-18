package vault

import (
	"fmt"
	"strings"
)

// LintRule defines a rule applied to secret keys or values.
type LintRule string

const (
	RuleNoEmptyValues  LintRule = "no-empty-values"
	RuleNoUppercaseKeys LintRule = "no-uppercase-keys"
	RuleNoSpacesInKeys LintRule = "no-spaces-in-keys"
)

// LintIssue describes a single lint violation.
type LintIssue struct {
	Path  string
	Key   string
	Rule  LintRule
	Detail string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("%s [%s] key=%q: %s", i.Path, i.Rule, i.Key, i.Detail)
}

// LintResult holds all issues found for a path.
type LintResult struct {
	Path   string
	Issues []LintIssue
}

func (r LintResult) Clean() bool {
	return len(r.Issues) == 0
}

// LintSecret runs lint rules against a secret's key/value pairs.
func LintSecret(path string, data map[string]string, rules []LintRule) LintResult {
	ruleSet := make(map[LintRule]bool, len(rules))
	for _, r := range rules {
		ruleSet[r] = true
	}

	result := LintResult{Path: path}

	for k, v := range data {
		if ruleSet[RuleNoEmptyValues] && strings.TrimSpace(v) == "" {
			result.Issues = append(result.Issues, LintIssue{
				Path: path, Key: k, Rule: RuleNoEmptyValues,
				Detail: "value is empty or whitespace",
			})
		}
		if ruleSet[RuleNoUppercaseKeys] && k != strings.ToLower(k) {
			result.Issues = append(result.Issues, LintIssue{
				Path: path, Key: k, Rule: RuleNoUppercaseKeys,
				Detail: "key contains uppercase characters",
			})
		}
		if ruleSet[RuleNoSpacesInKeys] && strings.Contains(k, " ") {
			result.Issues = append(result.Issues, LintIssue{
				Path: path, Key: k, Rule: RuleNoSpacesInKeys,
				Detail: "key contains spaces",
			})
		}
	}

	return result
}

// DefaultRules returns the default set of lint rules.
func DefaultRules() []LintRule {
	return []LintRule{RuleNoEmptyValues, RuleNoUppercaseKeys, RuleNoSpacesInKeys}
}
