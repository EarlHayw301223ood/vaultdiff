package cmd

import (
	"testing"
)

func TestBookmarkCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "bookmark" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'bookmark' command to be registered on root")
	}
}

func TestBookmarkCmd_HasSubcommands(t *testing.T) {
	var names []string
	for _, c := range bookmarkCmd.Commands() {
		names = append(names, c.Name())
	}
	expected := []string{"delete", "get", "set"}
	for _, e := range expected {
		found := false
		for _, n := range names {
			if n == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q not found", e)
		}
	}
}

func TestSetBookmarkCmd_RequiresThreeArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "bookmark", "set", "secret/app")
	if err == nil {
		t.Error("expected error with fewer than 3 args")
	}
}

func TestGetBookmarkCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "bookmark", "get", "secret/app")
	if err == nil {
		t.Error("expected error with fewer than 2 args")
	}
}

func TestDeleteBookmarkCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "bookmark", "delete")
	if err == nil {
		t.Error("expected error with no args")
	}
}

func TestBookmarkCmd_ShortDescription(t *testing.T) {
	if bookmarkCmd.Short == "" {
		t.Error("expected non-empty short description for bookmark command")
	}
}
