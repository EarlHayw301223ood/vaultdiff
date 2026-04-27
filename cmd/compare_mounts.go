package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var compareMountsPrefix string

func init() {
	compareMountsCmd := &cobra.Command{
		Use:   "compare-mounts <mountA> <mountB>",
		Short: "Compare all secrets across two KV v2 mounts",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompareMounts,
	}
	compareMountsCmd.Flags().StringVar(&compareMountsPrefix, "prefix", "", "Only compare paths that start with this prefix")
	rootCmd.AddCommand(compareMountsCmd)
}

func runCompareMounts(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	result, err := vault.CompareMounts(client.Logical(), args[0], args[1], compareMountsPrefix)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	printSection := func(label string, paths []string) {
		for _, p := range paths {
			fmt.Fprintf(w, "%s\t%s\n", label, p)
		}
	}

	printSection("ONLY-IN-"+args[0], result.OnlyInA)
	printSection("ONLY-IN-"+args[1], result.OnlyInB)
	printSection("DIVERGED", result.Diverged)
	printSection("IN-SYNC", result.InSync)

	fmt.Fprintln(cmd.OutOrStdout(), result.Summary())
	return nil
}
