package vault

import (
	"errors"
	"fmt"
)

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	SourcePath string
	DestPath   string
	Version    int
	Keys       int
}

// Copy reads the secret at sourcePath (at the given version ref) and writes it
// to destPath. Pass "latest" or "" for ref to copy the current version.
func Copy(client *Client, sourcePath, destPath, ref string) (*CopyResult, error) {
	if sourcePath == "" {
		return nil, errors.New("copy: source path must not be empty")
	}
	if destPath == "" {
		return nil, errors.New("copy: destination path must not be empty")
	}
	if sourcePath == destPath {
		return nil, errors.New("copy: source and destination paths must differ")
	}

	version, err := ResolveVersionRef(client, sourcePath, ref)
	if err != nil {
		return nil, fmt.Errorf("copy: resolve version: %w", err)
	}

	secret, err := FetchAtRef(client, sourcePath, version)
	if err != nil {
		return nil, fmt.Errorf("copy: fetch source: %w", err)
	}

	data := extractStringMap(secret)
	if len(data) == 0 {
		return nil, fmt.Errorf("copy: source secret at %q is empty or unreadable", sourcePath)
	}

	writePath := kvV2WritePath(destPath)
	_, err = client.Logical().Write(writePath, map[string]interface{}{
		"data": toInterfaceMap(data),
	})
	if err != nil {
		return nil, fmt.Errorf("copy: write destination: %w", err)
	}

	return &CopyResult{
		SourcePath: sourcePath,
		DestPath:   destPath,
		Version:    version,
		Keys:       len(data),
	}, nil
}
