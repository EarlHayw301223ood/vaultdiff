package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "Manage named version bookmarks for secret paths",
}

var setBookmarkCmd = &cobra.Command{
	Use:   "set <path> <name> <version>",
	Short: "Set a named bookmark pointing to a secret version",
	Args:  cobra.ExactArgs(3),
	RunE:  runSetBookmark,
}

var getBookmarkCmd = &cobra.Command{
	Use:   "get <path> <name>",
	Short: "Get the version number for a named bookmark",
	Args:  cobra.ExactArgs(2),
	RunE:  runGetBookmark,
}

var deleteBookmarkCmd = &cobra.Command{
	Use:   "delete <path> <name>",
	Short: "Delete a named bookmark",
	Args:  cobra.ExactArgs(2),
	RunE:  runDeleteBookmark,
}

func init() {
	bookmarkCmd.AddCommand(setBookmarkCmd)
	bookmarkCmd.AddCommand(getBookmarkCmd)
	bookmarkCmd.AddCommand(deleteBookmarkCmd)
	rootCmd.AddCommand(bookmarkCmd)
}

func runSetBookmark(cmd *cobra.Command, args []string) error {
	path, name, verStr := args[0], args[1], args[2]
	version, err := strconv.Atoi(verStr)
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", verStr, err)
	}

	client, err := vault.NewClient()
	if err != nil {
		return err
	}

	if err := vault.SetBookmark(client, path, name, version); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "bookmark %q set to version %d for %s\n", name, version, path)
	return nil
}

func runGetBookmark(cmd *cobra.Command, args []string) error {
	path, name := args[0], args[1]

	client, err := vault.NewClient()
	if err != nil {
		return err
	}

	version, err := vault.GetBookmark(client, path, name)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "bookmark %q → version %d\n", name, version)
	return nil
}

func runDeleteBookmark(cmd *cobra.Command, args []string) error {
	path, name := args[0], args[1]

	client, err := vault.NewClient()
	if err != nil {
		return err
	}

	if err := vault.DeleteBookmark(client, path, name); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "bookmark %q deleted from %s\n", name, path)
	return nil
}
