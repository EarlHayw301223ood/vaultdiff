package vault

import (
	"regexp"
	"strings"
)

// RedactRule defines a pattern and its replacement label.
type RedactRule struct {
	Name    string
	Pattern *regexp.Regexp
	Label   string
}

// RedactResult holds the outcome of redacting a secret map.
type RedactResult struct {
	Path        string
	Redacted    map[string]string
	RedactedKeys []string
}

// DefaultRedactRules returns a set of built-in sensitive key patterns.
func DefaultRedactRules() []RedactRule {
	return []RedactRule{
		{
			Name:    "password",
			Pattern: regexp.MustCompile(`(?i)(password|passwd|pwd)`),
			Label:   "[redacted:password]",
		},
		{
			Name:    "token",
			Pattern: regexp.MustCompile(`(?i)(token|secret|api_key|apikey)`),
			Label:   "[redacted:token]",
		},
		{
			Name:    "private_key",
			Pattern: regexp.MustCompile(`(?i)(private_key|privkey|pem)`),
			Label:   "[redacted:private_key]",
		},
	}
}

// RedactSecretData applies redaction rules to a secret key-value map.
// Keys matching any rule pattern have their values replaced with the rule label.
func RedactSecretData(path string, data map[string]string, rules []RedactRule) RedactResult {
	result := RedactResult{
		Path:    path,
		Redacted: make(map[string]string, len(data)),
	}

	for k, v := range data {
		matched := false
		for _, rule := range rules {
			if rule.Pattern.MatchString(k) {
				result.Redacted[k] = rule.Label
				result.RedactedKeys = append(result.RedactedKeys, k)
				matched = true
				break
			}
		}
		if !matched {
			result.Redacted[k] = v
		}
	}

	return result
}

// RedactKey returns true if the given key matches any of the provided rules.
func RedactKey(key string, rules []RedactRule) bool {
	for _, rule := range rules {
		if rule.Pattern.MatchString(key) {
			return true
		}
	}
	return false
}

// LabelForKey returns the redaction label for the first matching rule, or empty string.
func LabelForKey(key string, rules []RedactRule) string {
	for _, rule := range rules {
		if rule.Pattern.MatchString(strings.ToLower(key)) {
			return rule.Label
		}
	}
	return ""
}
