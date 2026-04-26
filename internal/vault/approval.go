package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// ApprovalStatus represents the current state of an approval request.
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

// ApprovalRequest holds metadata for a change-approval gate on a secret path.
type ApprovalRequest struct {
	Path      string         `json:"path"`
	Version   int            `json:"version"`
	Requestor string         `json:"requestor"`
	Reason    string         `json:"reason"`
	Status    ApprovalStatus `json:"status"`
	Reviewer  string         `json:"reviewer,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func approvalMetaPath(path string) string {
	path = strings.Trim(path, "/")
	segments := strings.SplitN(path, "/", 2)
	if len(segments) < 2 {
		return fmt.Sprintf("%s/metadata/_vaultdiff/approval/%s", segments[0], segments[0])
	}
	return fmt.Sprintf("%s/metadata/_vaultdiff/approval/%s", segments[0], segments[1])
}

// RequestApproval creates a pending approval record for the given secret path and version.
func RequestApproval(client *vaultapi.Client, path string, version int, requestor, reason string) (*ApprovalRequest, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if version <= 0 {
		return nil, errors.New("version must be a positive integer")
	}
	if requestor == "" {
		return nil, errors.New("requestor must not be empty")
	}

	now := time.Now().UTC()
	req := &ApprovalRequest{
		Path:      path,
		Version:   version,
		Requestor: requestor,
		Reason:    reason,
		Status:    ApprovalPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	data := map[string]interface{}{
		"path":       req.Path,
		"version":    req.Version,
		"requestor":  req.Requestor,
		"reason":     req.Reason,
		"status":     string(req.Status),
		"created_at": req.CreatedAt.Format(time.RFC3339),
		"updated_at": req.UpdatedAt.Format(time.RFC3339),
	}

	metaPath := approvalMetaPath(path)
	_, err := client.Logical().Write(metaPath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("write approval record: %w", err)
	}
	return req, nil
}

// ReviewApproval updates the status of an existing approval request.
func ReviewApproval(client *vaultapi.Client, path string, reviewer string, approved bool) (*ApprovalRequest, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if reviewer == "" {
		return nil, errors.New("reviewer must not be empty")
	}

	existing, err := GetApproval(client, path)
	if err != nil {
		return nil, err
	}

	if existing.Status != ApprovalPending {
		return nil, fmt.Errorf("approval is already %s", existing.Status)
	}

	status := ApprovalRejected
	if approved {
		status = ApprovalApproved
	}
	existing.Status = status
	existing.Reviewer = reviewer
	existing.UpdatedAt = time.Now().UTC()

	data := map[string]interface{}{
		"path":       existing.Path,
		"version":    existing.Version,
		"requestor":  existing.Requestor,
		"reason":     existing.Reason,
		"status":     string(existing.Status),
		"reviewer":   existing.Reviewer,
		"created_at": existing.CreatedAt.Format(time.RFC3339),
		"updated_at": existing.UpdatedAt.Format(time.RFC3339),
	}

	metaPath := approvalMetaPath(path)
	_, err = client.Logical().Write(metaPath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("update approval record: %w", err)
	}
	return existing, nil
}

// GetApproval retrieves the current approval record for a secret path.
func GetApproval(client *vaultapi.Client, path string) (*ApprovalRequest, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}

	metaPath := approvalMetaPath(path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read approval record: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no approval record found for path %q", path)
	}

	raw, _ := secret.Data["data"].(map[string]interface{})
	if raw == nil {
		raw = secret.Data
	}

	req := &ApprovalRequest{}
	req.Path, _ = raw["path"].(string)
	if v, ok := raw["version"].(float64); ok {
		req.Version = int(v)
	}
	req.Requestor, _ = raw["requestor"].(string)
	req.Reason, _ = raw["reason"].(string)
	if s, ok := raw["status"].(string); ok {
		req.Status = ApprovalStatus(s)
	}
	req.Reviewer, _ = raw["reviewer"].(string)
	if ts, ok := raw["created_at"].(string); ok {
		req.CreatedAt, _ = time.Parse(time.RFC3339, ts)
	}
	if ts, ok := raw["updated_at"].(string); ok {
		req.UpdatedAt, _ = time.Parse(time.RFC3339, ts)
	}
	return req, nil
}
