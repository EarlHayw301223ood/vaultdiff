package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestWatchCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "watch <path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("watch command not registered on root")
	}
}

func TestWatchCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "watch")
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestWatchCmd_HasIntervalFlag(t *testing.T) {
	f := watchCmd.Flags().Lookup("interval")
	if f == nil {
		t.Fatal("expected --interval flag")
	}
	if f.DefValue != "30" {
		t.Errorf("expected default 30, got %s", f.DefValue)
	}
}

func TestWatchCmd_ShortDescription(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	_, _ = executeCommand(rootCmd, "help", "watch")
	if !strings.Contains(buf.String(), "Poll a Vault KV path") {
		t.Errorf("expected short description in help output, got: %s", buf.String())
	}
}
