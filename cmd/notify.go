package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Dispatch notifications when secrets change",
}

var notifyFireCmd = &cobra.Command{
	Use:   "fire <path>",
	Short: "Fire a notify event for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runNotifyFire,
}

func init() {
	notifyFireCmd.Flags().StringP("operation", "o", "write", "operation name (write, delete, read)")
	notifyFireCmd.Flags().IntP("version", "v", 0, "secret version")
	notifyFireCmd.Flags().StringP("by", "b", "", "actor who triggered the event")
	notifyFireCmd.Flags().StringArrayP("target", "t", nil, "webhook target URL (repeatable)")
	notifyFireCmd.Flags().StringP("events", "e", "*", "comma-separated event filter")
	notifyFireCmd.Flags().StringP("template", "T", "", "body template with {{path}}, {{version}}, {{operation}}, {{by}}")
	notifyCmd.AddCommand(notifyFireCmd)
	rootCmd.AddCommand(notifyCmd)
}

func runNotifyFire(cmd *cobra.Command, args []string) error {
	path := args[0]
	op, _ := cmd.Flags().GetString("operation")
	ver, _ := cmd.Flags().GetInt("version")
	by, _ := cmd.Flags().GetString("by")
	targets, _ := cmd.Flags().GetStringArray("target")
	eventsRaw, _ := cmd.Flags().GetString("events")
	template, _ := cmd.Flags().GetString("template")

	if len(targets) == 0 {
		return fmt.Errorf("at least one --target URL is required")
	}

	events := strings.Split(eventsRaw, ",")
	configs := make([]vault.NotifyConfig, 0, len(targets))
	for _, tgt := range targets {
		configs = append(configs, vault.NotifyConfig{
			Channel:  vault.ChannelCustom,
			Target:   tgt,
			Events:   events,
			Template: template,
		})
	}

	d := vault.NewDispatcher(nil)
	event := vault.NotifyEvent{
		Path:      path,
		Version:   ver,
		Operation: op,
		By:        by,
		Triggered: time.Now().UTC(),
	}

	results := vault.Notify(d, event, configs)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return fmt.Errorf("encoding results: %w", err)
	}
	for _, r := range results {
		if !r.Success {
			return fmt.Errorf("one or more notifications failed")
		}
	}
	return nil
}
