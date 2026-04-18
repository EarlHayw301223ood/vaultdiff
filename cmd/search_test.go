package cmd

import (
	"testing"
)

func TestSearchCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "search" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("search command not registered on root")
	}
}

func TestSearchCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "search")
	if err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestSearchCmd_HasKeyFlag(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"search"})
	if err != nil || cmd == nil {
		t.Fatal("search command not found")
	}
	if cmd.Flags().Lookup("key") == nil {
		t.Fatal("expected --key flag")
	}
}

func TestSearchCmd_HasValueFlag(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"search"})
	if err != nil || cmd == nil {
		t.Fatal("search command not found")
	}
	if cmd.Flags().Lookup("value") == nil {
		t.Fatal("expected --value flag")
	}
}

func TestSearchCmd_ShortDescription(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"search"})
	if err != nil || cmd == nil {
		t.Fatal("search command not found")
	}
	if cmd.Short == "" {
		t.Fatal("expected non-empty short description")
	}
}
