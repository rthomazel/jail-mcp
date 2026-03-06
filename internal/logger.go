package internal

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	*slog.Logger
	f io.Closer
}

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
