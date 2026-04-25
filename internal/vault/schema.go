package vault

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines a validation rule for a secret key.
type SchemaRule struct {
	Key      string         // exact key name, or empty to match via Pattern
	Pattern  *regexp.Regexp // regex pattern to match key names
	Required bool
	Format   *regexp.Regexp // optional regex the value must satisfy
	MinLen   int
}

// Schema holds a collection of rules used to validate secret data.
type Schema struct {
	Rules []SchemaRule
}

// SchemaViolation describes a single rule violation.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

// ValidateSecret checks secret data against the schema and returns all violations.
func ValidateSecret(data map[string]string, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for _, rule := range schema.Rules {
		matched := collectMatchingKeys(data, rule)

		if rule.Required && len(matched) == 0 {
			name := rule.Key
			if name == "" && rule.Pattern != nil {
				name = rule.Pattern.String()
			}
			violations = append(violations, SchemaViolation{Key: name, Message: "required key is missing"})
			continue
		}

		for _, k := range matched {
			v := data[k]
			if rule.MinLen > 0 && len(v) < rule.MinLen {
				violations = append(violations, SchemaViolation{
					Key:     k,
					Message: fmt.Sprintf("value too short (min %d chars)", rule.MinLen),
				})
			}
			if rule.Format != nil && !rule.Format.MatchString(v) {
				violations = append(violations, SchemaViolation{
					Key:     k,
					Message: fmt.Sprintf("value does not match required format %q", rule.Format.String()),
				})
			}
		}
	}
	return violations
}

// ViolationsToError converts a slice of violations into a single joined error.
func ViolationsToError(vs []SchemaViolation) error {
	if len(vs) == 0 {
		return nil
	}
	msgs := make([]string, len(vs))
	for i, v := range vs {
		msgs[i] = v.Error()
	}
	return errors.New(strings.Join(msgs, "; "))
}

func collectMatchingKeys(data map[string]string, rule SchemaRule) []string {
	var keys []string
	for k := range data {
		if rule.Key != "" && k == rule.Key {
			keys = append(keys, k)
		} else if rule.Pattern != nil && rule.Key == "" && rule.Pattern.MatchString(k) {
			keys = append(keys, k)
		}
	}
	return keys
}
