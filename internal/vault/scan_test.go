package vault

import (
	"testing"
)

func TestFilterScanResults_KeepsMatches(t *testing.T) {
	input := []ScanResult{
		{Path: "secret/a", Mount: "secret", Versions: 3},
		{Path: "secret/b", Mount: "secret", Versions: 1},
		{Path: "secret/c", Mount: "secret", Versions: 5},
	}

	got := FilterScanResults(input, func(r ScanResult) bool {
		return r.Versions >= 3
	})

	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].Path != "secret/a" {
		t.Errorf("expected secret/a, got %s", got[0].Path)
	}
	if got[1].Path != "secret/c" {
		t.Errorf("expected secret/c, got %s", got[1].Path)
	}
}

func TestFilterScanResults_EmptyInput(t *testing.T) {
	got := FilterScanResults([]ScanResult{}, func(r ScanResult) bool { return true })
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(got))
	}
}

func TestFilterScanResults_NoneMatch(t *testing.T) {
	input := []ScanResult{
		{Path: "secret/x", Mount: "secret", Versions: 1},
	}
	got := FilterScanResults(input, func(r ScanResult) bool { return false })
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestFilterScanResults_AllMatch(t *testing.T) {
	input := []ScanResult{
		{Path: "a", Versions: 2},
		{Path: "b", Versions: 4},
	}
	got := FilterScanResults(input, func(r ScanResult) bool { return r.Versions > 0 })
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestScanMount_EmptyMount(t *testing.T) {
	_, err := ScanMount(nil, nil, "")
	if err == nil {
		t.Fatal("expected error for empty mount, got nil")
	}
}
