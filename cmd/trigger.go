package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Manage secret change triggers",
}

var setTriggerCmd = &cobra.Command{
	Use:   "set <path> <name>",
	Short: "Attach a trigger to a secret path",
	Args:  cobra.ExactArgs(2),
	RunE:  runSetTrigger,
}

var getTriggerCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get the trigger config for a secret path",
	Args:  cobra.ExactArgs(1),
	RunE:  runGetTrigger,
}

var triggerCondition string

func init() {
	setTriggerCmd.Flags().StringVar(&triggerCondition, "condition", "any",
		"Trigger condition: 'any' or 'version_gt:<n>'")
	triggerCmd.AddCommand(setTriggerCmd, getTriggerCmd)
	rootCmd.AddCommand(triggerCmd)
}

func runSetTrigger(cmd *cobra.Command, args []string) error {
	path, name := args[0], args[1]
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	cfg := vault.TriggerConfig{
		Name:      name,
		Path:      path,
		Condition: triggerCondition,
	}
	if err := vault.SetTrigger(client.Logical(), path, cfg); err != nil {
		return fmt.Errorf("set trigger: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "trigger %q set on %s (condition: %s)\n", name, path, triggerCondition)
	return nil
}

func runGetTrigger(cmd *cobra.Command, args []string) error {
	path := args[0]
	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	cfg, err := vault.GetTrigger(client.Logical(), path)
	if err != nil {
		return fmt.Errorf("get trigger: %w", err)
	}
	if cfg == nil {
		fmt.Fprintln(os.Stderr, "no trigger configured for", path)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "name:      %s\npath:      %s\ncondition: %s\n",
		cfg.Name, cfg.Path, cfg.Condition)
	return nil
}
