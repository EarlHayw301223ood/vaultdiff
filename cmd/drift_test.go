package cmd

import (
	"testing"
)

func TestDriftCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "drift <mountA> <mountB>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("drift command not registered on root")
	}
}

func TestDriftCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "drift")
	if err == nil {
		t.Error("expected error with no args")
	}

	_, err = executeCommand(rootCmd, "drift", "only-one")
	if err == nil {
		t.Error("expected error with one arg")
	}
}

func TestDriftCmd_HasPrefixFlag(t *testing.T) {
	flag := driftCmd.Flags().Lookup("prefix")
	if flag == nil {
		t.Fatal("expected --prefix flag to exist")
	}
	if flag.DefValue != "" {
		t.Errorf("expected empty default, got %q", flag.DefValue)
	}
}

func TestDriftCmd_ShortDescription(t *testing.T) {
	if driftCmd.Short == "" {
		t.Error("expected non-empty short description")
	}
}
