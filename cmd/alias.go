package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func init() {
	setAliasCmd := &cobra.Command{
		Use:   "set-alias <alias> <path>",
		Short: "Create or update a named alias pointing to a vault path",
		Args:  cobra.ExactArgs(2),
		RunE:  runSetAlias,
	}
	setAliasCmd.Flags().Int("version", 0, "Pin alias to a specific secret version")

	getAliasCmd := &cobra.Command{
		Use:   "get-alias <alias>",
		Short: "Resolve a named alias to its vault path",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetAlias,
	}

	rootCmd.AddCommand(setAliasCmd)
	rootCmd.AddCommand(getAliasCmd)
}

func runSetAlias(cmd *cobra.Command, args []string) error {
	alias := args[0]
	path := args[1]
	version, _ := cmd.Flags().GetInt("version")

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	if err := vault.SetAlias(client.Logical(), alias, path, version); err != nil {
		return fmt.Errorf("set alias: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "alias %q -> %s", alias, path)
	if version > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), " (v%s)", strconv.Itoa(version))
	}
	fmt.Fprintln(cmd.OutOrStdout())
	return nil
}

func runGetAlias(cmd *cobra.Command, args []string) error {
	alias := args[0]
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	entry, err := vault.GetAlias(client.Logical(), alias)
	if err != nil {
		return fmt.Errorf("get alias: %w", err)
	}
	if entry.Version > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%s@v%d\n", entry.Path, entry.Version)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), entry.Path)
	}
	return nil
}
