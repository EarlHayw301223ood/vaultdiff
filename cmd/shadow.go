package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/vaultdiff/internal/vault"
)

var shadowDryRun bool

func init() {
	shadowCmd := &cobra.Command{
		Use:   "shadow <source-path> <shadow-mount> <version>",
		Short: "Copy a specific secret version into a shadow mount for safe inspection",
		Args:  cobra.ExactArgs(3),
		RunE:  runShadow,
	}

	shadowCmd.Flags().BoolVar(&shadowDryRun, "dry-run", false, "Print what would be written without making changes")

	rootCmd.AddCommand(shadowCmd)
}

func runShadow(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]
	shadowMount := args[1]
	versionStr := args[2]

	version, err := strconv.Atoi(versionStr)
	if err != nil || version < 1 {
		return fmt.Errorf("version must be a positive integer, got %q", versionStr)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.Shadow(client, sourcePath, shadowMount, version, shadowDryRun)
	if err != nil {
		return err
	}

	if result.DryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] would shadow %s@v%d → %s\n",
			result.SourcePath, result.Version, result.ShadowPath)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "shadowed %s@v%d → %s\n",
		result.SourcePath, result.Version, result.ShadowPath)
	return nil
}
