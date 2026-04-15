package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

type mockLockClient struct {
	readSecret *api.Secret
	readErr    error
	writeErr   error
	writtenPath string
	writtenData map[string]interface{}
}

func (m *mockLockClient) ReadWithContext(_ context.Context, path string) (*api.Secret, error) {
	return m.readSecret, m.readErr
}

func (m *mockLockClient) WriteWithContext(_ context.Context, path string, data map[string]interface{}) (*api.Secret, error) {
	m.writtenPath = path
	m.writtenData = data
	return nil, m.writeErr
}

func TestAcquireLock_EmptyPath(t *testing.T) {
	_, err := AcquireLock(context.Background(), &mockLockClient{}, "", LockOptions{Owner: "ci", TTL: time.Minute})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAcquireLock_NoExistingLock(t *testing.T) {
	mc := &mockLockClient{readSecret: nil}
	res, err := AcquireLock(context.Background(), mc, "locks/deploy", LockOptions{Owner: "ci", TTL: time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Owner != "ci" {
		t.Errorf("expected owner ci, got %s", res.Owner)
	}
	if res.Path != "locks/deploy" {
		t.Errorf("unexpected path: %s", res.Path)
	}
}

func TestAcquireLock_ExistingExpiredLock(t *testing.T) {
	expired := time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339)
	mc := &mockLockClient{
		readSecret: &api.Secret{
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"owner":      "old-process",
					"expires_at": expired,
				},
			},
		},
	}
	res, err := AcquireLock(context.Background(), mc, "locks/deploy", LockOptions{Owner: "new", TTL: time.Minute})
	if err != nil {
		t.Fatalf("should acquire over expired lock: %v", err)
	}
	if res.Owner != "new" {
		t.Errorf("expected new owner, got %s", res.Owner)
	}
}

func TestAcquireLock_ActiveLockBlocked(t *testing.T) {
	future := time.Now().UTC().Add(10 * time.Minute).Format(time.RFC3339)
	mc := &mockLockClient{
		readSecret: &api.Secret{
			Data: map[string]interface{}{
				"data": map[string]interface{}{
					"owner":      "blocker",
					"expires_at": future,
				},
			},
		},
	}
	_, err := AcquireLock(context.Background(), mc, "locks/deploy", LockOptions{Owner: "me", TTL: time.Minute})
	if err == nil {
		t.Fatal("expected error when active lock exists")
	}
}

func TestAcquireLock_WriteError(t *testing.T) {
	mc := &mockLockClient{writeErr: fmt.Errorf("permission denied")}
	_, err := AcquireLock(context.Background(), mc, "locks/deploy", LockOptions{Owner: "ci", TTL: time.Minute})
	if err == nil {
		t.Fatal("expected write error to propagate")
	}
}

func TestReleaseLock_EmptyPath(t *testing.T) {
	err := ReleaseLock(context.Background(), &mockLockClient{}, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}
