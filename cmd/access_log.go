package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var accessLogCmd = &cobra.Command{
	Use:   "access-log",
	Short: "Append or read access log entries for a secret path",
}

var appendAccessLogCmd = &cobra.Command{
	Use:   "append <path>",
	Short: "Append an access log entry for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runAppendAccessLog,
}

var getAccessLogCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the latest access log entry for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetAccessLog,
}

func init() {
	accessLogCmd.AddCommand(appendAccessLogCmd)
	accessLogCmd.AddCommand(getAccessLogCmd)

	appendAccessLogCmd.Flags().String("operation", "read", "Operation name (read, promote, rollback, etc.)")
	appendAccessLogCmd.Flags().String("actor", "", "Actor performing the operation (defaults to current user)")
	appendAccessLogCmd.Flags().String("note", "", "Optional note to attach to the entry")
	appendAccessLogCmd.Flags().Int("version", 0, "Secret version the operation applies to")

	rootCmd.AddCommand(accessLogCmd)
}

func runAppendAccessLog(cmd *cobra.Command, args []string) error {
	path := args[0]

	operation, _ := cmd.Flags().GetString("operation")
	actor, _ := cmd.Flags().GetString("actor")
	note, _ := cmd.Flags().GetString("note")
	version, _ := cmd.Flags().GetInt("version")

	if actor == "" {
		actor = os.Getenv("USER")
		if actor == "" {
			actor = "unknown"
		}
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	entry := vault.AccessEntry{
		Path:      path,
		Version:   version,
		Operation: operation,
		Actor:     actor,
		Note:      note,
		Timestamp: time.Now().UTC(),
	}

	if err := vault.AppendAccessLog(client.Logical(), path, entry); err != nil {
		return fmt.Errorf("append access log: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "access log entry written for %s (version %s)\n",
		path, strconv.Itoa(version))
	return nil
}

func runGetAccessLog(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	entry, err := vault.GetAccessLog(client.Logical(), path)
	if err != nil {
		return fmt.Errorf("get access log: %w", err)
	}
	if entry == nil {
		fmt.Fprintf(cmd.OutOrStdout(), "no access log entry found for %s\n", path)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "path:      %s\n", entry.Path)
	fmt.Fprintf(cmd.OutOrStdout(), "version:   %d\n", entry.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "operation: %s\n", entry.Operation)
	fmt.Fprintf(cmd.OutOrStdout(), "actor:     %s\n", entry.Actor)
	if entry.Note != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "note:      %s\n", entry.Note)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "timestamp: %s\n", entry.Timestamp.Format(time.RFC3339))
	return nil
}
