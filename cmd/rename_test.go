package cmd

import (
	"testing"
)

func TestRenameCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			found = true
			break
		}
	}
	if !found {
		t.Error("rename command not registered on root")
	}
}

func TestRenameCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "rename", "only-one-arg")
	if err == nil {
		t.Error("expected error when only one arg provided")
	}
}

func TestRenameCmd_ShortDescription(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			if c.Short == "" {
				t.Error("rename command should have a short description")
			}
			return
		}
	}
}

func TestRenameCmd_AcceptsTwoArgs(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			if err := c.Args(c, []string{"src/key", "dst/key"}); err != nil {
				t.Errorf("expected two args to be accepted: %v", err)
			}
			return
		}
	}
}
