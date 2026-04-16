package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var copyRef string

var copyCmd = &cobra.Command{
	Use:   "copy <source-path> <dest-path>",
	Short: "Copy a secret from one path to another",
	Args:  cobra.ExactArgs(2),
	RunE:  runCopy,
}

func init() {
	copyCmd.Flags().StringVar(&copyRef, "ref", "latest", "Version ref to copy (number or 'latest')")
	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]
	destPath := args[1]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	result, err := vault.Copy(client, sourcePath, destPath, copyRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Copied %s@v%d → %s (%d keys)\n",
		result.SourcePath, result.Version, result.DestPath, result.Keys)
	return nil
}
