package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vaultdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp   time.Time         `json:"timestamp"`
	Environment string            `json:"environment"`
	Path        string            `json:"path"`
	FromVersion int               `json:"from_version"`
	ToVersion   int               `json:"to_version"`
	Changes     []diff.Change     `json:"changes"`
	Summary     Summary           `json:"summary"`
	User        string            `json:"user,omitempty"`
}

// Summary holds counts of change types for quick inspection.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Total    int `json:"total"`
}

// Logger writes audit entries to a destination.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger writing to the given writer.
// Pass nil to default to os.Stdout.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Record builds and writes an audit Entry as a JSON line.
func (l *Logger) Record(env, path string, from, to int, changes []diff.Change) error {
	entry := Entry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Path:        path,
		FromVersion: from,
		ToVersion:   to,
		Changes:     changes,
		Summary:     buildSummary(changes),
		User:        currentUser(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(l.out, "%s\n", data)
	return err
}

func buildSummary(changes []diff.Change) Summary {
	var s Summary
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified
	return s
}

func currentUser() string {
	if u := os.Getenv("VAULT_AUDIT_USER"); u != "" {
		return u
	}
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return ""
}
