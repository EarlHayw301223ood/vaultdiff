package cmd

import (
	"testing"
)

func TestNotifyCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "notify" {
			found = true
			break
		}
	}
	if !found {
		t.Error("notify command not registered on root")
	}
}

func TestNotifyCmd_HasSubcommands(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Name() == "notify" {
			for _, sub := range c.Commands() {
				if sub.Name() == "fire" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Error("notify fire subcommand not found")
	}
}

func TestNotifyFireCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "notify", "fire")
	if err == nil {
		t.Error("expected error when no path argument provided")
	}
}

func TestNotifyFireCmd_HasTargetFlag(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "notify" {
			for _, sub := range c.Commands() {
				if sub.Name() == "fire" {
					if sub.Flags().Lookup("target") == nil {
						t.Error("fire subcommand missing --target flag")
					}
					return
				}
			}
		}
	}
}

func TestNotifyFireCmd_HasTemplateFlag(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "notify" {
			for _, sub := range c.Commands() {
				if sub.Name() == "fire" {
					if sub.Flags().Lookup("template") == nil {
						t.Error("fire subcommand missing --template flag")
					}
					return
				}
			}
		}
	}
}

func TestNotifyCmd_ShortDescription(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "notify" {
			if c.Short == "" {
				t.Error("notify command should have a short description")
			}
			return
		}
	}
}
