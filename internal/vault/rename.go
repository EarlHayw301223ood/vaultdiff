package vault

import (
	"errors"
	"fmt"
)

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	SourcePath string
	DestPath   string
	Version    int
}

// Rename copies a secret from sourcePath to destPath, then deletes all
// versions at sourcePath by writing a tombstone marker.
// It returns an error if either path is empty or if source == dest.
func Rename(client *Client, sourcePath, destPath string) (*RenameResult, error) {
	if sourcePath == "" {
		return nil, errors.New("rename: source path must not be empty")
	}
	if destPath == "" {
		return nil, errors.New("rename: destination path must not be empty")
	}
	if sourcePath == destPath {
		return nil, errors.New("rename: source and destination paths must differ")
	}

	result, err := Copy(client, sourcePath, destPath)
	if err != nil {
		return nil, fmt.Errorf("rename: copy stage failed: %w", err)
	}

	destroyPath := kvV2WritePath(sourcePath)
	_, err = client.Logical().Delete(destroyPath, nil)
	if err != nil {
		return nil, fmt.Errorf("rename: delete stage failed: %w", err)
	}

	return &RenameResult{
		SourcePath: sourcePath,
		DestPath:   destPath,
		Version:    result.Version,
	}, nil
}
