package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var scanMinVersions int

var scanCmd = &cobra.Command{
	Use:   "scan <mount>",
	Short: "Scan all secrets under a mount and report version counts",
	Args:  cobra.ExactArgs(1),
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().IntVar(&scanMinVersions, "min-versions", 0, "Only show secrets with at least this many versions")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	mount := args[0]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	results, err := vault.ScanMount(cmd.Context(), client, mount)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	if scanMinVersions > 0 {
		results = vault.FilterScanResults(results, func(r vault.ScanResult) bool {
			return r.Versions >= scanMinVersions
		})
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No secrets found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tVERSIONS")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%d\n", r.Path, r.Versions)
	}
	return w.Flush()
}
