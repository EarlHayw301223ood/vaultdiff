package vault

import (
	"testing"
)

func TestTriggerMetaPath_Format(t *testing.T) {
	got := triggerMetaPath("secret/myapp/db")
	want := "vaultdiff/meta/secret/myapp/db/triggers"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTriggerMetaPath_TrimsSlashes(t *testing.T) {
	got := triggerMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/triggers"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetTrigger_EmptyPath(t *testing.T) {
	c := newStubClient(nil)
	err := SetTrigger(c, "", TriggerConfig{Name: "t1"})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetTrigger_EmptyName(t *testing.T) {
	c := newStubClient(nil)
	err := SetTrigger(c, "secret/app", TriggerConfig{})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSetTrigger_DefaultsConditionToAny(t *testing.T) {
	c := newStubClient(nil)
	err := SetTrigger(c, "secret/app", TriggerConfig{Name: "t1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTrigger_EmptyPath(t *testing.T) {
	c := newStubClient(nil)
	_, err := GetTrigger(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetTrigger_NilSecret(t *testing.T) {
	c := newStubClient(nil)
	cfg, err := GetTrigger(c, "secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Errorf("expected nil config for missing path")
	}
}

func TestEvaluateTrigger_Any(t *testing.T) {
	cfg := TriggerConfig{Condition: "any"}
	if !EvaluateTrigger(cfg, 1) {
		t.Error("expected true for 'any' condition")
	}
}

func TestEvaluateTrigger_VersionGt_Satisfied(t *testing.T) {
	cfg := TriggerConfig{Condition: "version_gt:3"}
	if !EvaluateTrigger(cfg, 4) {
		t.Error("expected true when version > threshold")
	}
}

func TestEvaluateTrigger_VersionGt_NotSatisfied(t *testing.T) {
	cfg := TriggerConfig{Condition: "version_gt:5"}
	if EvaluateTrigger(cfg, 3) {
		t.Error("expected false when version <= threshold")
	}
}

func TestEvaluateTrigger_UnknownCondition(t *testing.T) {
	cfg := TriggerConfig{Condition: "never"}
	if EvaluateTrigger(cfg, 99) {
		t.Error("expected false for unknown condition")
	}
}
