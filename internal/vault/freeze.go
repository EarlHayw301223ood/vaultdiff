package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const freezeMetaKeyPrefix = "vaultdiff/freeze"

// FreezeRecord holds the freeze state for a secret version.
type FreezeRecord struct {
	Frozen    bool      `json:"frozen"`
	Version   int       `json:"version"`
	Reason    string    `json:"reason"`
	FrozenBy  string    `json:"frozen_by"`
	FrozenAt  time.Time `json:"frozen_at"`
}

func freezeMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%s/%s", freezeMetaKeyPrefix, path)
}

// FreezeSecret marks a specific version of a secret as frozen, preventing
// automated overwrites or promotions.
func FreezeSecret(client LogicalClient, path string, version int, reason, actor string) (*FreezeRecord, error) {
	if path == "" {
		return nil, errors.New("freeze: path must not be empty")
	}
	if version <= 0 {
		return nil, errors.New("freeze: version must be a positive integer")
	}
	if reason == "" {
		return nil, errors.New("freeze: reason must not be empty")
	}

	rec := FreezeRecord{
		Frozen:   true,
		Version:  version,
		Reason:   reason,
		FrozenBy: actor,
		FrozenAt: time.Now().UTC(),
	}

	data := map[string]interface{}{
		"frozen":    rec.Frozen,
		"version":   rec.Version,
		"reason":    rec.Reason,
		"frozen_by": rec.FrozenBy,
		"frozen_at": rec.FrozenAt.Format(time.RFC3339),
	}

	metaPath := freezeMetaPath(path)
	_, err := client.Write(kvV2WritePath(metaPath), map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("freeze: write failed: %w", err)
	}
	return &rec, nil
}

// GetFreeze retrieves the freeze record for a secret path.
func GetFreeze(client LogicalClient, path string) (*FreezeRecord, error) {
	if path == "" {
		return nil, errors.New("freeze: path must not be empty")
	}

	metaPath := freezeMetaPath(path)
	secret, err := client.Read(kvV2ReadPath(metaPath))
	if err != nil {
		return nil, fmt.Errorf("freeze: read failed: %w", err)
	}
	if secret == nil {
		return &FreezeRecord{Frozen: false}, nil
	}

	data := extractStringMap(secret)

	var rec FreezeRecord
	rec.Frozen = data["frozen"] == "true"
	rec.Reason = data["reason"]
	rec.FrozenBy = data["frozen_by"]
	if v, ok := data["version"]; ok {
		fmt.Sscanf(v, "%d", &rec.Version)
	}
	if ts, ok := data["frozen_at"]; ok {
		rec.FrozenAt, _ = time.Parse(time.RFC3339, ts)
	}
	return &rec, nil
}

// UnfreezeSecret removes the freeze record for a secret path.
func UnfreezeSecret(client LogicalClient, path string) error {
	if path == "" {
		return errors.New("unfreeze: path must not be empty")
	}
	metaPath := freezeMetaPath(path)
	_, err := client.Write(kvV2WritePath(metaPath), map[string]interface{}{
		"data": map[string]interface{}{"frozen": "false"},
	})
	if err != nil {
		return fmt.Errorf("unfreeze: write failed: %w", err)
	}
	return nil
}
