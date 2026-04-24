package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Attach or retrieve named labels on secret versions",
}

var setLabelCmd = &cobra.Command{
	Use:   "set <path> <name> <value> <version>",
	Short: "Attach a label to a specific secret version",
	Args:  cobra.ExactArgs(4),
	RunE:  runSetLabel,
}

var getLabelCmd = &cobra.Command{
	Use:   "get <path> <name>",
	Short: "Retrieve a label attached to a secret",
	Args:  cobra.ExactArgs(2),
	RunE:  runGetLabel,
}

func init() {
	labelCmd.AddCommand(setLabelCmd)
	labelCmd.AddCommand(getLabelCmd)
	rootCmd.AddCommand(labelCmd)
}

func runSetLabel(cmd *cobra.Command, args []string) error {
	path := args[0]
	name := args[1]
	value := args[2]
	version, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[3], err)
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	entry, err := vault.SetLabel(client, path, name, value, version)
	if err != nil {
		return fmt.Errorf("setting label: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "label %q set on %s@v%d\n", entry.Name, entry.Path, entry.Version)
	return nil
}

func runGetLabel(cmd *cobra.Command, args []string) error {
	path := args[0]
	name := args[1]

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	entry, err := vault.GetLabel(client, path, name)
	if err != nil {
		return fmt.Errorf("getting label: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "label: %s\nvalue: %s\nversion: %d\nset_at: %s\n",
		entry.Name, entry.Value, entry.Version, entry.SetAt.Format("2006-01-02T15:04:05Z"))
	return nil
}
