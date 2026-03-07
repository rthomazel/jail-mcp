package internal

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func NewLogger(logPath string) (io.Closer, error) {
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(f, os.Stderr),
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)))

	return f, nil
}
