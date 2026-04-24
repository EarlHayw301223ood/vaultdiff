package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage and inspect Vault namespaces",
}

var listNamespaceCmd = &cobra.Command{
	Use:   "list [namespace-path]",
	Short: "List child namespaces under a given path",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runListNamespaces,
}

func init() {
	namespaceCmd.AddCommand(listNamespaceCmd)
	rootCmd.AddCommand(namespaceCmd)
}

func runListNamespaces(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("loading vault config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	namespacePath := ""
	if len(args) == 1 {
		namespacePath = vault.NamespacePath(args[0])
	}

	namespaces, err := vault.ListNamespaces(client, namespacePath)
	if err != nil {
		return fmt.Errorf("listing namespaces: %w", err)
	}

	if len(namespaces) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No namespaces found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tFULL PATH")
	for _, ns := range namespaces {
		fmt.Fprintf(w, "%s\t%s\n", ns.Path, ns.FullPath)
	}
	return w.Flush()
}
