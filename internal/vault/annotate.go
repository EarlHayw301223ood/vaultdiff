package vault

import (
	"errors"
	"fmt"
	"time"
)

// Annotation holds a user-defined note attached to a secret version.
type Annotation struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	Note      string    `json:"note"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// annotationMetaPath returns the KV path used to store an annotation.
func annotationMetaPath(path string, version int) string {
	return fmt.Sprintf(".annotations/%s/v%d", path, version)
}

// SetAnnotation writes a note for a given secret path and version.
func SetAnnotation(client *Client, path string, version int, note, author string) (*Annotation, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if version <= 0 {
		return nil, errors.New("version must be a positive integer")
	}
	if note == "" {
		return nil, errors.New("note must not be empty")
	}

	a := &Annotation{
		Path:      path,
		Version:   version,
		Note:      note,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	}

	data := map[string]interface{}{
		"path":       a.Path,
		"version":    a.Version,
		"note":       a.Note,
		"author":     a.Author,
		"created_at": a.CreatedAt.Format(time.RFC3339),
	}

	metaPath := annotationMetaPath(path, version)
	writePath := kvV2WritePath(metaPath)
	_, err := client.Logical().Write(writePath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("write annotation: %w", err)
	}
	return a, nil
}

// GetAnnotation retrieves the annotation for a given secret path and version.
func GetAnnotation(client *Client, path string, version int) (*Annotation, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if version <= 0 {
		return nil, errors.New("version must be a positive integer")
	}

	metaPath := annotationMetaPath(path, version)
	readPath := kvV2ReadPath(metaPath)
	secret, err := client.Logical().Read(readPath)
	if err != nil {
		return nil, fmt.Errorf("read annotation: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no annotation found for %s@v%d", path, version)
	}

	data := extractStringMap(secret)
	createdAt, _ := time.Parse(time.RFC3339, data["created_at"])
	return &Annotation{
		Path:      data["path"],
		Version:   version,
		Note:      data["note"],
		Author:    data["author"],
		CreatedAt: createdAt,
	}, nil
}
