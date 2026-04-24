package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/internal/vault"
)

var envCmd = &cobra.Command{
	Use:   "env <path>",
	Short: "Export secret as environment variables",
	Long:  "Fetch a secret and print its keys as shell-compatible KEY=VALUE lines.",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnv,
}

func init() {
	envCmd.Flags().StringP("ref", "r", "latest", "version ref to export (e.g. latest, 3)")
	envCmd.Flags().BoolP("sorted", "s", false, "sort output lines alphabetically")
	envCmd.Flags().StringP("prefix", "p", "", "prefix to prepend to every key")
	rootCmd.AddCommand(envCmd)
}

func runEnv(cmd *cobra.Command, args []string) error {
	path := args[0]
	ref, _ := cmd.Flags().GetString("ref")
	sorted, _ := cmd.Flags().GetBool("sorted")
	prefix, _ := cmd.Flags().GetString("prefix")

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.ExportEnv(client, path, ref)
	if err != nil {
		return fmt.Errorf("export env: %w", err)
	}

	lines := result.Lines
	if prefix != "" {
		pfx := strings.ToUpper(prefix)
		for i, l := range lines {
			lines[i] = pfx + l
		}
	}

	if sorted {
		sort.Strings(lines)
	}

	for _, l := range lines {
		fmt.Fprintln(os.Stdout, l)
	}

	fmt.Fprintf(os.Stderr, "# exported %d variable(s) from %s@%s\n", result.Count, result.Path, ref)
	return nil
}
