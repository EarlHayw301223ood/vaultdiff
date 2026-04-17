package cmd

import (
	"testing"
)

func TestPinCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "pin" {
			return
		}
	}
	t.Error("pin command not registered on root")
}

func TestSetPinCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "pin", "set", "secret/app")
	if err == nil {
		t.Error("expected error with one arg, got nil")
	}
}

func TestGetPinCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "pin", "get")
	if err == nil {
		t.Error("expected error with no args, got nil")
	}
}

func TestClearPinCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "pin", "clear")
	if err == nil {
		t.Error("expected error with no args, got nil")
	}
}

func TestPinCmd_ShortDescription(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "pin" {
			if sub.Short == "" {
				t.Error("pin command missing short description")
			}
			return
		}
	}
}
