package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
)

func newRollbackTestClient(t *testing.T, handler http.Handler) *api.Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	cfg := api.DefaultConfig()
	cfg.Address = ts.URL
	c, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	c.SetToken("test-token")
	return c
}

func TestRollback_EmptyPath(t *testing.T) {
	client := newRollbackTestClient(t, http.NotFoundHandler())
	_, err := Rollback(context.Background(), client, "", 1)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRollback_InvalidVersion(t *testing.T) {
	client := newRollbackTestClient(t, http.NotFoundHandler())
	_, err := Rollback(context.Background(), client, "secret/data/foo", 0)
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestRollback_FetchError(t *testing.T) {
	// Server returns 500 to simulate a fetch failure.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	client := newRollbackTestClient(t, handler)
	_, err := Rollback(context.Background(), client, "secret/myapp", 2)
	if err == nil {
		t.Fatal("expected error when fetch fails")
	}
}

func TestRollbackResult_Fields(t *testing.T) {
	r := &RollbackResult{
		Path:        "secret/myapp",
		FromVersion: 5,
		ToVersion:   3,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.FromVersion != 5 {
		t.Errorf("unexpected FromVersion: %d", r.FromVersion)
	}
	if r.ToVersion != 3 {
		t.Errorf("unexpected ToVersion: %d", r.ToVersion)
	}
}
