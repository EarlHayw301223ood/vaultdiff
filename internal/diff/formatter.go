package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// FormatOptions controls rendering behaviour.
type FormatOptions struct {
	Color   bool
	MaskValues bool
}

// Render writes a human-readable diff to w.
func Render(w io.Writer, result *Result, opts FormatOptions) {
	header := fmt.Sprintf("--- %s (v%d)\n+++ %s (v%d)\n",
		result.Path, result.FromVersion,
		result.Path, result.ToVersion,
	)
	fmt.Fprint(w, header)

	if !result.HasChanges() {
		fmt.Fprintln(w, "  (no changes)")
		return
	}

	for _, c := range result.Changes {
		switch c.Change {
		case ChangeAdded:
			line := fmt.Sprintf("+ %s = %s", c.Key, maskIfNeeded(c.NewValue, opts.MaskValues))
			fmt.Fprintln(w, colorize(line, colorGreen, opts.Color))
		case ChangeRemoved:
			line := fmt.Sprintf("- %s = %s", c.Key, maskIfNeeded(c.OldValue, opts.MaskValues))
			fmt.Fprintln(w, colorize(line, colorRed, opts.Color))
		case ChangeModified:
			oldLine := fmt.Sprintf("- %s = %s", c.Key, maskIfNeeded(c.OldValue, opts.MaskValues))
			newLine := fmt.Sprintf("+ %s = %s", c.Key, maskIfNeeded(c.NewValue, opts.MaskValues))
			fmt.Fprintln(w, colorize(oldLine, colorRed, opts.Color))
			fmt.Fprintln(w, colorize(newLine, colorGreen, opts.Color))
		case ChangeUnchanged:
			fmt.Fprintf(w, "  %s = %s\n", c.Key, maskIfNeeded(c.NewValue, opts.MaskValues))
		}
	}

	summaryLine := fmt.Sprintf("\n@ %s", result.Summary())
	fmt.Fprintln(w, colorize(summaryLine, colorCyan, opts.Color))
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}

func maskIfNeeded(value string, mask bool) string {
	if !mask || value == "" {
		return value
	}
	return strings.Repeat("*", 8)
}
