package vault

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// dependencyMetaPath returns the metadata path for dependency records.
func dependencyMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/dependencies", path)
}

// Dependency represents a declared dependency between two secret paths.
type Dependency struct {
	SourcePath string    `json:"source_path"`
	TargetPath string    `json:"target_path"`
	Label      string    `json:"label,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// DependencyList holds all dependencies registered for a path.
type DependencyList struct {
	Dependencies []Dependency `json:"dependencies"`
}

// AddDependency records that sourcePath depends on targetPath.
func AddDependency(client LogicalClient, sourcePath, targetPath, label string) error {
	if sourcePath == "" {
		return fmt.Errorf("source path must not be empty")
	}
	if targetPath == "" {
		return fmt.Errorf("target path must not be empty")
	}
	if sourcePath == targetPath {
		return fmt.Errorf("source and target paths must differ")
	}

	metaPath := dependencyMetaPath(sourcePath)
	existing, err := getDependencies(client, metaPath)
	if err != nil {
		return err
	}

	for _, d := range existing.Dependencies {
		if d.TargetPath == targetPath {
			return fmt.Errorf("dependency on %q already exists", targetPath)
		}
	}

	existing.Dependencies = append(existing.Dependencies, Dependency{
		SourcePath: sourcePath,
		TargetPath: targetPath,
		Label:      label,
		CreatedAt:  time.Now().UTC(),
	})

	_, err = client.Write(kvV2WritePath(metaPath), map[string]interface{}{
		"data": map[string]interface{}{
			"dependencies": marshalDependencies(existing.Dependencies),
		},
	})
	return err
}

// GetDependencies returns all dependencies declared for the given source path.
func GetDependencies(client LogicalClient, sourcePath string) (*DependencyList, error) {
	if sourcePath == "" {
		return nil, fmt.Errorf("source path must not be empty")
	}
	return getDependencies(client, dependencyMetaPath(sourcePath))
}

func getDependencies(client LogicalClient, metaPath string) (*DependencyList, error) {
	secret, err := client.Read(kvV2ReadPath(metaPath))
	if err != nil || secret == nil {
		return &DependencyList{}, nil
	}
	data := extractStringMap(secret)
	list := &DependencyList{}
	raw, ok := data["dependencies"]
	if !ok {
		return list, nil
	}
	for _, entry := range strings.Split(raw, "|") {
		parts := strings.SplitN(entry, ":", 3)
		if len(parts) < 2 {
			continue
		}
		d := Dependency{TargetPath: parts[0], Label: parts[1]}
		if len(parts) == 3 {
			d.CreatedAt, _ = time.Parse(time.RFC3339, parts[2])
		}
		list.Dependencies = append(list.Dependencies, d)
	}
	sort.Slice(list.Dependencies, func(i, j int) bool {
		return list.Dependencies[i].TargetPath < list.Dependencies[j].TargetPath
	})
	return list, nil
}

func marshalDependencies(deps []Dependency) string {
	parts := make([]string, 0, len(deps))
	for _, d := range deps {
		parts = append(parts, fmt.Sprintf("%s:%s:%s", d.TargetPath, d.Label, d.CreatedAt.Format(time.RFC3339)))
	}
	return strings.Join(parts, "|")
}
