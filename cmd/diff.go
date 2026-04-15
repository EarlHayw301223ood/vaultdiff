package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/audit"
	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <path-a> <path-b>",
	Short: "Diff two Vault secret paths (optionally at specific versions)",
	Example: `  vaultdiff diff secret/data/prod/app secret/data/staging/app
  vaultdiff diff secret/data/app@3 secret/data/app@5
  vaultdiff diff secret/data/app@latest secret/data/app@2 --output json`,
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	diffCmd.Flags().String("export", "", "Write report to this file path")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}

	if addr, _ := cmd.Flags().GetString("vault-addr"); addr != "" {
		cfg.Address = addr
	}
	if token, _ := cmd.Flags().GetString("vault-token"); token != "" {
		cfg.Token = token
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	refA, refB := args[0], args[1]
	pair, err := vault.FetchPairAtRefs(cmd.Context(), client, refA, refB)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	changes := diff.Compare(pair.A.Data, pair.B.Data)
	mask, _ := cmd.Flags().GetBool("mask")
	fmt.Print(diff.Render(changes, mask))

	report := diff.NewReport(refA, refB, changes, pair.A.Metadata, pair.B.Metadata)

	if exportPath, _ := cmd.Flags().GetString("export"); exportPath != "" {
		format, _ := cmd.Flags().GetString("output")
		fmt := diff.ParseFormat(format)
		if err := diff.Export(report, fmt, exportPath); err != nil {
			return fmt.Errorf("export: %w", err)
		}
	}

	logger, err := resolveAuditLogger(cmd)
	if err != nil {
		return err
	}
	return logger.Record(refA, refB, changes)
}

func resolveAuditLogger(cmd *cobra.Command) (audit.Logger, error) {
	path, _ := cmd.Flags().GetString("audit-log")
	if path == "" {
		return audit.NewLogger(os.Stdout), nil
	}
	return audit.NewFileLogger(path)
}
