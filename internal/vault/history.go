package vault

import (
	"fmt"
	"sort"
	"time"
)

// VersionHistory holds metadata for all versions of a secret path.
type VersionHistory struct {
	Path     string
	Versions []VersionMeta
}

// VersionMeta holds metadata for a single version.
type VersionMeta struct {
	Version     int
	CreatedTime time.Time
	DeletedTime *time.Time
	Destroyed   bool
}

// GetHistory returns version metadata for all versions of the given path.
func GetHistory(client *Client, path string) (*VersionHistory, error) {
	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}

	versions, err := ListVersions(client, path)
	if err != nil {
		return nil, fmt.Errorf("list versions: %w", err)
	}

	metas := make([]VersionMeta, 0, len(versions))
	for _, v := range versions {
		meta := VersionMeta{
			Version:     v.Version,
			CreatedTime: v.CreatedTime,
			Destroyed:   v.Destroyed,
		}
		if !v.DeletionTime.IsZero() {
			t := v.DeletionTime
			meta.DeletedTime = &t
		}
		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	return &VersionHistory{
		Path:     path,
		Versions: metas,
	}, nil
}

// ActiveVersions returns only versions that are neither deleted nor destroyed.
func (h *VersionHistory) ActiveVersions() []VersionMeta {
	var active []VersionMeta
	for _, v := range h.Versions {
		if v.DeletedTime == nil && !v.Destroyed {
			active = append(active, v)
		}
	}
	return active
}

// Latest returns the highest version number in the history.
func (h *VersionHistory) Latest() (VersionMeta, bool) {
	if len(h.Versions) == 0 {
		return VersionMeta{}, false
	}
	return h.Versions[len(h.Versions)-1], true
}
