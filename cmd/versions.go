package cmd

import (
	"fmt"
	/tabwriter"
	"os"

	"github.com/spf13/cobra"

	"github.comdiff/internal/vault"
)

var versionsCmd = &cobra.Command{
	Use:   "versions <path>",
	Short: "List available versions of a Vault KV secret",
	Example: `  vaultdiff versions secret/data/prod/app`,cobra.ExactArgs(1),
	RunE:  runVersions,
}

func init() {
	rootCmd.AddCommand(versionsCmd)
}

func runVersions(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	path := args[0]
	versions, err := vault.ListVersions(cmd.Context(), client, path)
	if err != nil {
		return fmt.Errorf("list versions: %w", err)
	}

	if len(versions) == 0 {
		fmt.Println("No versions found for", path)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tDESTROYED\tDELETED")
	for _, v := range versions {
		destroyed := "-"
		if v.Destroyed {
			destroyed = "yes"
		}
		deleted := "-"
		if !v.DeletionTime.IsZero() {
			deleted = v.DeletionTime.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			v.Version,
			v.CreatedTime.Format("2006-01-02 15:04:05"),
			destroyed,
			deleted,
		)
	}
	return w.Flush()
}
