package vault

import (
	"errors"
	"fmt"
)

// RestoreResult holds the outcome of a restore operation.
type RestoreResult struct {
	SourcePath  string
	DestPath    string
	Version     int
	RestoredTo  string
}

// Restore copies a specific version of a secret at srcPath into dstPath,
// effectively restoring an older state to a new (or same) location.
func Restore(client *Client, srcPath, dstPath string, version int) (*RestoreResult, error) {
	if srcPath == "" {
		return nil, errors.New("source path must not be empty")
	}
	if dstPath == "" {
		return nil, errors.New("destination path must not be empty")
	}
	if version < 1 {
		return nil, fmt.Errorf("version must be >= 1, got %d", version)
	}

	ref := fmt.Sprintf("%d", version)
	secret, err := FetchAtRef(client, srcPath, ref)
	if err != nil {
		return nil, fmt.Errorf("fetch version %d of %q: %w", version, srcPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %q version %d", srcPath, version)
	}

	data := extractStringMap(secret)
	if len(data) == 0 {
		return nil, fmt.Errorf("secret at %q version %d contains no data", srcPath, version)
	}

	writePath := kvV2WritePath(dstPath)
	_, err = client.Logical().Write(writePath, toInterfaceMap(data))
	if err != nil {
		return nil, fmt.Errorf("write restored secret to %q: %w", dstPath, err)
	}

	return &RestoreResult{
		SourcePath: srcPath,
		DestPath:   dstPath,
		Version:    version,
		RestoredTo: dstPath,
	}, nil
}
