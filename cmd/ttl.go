package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/vaultdiff/internal/vault"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Manage TTLs on secret versions",
}

var setTTLCmd = &cobra.Command{
	Use:   "set <path> <version> <duration>",
	Short: "Set a TTL on a secret version (e.g. 24h, 30m)",
	Args:  cobra.ExactArgs(3),
	RunE:  runSetTTL,
}

var getTTLCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the TTL record for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetTTL,
}

func init() {
	ttlCmd.AddCommand(setTTLCmd)
	ttlCmd.AddCommand(getTTLCmd)
	rootCmd.AddCommand(ttlCmd)
}

func runSetTTL(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil || version <= 0 {
		return fmt.Errorf("version must be a positive integer, got %q", args[1])
	}

	ttl, err := time.ParseDuration(args[2])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", args[2], err)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	setBy, _ := cmd.Flags().GetString("set-by")
	if setBy == "" {
		setBy = "vaultdiff"
	}

	if err := vault.SetTTL(client.Logical(), path, version, ttl, setBy); err != nil {
		return fmt.Errorf("set ttl: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "TTL set: path=%s version=%d ttl=%s\n", path, version, ttl)
	return nil
}

func runGetTTL(cmd *cobra.Command, args []string) error {
	path := args[0]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	rec, err := vault.GetTTL(client.Logical(), path)
	if err != nil {
		return fmt.Errorf("get ttl: %w", err)
	}
	if rec == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "no TTL set for", path)
		return nil
	}

	status := "active"
	if rec.IsExpired() {
		status = "EXPIRED"
	}

	fmt.Fprintf(cmd.OutOrStdout(),
		"path=%s version=%d expires_at=%s remaining=%s status=%s set_by=%s\n",
		rec.Path, rec.Version,
		rec.ExpiresAt.Format(time.RFC3339),
		rec.RemainingTTL(),
		status,
		rec.SetBy,
	)
	return nil
}
