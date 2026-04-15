package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// WatchEventRecord is a serialisable audit record for a WatchEvent.
type WatchEventRecord struct {
	Timestamp  time.Time `json:"timestamp"`
	Path       string    `json:"path"`
	OldVersion int       `json:"old_version"`
	NewVersion int       `json:"new_version"`
}

// NewWatchEventRecord converts a WatchEvent into an audit record.
func NewWatchEventRecord(evt WatchEvent) WatchEventRecord {
	return WatchEventRecord{
		Timestamp:  evt.DetectedAt,
		Path:       evt.Path,
		OldVersion: evt.OldVersion,
		NewVersion: evt.NewVersion,
	}
}

// WriteJSON serialises the record as a single JSON line to w.
func (r WatchEventRecord) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
}

// String returns a human-readable summary of the event.
func (r WatchEventRecord) String() string {
	return fmt.Sprintf("[%s] %s: v%d → v%d",
		r.Timestamp.Format(time.RFC3339),
		r.Path,
		r.OldVersion,
		r.NewVersion,
	)
}
