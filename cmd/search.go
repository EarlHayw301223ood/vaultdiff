package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var (
	searchKeyPattern   string
	searchValuePattern string
	searchCaseSensitive bool
)

var searchCmd = &cobra.Command{
	Use:   "search <mount>",
	Short: "Search secrets by key or value pattern",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().StringVarP(&searchKeyPattern, "key", "k", "", "key pattern to search")
	searchCmd.Flags().StringVarP(&searchValuePattern, "value", "v", "", "value pattern to search")
	searchCmd.Flags().BoolVar(&searchCaseSensitive, "case-sensitive", false, "enable case-sensitive matching")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	mount := args[0]

	cfg, err := vault.ConfigFromEnv()
	if err != nil {
		return err
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	opts := vault.SearchOptions{
		KeyPattern:    searchKeyPattern,
		ValuePattern:  searchValuePattern,
		CaseSensitive: searchCaseSensitive,
	}

	results, err := vault.SearchMount(client, mount, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no matches found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tMATCHED KEYS")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%v\n", r.Path, r.MatchedKeys)
	}
	return w.Flush()
}
