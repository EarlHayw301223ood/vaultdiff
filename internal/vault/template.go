package vault

import (
	"errors"
	"fmt"
	"strings"
)

// TemplateResult holds the rendered output of a secret template.
type TemplateResult struct {
	Path     string
	Rendered string
}

// RenderTemplate substitutes placeholders in a template string with values
// from the secret at the given path and version ref.
//
// Placeholders use the form {{ key }}.
func RenderTemplate(client *Client, path, ref, tmpl string) (*TemplateResult, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if tmpl == "" {
		return nil, errors.New("template must not be empty")
	}

	secret, err := FetchAtRef(client, path, ref)
	if err != nil {
		return nil, fmt.Errorf("fetch secret: %w", err)
	}

	data := extractStringMap(secret)
	if data == nil {
		return nil, fmt.Errorf("no data found at path %q", path)
	}

	rendered, err := interpolate(tmpl, data)
	if err != nil {
		return nil, fmt.Errorf("interpolate: %w", err)
	}

	return &TemplateResult{Path: path, Rendered: rendered}, nil
}

// interpolate replaces {{ key }} placeholders with values from data.
func interpolate(tmpl string, data map[string]string) (string, error) {
	result := tmpl
	var missing []string

	for key, val := range data {
		placeholder := "{{ " + key + " }}"
		result = strings.ReplaceAll(result, placeholder, val)
	}

	// detect any unreplaced placeholders
	for {
		start := strings.Index(result, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}}")
		if end == -1 {
			break
		}
		token := strings.TrimSpace(result[start+2 : start+end])
		missing = append(missing, token)
		// remove to avoid infinite loop
		result = result[:start] + result[start+end+2:]
	}

	if len(missing) > 0 {
		return "", fmt.Errorf("unresolved placeholders: %s", strings.Join(missing, ", "))
	}

	return result, nil
}
