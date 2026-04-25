package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var auditTrailCmd = &cobra.Command{
	Use:   "audit-trail",
	Short: "Manage audit trail entries for secrets",
}

var appendAuditTrailCmd = &cobra.Command{
	Use:   "append <path>",
	Short: "Append an entry to the audit trail for a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runAppendAuditTrail,
}

var getAuditTrailCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the full audit trail for a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetAuditTrail,
}

var (
	auditOperation string
	auditActor     string
	auditVersion   int
	auditNote      string
)

func init() {
	appendAuditTrailCmd.Flags().StringVar(&auditOperation, "operation", "write", "Operation name (e.g. write, read, promote)")
	appendAuditTrailCmd.Flags().StringVar(&auditActor, "actor", "", "Actor performing the operation")
	appendAuditTrailCmd.Flags().IntVar(&auditVersion, "version", 1, "Secret version the operation applies to")
	appendAuditTrailCmd.Flags().StringVar(&auditNote, "note", "", "Optional note")

	auditTrailCmd.AddCommand(appendAuditTrailCmd)
	auditTrailCmd.AddCommand(getAuditTrailCmd)
	rootCmd.AddCommand(auditTrailCmd)
}

func runAppendAuditTrail(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	if err := vault.AppendAuditTrail(client, args[0], auditOperation, auditActor, auditVersion, auditNote); err != nil {
		return fmt.Errorf("append audit trail: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "audit trail entry recorded for %s\n", args[0])
	return nil
}

func runGetAuditTrail(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	trail, err := vault.GetAuditTrail(client, args[0])
	if err != nil {
		return fmt.Errorf("get audit trail: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tACTOR\tVERSION\tNOTE")
	for _, e := range trail.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Operation, e.Actor, e.Version, e.Note)
	}
	w.Flush()

	if len(trail.Entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no audit trail entries found")
	}
	return nil
}
