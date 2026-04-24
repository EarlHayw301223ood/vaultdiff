package vault

import (
	"testing"
	"time"
)

func TestGarbageCollect_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := GarbageCollect(c, "", GCOptions{KeepLast: 1})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGarbageCollect_KeepLastEnforcedMinimum(t *testing.T) {
	// opts.KeepLast = 0 should be clamped to 1
	c := newStubClient()
	c.stubVersions["secret/data/app"] = []int{1, 2, 3}
	c.stubMeta["secret/data/app"] = map[int]time.Time{
		1: time.Now().Add(-72 * time.Hour),
		2: time.Now().Add(-48 * time.Hour),
		3: time.Now().Add(-1 * time.Hour),
	}
	res, err := GarbageCollect(c, "secret/app", GCOptions{KeepLast: 0, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.RetainedVersions) < 1 {
		t.Error("expected at least 1 retained version")
	}
}

func TestGarbageCollect_DryRunDoesNotDestroy(t *testing.T) {
	c := newStubClient()
	c.stubVersions["secret/data/app"] = []int{1, 2, 3, 4, 5}
	now := time.Now()
	c.stubMeta["secret/data/app"] = map[int]time.Time{
		1: now.Add(-200 * time.Hour),
		2: now.Add(-150 * time.Hour),
		3: now.Add(-100 * time.Hour),
		4: now.Add(-50 * time.Hour),
		5: now.Add(-1 * time.Hour),
	}
	res, err := GarbageCollect(c, "secret/app", GCOptions{
		KeepLast: 2,
		MaxAge:   24 * time.Hour,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if c.destroyCalled {
		t.Error("destroy should not be called in dry-run mode")
	}
}

func TestGarbageCollect_DeletesOldVersions(t *testing.T) {
	c := newStubClient()
	c.stubVersions["secret/data/app"] = []int{1, 2, 3}
	now := time.Now()
	c.stubMeta["secret/data/app"] = map[int]time.Time{
		1: now.Add(-200 * time.Hour),
		2: now.Add(-100 * time.Hour),
		3: now.Add(-1 * time.Hour),
	}
	res, err := GarbageCollect(c, "secret/app", GCOptions{
		KeepLast: 1,
		MaxAge:   24 * time.Hour,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.DeletedVersions) == 0 {
		t.Error("expected at least one deleted version")
	}
	if !c.destroyCalled {
		t.Error("expected destroy to be called")
	}
}

func TestGCResult_Fields(t *testing.T) {
	r := &GCResult{
		Path:             "secret/app",
		DeletedVersions:  []int{1, 2},
		RetainedVersions: []int{3},
		DryRun:           false,
	}
	if r.Path != "secret/app" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if len(r.DeletedVersions) != 2 {
		t.Errorf("expected 2 deleted, got %d", len(r.DeletedVersions))
	}
}
