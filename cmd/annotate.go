package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

var annotateAuthor string

func init() {
	annotateSetCmd := &cobra.Command{
		Use:   "set <path> <version> <note>",
		Short: "Attach a note to a secret version",
		Args:  cobra.ExactArgs(3),
		RunE:  runSetAnnotation,
	}
	annotateSetCmd.Flags().StringVar(&annotateAuthor, "author", "", "Author of the annotation (defaults to current user)")

	annotateGetCmd := &cobra.Command{
		Use:   "get <path> <version>",
		Short: "Retrieve the note attached to a secret version",
		Args:  cobra.ExactArgs(2),
		RunE:  runGetAnnotation,
	}

	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Manage annotations on secret versions",
	}
	annotateCmd.AddCommand(annotateSetCmd)
	annotateCmd.AddCommand(annotateGetCmd)
	rootCmd.AddCommand(annotateCmd)
}

func runSetAnnotation(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}
	note := args[2]

	author := annotateAuthor
	if author == "" {
		author = os.Getenv("USER")
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	a, err := vault.SetAnnotation(client, path, version, note, author)
	if err != nil {
		return err
	}
	fmt.Printf("Annotation set for %s@v%d by %s\n", a.Path, a.Version, a.Author)
	return nil
}

func runGetAnnotation(cmd *cobra.Command, args []string) error {
	path := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", args[1], err)
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	a, err := vault.GetAnnotation(client, path, version)
	if err != nil {
		return err
	}
	fmt.Printf("Path:    %s\nVersion: %d\nAuthor:  %s\nDate:    %s\nNote:    %s\n",
		a.Path, a.Version, a.Author, a.CreatedAt.Format("2006-01-02 15:04:05"), a.Note)
	return nil
}
