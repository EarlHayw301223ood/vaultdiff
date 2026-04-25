package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

var compareVersionsCmd = &cobra.Command{
	Use:   "compare-versions <path> <versionA> <versionB>",
	Short: "Diff two versions of the same secret",
	Args:  cobra.ExactArgs(3),
	RunE:  runCompareVersions,
}

func init() {
	compareVersionsCmd.Flags().Bool("mask", true, "mask secret values in output")
	compareVersionsCmd.Flags().String("format", "text", "output format: text or json")
	rootCmd.AddCommand(compareVersionsCmd)
}

func runCompareVersions(cmd *cobra.Command, args []string) error {
	path := args[0]

	versionA, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid versionA %q: %w", args[1], err)
	}
	versionB, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid versionB %q: %w", args[2], err)
	}

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	cmp, err := vault.CompareVersions(client, path, versionA, versionB)
	if err != nil {
		return err
	}

	changes := diff.Compare(cmp.DataA, cmp.DataB)

	mask, _ := cmd.Flags().GetBool("mask")
	fmt.Fprintf(cmd.OutOrStdout(), "Comparing %s v%d → v%d\n", cmp.Path, cmp.VersionA, cmp.VersionB)
	fmt.Fprint(cmd.OutOrStdout(), diff.Render(changes, mask))
	return nil
}
