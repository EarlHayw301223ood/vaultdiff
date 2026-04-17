package cmd

import (
	"testing"
)

func TestSetProtectCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "set-protect <path> [reason]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("set-protect command not registered on root")
	}
}

func TestGetProtectCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "get-protect <path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("get-protect command not registered on root")
	}
}

func TestClearProtectCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "clear-protect <path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("clear-protect command not registered on root")
	}
}

func TestSetProtectCmd_RequiresAtLeastOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "set-protect")
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestGetProtectCmd_RequiresExactlyOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "get-protect")
	if err == nil {
		t.Error("expected error when no args provided")
	}
}
