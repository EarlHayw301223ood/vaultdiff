package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileLogger wraps Logger with a file-backed writer that rotates by date.
type FileLogger struct {
	*Logger
	dir  string
	file *os.File
	date string
}

// NewFileLogger creates a FileLogger that appends entries to
// <dir>/vaultdiff-audit-YYYY-MM-DD.jsonl, rotating daily.
func NewFileLogger(dir string) (*FileLogger, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}

	fl := &FileLogger{dir: dir}
	if err := fl.rotate(); err != nil {
		return nil, err
	}
	fl.Logger = NewLogger(fl.file)
	return fl, nil
}

// Close flushes and closes the underlying log file.
func (fl *FileLogger) Close() error {
	if fl.file != nil {
		return fl.file.Close()
	}
	return nil
}

// rotate opens (or creates) the log file for today's date.
func (fl *FileLogger) rotate() error {
	today := time.Now().UTC().Format("2006-01-02")
	if fl.date == today && fl.file != nil {
		return nil
	}

	name := filepath.Join(fl.dir, fmt.Sprintf("vaultdiff-audit-%s.jsonl", today))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}

	if fl.file != nil {
		_ = fl.file.Close()
	}
	fl.file = f
	fl.date = today
	if fl.Logger != nil {
		fl.Logger.out = f
	}
	return nil
}

// LogPath returns the current active log file path.
func (fl *FileLogger) LogPath() string {
	if fl.file == nil {
		return ""
	}
	return fl.file.Name()
}
