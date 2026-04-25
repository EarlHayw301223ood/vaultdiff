package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var dependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Manage secret dependencies",
}

var addDependencyCmd = &cobra.Command{
	Use:   "add <source-path> <target-path>",
	Short: "Declare that source-path depends on target-path",
	Args:  cobra.ExactArgs(2),
	RunE:  runAddDependency,
}

var listDependencyCmd = &cobra.Command{
	Use:   "list <path>",
	Short: "List all dependencies declared for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runListDependencies,
}

var depLabel string

func init() {
	addDependencyCmd.Flags().StringVar(&depLabel, "label", "", "Optional label describing the dependency relationship")
	dependencyCmd.AddCommand(addDependencyCmd)
	dependencyCmd.AddCommand(listDependencyCmd)
	rootCmd.AddCommand(dependencyCmd)
}

func runAddDependency(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}
	if err := vault.AddDependency(client.Logical(), args[0], args[1], depLabel); err != nil {
		return fmt.Errorf("add dependency: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "dependency added: %s → %s\n", args[0], args[1])
	return nil
}

func runListDependencies(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}
	list, err := vault.GetDependencies(client.Logical(), args[0])
	if err != nil {
		return fmt.Errorf("list dependencies: %w", err)
	}
	if len(list.Dependencies) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no dependencies registered")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TARGET\tLABEL\tCREATED")
	for _, d := range list.Dependencies {
		fmt.Fprintf(w, "%s\t%s\t%s\n", d.TargetPath, d.Label, d.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
