package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

var schemaCmd = &cobra.Command{
	Use:   "schema <path>",
	Short: "Validate a secret against a JSON schema file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSchema,
}

var schemaFile string

func init() {
	schemaCmd.Flags().StringVar(&schemaFile, "schema", "", "path to JSON schema file (required)")
	_ = schemaCmd.MarkFlagRequired("schema")
	rootCmd.AddCommand(schemaCmd)
}

type jsonSchemaRule struct {
	Key      string `json:"key"`
	Pattern  string `json:"pattern"`
	Required bool   `json:"required"`
	Format   string `json:"format"`
	MinLen   int    `json:"min_len"`
}

func runSchema(cmd *cobra.Command, args []string) error {
	path := args[0]

	raw, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var jsonRules []jsonSchemaRule
	if err := json.Unmarshal(raw, &jsonRules); err != nil {
		return fmt.Errorf("parsing schema file: %w", err)
	}

	schema, err := buildSchema(jsonRules)
	if err != nil {
		return err
	}

	client, err := vault.NewClient(vault.ConfigFromEnv())
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secret, err := client.Logical().Read(path)
	if err != nil {
		return fmt.Errorf("reading secret: %w", err)
	}

	data := vault.ExtractStringMap(secret)
	violations := vault.ValidateSecret(data, schema)

	if len(violations) == 0 {
		fmt.Println("✓ secret passes schema validation")
		return nil
	}

	fmt.Fprintf(os.Stderr, "schema violations for %s:\n", path)
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v.Error())
	}
	return fmt.Errorf("%d violation(s) found", len(violations))
}

func buildSchema(rules []jsonSchemaRule) (vault.Schema, error) {
	var s vault.Schema
	for _, r := range rules {
		sr := vault.SchemaRule{
			Key:      r.Key,
			Required: r.Required,
			MinLen:   r.MinLen,
		}
		if r.Pattern != "" {
			pat, err := regexp.Compile(r.Pattern)
			if err != nil {
				return s, fmt.Errorf("invalid pattern %q: %w", r.Pattern, err)
			}
			sr.Pattern = pat
		}
		if r.Format != "" {
			fmt2, err := regexp.Compile(r.Format)
			if err != nil {
				return s, fmt.Errorf("invalid format %q: %w", r.Format, err)
			}
			sr.Format = fmt2
		}
		s.Rules = append(s.Rules, sr)
	}
	return s, nil
}
