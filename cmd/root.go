package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit HashiCorp Vault secret versions across environments",
	Long: `vaultdiff compares Vault KV secret versions, renders diffs, and writes
audit logs so teams can track what changed, when, and by whom.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().String("vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().Bool("mask", true, "Mask secret values in output")
	rootCmd.PersistentFlags().String("output", "text", "Output format: text or json")
	rootCmd.PersistentFlags().String("audit-log", "", "Path to audit log file (default: stdout)")
}
