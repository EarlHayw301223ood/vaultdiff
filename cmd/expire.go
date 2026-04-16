package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/vaultdiff/internal/vault"
)

var expireTTL string

var expireCmd = &cobra.Command{
	Use:   "expire <mount>",
	Short: "Scan a mount for secrets whose latest version exceeds a TTL",
	Args:  cobra.ExactArgs(1),
	RunE:  runExpire,
}

func init() {
	expireCmd.Flags().StringVar(&expireTTL, "ttl", "720h", "Maximum allowed age for a secret version (e.g. 24h, 30d)")
	rootCmd.AddCommand(expireCmd)
}

func runExpire(cmd *cobra.Command, args []string) error {
	mount := args[0]

	ttl, err := time.ParseDuration(expireTTL)
	if err != nil {
		return fmt.Errorf("invalid ttl %q: %w", expireTTL, err)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	results, err := vault.ScanExpired(client, mount, ttl)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No expired secrets found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tVERSION\tCREATED\tEXPIRED AT")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			r.Path,
			r.Version,
			r.CreatedAt.Format(time.RFC3339),
			r.ExpiredAt.Format(time.RFC3339),
		)
	}
	w.Flush()
	fmt.Printf("\n%d expired secret(s) found.\n", len(results))
	return nil
}
