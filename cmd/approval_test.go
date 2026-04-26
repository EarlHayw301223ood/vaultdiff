package cmd

import (
	"testing"
)

func TestApprovalCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "approval" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'approval' command to be registered on root")
	}
}

func TestApprovalCmd_HasSubcommands(t *testing.T) {
	var cmd = approvalCmd
	subs := map[string]bool{}
	for _, c := range cmd.Commands() {
		subs[c.Use] = true
	}
	for _, want := range []string{
		"request <path> <version> <requestor>",
		"review <path> <reviewer> <approve|reject>",
		"get <path>",
	} {
		if !subs[want] {
			t.Errorf("expected subcommand %q to be registered", want)
		}
	}
}

func TestRequestApprovalCmd_RequiresThreeArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "approval", "request", "secret/foo")
	if err == nil {
		t.Error("expected error when fewer than 3 args provided")
	}
}

func TestReviewApprovalCmd_RequiresThreeArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "approval", "review", "secret/foo")
	if err == nil {
		t.Error("expected error when fewer than 3 args provided")
	}
}

func TestGetApprovalCmd_RequiresOneArg(t *testing.T) {
	_, err := executeCommand(rootCmd, "approval", "get")
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestApprovalCmd_ShortDescription(t *testing.T) {
	if approvalCmd.Short == "" {
		t.Error("expected non-empty short description for approval command")
	}
}

func TestRequestApprovalCmd_HasReasonFlag(t *testing.T) {
	f := requestApprovalCmd.Flags().Lookup("reason")
	if f == nil {
		t.Error("expected --reason flag on request subcommand")
	}
}
