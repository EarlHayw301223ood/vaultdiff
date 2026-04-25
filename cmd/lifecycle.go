package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "Manage secret lifecycle stages",
}

var setLifecycleCmd = &cobra.Command{
	Use:   "set <path> <version> <stage>",
	Short: "Set the lifecycle stage for a secret version",
	Args:  cobra.ExactArgs(3),
	RunE:  runSetLifecycle,
}

var getLifecycleCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the current lifecycle stage for a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetLifecycle,
}

var lifecycleReason string
var lifecycleBy string

func init() {
	setLifecycleCmd.Flags().StringVar(&lifecycleReason, "reason", "", "Reason for the lifecycle change")
	setLifecycleCmd.Flags().StringVar(&lifecycleBy, "by", "", "Identity making the change (defaults to current user)")
	lifecycleCmd.AddCommand(setLifecycleCmd)
	lifecycleCmd.AddCommand(getLifecycleCmd)
	rootCmd.AddCommand(lifecycleCmd)
}

func runSetLifecycle(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}
	stage := vault.LifecycleStage(args[2])

	by := lifecycleBy
	if by == "" {
		by, _ = os.Hostname()
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}

	rec, err := vault.SetLifecycle(client.Logical(), path, version, stage, by, lifecycleReason)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "lifecycle set: path=%s version=%d stage=%s changed_by=%s\n",
		rec.Path, rec.Version, rec.Stage, rec.ChangedBy)
	return nil
}

func runGetLifecycle(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return err
	}

	rec, err := vault.GetLifecycle(client.Logical(), args[0])
	if err != nil {
		return err
	}
	if rec == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "no lifecycle record found")
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "path=%s version=%d stage=%s changed_by=%s changed_at=%s reason=%s\n",
		rec.Path, rec.Version, rec.Stage, rec.ChangedBy, rec.ChangedAt.Format("2006-01-02T15:04:05Z"), rec.Reason)
	return nil
}
