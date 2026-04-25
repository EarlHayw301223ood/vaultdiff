package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/vaultdiff/internal/vault"
)

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Manage write quotas for secret paths",
}

var setQuotaCmd = &cobra.Command{
	Use:   "set <path> <max-writes> <window-sec>",
	Short: "Set a write quota on a secret path",
	Args:  cobra.ExactArgs(3),
	RunE:  runSetQuota,
}

var getQuotaCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the current quota for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetQuota,
}

func init() {
	quotaCmd.AddCommand(setQuotaCmd)
	quotaCmd.AddCommand(getQuotaCmd)
	rootCmd.AddCommand(quotaCmd)
}

func runSetQuota(cmd *cobra.Command, args []string) error {
	path := args[0]

	maxWrites, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid max-writes value: %w", err)
	}

	windowSec, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid window-sec value: %w", err)
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	if err := vault.SetQuota(client.Logical(), path, maxWrites, windowSec); err != nil {
		return fmt.Errorf("set quota: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "quota set: path=%s max_writes=%d window_sec=%d\n", path, maxWrites, windowSec)
	return nil
}

func runGetQuota(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	record, err := vault.GetQuota(client.Logical(), path)
	if err != nil {
		return fmt.Errorf("get quota: %w", err)
	}

	if record == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "no quota set for path:", path)
		return nil
	}

	status := "ok"
	if record.Exceeded() {
		status = "EXCEEDED"
	}

	fmt.Fprintf(cmd.OutOrStdout(),
		"path=%s writes=%d/%d window=%ds reset_at=%s status=%s remaining=%d\n",
		record.Path, record.Writes, record.MaxWrites, record.WindowSec,
		record.ResetAt.Format("2006-01-02T15:04:05Z"), status, record.Remaining(),
	)
	return nil
}
