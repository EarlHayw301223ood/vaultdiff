package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var retentionCmd = &cobra.Command{
	Use:   "retention",
	Short: "Manage retention policies for secret versions",
}

var setRetentionCmd = &cobra.Command{
	Use:   "set <path> <version> <max-age> <max-versions>",
	Short: "Set a retention policy on a secret",
	Args:  cobra.ExactArgs(4),
	RunE:  runSetRetention,
}

var getRetentionCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the retention policy for a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetRetention,
}

var retentionSetBy string

func init() {
	setRetentionCmd.Flags().StringVar(&retentionSetBy, "set-by", "", "Identity of the user setting the policy")
	retentionCmd.AddCommand(setRetentionCmd)
	retentionCmd.AddCommand(getRetentionCmd)
	rootCmd.AddCommand(retentionCmd)
}

func runSetRetention(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}
	maxAge, err := time.ParseDuration(args[2])
	if err != nil {
		return fmt.Errorf("invalid max-age %q: %w", args[2], err)
	}
	maxVersions, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("invalid max-versions %q: %w", args[3], err)
	}

	cfg, err := vaultConfigFromCmd(cmd)
	if err != nil {
		return err
	}
	client, err := newVaultClient(cfg)
	if err != nil {
		return err
	}

	setBy := retentionSetBy
	if setBy == "" {
		setBy, _ = os.Hostname()
	}

	if err := vault.SetRetention(client.Logical(), path, version, maxAge, maxVersions, setBy); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "retention policy set for %s (v%d)\n", path, version)
	return nil
}

func runGetRetention(cmd *cobra.Command, args []string) error {
	cfg, err := vaultConfigFromCmd(cmd)
	if err != nil {
		return err
	}
	client, err := newVaultClient(cfg)
	if err != nil {
		return err
	}

	policy, err := vault.GetRetention(client.Logical(), args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "path:         %s\n", policy.Path)
	fmt.Fprintf(cmd.OutOrStdout(), "version:      %d\n", policy.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "max_age:      %s\n", policy.MaxAge)
	fmt.Fprintf(cmd.OutOrStdout(), "max_versions: %d\n", policy.MaxVersions)
	fmt.Fprintf(cmd.OutOrStdout(), "set_by:       %s\n", policy.SetBy)
	fmt.Fprintf(cmd.OutOrStdout(), "set_at:       %s\n", policy.SetAt.Format(time.RFC3339))
	return nil
}
