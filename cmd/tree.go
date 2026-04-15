package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var treeMount string

var treeCmd = &cobra.Command{
	Use:   "tree <path-prefix>",
	Short: "List all secret paths under a KV v2 prefix",
	Args:  cobra.ExactArgs(1),
	RunE:  runTree,
}

func init() {
	treeCmd.Flags().StringVar(&treeMount, "mount", "secret", "KV v2 mount name")
	rootCmd.AddCommand(treeCmd)
}

func runTree(cmd *cobra.Command, args []string) error {
	prefix := args[0]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	ctx := cmd.Context()
	paths, err := vault.ListTree(ctx, client, treeMount, prefix)
	if err != nil {
		return fmt.Errorf("listing tree: %w", err)
	}

	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "no secrets found under prefix:", prefix)
		return nil
	}

	for _, p := range paths {
		fmt.Println(p)
	}
	fmt.Fprintf(os.Stderr, "\n%d secret(s) found\n", len(paths))
	return nil
}
