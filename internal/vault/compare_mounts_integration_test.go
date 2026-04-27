package vault

import (
	"testing"
)

func TestCompareMounts_IdenticalSnapshots(t *testing.T) {
	data := map[string]map[string]interface{}{
		"secret/data/app/db": {"pass": "abc"},
		"secret/data/app/api": {"key": "xyz"},
	}
	c := newStubClient(data)

	// Both mounts point at the same stub — every path should be in-sync.
	result, err := CompareMounts(c, "secret", "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diverged) != 0 {
		t.Errorf("expected 0 diverged, got %d", len(result.Diverged))
	}
	if len(result.OnlyInA) != 0 {
		t.Errorf("expected 0 only-in-A, got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 0 {
		t.Errorf("expected 0 only-in-B, got %d", len(result.OnlyInB))
	}
}

func TestCompareMounts_SummaryNonEmpty(t *testing.T) {
	r := &MountCompareResult{
		MountA:  "dev",
		MountB:  "prod",
		InSync:  []string{"app/db"},
		Diverged: []string{"app/api"},
	}
	if r.Summary() == "" {
		t.Fatal("Summary must not be empty")
	}
}

func TestCompareMounts_PrefixFiltersResults(t *testing.T) {
	paths := []string{"app/db", "app/api", "infra/net", "infra/dns"}
	got := filterPaths(paths, "infra/")
	if len(got) != 2 {
		t.Fatalf("expected 2 infra paths, got %d", len(got))
	}
	for _, p := range got {
		if p != "infra/net" && p != "infra/dns" {
			t.Errorf("unexpected path %q after prefix filter", p)
		}
	}
}
