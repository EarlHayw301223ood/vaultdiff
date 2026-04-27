package vault

import (
	"fmt"
	"strings"
	"time"
)

// ComplianceStatus represents the compliance state of a secret.
type ComplianceStatus string

const (
	CompliancePass    ComplianceStatus = "pass"
	ComplianceFail    ComplianceStatus = "fail"
	ComplianceWarning ComplianceStatus = "warning"
)

// ComplianceRule defines a single compliance check.
type ComplianceRule struct {
	Name        string
	Description string
	Check       func(data map[string]string) (bool, string)
}

// ComplianceResult holds the outcome of a compliance evaluation.
type ComplianceResult struct {
	Path      string
	Version   int
	Status    ComplianceStatus
	Violations []ComplianceViolation
	CheckedAt time.Time
}

// ComplianceViolation describes a single rule failure.
type ComplianceViolation struct {
	Rule    string
	Message string
}

// DefaultComplianceRules returns the built-in compliance rules.
func DefaultComplianceRules() []ComplianceRule {
	return []ComplianceRule{
		{
			Name:        "no-empty-values",
			Description: "All secret keys must have non-empty values",
			Check: func(data map[string]string) (bool, string) {
				for k, v := range data {
					if strings.TrimSpace(v) == "" {
						return false, fmt.Sprintf("key %q has an empty value", k)
					}
				}
				return true, ""
			},
		},
		{
			Name:        "no-plaintext-password-key",
			Description: "Keys named 'password' should not contain plaintext indicators",
			Check: func(data map[string]string) (bool, string) {
				for k, v := range data {
					if strings.EqualFold(k, "password") && strings.EqualFold(v, "password") {
						return false, "key 'password' appears to contain a literal plaintext password"
					}
				}
				return true, ""
			},
		},
		{
			Name:        "min-secret-count",
			Description: "Secrets must contain at least one key-value pair",
			Check: func(data map[string]string) (bool, string) {
				if len(data) == 0 {
					return false, "secret contains no key-value pairs"
				}
				return true, ""
			},
		},
	}
}

// CheckCompliance evaluates a secret's data against the provided rules.
func CheckCompliance(path string, version int, data map[string]string, rules []ComplianceRule) ComplianceResult {
	result := ComplianceResult{
		Path:      path,
		Version:   version,
		Status:    CompliancePass,
		CheckedAt: time.Now().UTC(),
	}

	for _, rule := range rules {
		ok, msg := rule.Check(data)
		if !ok {
			result.Violations = append(result.Violations, ComplianceViolation{
				Rule:    rule.Name,
				Message: msg,
			})
		}
	}

	if len(result.Violations) > 0 {
		result.Status = ComplianceFail
	}

	return result
}
