package vault

import (
	"fmt"
	"strings"
	"time"
)

// TraceEntry records a single operation performed on a secret path.
type TraceEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Version   int       `json:"version,omitempty"`
	User      string    `json:"user,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// TraceLog holds an ordered list of trace entries for a secret path.
type TraceLog struct {
	Path    string       `json:"path"`
	Entries []TraceEntry `json:"entries"`
}

// Latest returns the most recent trace entry, or nil if empty.
func (t *TraceLog) Latest() *TraceEntry {
	if len(t.Entries) == 0 {
		return nil
	}
	return &t.Entries[len(t.Entries)-1]
}

// FilterByOperation returns entries matching the given operation (case-insensitive).
func (t *TraceLog) FilterByOperation(op string) []TraceEntry {
	op = strings.ToLower(op)
	var out []TraceEntry
	for _, e := range t.Entries {
		if strings.ToLower(e.Operation) == op {
			out = append(out, e)
		}
	}
	return out
}

func traceMetaPath(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("vaultdiff-meta/%s/__trace__", path)
}
