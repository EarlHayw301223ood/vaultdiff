package vault

import (
	"errors"
	"fmt"
	"strings"
)

// EnvExportResult holds the rendered environment variable output for a secret path.
type EnvExportResult struct {
	Path   string
	Lines  []string
	Count  int
}

// ExportEnv fetches the secret at path/ref and returns KEY=VALUE lines
// suitable for use as shell environment variables or a .env file.
func ExportEnv(client *Client, path, ref string) (*EnvExportResult, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("path must not be empty")
	}

	secret, err := FetchAtRef(client, path, ref)
	if err != nil {
		return nil, fmt.Errorf("fetch secret: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %q", path)
	}

	data, err := extractStringMap(secret)
	if err != nil {
		return nil, fmt.Errorf("extract secret data: %w", err)
	}

	lines := make([]string, 0, len(data))
	for k, v := range data {
		key := toEnvKey(k)
		lines = append(lines, fmt.Sprintf("%s=%s", key, shellEscape(v)))
	}

	return &EnvExportResult{
		Path:  path,
		Lines: lines,
		Count: len(lines),
	}, nil
}

// toEnvKey converts a secret key to a canonical uppercase env var name,
// replacing non-alphanumeric characters with underscores.
func toEnvKey(k string) string {
	k = strings.ToUpper(k)
	var sb strings.Builder
	for _, r := range k {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('_')
		}
	}
	return sb.String()
}

// shellEscape wraps a value in single quotes, escaping any embedded single quotes.
func shellEscape(v string) string {
	escaped := strings.ReplaceAll(v, "'", "'\\''")
	return "'" + escaped + "'"
}
