package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var (
	mergeStrategy string
	mergeDryRun   bool
)

var mergeCmd = &cobra.Command{
	Use:   "merge <source> <dest>",
	Short: "Merge secrets from source path into destination path",
	Args:  cobra.ExactArgs(2),
	RunE:  runMerge,
}

func init() {
	mergeCmd.Flags().StringVar(&mergeStrategy, "strategy", "ours",
		"Conflict resolution strategy: ours, theirs, union")
	mergeCmd.Flags().BoolVar(&mergeDryRun, "dry-run", false,
		"Preview merge without writing changes")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	strategy := vault.MergeStrategy(strings.ToLower(mergeStrategy))
	switch strategy {
	case vault.MergeStrategyOurs, vault.MergeStrategyTheirs, vault.MergeStrategyUnion:
		// valid
	default:
		return fmt.Errorf("unknown strategy %q: choose ours, theirs, or union", mergeStrategy)
	}

	res, err := vault.Merge(client.Logical(), args[0], args[1], vault.MergeOptions{
		Strategy: strategy,
		DryRun:   mergeDryRun,
	})
	if err != nil {
		return err
	}

	if mergeDryRun {
		fmt.Fprintln(os.Stdout, "[dry-run] merge preview:")
	}

	if len(res.Added) > 0 {
		fmt.Fprintf(os.Stdout, "  added (%d):      %s\n", len(res.Added), strings.Join(res.Added, ", "))
	}
	if len(res.Overwritten) > 0 {
		fmt.Fprintf(os.Stdout, "  overwritten (%d): %s\n", len(res.Overwritten), strings.Join(res.Overwritten, ", "))
	}
	if len(res.Conflicts) > 0 {
		fmt.Fprintf(os.Stdout, "  conflicts (%d):  %s\n", len(res.Conflicts), strings.Join(res.Conflicts, ", "))
	}

	if !mergeDryRun {
		fmt.Fprintf(os.Stdout, "merged %s → %s (strategy: %s)\n", args[0], args[1], strategy)
	}
	return nil
}
