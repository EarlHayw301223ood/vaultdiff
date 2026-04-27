package vault

import (
	"fmt"
	"strings"
	"unicode"
)

// SanitizeResult holds the outcome of sanitizing a secret's data.
type SanitizeResult struct {
	Path        string
	Original    map[string]string
	Sanitized   map[string]string
	ChangedKeys []string
}

// SanitizeOptions controls how sanitization is applied.
type SanitizeOptions struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// NormalizeKeys lowercases all keys and replaces spaces with underscores.
	NormalizeKeys bool
	// StripNonPrintable removes non-printable characters from values.
	StripNonPrintable bool
}

// DefaultSanitizeOptions returns a sensible default sanitization config.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		TrimSpace:         true,
		NormalizeKeys:     false,
		StripNonPrintable: true,
	}
}

// SanitizeSecret reads a secret at path, applies sanitization rules, and
// returns a SanitizeResult describing what changed. It does NOT write back;
// callers decide whether to persist the result.
func SanitizeSecret(client LogicalClient, path string, opts SanitizeOptions) (*SanitizeResult, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("sanitize: path must not be empty")
	}

	secret, err := client.Read(kvV2DataPath(path))
	if err != nil {
		return nil, fmt.Errorf("sanitize: read %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("sanitize: secret not found at %s", path)
	}

	original, err := extractStringMap(secret)
	if err != nil {
		return nil, fmt.Errorf("sanitize: extract data: %w", err)
	}

	sanitized := make(map[string]string, len(original))
	var changed []string

	for k, v := range original {
		newKey := k
		newVal := v

		if opts.NormalizeKeys {
			newKey = strings.ToLower(strings.ReplaceAll(k, " ", "_"))
		}
		if opts.TrimSpace {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.StripNonPrintable {
			newVal = strings.Map(func(r rune) rune {
				if unicode.IsPrint(r) {
					return r
				}
				return -1
			}, newVal)
		}

		sanitized[newKey] = newVal
		if newKey != k || newVal != v {
			changed = append(changed, k)
		}
	}

	return &SanitizeResult{
		Path:        path,
		Original:    original,
		Sanitized:   sanitized,
		ChangedKeys: changed,
	}, nil
}
