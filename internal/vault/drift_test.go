package vault

import (
	"testing"
)

func makeDriftSnapshot(data map[string]map[string]string) *Snapshot {
	return makeSnapshot(data)
}

func TestDetectDrift_AllInSync(t *testing.T) {
	data := map[string]map[string]string{
		"app/config": {"key": "value"},
	}
	snapA := makeDriftSnapshot(data)
	snapB := makeDriftSnapshot(data)

	report, err := DetectDrift("secret/", "secret-prod/", snapA, snapB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(report.Results))
	}
	if !report.Results[0].InSync {
		t.Errorf("expected InSync=true")
	}
}

func TestDetectDrift_Diverged(t *testing.T) {
	snapA := makeDriftSnapshot(map[string]map[string]string{
		"app/config": {"key": "old"},
	})
	snapB := makeDriftSnapshot(map[string]map[string]string{
		"app/config": {"key": "new"},
	})

	report, err := DetectDrift("secret/", "secret-prod/", snapA, snapB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Results[0].Diverged {
		t.Errorf("expected Diverged=true")
	}
}

func TestDetectDrift_OnlyInA(t *testing.T) {
	snapA := makeDriftSnapshot(map[string]map[string]string{
		"app/only-a": {"x": "1"},
	})
	snapB := makeDriftSnapshot(map[string]map[string]string{})

	report, err := DetectDrift("a/", "b/", snapA, snapB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Results[0].OnlyInA {
		t.Errorf("expected OnlyInA=true")
	}
}

func TestDetectDrift_OnlyInB(t *testing.T) {
	snapA := makeDriftSnapshot(map[string]map[string]string{})
	snapB := makeDriftSnapshot(map[string]map[string]string{
		"app/only-b": {"x": "1"},
	})

	report, err := DetectDrift("a/", "b/", snapA, snapB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.Results[0].OnlyInB {
		t.Errorf("expected OnlyInB=true")
	}
}

func TestDetectDrift_NilSnapshot(t *testing.T) {
	_, err := DetectDrift("a/", "b/", nil, nil)
	if err == nil {
		t.Error("expected error for nil snapshots")
	}
}

func TestDriftReport_Summary(t *testing.T) {
	report := &DriftReport{
		Results: []DriftResult{
			{InSync: true},
			{Diverged: true},
			{OnlyInA: true},
			{OnlyInB: true},
			{InSync: true},
		},
	}
	summary := report.Summary()
	if summary["in_sync"] != 2 {
		t.Errorf("expected 2 in_sync, got %d", summary["in_sync"])
	}
	if summary["diverged"] != 1 {
		t.Errorf("expected 1 diverged, got %d", summary["diverged"])
	}
	if summary["only_in_a"] != 1 {
		t.Errorf("expected 1 only_in_a, got %d", summary["only_in_a"])
	}
	if summary["only_in_b"] != 1 {
		t.Errorf("expected 1 only_in_b, got %d", summary["only_in_b"])
	}
}
