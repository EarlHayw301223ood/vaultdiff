package cmd

import (
	"testing"
)

func TestBaselineCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "baseline" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'baseline' command to be registered on root")
	}
}

func TestSaveBaselineCmd_RequiresThreeArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "baseline", "save", "secret/myapp")
	if err == nil {
		t.Error("expected error with fewer than 3 args")
	}
}

func TestGetBaselineCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "baseline", "get", "secret/myapp")
	if err == nil {
		t.Error("expected error with fewer than 2 args")
	}
}

func TestBaselineCmd_HasSubcommands(t *testing.T) {
	var found []string
	for _, c := range rootCmd.Commands() {
		if c.Name() == "baseline" {
			for _, sub := range c.Commands() {
				found = append(found, sub.Name())
			}
		}
	}
	expected := map[string]bool{"save": false, "get": false}
	for _, name := range found {
		expected[name] = true
	}
	for sub, ok := range expected {
		if !ok {
			t.Errorf("expected subcommand %q to be registered under baseline", sub)
		}
	}
}

func TestBaselineCmd_ShortDescription(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "baseline" {
			if c.Short == "" {
				t.Error("expected non-empty short description for baseline command")
			}
			return
		}
	}
}
