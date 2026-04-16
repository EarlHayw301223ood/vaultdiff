package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/internal/vault"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <path> <version>",
	Short: "Archive a specific version of a secret to a safe location",
	Args:  cobra.ExactArgs(2),
	RunE:  runArchive,
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	result, err := vault.Archive(client.Logical(), path, version)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "archived %s v%d → %s (at %s)\n",
		path, result.Version, result.Path, result.ArchivedAt.Format("2006-01-02T15:04:05Z"))
	return nil
}
