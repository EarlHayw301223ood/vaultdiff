package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback <path> <version>",
	Short: "Restore a secret to a previous version",
	Long: `Reads the data at the specified version of a KV v2 secret and
writes it back as a new version, effectively rolling back the secret.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runRollback,
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, args []string) error {
	path := args[0]
	versionStr := args[1]

	targetVersion, err := strconv.Atoi(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version %q: must be a positive integer", versionStr)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.Rollback(cmd.Context(), client, path, targetVersion)
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(),
		"Rolled back %q from version %d → version %d (written as new version)\n",
		result.Path, result.FromVersion, result.ToVersion,
	)
	return nil
}
