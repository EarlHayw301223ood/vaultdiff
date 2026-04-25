package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// LifecycleStage represents a named stage in a secret's lifecycle.
type LifecycleStage string

const (
	StageActive     LifecycleStage = "active"
	StageDeprecated LifecycleStage = "deprecated"
	StageRetired    LifecycleStage = "retired"
	StageReview     LifecycleStage = "review"
)

// LifecycleRecord holds the current lifecycle state of a secret version.
type LifecycleRecord struct {
	Path      string         `json:"path"`
	Version   int            `json:"version"`
	Stage     LifecycleStage `json:"stage"`
	ChangedBy string         `json:"changed_by"`
	ChangedAt time.Time      `json:"changed_at"`
	Reason    string         `json:"reason,omitempty"`
}

func lifecycleMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/lifecycle", path)
}

// SetLifecycle writes a lifecycle record for the given secret version.
func SetLifecycle(client LogicalClient, path string, version int, stage LifecycleStage, changedBy, reason string) (*LifecycleRecord, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if version < 1 {
		return nil, errors.New("version must be >= 1")
	}
	validStages := map[LifecycleStage]bool{
		StageActive: true, StageDeprecated: true, StageRetired: true, StageReview: true,
	}
	if !validStages[stage] {
		return nil, fmt.Errorf("unknown lifecycle stage: %q", stage)
	}

	rec := &LifecycleRecord{
		Path:      path,
		Version:   version,
		Stage:     stage,
		ChangedBy: changedBy,
		ChangedAt: time.Now().UTC(),
		Reason:    reason,
	}

	data := map[string]interface{}{
		"path":       rec.Path,
		"version":    rec.Version,
		"stage":      string(rec.Stage),
		"changed_by": rec.ChangedBy,
		"changed_at": rec.ChangedAt.Format(time.RFC3339),
		"reason":     rec.Reason,
	}

	metaPath := lifecycleMetaPath(path)
	_, err := client.Write(metaPath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("write lifecycle: %w", err)
	}
	return rec, nil
}

// GetLifecycle retrieves the lifecycle record for a secret path.
func GetLifecycle(client LogicalClient, path string) (*LifecycleRecord, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	metaPath := lifecycleMetaPath(path)
	secret, err := client.Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read lifecycle: %w", err)
	}
	if secret == nil {
		return nil, nil
	}

	data := extractStringMap(secret)
	ts, _ := time.Parse(time.RFC3339, data["changed_at"])
	ver := 0
	fmt.Sscanf(data["version"], "%d", &ver)

	return &LifecycleRecord{
		Path:      data["path"],
		Version:   ver,
		Stage:     LifecycleStage(data["stage"]),
		ChangedBy: data["changed_by"],
		ChangedAt: ts,
		Reason:    data["reason"],
	}, nil
}
