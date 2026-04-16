package vault

import (
	"errors"
	"fmt"
	"time"
)

// ArchiveResult holds metadata about an archived secret path.
type ArchiveResult struct {
	Path      string
	Version   int
	ArchivedAt time.Time
}

// Archive writes a snapshot of the secret at path to a designated archive
// location (e.g. <mount>/archive/<path>) so it can be retrieved later.
func Archive(client LogicalClient, path string, version int) (*ArchiveResult, error) {
	if path == "" {
		return nil, errors.New("archive: path must not be empty")
	}
	if version < 1 {
		return nil, errors.New("archive: version must be >= 1")
	}

	secret, err := FetchAtRef(client, path, version)
	if err != nil {
		return nil, fmt.Errorf("archive: fetch failed: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("archive: no secret found at %s version %d", path, version)
	}

	archivePath := archiveDestPath(path, version)
	data := toInterfaceMap(secret)

	_, err = client.Write(kvV2WritePath(archivePath), map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("archive: write failed: %w", err)
	}

	return &ArchiveResult{
		Path:       archivePath,
		Version:    version,
		ArchivedAt: time.Now().UTC(),
	}, nil
}

func archiveDestPath(path string, version int) string {
	return fmt.Sprintf("archive/%s/v%d", path, version)
}
