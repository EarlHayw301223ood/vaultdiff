package vault

import (
	"testing"
)

func TestApprovalMetaPath_Standard(t *testing.T) {
	got := approvalMetaPath("secret/myapp/db")
	want := "secret/metadata/_vaultdiff/approval/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApprovalMetaPath_TrimsSlashes(t *testing.T) {
	got := approvalMetaPath("/secret/myapp/db/")
	want := "secret/metadata/_vaultdiff/approval/myapp/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApprovalMetaPath_SingleSegment(t *testing.T) {
	got := approvalMetaPath("secret")
	want := "secret/metadata/_vaultdiff/approval/secret"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRequestApproval_EmptyPath(t *testing.T) {
	_, err := RequestApproval(nil, "", 1, "alice", "deploy")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestRequestApproval_ZeroVersion(t *testing.T) {
	_, err := RequestApproval(nil, "secret/foo", 0, "alice", "deploy")
	if err == nil || err.Error() != "version must be a positive integer" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestRequestApproval_NegativeVersion(t *testing.T) {
	_, err := RequestApproval(nil, "secret/foo", -3, "alice", "deploy")
	if err == nil {
		t.Error("expected error for negative version")
	}
}

func TestRequestApproval_EmptyRequestor(t *testing.T) {
	_, err := RequestApproval(nil, "secret/foo", 1, "", "deploy")
	if err == nil || err.Error() != "requestor must not be empty" {
		t.Errorf("expected requestor error, got %v", err)
	}
}

func TestReviewApproval_EmptyPath(t *testing.T) {
	_, err := ReviewApproval(nil, "", "bob", true)
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestReviewApproval_EmptyReviewer(t *testing.T) {
	_, err := ReviewApproval(nil, "secret/foo", "", true)
	if err == nil || err.Error() != "reviewer must not be empty" {
		t.Errorf("expected reviewer error, got %v", err)
	}
}

func TestGetApproval_EmptyPath(t *testing.T) {
	_, err := GetApproval(nil, "")
	if err == nil || err.Error() != "path must not be empty" {
		t.Errorf("expected empty path error, got %v", err)
	}
}

func TestApprovalStatus_Constants(t *testing.T) {
	if ApprovalPending != "pending" {
		t.Errorf("unexpected pending value: %s", ApprovalPending)
	}
	if ApprovalApproved != "approved" {
		t.Errorf("unexpected approved value: %s", ApprovalApproved)
	}
	if ApprovalRejected != "rejected" {
		t.Errorf("unexpected rejected value: %s", ApprovalRejected)
	}
}

func TestApprovalRequest_Fields(t *testing.T) {
	req := &ApprovalRequest{
		Path:      "secret/app",
		Version:   2,
		Requestor: "alice",
		Reason:    "hotfix",
		Status:    ApprovalPending,
	}
	if req.Path != "secret/app" {
		t.Errorf("unexpected path: %s", req.Path)
	}
	if req.Version != 2 {
		t.Errorf("unexpected version: %d", req.Version)
	}
	if req.Status != ApprovalPending {
		t.Errorf("unexpected status: %s", req.Status)
	}
}
