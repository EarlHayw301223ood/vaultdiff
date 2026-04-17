package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Save and retrieve named baselines for secret paths",
}

var saveBaselineCmd = &cobra.Command{
	Use:   "save <path> <name> <version>",
	Short: "Save a named baseline for a secret path at a specific version",
	Args:  cobra.ExactArgs(3),
	RunE:  runSaveBaseline,
}

var getBaselineCmd = &cobra.Command{
	Use:   "get <path> <name>",
	Short: "Retrieve a named baseline for a secret path",
	Args:  cobra.ExactArgs(2),
	RunE:  runGetBaseline,
}

func init() {
	baselineCmd.AddCommand(saveBaselineCmd)
	baselineCmd.AddCommand(getBaselineCmd)
	rootCmd.AddCommand(baselineCmd)
}

func runSaveBaseline(cmd *cobra.Command, args []string) error {
	path, name, verStr := args[0], args[1], args[2]
	version, err := strconv.Atoi(verStr)
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", verStr, err)
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	bl, err := vault.SaveBaseline(client, path, name, version)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Baseline %q saved for %s@v%d (created: %s)\n",
		bl.Name, bl.Path, bl.Version, bl.CreatedAt.Format("2006-01-02T15:04:05Z"))
	return nil
}

func runGetBaseline(cmd *cobra.Command, args []string) error {
	path, name := args[0], args[1]

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	bl, err := vault.GetBaseline(client, path, name)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Baseline: %s | Path: %s | Keys: %d\n", bl.Name, bl.Path, len(bl.Data))
	for k, v := range bl.Data {
		fmt.Fprintf(os.Stdout, "  %s = %s\n", k, v)
	}
	return nil
}
