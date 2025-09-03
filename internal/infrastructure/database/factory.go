// Package database provides database abstraction and connection management.
package database

import (
	"fmt"
	"log/slog"
)

// Factory creates database instances based on configuration
type Factory struct {
	logger *slog.Logger
}

// NewFactory creates a new database factory
func NewFactory(logger *slog.Logger) *Factory {
	return &Factory{
		logger: logger,
	}
}

// Create creates a new database connection based on the provided configuration
func (f *Factory) Create(cfg *Config) (Database, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config is nil")
	}

	switch cfg.Type {
	case TypePostgreSQL:
		return NewPostgreSQL(cfg, f.logger)
	case TypeSQLite:
		return NewSQLite(cfg, f.logger)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}

// CreateFromURL creates a database connection from a connection URL
// It auto-detects the database type from the URL scheme
func (f *Factory) CreateFromURL(url string) (Database, error) {
	dbType, err := detectTypeFromURL(url)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig(dbType)
	cfg.URL = url

	return f.Create(cfg)
}

// detectTypeFromURL detects the database type from the connection URL
func detectTypeFromURL(url string) (Type, error) {
	switch {
	case len(url) >= 11 && url[:11] == "postgresql:":
		return TypePostgreSQL, nil
	case len(url) >= 9 && url[:9] == "postgres:":
		return TypePostgreSQL, nil
	case len(url) >= 4 && url[:4] == "pg:":
		return TypePostgreSQL, nil
	case len(url) >= 7 && url[:7] == "sqlite:":
		return TypeSQLite, nil
	case len(url) >= 8 && url[:8] == "sqlite3:":
		return TypeSQLite, nil
	case len(url) >= 7 && url[:7] == "file://":
		return TypeSQLite, nil
	case len(url) >= 10 && url[:10] == ":memory:":
		return TypeSQLite, nil
	default:
		return "", fmt.Errorf("cannot detect database type from URL: %s", url)
	}
}
