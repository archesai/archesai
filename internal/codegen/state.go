// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"fmt"
	"text/template"

	"github.com/speakeasy-api/openapi/openapi"

	"github.com/archesai/archesai/internal/parsers"
)

// GlobalState holds the global state for code generation, similar to oapi-codegen
type GlobalState struct {
	// Core configuration
	Options Configuration

	// OpenAPI spec and parsed data
	ProcessedSchemas map[string]*parsers.ProcessedSchema
	Operations       []parsers.OperationDef
	XCodegenMap      map[string]*parsers.XCodegenExtension // Maps schema name to x-codegen extension

	// Templates
	Templates map[string]*template.Template

	// File writer
	FileWriter *FileWriter

	// Import mappings
	ImportMapping map[string]string

	// Track what's been generated to avoid duplicates
	GeneratedTypes map[string]bool
	GeneratedFiles map[string]bool
}

// NewGlobalState creates a new global state
func NewGlobalState() *GlobalState {
	return &GlobalState{
		ImportMapping:  make(map[string]string),
		GeneratedTypes: make(map[string]bool),
		GeneratedFiles: make(map[string]bool),
		XCodegenMap:    make(map[string]*parsers.XCodegenExtension),
	}
}

// Initialize sets up the global state with all necessary data
func (gs *GlobalState) Initialize(
	doc *openapi.OpenAPI,
	templates map[string]*template.Template,
	fileWriter *FileWriter,
	options Configuration,
) error {
	gs.Operations = parsers.ExtractOperations(doc)

	// Process all schemas using the parser
	processedSchemas, err := parsers.ProcessAllSchemas(doc)
	if err != nil {
		return fmt.Errorf("failed to process schemas: %w", err)
	}
	gs.ProcessedSchemas = processedSchemas

	// Extract x-codegen map for quick lookup
	gs.XCodegenMap = make(map[string]*parsers.XCodegenExtension)
	for name, processed := range processedSchemas {
		if processed.XCodegen != nil {
			gs.XCodegenMap[name] = processed.XCodegen
		}
	}

	gs.Templates = templates
	gs.FileWriter = fileWriter
	gs.Options = options

	return nil
}

// HasGeneratedType checks if a type has already been generated
func (gs *GlobalState) HasGeneratedType(typeName string) bool {
	return gs.GeneratedTypes[typeName]
}

// MarkTypeGenerated marks a type as generated
func (gs *GlobalState) MarkTypeGenerated(typeName string) {
	gs.GeneratedTypes[typeName] = true
}

// HasGeneratedFile checks if a file has already been generated
func (gs *GlobalState) HasGeneratedFile(filePath string) bool {
	return gs.GeneratedFiles[filePath]
}

// MarkFileGenerated marks a file as generated
func (gs *GlobalState) MarkFileGenerated(filePath string) {
	gs.GeneratedFiles[filePath] = true
}

// GetTemplate retrieves a template by name
func (gs *GlobalState) GetTemplate(name string) (*template.Template, bool) {
	tmpl, ok := gs.Templates[name]
	return tmpl, ok
}

// globalState is the singleton instance
var globalState = NewGlobalState()
