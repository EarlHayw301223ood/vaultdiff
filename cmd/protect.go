package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/vaultdiff/internal/vault"
)

var protectVersion int

func init() {
	setCmd := &cobra.Command{
		Use:   "set-protect <path> [reason]",
		Short: "Mark a secret path as protected",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runSetProtect,
	}

	getCmd := &cobra.Command{
		Use:   "get-protect <path>",
		Short: "Get protection status of a secret path",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetProtect,
	}

	clearCmd := &cobra.Command{
		Use:   "clear-protect <path>",
		Short: "Remove protection from a secret path",
		Args:  cobra.ExactArgs(1),
		RunE:  runClearProtect,
	}

	rootCmd.AddCommand(setCmd, getCmd, clearCmd)
}

func runSetProtect(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}
	reason := ""
	if len(args) == 2 {
		reason = args[1]
	}
	if err := vault.SetProtection(client.Logical(), args[0], reason); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Protection set on %q\n", args[0])
	return nil
}

func runGetProtect(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}
	protected, reason, err := vault.GetProtection(client.Logical(), args[0])
	if err != nil {
		return err
	}
	if protected {
		fmt.Fprintf(os.Stdout, "%s is protected: %s\n", args[0], reason)
	} else {
		fmt.Fprintf(os.Stdout, "%s is not protected\n", args[0])
	}
	return nil
}

func runClearProtect(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}
	if err := vault.ClearProtection(client.Logical(), args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Protection cleared on %q\n", args[0])
	return nil
}
