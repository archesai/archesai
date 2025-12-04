// Package logger provides centralized logging configuration and utilities
package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

// New creates a configured logger with stdout output.
func New(cfg Config) *slog.Logger {
	// FIXME should i do this?
	if cfg.Pretty {
		return NewPretty(cfg)
	}

	return NewDefault(cfg)
}

// NewDefault creates a default logger with info level and stdout output.
func NewDefault(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

// NewPretty creates a pretty-printed logger for development use.
func NewPretty(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	handler := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: level == slog.LevelDebug,
		ReportCaller:    level == slog.LevelDebug,
		TimeFormat:      time.Kitchen,
		Level:           log.Level(level),
	})
	return slog.New(handler)
}

// NewDiscard creates a test logger that discards all output.
func NewDiscard() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

// NewWithWriter creates a logger with a custom writer.
func NewWithWriter(w io.Writer, cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

// parseLevel converts string log level to slog.Level.
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
		return slog.LevelError + 1
	default:
		return slog.LevelInfo
	}
}
