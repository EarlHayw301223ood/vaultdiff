package vault

import (
	"testing"
)

func TestBaselineMetaPath_Format(t *testing.T) {
	got := baselineMetaPath("secret/myapp", "v1")
	want := "vaultdiff/baselines/secret/myapp/v1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBaselineMetaPath_TrimsSlashes(t *testing.T) {
	got := baselineMetaPath("/secret/myapp/", "prod")
	want := "vaultdiff/baselines/secret/myapp/prod"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSaveBaseline_EmptyPath(t *testing.T) {
	_, err := SaveBaseline(nil, "", "v1", 1)
	if err == nil || err.Error() != "path is required" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestSaveBaseline_EmptyName(t *testing.T) {
	_, err := SaveBaseline(nil, "secret/myapp", "", 1)
	if err == nil || err.Error() != "baseline name is required" {
		t.Errorf("expected name error, got %v", err)
	}
}

func TestSaveBaseline_InvalidVersion(t *testing.T) {
	_, err := SaveBaseline(nil, "secret/myapp", "v1", 0)
	if err == nil || err.Error() != "version must be >= 1" {
		t.Errorf("expected version error, got %v", err)
	}
}

func TestGetBaseline_EmptyPath(t *testing.T) {
	_, err := GetBaseline(nil, "", "v1")
	if err == nil || err.Error() != "path is required" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestGetBaseline_EmptyName(t *testing.T) {
	_, err := GetBaseline(nil, "secret/myapp", "")
	if err == nil || err.Error() != "baseline name is required" {
		t.Errorf("expected name error, got %v", err)
	}
}

func TestBaseline_Fields(t *testing.T) {
	bl := &Baseline{
		Name:    "prod-snapshot",
		Path:    "secret/myapp",
		Version: 3,
		Data:    map[string]string{"KEY": "value"},
	}
	if bl.Name != "prod-snapshot" {
		t.Errorf("unexpected Name: %s", bl.Name)
	}
	if bl.Version != 3 {
		t.Errorf("unexpected Version: %d", bl.Version)
	}
	if bl.Data["KEY"] != "value" {
		t.Errorf("unexpected Data value")
	}
}
