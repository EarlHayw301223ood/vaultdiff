package vault

import (
	"fmt"
	"strings"
)

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourcePath string
	DestPath   string
	Version    int
	Skipped    bool
	SkipReason string
}

// CloneOptions configures the behaviour of Clone.
type CloneOptions struct {
	// OverwriteExisting controls whether an existing secret at DestPath is overwritten.
	OverwriteExisting bool
}

// Clone copies the secret at sourcePath (at the given version ref) to destPath.
// It returns a CloneResult describing what happened.
func Clone(logical LogicalClient, sourcePath, destPath, ref string, opts CloneOptions) (CloneResult, error) {
	result := CloneResult{
		SourcePath: sourcePath,
		DestPath:   destPath,
	}

	if strings.TrimSpace(sourcePath) == "" {
		return result, fmt.Errorf("clone: sourcePath must not be empty")
	}
	if strings.TrimSpace(destPath) == "" {
		return result, fmt.Errorf("clone: destPath must not be empty")
	}

	version, err := ResolveVersionRef(logical, sourcePath, ref)
	if err != nil {
		return result, fmt.Errorf("clone: resolve version: %w", err)
	}
	result.Version = version

	secret, err := FetchAtRef(logical, sourcePath, version)
	if err != nil {
		return result, fmt.Errorf("clone: fetch source: %w", err)
	}

	if !opts.OverwriteExisting {
		existing, _ := FetchAtRef(logical, destPath, 0)
		if existing != nil && len(existing) > 0 {
			result.Skipped = true
			result.SkipReason = "destination already exists (use --overwrite to replace)"
			return result, nil
		}
	}

	writePath := kvV2WritePath(destPath)
	_, err = logical.Write(writePath, toInterfaceMap(secret))
	if err != nil {
		return result, fmt.Errorf("clone: write destination: %w", err)
	}

	return result, nil
}
