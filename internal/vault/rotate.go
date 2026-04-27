package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// RotateResult holds the outcome of a secret rotation operation.
type RotateResult struct {
	Path       string
	OldVersion int
	NewVersion int
	RotatedAt  time.Time
}

// RotateOptions controls how rotation behaves.
type RotateOptions struct {
	// Transform is an optional function applied to each value before writing.
	// If nil, values are copied as-is.
	Transform func(key, value string) (string, error)
	// DryRun skips the write and returns a simulated result.
	DryRun bool
}

// Rotate reads the latest version of the secret at path, optionally
// transforms its values, and writes a new version back. It returns the
// old and new version numbers so callers can audit the change.
func Rotate(client LogicalClient, path string, opts RotateOptions) (*RotateResult, error) {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil, errors.New("rotate: path must not be empty")
	}

	reaPath := kvV2DataPath(path)
	secret, err := client.Read(reaPath)
	if err != nil {
		return nil, fmt.Errorf("rotate: read %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("rotate: secret not found at %s", path)
	}

	data := extractStringMap(secret)
	meta := extractMetadata(secret)
	oldVersion := meta.Version

	newData := make(map[string]interface{}, len(data))
	for k, v := range data {
		if opts.Transform != nil {
			transformed, terr := opts.Transform(k, v)
			if terr != nil {
				return nil, fmt.Errorf("rotate: transform key %q: %w", k, terr)
			}
			newData[k] = transformed
		} else {
			newData[k] = v
		}
	}

	if opts.DryRun {
		return &RotateResult{
			Path:       path,
			OldVersion: oldVersion,
			NewVersion: oldVersion + 1,
			RotatedAt:  time.Now().UTC(),
		}, nil
	}

	writePath := kvV2WritePath(path)
	_, err = client.Write(writePath, map[string]interface{}{"data": newData})
	if err != nil {
		return nil, fmt.Errorf("rotate: write %s: %w", path, err)
	}

	return &RotateResult{
		Path:       path,
		OldVersion: oldVersion,
		NewVersion: oldVersion + 1,
		RotatedAt:  time.Now().UTC(),
	}, nil
}
