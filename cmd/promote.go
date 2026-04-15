package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var promoteRef string

var promoteCmd = &cobra.Command{
	Use:   "promote <source-path> <dest-path>",
	Short: "Promote a secret version from one path to another",
	Long: `Copies a secret (at an optional version ref) from source-path to dest-path,
writing it as a new version at the destination. Useful for promoting secrets
across environments (e.g. staging → production).`,
	Args: cobra.ExactArgs(2),
	RunE: runPromote,
}

func init() {
	promoteCmd.Flags().StringVar(&promoteRef, "ref", "latest", "Source version ref (latest or version number)")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
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

	ctx := cmd.Context()

	result, err := vault.Promote(ctx, client, sourcePath, destPath, promoteRef)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Promoted %s@v%d → %s@v%d (%d keys)\n",
		result.SourcePath, result.SourceVersion,
		result.DestPath, result.DestVersion,
		result.Keys,
	)
	return nil
}
