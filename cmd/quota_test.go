package cmd

import (
	"testing"
)

func TestQuotaCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "quota" {
			found = true
			break
		}
	}
	if !found {
		t.Error("quota command not registered on root")
	}
}

func TestQuotaCmd_HasSubcommands(t *testing.T) {
	var names []string
	for _, c := range quotaCmd.Commands() {
		names = append(names, c.Name())
	}
	expected := map[string]bool{"set": false, "get": false}
	for _, n := range names {
		expected[n] = true
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}

func TestSetQuotaCmd_RequiresThreeArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "quota", "set", "secret/app")
	if err == nil {
		t.Error("expected error when fewer than 3 args given")
	}
}

func TestGetQuotaCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "quota", "get")
	if err == nil {
		t.Error("expected error when no args given")
	}
}

func TestQuotaCmd_ShortDescription(t *testing.T) {
	if quotaCmd.Short == "" {
		t.Error("quota command should have a short description")
	}
}

func TestSetQuotaCmd_InvalidMaxWrites(t *testing.T) {
	_, err := executeCommand(rootCmd, "quota", "set", "secret/app", "notanumber", "60")
	if err == nil {
		t.Error("expected error for non-numeric max-writes")
	}
}
