package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const immutableMetaKey = "immutable"

type ImmutableRecord struct {
	Enabled   bool      `json:"enabled"`
	SetBy     string    `json:"set_by"`
	SetAt     time.Time `json:"set_at"`
	Reason    string    `json:"reason,omitempty"`
}

func immutableMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/immutable", path)
}

// SetImmutable marks a secret path as immutable, preventing future writes.
func SetImmutable(client LogicalClient, path, setBy, reason string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("path must not be empty")
	}
	if strings.TrimSpace(setBy) == "" {
		return errors.New("setBy must not be empty")
	}

	rec := ImmutableRecord{
		Enabled: true,
		SetBy:   setBy,
		SetAt:   time.Now().UTC(),
		Reason:  reason,
	}

	metaPath := immutableMetaPath(path)
	_, err := client.Write(metaPath, map[string]interface{}{
		"data": map[string]interface{}{
			"enabled": rec.Enabled,
			"set_by":  rec.SetBy,
			"set_at":  rec.SetAt.Format(time.RFC3339),
			"reason":  rec.Reason,
		},
	})
	return err
}

// GetImmutable retrieves the immutability record for a secret path.
func GetImmutable(client LogicalClient, path string) (*ImmutableRecord, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("path must not be empty")
	}

	metaPath := immutableMetaPath(path)
	secret, err := client.Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read immutable meta: %w", err)
	}
	if secret == nil {
		return &ImmutableRecord{Enabled: false}, nil
	}

	data := extractStringMap(secret)
	rec := &ImmutableRecord{
		Enabled: data["enabled"] == "true",
		SetBy:   data["set_by"],
		Reason:  data["reason"],
	}
	if raw := data["set_at"]; raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			rec.SetAt = t
		}
	}
	return rec, nil
}

// ClearImmutable removes the immutability flag from a secret path.
func ClearImmutable(client LogicalClient, path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("path must not be empty")
	}
	metaPath := immutableMetaPath(path)
	_, err := client.Delete(metaPath)
	return err
}
