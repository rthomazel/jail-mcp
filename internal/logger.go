package internal

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Logger wraps slog and holds the file handle so we can close it cleanly on shutdown.
type Logger struct {
	*slog.Logger
	f io.Closer
}

// NewLogger opens (or creates) logPath for appending and returns a text-format
// slog logger writing to it. Text format is chosen over JSON for human readability
// when tailing or grepping the log file directly.
//
// Example line:
//
//	time=2026-03-05T14:32:01Z level=INFO msg="exec done" cmd="go build ./..." exit_code=0 duration=1.82s
func NewLogger(logPath string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	w := io.MultiWriter(f, os.Stderr)
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		Logger: slog.New(handler),
		f:      f,
	}, nil
}

func (l *Logger) Close() {
	_ = l.f.Close()
}
