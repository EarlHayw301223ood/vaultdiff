package cmd

import (
	"testing"
)

func TestTriggerCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "trigger" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'trigger' command to be registered on root")
	}
}

func TestTriggerCmd_HasSubcommands(t *testing.T) {
	var names []string
	for _, c := range triggerCmd.Commands() {
		names = append(names, c.Use)
	}
	want := map[string]bool{"set <path> <name>": true, "get <path>": true}
	for _, n := range names {
		delete(want, n)
	}
	if len(want) > 0 {
		t.Errorf("missing subcommands: %v", want)
	}
}

func TestSetTriggerCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "trigger", "set", "secret/app")
	if err == nil {
		t.Error("expected error when only one arg provided")
	}
}

func TestGetTriggerCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "trigger", "get")
	if err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestSetTriggerCmd_HasConditionFlag(t *testing.T) {
	f := setTriggerCmd.Flags().Lookup("condition")
	if f == nil {
		t.Fatal("expected --condition flag on set-trigger command")
	}
	if f.DefValue != "any" {
		t.Errorf("expected default 'any', got %q", f.DefValue)
	}
}

func TestTriggerCmd_ShortDescription(t *testing.T) {
	if triggerCmd.Short == "" {
		t.Error("trigger command should have a short description")
	}
}
