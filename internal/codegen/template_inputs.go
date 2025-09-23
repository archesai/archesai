package codegen

import (
	"github.com/archesai/archesai/internal/templates"
)

// =============================================================================
// Specialized Template Input Structures
// These structures are inputs to specific template generators.
// Base structures remain in the templates package.
// =============================================================================

// TypesTemplateInput is the input for type generation templates.
type TypesTemplateInput struct {
	templates.TemplateData
	Schemas     []templates.SchemaData    // All schemas to generate
	Constants   []templates.ConstantDef   // Enums and constants
	TypeAliases []templates.TypeAliasData // Type definitions
}

// HandlerTemplateInput is the input for HTTP handler templates.
type HandlerTemplateInput struct {
	templates.TemplateData
	Operations    []templates.OperationData // HTTP operations to handle
	Middleware    []string                  // Required middleware
	HasFileUpload bool                      // Whether multipart handling needed
	ServiceName   string                    // Name of the service to inject
}

// ServiceTemplateInput is the input for service layer templates.
type ServiceTemplateInput struct {
	templates.TemplateData
	Entities []templates.EntityData // Entities to create services for
	Methods  []templates.MethodData // Service methods to generate
}

// RepositoryTemplateInput is the input for repository interface templates.
type RepositoryTemplateInput struct {
	templates.TemplateData
	Entities      []templates.EntityData // Entities needing repositories
	Operations    []string               // CRUD operations to include
	CustomMethods []templates.MethodData // Additional repository methods
	DatabaseType  string                 // postgres, sqlite, etc.
}

// EventsTemplateInput is the input for event publisher templates.
type EventsTemplateInput struct {
	templates.TemplateData
	Events     []templates.EventData // Events to publish
	EventTypes []string              // Types of events (created, updated, deleted)
}

// CacheTemplateInput is the input for cache layer generation templates.
type CacheTemplateInput struct {
	templates.TemplateData
	Entities  []templates.EntityData // Entities to cache
	CacheType string                 // redis, memory, etc.
	TTL       int                    // Default TTL in seconds
}

// SQLTemplateInput is the input for SQL schema generation templates.
type SQLTemplateInput struct {
	templates.TemplateData
	Tables       []templates.TableData // Tables to generate
	DatabaseType string                // postgres, sqlite, mysql
}
