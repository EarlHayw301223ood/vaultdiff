package vault

import (
	"fmt"
	"sort"
	"time"
)

// ReplayEntry represents a single replayed version of a secret.
type ReplayEntry struct {
	Version  int
	Data     map[string]string
	CreatedAt time.Time
}

// ReplayLog is an ordered sequence of secret versions reconstructed from history.
type ReplayLog struct {
	Path    string
	Entries []ReplayEntry
}

// At returns the ReplayEntry for the given version, or an error if not found.
func (r *ReplayLog) At(version int) (ReplayEntry, error) {
	for _, e := range r.Entries {
		if e.Version == version {
			return e, nil
		}
	}
	return ReplayEntry{}, fmt.Errorf("replay: version %d not found in log for %q", version, r.Path)
}

// Latest returns the highest-version entry in the log.
func (r *ReplayLog) Latest() (ReplayEntry, bool) {
	if len(r.Entries) == 0 {
		return ReplayEntry{}, false
	}
	return r.Entries[len(r.Entries)-1], true
}

// ReplaySecret fetches all active versions of a secret at path and returns
// a ReplayLog ordered by ascending version number.
func ReplaySecret(client LogicalClient, path string, mount string) (*ReplayLog, error) {
	if path == "" {
		return nil, fmt.Errorf("replay: path must not be empty")
	}
	if mount == "" {
		return nil, fmt.Errorf("replay: mount must not be empty")
	}

	versions, err := ListVersions(client, mount, path)
	if err != nil {
		return nil, fmt.Errorf("replay: listing versions for %q: %w", path, err)
	}

	log := &ReplayLog{Path: path}

	for _, v := range versions {
		secret, err := FetchAtRef(client, mount, path, fmt.Sprintf("%d", v.Version))
		if err != nil || secret == nil {
			continue
		}
		data, _ := extractStringMap(secret)
		log.Entries = append(log.Entries, ReplayEntry{
			Version:   v.Version,
			Data:      data,
			CreatedAt: v.CreatedAt,
		})
	}

	sort.Slice(log.Entries, func(i, j int) bool {
		return log.Entries[i].Version < log.Entries[j].Version
	})

	return log, nil
}
