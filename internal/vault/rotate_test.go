package vault

import (
	"errors"
	"strings"
	"testing"
)

func TestRotate_EmptyPath(t *testing.T) {
	client := newStubClient(nil)
	_, err := Rotate(client, "", RotateOptions{})
	if err == nil || !strings.Contains(err.Error(), "path must not be empty") {
		t.Fatalf("expected empty-path error, got %v", err)
	}
}

func TestRotate_SecretNotFound(t *testing.T) {
	client := newStubClient(nil)
	_, err := Rotate(client, "secret/missing", RotateOptions{})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not-found error, got %v", err)
	}
}

func TestRotate_DryRunDoesNotWrite(t *testing.T) {
	client := newStubClient(map[string]interface{}{
		"data": map[string]interface{}{"API_KEY": "old-value"},
		"metadata": map[string]interface{}{"version": float64(3)},
	})

	res, err := Rotate(client, "secret/svc", RotateOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.WriteCalls() != 0 {
		t.Fatalf("dry-run must not write; got %d write calls", client.WriteCalls())
	}
	if res.OldVersion != 3 {
		t.Errorf("expected OldVersion 3, got %d", res.OldVersion)
	}
	if res.NewVersion != 4 {
		t.Errorf("expected NewVersion 4, got %d", res.NewVersion)
	}
}

func TestRotate_WritesNewVersion(t *testing.T) {
	client := newStubClient(map[string]interface{}{
		"data": map[string]interface{}{"TOKEN": "abc"},
		"metadata": map[string]interface{}{"version": float64(1)},
	})

	res, err := Rotate(client, "secret/app", RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.WriteCalls() != 1 {
		t.Fatalf("expected 1 write call, got %d", client.WriteCalls())
	}
	if res.Path != "secret/app" {
		t.Errorf("unexpected path %q", res.Path)
	}
	if res.RotatedAt.IsZero() {
		t.Error("RotatedAt must not be zero")
	}
}

func TestRotate_TransformApplied(t *testing.T) {
	client := newStubClient(map[string]interface{}{
		"data": map[string]interface{}{"PASS": "plain"},
		"metadata": map[string]interface{}{"version": float64(2)},
	})

	transform := func(_, v string) (string, error) { return strings.ToUpper(v), nil }
	_, err := Rotate(client, "secret/db", RotateOptions{Transform: transform})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	written := client.LastWritten()
	dataMap, _ := written["data"].(map[string]interface{})
	if dataMap["PASS"] != "PLAIN" {
		t.Errorf("expected transformed value PLAIN, got %v", dataMap["PASS"])
	}
}

func TestRotate_TransformError(t *testing.T) {
	client := newStubClient(map[string]interface{}{
		"data": map[string]interface{}{"KEY": "val"},
		"metadata": map[string]interface{}{"version": float64(1)},
	})

	transform := func(_, _ string) (string, error) { return "", errors.New("transform failed") }
	_, err := Rotate(client, "secret/x", RotateOptions{Transform: transform})
	if err == nil || !strings.Contains(err.Error(), "transform failed") {
		t.Fatalf("expected transform error, got %v", err)
	}
}
