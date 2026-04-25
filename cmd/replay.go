package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var replayCmd = &cobra.Command{
	Use:   "replay <path>",
	Short: "Replay all active versions of a secret in chronological order",
	Args:  cobra.ExactArgs(1),
	RunE:  runReplay,
}

var replayMount string

func init() {
	replayCmd.Flags().StringVar(&replayMount, "mount", "secret", "KV v2 mount name")
	rootCmd.AddCommand(replayCmd)
}

func runReplay(cmd *cobra.Command, args []string) error {
	path := args[0]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("replay: config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("replay: client: %w", err)
	}

	log, err := vault.ReplaySecret(client.Logical(), path, replayMount)
	if err != nil {
		return err
	}

	if len(log.Entries) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "no active versions found for %q\n", path)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tKEYS")
	for _, e := range log.Entries {
		fmt.Fprintf(w, "%d\t%s\t%d\n",
			e.Version,
			e.CreatedAt.Format("2006-01-02T15:04:05Z"),
			len(e.Data),
		)
	}
	return w.Flush()
}
