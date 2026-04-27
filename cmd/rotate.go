package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate <path>",
	Short: "Re-write a secret to create a new version (optionally transforming values)",
	Args:  cobra.ExactArgs(1),
	RunE:  runRotate,
}

var (
	rotateDryRun    bool
	rotateUppercase bool
)

func init() {
	rotateCmd.Flags().BoolVar(&rotateDryRun, "dry-run", false, "Simulate rotation without writing")
	rotateCmd.Flags().BoolVar(&rotateUppercase, "uppercase-values", false, "Transform all values to uppercase")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	opts := vault.RotateOptions{
		DryRun: rotateDryRun,
	}

	if rotateUppercase {
		opts.Transform = func(_, v string) (string, error) {
			return strings.ToUpper(v), nil
		}
	}

	res, err := vault.Rotate(client.Logical(), args[0], opts)
	if err != nil {
		return err
	}

	if rotateDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] %s: version %d → %d\n",
			res.Path, res.OldVersion, res.NewVersion)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "rotated %s: version %d → %d (at %s)\n",
			res.Path, res.OldVersion, res.NewVersion,
			res.RotatedAt.Format("2006-01-02T15:04:05Z"))
	}
	return nil
}
