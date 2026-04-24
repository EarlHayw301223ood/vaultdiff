package vault

import (
	"fmt"
	"time"
)

type versionMeta struct {
	CreatedTime time.Time
	Destroyed   bool
}

// fetchVersionMeta reads the creation time of a specific version from KV v2
// metadata. It is a thin wrapper around the logical client.
func fetchVersionMeta(l LogicalClient, path string, version int) (versionMeta, error) {
	metaPath := fmt.Sprintf("%s/metadata/%s", mountPrefix(path), stripMount(path))
	secret, err := l.Read(metaPath)
	if err != nil {
		return versionMeta{}, err
	}
	if secret == nil || secret.Data == nil {
		return versionMeta{}, fmt.Errorf("no metadata at %q", path)
	}
	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok {
		return versionMeta{}, fmt.Errorf("unexpected metadata shape at %q", path)
	}
	key := fmt.Sprintf("%d", version)
	entry, ok := versions[key].(map[string]interface{})
	if !ok {
		return versionMeta{}, fmt.Errorf("version %d not found in metadata for %q", version, path)
	}
	var meta versionMeta
	if ct, ok := entry["created_time"].(string); ok {
		meta.CreatedTime, _ = time.Parse(time.RFC3339Nano, ct)
	}
	if d, ok := entry["destroyed"].(bool); ok {
		meta.Destroyed = d
	}
	return meta, nil
}

// destroyVersions permanently destroys the listed versions at path.
func destroyVersions(l LogicalClient, path string, versions []int) error {
	data := make([]interface{}, len(versions))
	for i, v := range versions {
		data[i] = v
	}
	destroyPath := fmt.Sprintf("%s/destroy/%s", mountPrefix(path), stripMount(path))
	_, err := l.Write(destroyPath, map[string]interface{}{"versions": data})
	return err
}

// stripMount removes the first path segment (the mount) from path.
func stripMount(path string) string {
	for i, c := range path {
		if c == '/' && i > 0 {
			return path[i+1:]
		}
	}
	return path
}
