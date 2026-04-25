package vault

import (
	"testing"
)

func TestShadowMetaPath_Format(t *testing.T) {
	got := shadowMetaPath("secret/data/myapp/config")
	want := "vaultdiff/shadow/secret/data/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestShadowMetaPath_TrimsSlashes(t *testing.T) {
	got := shadowMetaPath("/secret/data/myapp/")
	want := "vaultdiff/shadow/secret/data/myapp"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestShadow_EmptySourcePath(t *testing.T) {
	_, err := Shadow(nil, "", "shadow", 1, false)
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestShadow_EmptyShadowMount(t *testing.T) {
	_, err := Shadow(nil, "secret/data/app", "", 1, false)
	if err == nil {
		t.Fatal("expected error for empty shadow mount")
	}
}

func TestShadow_InvalidVersion(t *testing.T) {
	_, err := Shadow(nil, "secret/data/app", "shadow", 0, false)
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestShadow_NegativeVersion(t *testing.T) {
	_, err := Shadow(nil, "secret/data/app", "shadow", -3, false)
	if err == nil {
		t.Fatal("expected error for negative version")
	}
}

func TestShadowResult_Fields(t *testing.T) {
	r := &ShadowResult{
		SourcePath: "secret/data/app",
		ShadowPath: "shadow/data/app",
		Version:    2,
		DryRun:     true,
	}
	if r.SourcePath != "secret/data/app" {
		t.Errorf("unexpected SourcePath: %q", r.SourcePath)
	}
	if r.Version != 2 {
		t.Errorf("unexpected Version: %d", r.Version)
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestFirstSegment_Standard(t *testing.T) {
	got := firstSegment("secret/data/app/config")
	if got != "secret" {
		t.Errorf("got %q, want %q", got, "secret")
	}
}

func TestFirstSegment_SingleSegment(t *testing.T) {
	got := firstSegment("secret")
	if got != "secret" {
		t.Errorf("got %q, want %q", got, "secret")
	}
}

func TestFirstSegment_LeadingSlash(t *testing.T) {
	got := firstSegment("/kv/data/foo")
	if got != "kv" {
		t.Errorf("got %q, want %q", got, "kv")
	}
}
