// Package logger provides centralized logging configuration and utilities
package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// New creates a configured logger with stdout output
func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

// NewPretty creates a pretty-printed logger for development use
func NewPretty(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
		AddSource:  level == slog.LevelDebug,
	}))
}

// NewTest creates a test logger that discards all output
func NewTest() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// NewWithWriter creates a logger with a custom writer
func NewWithWriter(w io.Writer, cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	}))
}

// parseLevel converts string log level to slog.Level
func parseLevel(level string) slog.Level {
	switch level {
	case "debug", "trace":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error", "fatal":
		return slog.LevelError
	case "silent":
		// Set to highest level to suppress all logs
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}
