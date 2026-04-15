package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var tagMount string

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage named tags pointing to secret versions",
	}

	setCmd := &cobra.Command{
		Use:   "set <path> <tag> <version>",
		Short: "Create or update a tag pointing to a specific secret version",
		Args:  cobra.ExactArgs(3),
		RunE:  runSetTag,
	}
	setCmd.Flags().StringVar(&tagMount, "mount", "secret", "KV v2 mount name")

	getCmd := &cobra.Command{
		Use:   "get <path> <tag>",
		Short: "Retrieve the version a tag points to",
		Args:  cobra.ExactArgs(2),
		RunE:  runGetTag,
	}
	getCmd.Flags().StringVar(&tagMount, "mount", "secret", "KV v2 mount name")

	tagCmd.AddCommand(setCmd, getCmd)
	rootCmd.AddCommand(tagCmd)
}

func runSetTag(cmd *cobra.Command, args []string) error {
	secretPath := args[0]
	tagName := args[1]

	version, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid version %q: must be an integer", args[2])
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	if err := vault.SetTag(client, tagMount, secretPath, tagName, version); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Tag %q set on %s → version %d\n", tagName, secretPath, version)
	return nil
}

func runGetTag(cmd *cobra.Command, args []string) error {
	secretPath := args[0]
	tagName := args[1]

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	tag, err := vault.GetTag(client, tagMount, secretPath, tagName)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "tag=%s path=%s version=%d\n", tag.Name, tag.Path, tag.Version)
	return nil
}
