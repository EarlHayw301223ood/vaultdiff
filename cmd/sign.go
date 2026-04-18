package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign and verify secret versions using HMAC-SHA256",
}

var signSetCmd = &cobra.Command{
	Use:   "set <path> <version>",
	Short: "Compute and store an HMAC-SHA256 signature for a secret version",
	Args:  cobra.ExactArgs(2),
	RunE:  runSignSet,
}

var signVerifyCmd = &cobra.Command{
	Use:   "verify <path> <version>",
	Short: "Verify the stored signature for a secret version",
	Args:  cobra.ExactArgs(2),
	RunE:  runSignVerify,
}

func init() {
	signSetCmd.Flags().String("hmac-key", "", "HMAC key (or set VAULTDIFF_HMAC_KEY)")
	signVerifyCmd.Flags().String("hmac-key", "", "HMAC key (or set VAULTDIFF_HMAC_KEY)")
	signCmd.AddCommand(signSetCmd, signVerifyCmd)
	rootCmd.AddCommand(signCmd)
}

func resolveHMACKey(cmd *cobra.Command) string {
	if k, _ := cmd.Flags().GetString("hmac-key"); k != "" {
		return k
	}
	return os.Getenv("VAULTDIFF_HMAC_KEY")
}

func runSignSet(cmd *cobra.Command, args []string) error {
	path := args[0]
	v, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}
	key := resolveHMACKey(cmd)
	if key == "" {
		return fmt.Errorf("hmac key required: use --hmac-key or VAULTDIFF_HMAC_KEY")
	}

	client, err := vault.NewClient()
	if err != nil {
		return err
	}

	res, err := vault.SignSecret(client, path, v, key)
	if err != nil {
		return err
	}
	fmt.Printf("signed %s@v%d: %s\n", res.Path, res.Version, res.Signature)
	return nil
}

func runSignVerify(cmd *cobra.Command, args []string) error {
	path := args[0]
	v, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}
	key := resolveHMACKey(cmd)
	if key == "" {
		return fmt.Errorf("hmac key required: use --hmac-key or VAULTDIFF_HMAC_KEY")
	}

	client, err := vault.NewClient()
	if err != nil {
		return err
	}

	ok, err := vault.VerifySecret(client, path, v, key)
	if err != nil {
		return err
	}
	if ok {
		fmt.Printf("✓ signature valid for %s@v%d\n", path, v)
	} else {
		fmt.Printf("✗ signature mismatch for %s@v%d\n", path, v)
		os.Exit(1)
	}
	return nil
}
