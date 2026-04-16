package cmd

import (
	"testing"
)

func TestSetAliasCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "set-alias <alias> <path>" {
			found = true
		}
	}
	if !found {
		t.Error("set-alias command not registered on root")
	}
}

func TestGetAliasCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "get-alias <alias>" {
			found = true
		}
	}
	if !found {
		t.Error("get-alias command not registered on root")
	}
}

func TestSetAliasCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "set-alias", "only-one")
	if err == nil {
		t.Error("expected error with one arg")
	}
}

func TestGetAliasCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "get-alias")
	if err == nil {
		t.Error("expected error with no args")
	}
}

func TestSetAliasCmd_HasVersionFlag(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "set-alias <alias> <path>" {
			if c.Flags().Lookup("version") == nil {
				t.Error("expected --version flag on set-alias")
			}
			return
		}
	}
	t.Error("set-alias command not found")
}
