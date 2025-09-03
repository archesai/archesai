// Package database provides database abstraction and connection management.
package database

import (
	"fmt"
	"log/slog"

	"github.com/archesai/archesai/internal/generated/api"
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
func (f *Factory) Create(cfg *api.DatabaseConfig) (Database, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config is nil")
	}

	// Determine database type from config or URL
	var dbType Type
	if cfg.Type != "" {
		dbType = ParseTypeFromString(string(cfg.Type))
	} else {
		dbType = DetectTypeFromURL(cfg.Url)
	}

	switch dbType {
	case TypePostgreSQL:
		return NewPostgreSQL(cfg, f.logger)
	case TypeSQLite:
		return NewSQLite(cfg, f.logger)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// CreateFromURL creates a database connection from a connection URL
// It auto-detects the database type from the URL scheme
func (f *Factory) CreateFromURL(url string) (Database, error) {
	dbType := DetectTypeFromURL(url)
	cfg := &api.DatabaseConfig{
		Enabled: true,
		Url:     url,
		Type:    api.DatabaseConfigType(dbType),
	}
	return f.Create(cfg)
}
