package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"context"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch <path>",
	Short: "Poll a Vault KV path and print an alert when the version changes",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 30, "Poll interval in seconds")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	path := args[0]

	if watchInterval <= 0 {
		return fmt.Errorf("interval must be a positive number of seconds, got %d", watchInterval)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("vault config: %w", err)
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	interval := time.Duration(watchInterval) * time.Second
	watcher := vault.NewWatcher(client.Logical(), path, interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s every %s — press Ctrl+C to stop\n", path, interval)

	ch, err := watcher.Watch(ctx)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	for evt := range ch {
		fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s changed: v%d → v%d\n",
			evt.DetectedAt.Format(time.RFC3339),
			evt.Path,
			evt.OldVersion,
			evt.NewVersion,
		)
	}
	return nil
}
