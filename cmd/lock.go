package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var (
	lockTTL   time.Duration
	lockOwner string
)

func init() {
	lockCmd := &cobra.Command{
		Use:   "lock <subcommand>",
		Short: "Manage advisory locks on Vault KV paths",
	}

	acquireCmd := &cobra.Command{
		Use:   "acquire <path>",
		Short: "Acquire an advisory lock at the given KV path",
		Args:  cobra.ExactArgs(1),
		RunE:  runAcquireLock,
	}
	acquireCmd.Flags().DurationVar(&lockTTL, "ttl", 5*time.Minute, "Duration the lock should be held")
	acquireCmd.Flags().StringVar(&lockOwner, "owner", "", "Owner identifier (defaults to current user)")

	releaseCmd := &cobra.Command{
		Use:   "release <path>",
		Short: "Release an advisory lock at the given KV path",
		Args:  cobra.ExactArgs(1),
		RunE:  runReleaseLock,
	}

	lockCmd.AddCommand(acquireCmd, releaseCmd)
	rootCmd.AddCommand(lockCmd)
}

func runAcquireLock(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	owner := lockOwner
	if owner == "" {
		owner, _ = os.Hostname()
	}

	res, err := vault.AcquireLock(context.Background(), client.Logical(), args[0], vault.LockOptions{
		TTL:   lockTTL,
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Lock acquired: path=%s owner=%s acquired_at=%s\n",
		res.Path, res.Owner, res.AcquiredAt.Format(time.RFC3339))
	return nil
}

func runReleaseLock(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	if err := vault.ReleaseLock(context.Background(), client.Logical(), args[0]); err != nil {
		return fmt.Errorf("release lock: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Lock released: path=%s\n", args[0])
	return nil
}
