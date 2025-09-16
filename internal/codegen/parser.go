// Package codegen provides schema parsing using the Speakeasy OpenAPI library.
// This parser offers better OpenAPI 3.1.1 support and improved x-codegen extraction.
package codegen

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/speakeasy-api/openapi/extensions"
	"github.com/speakeasy-api/openapi/jsonschema/oas3"
	"github.com/speakeasy-api/openapi/openapi"
)

// ParsedSchema wraps the OpenAPI schema with code generation metadata.
type ParsedSchema struct {
	// The underlying OpenAPI schema
	*oas3.Schema

	// Schema name (e.g., "User", "Organization")
	Name string

	// Domain this schema belongs to (e.g., "auth", "organizations")
	Domain string

	// x-codegen extension containing all generation configuration
	XCodegen *XCodegen `yaml:"x-codegen,omitempty" json:"x-codegen,omitempty"`
}

// Parser handles parsing of OpenAPI schemas with x-codegen extensions.
// This implementation uses the Speakeasy OpenAPI library for better OpenAPI 3.1.1 support.
type Parser struct {
	// Parsed schemas cache
	schemas map[string]*ParsedSchema

	// Base directory for relative paths
	baseDir string

	// Warnings accumulated during parsing
	warnings []string

	// OpenAPI document
	doc *openapi.OpenAPI
}

// NewParser creates a new schema parser using the Speakeasy OpenAPI library.
func NewParser(baseDir string) *Parser {
	return &Parser{
		schemas:  make(map[string]*ParsedSchema),
		baseDir:  baseDir,
		warnings: []string{},
	}
}

// ParseOpenAPISpec parses a complete OpenAPI specification file using Speakeasy.
func (p *Parser) ParseOpenAPISpec(specPath string) (map[string]*ParsedSchema, error) {
	ctx := context.Background()

	// Open the spec file
	f, err := os.Open(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open spec file: %w", err)
	}
	defer func() { _ = f.Close() }()

	// Parse and validate the OpenAPI document
	doc, validationErrs, err := openapi.Unmarshal(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OpenAPI spec: %w", err)
	}

	// Store document for reference
	p.doc = doc

	// Collect validation warnings
	for _, validErr := range validationErrs {
		p.warnings = append(p.warnings, validErr.Error())
	}

	// Check if we have components and schemas
	if doc.Components == nil || doc.Components.Schemas == nil {
		return nil, fmt.Errorf("no components.schemas found in spec")
	}

	// Parse all schemas in components
	schemas := make(map[string]*ParsedSchema)
	for schemaName := range doc.Components.Schemas.Keys() {
		schemaRef := doc.Components.Schemas.GetOrZero(schemaName)
		if schemaRef != nil && schemaRef.IsLeft() {
			schema := schemaRef.GetLeft()
			parsed := p.parseSchema(schemaName, schema)
			if parsed != nil {
				schemas[schemaName] = parsed
			}
		}
	}

	// Cache the parsed schemas
	p.schemas = schemas

	return schemas, nil
}

// parseSchema converts a Speakeasy schema to our ParsedSchema format.
func (p *Parser) parseSchema(name string, schema *oas3.Schema) *ParsedSchema {
	// Create parsed schema
	parsed := &ParsedSchema{
		Schema: schema,
		Name:   name,
		Domain: p.inferDomain("", name),
	}

	// Extract x-codegen extension
	if schema.Extensions != nil {
		xcodegen := p.extractXCodegen(schema.Extensions)
		if xcodegen != nil {
			parsed.XCodegen = xcodegen

			// Events are available in xcodegen.Events if needed
		}
	}

	// The schema already contains all properties, allOf, defaults, and enums
	// We don't need to duplicate them

	return parsed
}

// extractXCodegen extracts the x-codegen extension from schema extensions.
func (p *Parser) extractXCodegen(ext *extensions.Extensions) *XCodegen {
	raw, err := extensions.GetExtensionValue[interface{}](ext, "x-codegen")
	if err != nil || raw == nil {
		return nil
	}

	// Marshal to JSON then unmarshal to our type
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		p.warnings = append(p.warnings, fmt.Sprintf("failed to marshal x-codegen: %v", err))
		return nil
	}

	var xcodegen XCodegen
	if err := json.Unmarshal(jsonBytes, &xcodegen); err != nil {
		p.warnings = append(p.warnings, fmt.Sprintf("failed to unmarshal x-codegen: %v", err))
		return nil
	}

	return &xcodegen
}

// GetWarnings returns any warnings accumulated during parsing.
func (p *Parser) GetWarnings() []string {
	return p.warnings
}

// inferDomain attempts to infer the domain from schema name.
func (p *Parser) inferDomain(_, schemaName string) string {
	// Simple rule: use lowercase plural form of the entity name
	// This creates a one-to-one mapping between entities and packages

	// Convert schema name to lowercase
	lower := strings.ToLower(schemaName)

	// Simple pluralization rules
	if strings.HasSuffix(lower, "s") {
		// Already plural (e.g., "sessions" stays "sessions")
		return lower
	}
	if strings.HasSuffix(lower, "y") {
		// Change y to ies (e.g., "entity" -> "entities")
		return lower[:len(lower)-1] + "ies"
	}
	// Just add s (e.g., "user" -> "users", "artifact" -> "artifacts")
	return lower + "s"
}

// GetDefaultValues extracts all default values from a schema.
func (p *Parser) GetDefaultValues(schemaName string) (map[string]interface{}, error) {
	schema, exists := p.schemas[schemaName]
	if !exists {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}

	defaults := make(map[string]interface{})

	// Extract defaults from properties
	if schema.Properties != nil {
		for propName := range schema.Properties.Keys() {
			propRef := schema.Properties.GetOrZero(propName)
			if propRef != nil && propRef.IsLeft() {
				prop := propRef.GetLeft()
				if prop.Default != nil {
					var defaultValue any
					if err := prop.Default.Decode(&defaultValue); err == nil {
						defaults[propName] = defaultValue
					}
				}
			}
		}
	}

	return defaults, nil
}

// ParseFile is a compatibility method that delegates to ParseOpenAPISpec.
func (p *Parser) ParseFile(filePath string) (map[string]*ParsedSchema, error) {
	return p.ParseOpenAPISpec(filePath)
}

// Helper functions for code generation compatibility

// HasXCodegen checks if a schema has x-codegen extensions.
func HasXCodegen(schema *ParsedSchema) bool {
	return schema != nil && schema.XCodegen != nil
}

// NeedsRepository checks if a schema needs repository generation.
func NeedsRepository(schema *ParsedSchema) bool {
	if schema == nil || schema.XCodegen == nil {
		return false
	}
	// Check if repository operations are defined
	return len(schema.XCodegen.Repository.Operations) > 0
}

// NeedsCache checks if a schema needs cache generation.
func NeedsCache(schema *ParsedSchema) bool {
	if schema == nil || schema.XCodegen == nil {
		return false
	}
	return schema.XCodegen.Cache.Enabled
}

// NeedsEvents checks if a schema needs event generation.
func NeedsEvents(schema *ParsedSchema) bool {
	if schema == nil || schema.XCodegen == nil {
		return false
	}
	return len(schema.XCodegen.Events) > 0
}

// NeedsService checks if a schema needs service generation.
func NeedsService(schema *ParsedSchema) bool {
	if schema == nil || schema.XCodegen == nil {
		return false
	}
	// Check if repository operations are defined
	return len(schema.XCodegen.Repository.Operations) > 0
}

// NeedsHandler checks if a schema needs handler generation.
func NeedsHandler(schema *ParsedSchema) bool {
	if schema == nil || schema.XCodegen == nil {
		return false
	}
	// Check if repository operations are defined
	return len(schema.XCodegen.Repository.Operations) > 0
}
