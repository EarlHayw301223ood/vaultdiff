package diff

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Format represents an output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat converts a string to a Format, returning an error for unknowns.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be text or json", s)
	}
}

// ExportOptions controls how a report is exported.
type ExportOptions struct {
	Format      Format
	OutputPath  string // empty means stdout
	MaskSecrets bool
}

// Export writes the report to the destination described by opts.
func Export(r Report, opts ExportOptions) error {
	var w io.Writer = os.Stdout

	if opts.OutputPath != "" {
		path := resolveOutputPath(opts.OutputPath, opts.Format)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if err := write(r, w, opts); err != nil {
		return err
	}
	return nil
}

// write dispatches to the appropriate serialisation method based on opts.Format.
func write(r Report, w io.Writer, opts ExportOptions) error {
	switch opts.Format {
	case FormatJSON:
		return r.WriteJSON(w)
	default:
		return r.WriteText(w, opts.MaskSecrets)
	}
}

// resolveOutputPath appends a datestamped suffix and extension when the path
// looks like a directory (no extension provided).
func resolveOutputPath(path string, f Format) string {
	if filepath.Ext(path) != "" {
		return path
	}
	ext := ".txt"
	if f == FormatJSON {
		ext = ".json"
	}
	stamp := time.Now().UTC().Format("20060102-150405")
	return filepath.Join(path, fmt.Sprintf("vaultdiff-%s%s", stamp, ext))
}
