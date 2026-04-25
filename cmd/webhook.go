package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Dispatch a webhook event for a secret path",
}

var webhookFireCmd = &cobra.Command{
	Use:   "fire <path>",
	Short: "Send a webhook event to one or more URLs",
	Args:  cobra.ExactArgs(1),
	RunE:  runWebhookFire,
}

var (
	webhookURLs      []string
	webhookOperation string
	webhookVersion   int
	webhookHeader    []string
)

func init() {
	webhookFireCmd.Flags().StringArrayVar(&webhookURLs, "url", nil, "Webhook URL (repeatable)")
	webhookFireCmd.Flags().StringVar(&webhookOperation, "operation", "write", "Operation name to include in payload")
	webhookFireCmd.Flags().IntVar(&webhookVersion, "version", 0, "Secret version to include in payload")
	webhookFireCmd.Flags().StringArrayVar(&webhookHeader, "header", nil, "Extra header in Key:Value format (repeatable)")
	_ = webhookFireCmd.MarkFlagRequired("url")

	webhookCmd.AddCommand(webhookFireCmd)
	rootCmd.AddCommand(webhookCmd)
}

func runWebhookFire(cmd *cobra.Command, args []string) error {
	path := args[0]

	headers := make(map[string]string)
	for _, h := range webhookHeader {
		for i := 0; i < len(h); i++ {
			if h[i] == ':' {
				headers[h[:i]] = h[i+1:]
				break
			}
		}
	}

	configs := make([]vault.WebhookConfig, len(webhookURLs))
	for i, u := range webhookURLs {
		configs[i] = vault.WebhookConfig{
			URL:     u,
			Headers: headers,
			Timeout: 10 * time.Second,
		}
	}

	event := vault.WebhookEvent{
		Path:      path,
		Operation: webhookOperation,
		Version:   webhookVersion,
		Timestamp: time.Now().UTC(),
	}

	d := vault.NewDispatcher()
	results := d.Send(context.Background(), event, configs)

	hasErr := false
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "ERROR %s: %v\n", r.URL, r.Err)
			hasErr = true
		} else {
			fmt.Printf("OK    %s [%d]\n", r.URL, r.StatusCode)
		}
	}
	if hasErr {
		return fmt.Errorf("one or more webhook deliveries failed")
	}
	return nil
}
