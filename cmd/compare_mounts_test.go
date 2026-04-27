package cmd

import (
	"testing"
)

func TestCompareMountsCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "compare-mounts <mountA> <mountB>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("compare-mounts command not registered on root")
	}
}

func TestCompareMountsCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "compare-mounts")
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestCompareMountsCmd_RequiresExactlyTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "compare-mounts", "only-one")
	if err == nil {
		t.Fatal("expected error when only one arg provided")
	}
}

func TestCompareMountsCmd_HasPrefixFlag(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "compare-mounts <mountA> <mountB>" {
			if c.Flags().Lookup("prefix") != nil {
				found = true
			}
			break
		}
	}
	if !found {
		t.Fatal("compare-mounts command missing --prefix flag")
	}
}

func TestCompareMountsCmd_ShortDescription(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "compare-mounts <mountA> <mountB>" {
			if c.Short == "" {
				t.Fatal("compare-mounts command has empty short description")
			}
			return
		}
	}
	t.Fatal("compare-mounts command not found")
}
