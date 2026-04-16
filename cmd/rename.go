package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var renameCmd = &cobra.Command{
	Use:   "rename <source-path> <dest-path>",
	Short: "Rename a secret by copying it to a new path and removing the original",
	Args:  cobra.ExactArgs(2),
	RunE:  runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]
	destPath := args[1]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.Rename(client, sourcePath, destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Renamed %s → %s (version %d written to dest)\n",
		result.SourcePath, result.DestPath, result.Version)
	return nil
}
