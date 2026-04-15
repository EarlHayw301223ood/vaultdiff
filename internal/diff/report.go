package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Report holds the full diff result with metadata for export.
type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	SourcePath  string    `json:"source_path"`
	TargetPath  string    `json:"target_path"`
	SourceEnv   string    `json:"source_env"`
	TargetEnv   string    `json:"target_env"`
	Changes     []Change  `json:"changes"`
	Summary     Summary   `json:"summary"`
}

// Summary aggregates counts of each change type.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Total    int `json:"total"`
}

// NewReport constructs a Report from a slice of Changes and metadata.
func NewReport(sourceEnv, targetEnv, sourcePath, targetPath string, changes []Change) Report {
	s := Summary{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Modified:
			s.Modified++
		}
	}
	return Report{
		GeneratedAt: time.Now().UTC(),
		SourcePath:  sourcePath,
		TargetPath:  targetPath,
		SourceEnv:   sourceEnv,
		TargetEnv:   targetEnv,
		Changes:     changes,
		Summary:     s,
	}
}

// WriteJSON serialises the report as indented JSON to w.
func (r Report) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// WriteText writes a human-readable report to w.
func (r Report) WriteText(w io.Writer, maskSecrets bool) error {
	fmt.Fprintf(w, "VaultDiff Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Source : [%s] %s\n", r.SourceEnv, r.SourcePath)
	fmt.Fprintf(w, "Target : [%s] %s\n\n", r.TargetEnv, r.TargetPath)
	if err := Render(w, r.Changes, maskSecrets); err != nil {
		return err
	}
	fmt.Fprintf(w, "\nSummary — added: %d  removed: %d  modified: %d  total: %d\n",
		r.Summary.Added, r.Summary.Removed, r.Summary.Modified, r.Summary.Total)
	return nil
}
