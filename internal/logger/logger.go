package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

type loggerContextKey struct{}

var key = loggerContextKey{}

func Context(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		log,
	)
}

func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(key).(*slog.Logger)
	if !ok {
		panic("no logger in context")
	}
	return log
}

func NewLogger(cfg Config) (*Logger, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level")
	}

	if err := os.MkdirAll(cfg.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	filePath := filepath.Join(cfg.Folder, fmt.Sprintf("%s.log", timestamp))

	file, err := os.OpenFile(
		filePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	writer := io.MultiWriter(os.Stdout, file)

	log := slog.New(slog.NewJSONHandler(
		writer,
		&slog.HandlerOptions{
			Level: level,
		},
	))

	return &Logger{log, file}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}
