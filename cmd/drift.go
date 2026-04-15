package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"vaultdiff/internal/vault"
)

var driftCmd = &cobra.Command{
	Use:   "drift <mountA> <mountB>",
	Short: "Detect secret drift between two Vault mounts",
	Args:  cobra.ExactArgs(2),
	RunE:  runDrift,
}

func init() {
	driftCmd.Flags().StringP("prefix", "p", "", "only compare paths with this prefix")
	rootCmd.AddCommand(driftCmd)
}

func runDrift(cmd *cobra.Command, args []string) error {
	mountA := args[0]
	mountB := args[1]
	prefix, _ := cmd.Flags().GetString("prefix")

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("drift: config error: %w", err)
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("drift: client error: %w", err)
	}

	snapA, err := vault.SnapshotMount(client, mountA, prefix)
	if err != nil {
		return fmt.Errorf("drift: snapshot %s: %w", mountA, err)
	}
	snapB, err := vault.SnapshotMount(client, mountB, prefix)
	if err != nil {
		return fmt.Errorf("drift: snapshot %s: %w", mountB, err)
	}

	report, err := vault.DetectDrift(mountA, mountB, snapA, snapB)
	if err != nil {
		return fmt.Errorf("drift: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "PATH\tSTATUS\n")
	for _, r := range report.Results {
		status := driftStatus(r)
		fmt.Fprintf(w, "%s\t%s\n", r.Path, status)
	}
	w.Flush()

	summary := report.Summary()
	fmt.Printf("\nSummary: %d in-sync, %d diverged, %d only-in-%s, %d only-in-%s\n",
		summary["in_sync"], summary["diverged"],
		summary["only_in_a"], mountA,
		summary["only_in_b"], mountB,
	)
	return nil
}

func driftStatus(r vault.DriftResult) string {
	switch {
	case r.InSync:
		return "in-sync"
	case r.Diverged:
		return "diverged"
	case r.OnlyInA:
		return "only-in-a"
	case r.OnlyInB:
		return "only-in-b"
	default:
		return "unknown"
	}
}
