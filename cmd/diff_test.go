package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return buf.String(), err
}

func TestDiffCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "diff <path-a> <path-b>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'diff' subcommand to be registered on root")
	}
}

func TestDiffCmd_RequiresTwoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "diff", "only-one-arg")
	if err == nil {
		t.Error("expected error when fewer than 2 args provided")
	}
}

func TestDiffCmd_HasExportFlag(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "diff <path-a> <path-b>" {
			if sub.Flags().Lookup("export") == nil {
				t.Error("expected --export flag on diff command")
			}
			return
		}
	}
	t.Error("diff command not found")
}

func TestRootCmd_HasPersistentFlags(t *testing.T) {
	expected := []string{"vault-addr", "vault-token", "mask", "output", "audit-log"}
	for _, name := range expected {
		if rootCmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("expected persistent flag --%s on root command", name)
		}
	}
}
