package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var gcCmd = &cobra.Command{
	Use:   "gc <path>",
	Short: "Garbage-collect old secret versions at a KV v2 path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGC,
}

func init() {
	gcCmd.Flags().Int("keep-last", 5, "Number of most-recent versions to always retain")
	gcCmd.Flags().Duration("max-age", 0, "Delete versions older than this duration (e.g. 720h). Zero disables age check.")
	gcCmd.Flags().Bool("dry-run", false, "Print what would be deleted without making changes")
	rootCmd.AddCommand(gcCmd)
}

func runGC(cmd *cobra.Command, args []string) error {
	path := args[0]
	keepLast, _ := cmd.Flags().GetInt("keep-last")
	maxAge, _ := cmd.Flags().GetDuration("max-age")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if keepLast < 1 {
		return fmt.Errorf("--keep-last must be at least 1")
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	opts := vault.GCOptions{
		KeepLast: keepLast,
		MaxAge:   maxAge,
		DryRun:   dryRun,
	}

	result, err := vault.GarbageCollect(client.Logical(), path, opts)
	if err != nil {
		return err
	}

	prefix := ""
	if result.DryRun {
		prefix = "[dry-run] "
	}

	if len(result.DeletedVersions) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "%sNothing to collect at %s\n", prefix, path)
		return nil
	}

	verStrs := make([]string, len(result.DeletedVersions))
	for i, v := range result.DeletedVersions {
		verStrs[i] = fmt.Sprintf("%d", v)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%sCollected %d version(s) at %s: [%s]\n",
		prefix,
		len(result.DeletedVersions),
		path,
		strings.Join(verStrs, ", "),
	)
	fmt.Fprintf(cmd.OutOrStdout(), "Retained %d version(s). Max-age: %s\n",
		len(result.RetainedVersions),
		fmtDuration(opts.MaxAge),
	)
	return nil
}

func fmtDuration(d time.Duration) string {
	if d == 0 {
		return "none"
	}
	return d.String()
}
