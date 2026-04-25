package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// TriggerEvent represents a condition that fires when a secret version changes.
type TriggerEvent struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	FiredAt   time.Time `json:"fired_at"`
	Condition string    `json:"condition"`
}

// TriggerConfig defines a named trigger attached to a secret path.
type TriggerConfig struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Condition string `json:"condition"` // "any", "version_gt:<n>"
}

func triggerMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff/meta/%s/triggers", path)
}

// SetTrigger persists a TriggerConfig as metadata at the canonical meta path.
func SetTrigger(c LogicalClient, path string, cfg TriggerConfig) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if cfg.Name == "" {
		return errors.New("trigger name must not be empty")
	}
	if cfg.Condition == "" {
		cfg.Condition = "any"
	}
	metaPath := triggerMetaPath(path)
	_, err := c.Write(metaPath, map[string]interface{}{
		"name":      cfg.Name,
		"path":      path,
		"condition": cfg.Condition,
	})
	return err
}

// GetTrigger retrieves the TriggerConfig stored for a path.
func GetTrigger(c LogicalClient, path string) (*TriggerConfig, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	metaPath := triggerMetaPath(path)
	secret, err := c.Read(metaPath)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}
	cfg := &TriggerConfig{
		Name:      fmt.Sprintf("%v", secret.Data["name"]),
		Path:      fmt.Sprintf("%v", secret.Data["path"]),
		Condition: fmt.Sprintf("%v", secret.Data["condition"]),
	}
	return cfg, nil
}

// EvaluateTrigger checks whether the trigger condition is satisfied for the given version.
func EvaluateTrigger(cfg TriggerConfig, version int) bool {
	switch {
	case cfg.Condition == "any":
		return true
	case strings.HasPrefix(cfg.Condition, "version_gt:"):
		var threshold int
		fmt.Sscanf(strings.TrimPrefix(cfg.Condition, "version_gt:"), "%d", &threshold)
		return version > threshold
	}
	return false
}
