package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage pinned versions of secrets",
}

var setPinCmd = &cobra.Command{
	Use:   "set <path> <version>",
	Short: "Pin a secret path to a specific version",
	Args:  cobra.ExactArgs(2),
	RunE:  runSetPin,
}

var getPinCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the pinned version for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetPin,
}

var clearPinCmd = &cobra.Command{
	Use:   "clear <path>",
	Short: "Clear the pin for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runClearPin,
}

func init() {
	pinCmd.AddCommand(setPinCmd, getPinCmd, clearPinCmd)
	rootCmd.AddCommand(pinCmd)
}

func runSetPin(cmd *cobra.Command, args []string) error {
	path := args[0]
	v, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}
	res, err := vault.SetPin(client.Logical().(*vault.LogicalWrapper).Client(), path, v)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "pinned %s @ v%d\n", res.Path, res.Version)
	return nil
}

func runGetPin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}
	res, err := vault.GetPin(client.Logical().(*vault.LogicalWrapper).Client(), args[0])
	if err != nil {
		return err
	}
	if !res.Pinned {
		fmt.Fprintf(cmd.OutOrStdout(), "%s is not pinned\n", res.Path)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s pinned @ v%d\n", res.Path, res.Version)
	return nil
}

func runClearPin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}
	if err := vault.ClearPin(client.Logical().(*vault.LogicalWrapper).Client(), args[0]); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "pin cleared for %s\n", args[0])
	return nil
}
