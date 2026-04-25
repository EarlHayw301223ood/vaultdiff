package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/vaultdiff/internal/vault"
)

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "View or append operation traces for a secret path",
}

var traceGetCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Display the trace log for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetTrace,
}

var traceAddCmd = &cobra.Command{
	Use:   "add <path> <operation>",
	Short: "Append a trace entry for a secret path",
	Args:  cobra.ExactArgs(2),
	RunE:  runAddTrace,
}

var traceFilterOp string
var traceVersion int
var traceUser string
var traceNote string

func init() {
	traceGetCmd.Flags().StringVar(&traceFilterOp, "op", "", "Filter entries by operation")
	traceAddCmd.Flags().IntVar(&traceVersion, "version", 0, "Secret version this trace entry relates to")
	traceAddCmd.Flags().StringVar(&traceUser, "user", "", "User performing the operation")
	traceAddCmd.Flags().StringVar(&traceNote, "note", "", "Optional note to attach")

	traceCmd.AddCommand(traceGetCmd)
	traceCmd.AddCommand(traceAddCmd)
	rootCmd.AddCommand(traceCmd)
}

func runGetTrace(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	log, err := vault.GetTrace(client, args[0])
	if err != nil {
		return err
	}

	entries := log.Entries
	if traceFilterOp != "" {
		entries = log.FilterByOperation(traceFilterOp)
	}

	if len(entries) == 0 {
		fmt.Println("no trace entries found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tVERSION\tUSER\tNOTE")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Operation, e.Version, e.User, e.Note)
	}
	return w.Flush()
}

func runAddTrace(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	if err := vault.AppendTrace(client, args[0], args[1], traceUser, traceNote, traceVersion); err != nil {
		return err
	}

	fmt.Printf("trace entry recorded for %s (op=%s)\n", args[0], args[1])
	return nil
}
