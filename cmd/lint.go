package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/internal/vault"
)

var lintRules []string

var lintCmd = &cobra.Command{
	Use:   "lint <mount>",
	Short: "Lint secrets under a mount for common issues",
	Args:  cobra.ExactArgs(1),
	RunE:  runLint,
}

func init() {
	lintCmd.Flags().StringSliceVar(&lintRules, "rules", nil, "comma-separated rules to apply (default: all)")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	mount := args[0]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	rules := vault.DefaultRules()
	if len(lintRules) > 0 {
		rules = make([]vault.LintRule, len(lintRules))
		for i, r := range lintRules {
			rules[i] = vault.LintRule(r)
		}
	}

	results, err := vault.LintMount(cmd.Context(), client, mount, rules)
	if err != nil {
		return fmt.Errorf("lint: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ No issues found.")
		return nil
	}

	for _, r := range results {
		for _, issue := range r.Issues {
			fmt.Fprintln(cmd.OutOrStdout(), issue.String())
		}
	}

	summary := vault.LintSummary(results)
	fmt.Fprintf(cmd.OutOrStdout(), "\n%d path(s) with issues\n", len(results))
	for rule, count := range summary {
		fmt.Fprintf(cmd.OutOrStdout(), "  %-30s %d\n", rule, count)
	}

	os.Exit(1)
	return nil
}
