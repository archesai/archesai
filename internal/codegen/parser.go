// Package codegen provides schema parsing using the Speakeasy OpenAPI library.
// This parser offers better OpenAPI 3.1.1 support and improved x-codegen extraction.
package codegen

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/speakeasy-api/openapi/extensions"
	"github.com/speakeasy-api/openapi/jsonschema/oas3"
	"github.com/speakeasy-api/openapi/openapi"
)

// ParsedSchema represents a fully parsed and analyzed schema with x-codegen metadata.
type ParsedSchema struct {
	Schema

	// Domain this schema belongs to (e.g., "auth", "organizations")
	Domain string

	// File path where schema was defined
	SourceFile string

	// Timestamp when parsed
	ParsedAt time.Time

	// Any parsing warnings
	Warnings []string

	// Primary key field name (typically "ID" or "Id")
	PrimaryKey string

	// Type is the entity type name (same as Name, for template compatibility)
	Type string

	// Events extracted and formatted from XCodegen.Events
	Events []Event
}

// Event represents a domain event for code generation.
type Event struct {
	Type        string // e.g., "UserCreated"
	Description string // e.g., "User created event"
}

// Schema represents a parsed OpenAPI schema with x-codegen extensions.
type Schema struct {
	// Schema name (e.g., "User", "Organization")
	Name string

	// OpenAPI schema type
	Type string `yaml:"type" json:"type"`

	// Schema description
	Description string `yaml:"description,omitempty" json:"description,omitempty"`

	// Required fields
	Required []string `yaml:"required,omitempty" json:"required,omitempty"`

	// Schema properties
	Properties map[string]Property `yaml:"properties,omitempty" json:"properties,omitempty"`

	// x-codegen extension at schema level
	XCodegen *XCodegen `yaml:"x-codegen,omitempty" json:"x-codegen,omitempty"`

	// AllOf references (for composition) - using interface{} to handle both refs and inline schemas
	AllOf []interface{} `yaml:"allOf,omitempty" json:"allOf,omitempty"`

	// Default value for the entire object
	Default interface{} `yaml:"default,omitempty" json:"default,omitempty"`

	// Enum values if applicable
	Enum []interface{} `yaml:"enum,omitempty" json:"enum,omitempty"`
}

// Property represents a schema property with potential x-codegen extensions.
type Property struct {
	// Property type
	Type string `yaml:"type" json:"type"`

	// Format hint (e.g., "uuid", "email", "date-time")
	Format string `yaml:"format,omitempty" json:"format,omitempty"`

	// Property description
	Description string `yaml:"description,omitempty" json:"description,omitempty"`

	// Default value
	Default interface{} `yaml:"default,omitempty" json:"default,omitempty"`

	// Enum values
	Enum []interface{} `yaml:"enum,omitempty" json:"enum,omitempty"`

	// Reference to another schema
	Ref string `yaml:"$ref,omitempty" json:"$ref,omitempty"`

	// Array items type
	Items *Property `yaml:"items,omitempty" json:"items,omitempty"`

	// Nested object properties
	Properties map[string]Property `yaml:"properties,omitempty" json:"properties,omitempty"`

	// x-codegen extension at property level
	XCodegen *PropertyXCodegen `yaml:"x-codegen,omitempty" json:"x-codegen,omitempty"`

	// Required fields for object types
	Required []string `yaml:"required,omitempty" json:"required,omitempty"`
}

// PropertyXCodegen represents x-codegen at the property level.
type PropertyXCodegen struct {
	// Create unique constraint
	Unique *bool `yaml:"unique,omitempty" json:"unique,omitempty"`

	// Create database index
	Index *bool `yaml:"index,omitempty" json:"index,omitempty"`

	// Field is searchable (full-text search)
	Searchable *bool `yaml:"searchable,omitempty" json:"searchable,omitempty"`

	// Custom validation rule
	Validation *XCodegenValidation `yaml:"validation,omitempty" json:"validation,omitempty"`

	// Mark as primary key (legacy field)
	PrimaryKey bool `yaml:"primary-key,omitempty" json:"primary-key,omitempty"`

	// Field is immutable after creation
	Immutable bool `yaml:"immutable,omitempty" json:"immutable,omitempty"`

	// Database column name (if different from property name)
	ColumnName string `yaml:"column-name,omitempty" json:"column-name,omitempty"`

	// Default value expression
	DefaultValue string `yaml:"default-value,omitempty" json:"default-value,omitempty"`

	// Auto-generate value (e.g., "uuid", "timestamp")
	AutoGenerate string `yaml:"auto-generate,omitempty" json:"auto-generate,omitempty"`
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
	// Create base parsed schema
	parsed := &ParsedSchema{
		Schema: Schema{
			Name: name,
		},
		Domain:     p.inferDomain("", name),
		SourceFile: "openapi.yaml", // Since it's from the bundled spec
		ParsedAt:   time.Now(),
		Warnings:   []string{},
		PrimaryKey: "Id", // Default primary key field
		Type:       name, // For template compatibility
	}

	// Extract basic schema info
	if schema.Type != nil {
		types := schema.GetType()
		if len(types) > 0 {
			parsed.Schema.Type = string(types[0])
		}
	}

	if schema.Description != nil {
		parsed.Description = *schema.Description
	}

	parsed.Required = schema.Required

	// Extract x-codegen extension
	if schema.Extensions != nil {
		xcodegen := p.extractXCodegen(schema.Extensions)
		if xcodegen != nil {
			parsed.XCodegen = xcodegen

			// Extract events if present
			if len(xcodegen.Events) > 0 {
				parsed.Events = make([]Event, 0, len(xcodegen.Events))
				for _, eventName := range xcodegen.Events {
					parsed.Events = append(parsed.Events, Event{
						Type:        eventName,
						Description: fmt.Sprintf("%s event", eventName),
					})
				}
			}
		}
	}

	// Parse properties
	if schema.Properties != nil {
		parsed.Properties = make(map[string]Property)
		for propName := range schema.Properties.Keys() {
			propRef := schema.Properties.GetOrZero(propName)
			if propRef != nil && propRef.IsLeft() {
				prop := propRef.GetLeft()
				parsed.Properties[propName] = p.parseProperty(propName, prop)
			}
		}

		// Properties are already stored at top level
	}

	// Extract default values if present
	if schema.Default != nil {
		var defaultValue interface{}
		if err := schema.Default.Decode(&defaultValue); err == nil {
			parsed.Default = defaultValue
		}
	}

	// Extract enum values if present
	if len(schema.Enum) > 0 {
		parsed.Enum = make([]interface{}, len(schema.Enum))
		for i, enumVal := range schema.Enum {
			var value interface{}
			if err := enumVal.Decode(&value); err == nil {
				parsed.Enum[i] = value
			} else {
				parsed.Enum[i] = enumVal
			}
		}
	}

	return parsed
}

// parseProperty converts a Speakeasy property to our Property format.
func (p *Parser) parseProperty(_ string, prop *oas3.Schema) Property {
	result := Property{}

	// Extract type
	if prop.Type != nil {
		types := prop.GetType()
		if len(types) > 0 {
			result.Type = string(types[0])
		}
	}

	// Extract format
	if prop.Format != nil {
		result.Format = *prop.Format
	}

	// Extract description
	if prop.Description != nil {
		result.Description = *prop.Description
	}

	// Extract default value
	if prop.Default != nil {
		var defaultValue interface{}
		if err := prop.Default.Decode(&defaultValue); err == nil {
			result.Default = defaultValue
		}
	}

	// Extract enum values
	if len(prop.Enum) > 0 {
		result.Enum = make([]interface{}, len(prop.Enum))
		for i, enumVal := range prop.Enum {
			var value interface{}
			if err := enumVal.Decode(&value); err == nil {
				result.Enum[i] = value
			}
		}
	}

	// Extract x-codegen extension at property level
	if prop.Extensions != nil {
		xcodegen := p.extractPropertyXCodegen(prop.Extensions)
		if xcodegen != nil {
			result.XCodegen = xcodegen
		}
	}

	// Handle array items
	if result.Type == "array" && prop.Items != nil && prop.Items.IsLeft() {
		itemSchema := prop.Items.GetLeft()
		itemProp := p.parseProperty("item", itemSchema)
		result.Items = &itemProp
	}

	// Handle nested object properties
	if result.Type == "object" && prop.Properties != nil {
		result.Properties = make(map[string]Property)
		for subPropName := range prop.Properties.Keys() {
			subPropRef := prop.Properties.GetOrZero(subPropName)
			if subPropRef != nil && subPropRef.IsLeft() {
				subProp := subPropRef.GetLeft()
				result.Properties[subPropName] = p.parseProperty(subPropName, subProp)
			}
		}
	}

	return result
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

// extractPropertyXCodegen extracts property-level x-codegen extension.
func (p *Parser) extractPropertyXCodegen(ext *extensions.Extensions) *PropertyXCodegen {
	raw, err := extensions.GetExtensionValue[interface{}](ext, "x-codegen")
	if err != nil || raw == nil {
		return nil
	}

	// Marshal to JSON then unmarshal to our type
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return nil
	}

	var xcodegen PropertyXCodegen
	if err := json.Unmarshal(jsonBytes, &xcodegen); err != nil {
		return nil
	}

	return &xcodegen
}

// ParseWithTags parses schemas associated with specific OpenAPI tags.
func (p *Parser) ParseWithTags(specPath string, tags []string) (map[string]*ParsedSchema, error) {
	// First parse the entire spec
	allSchemas, err := p.ParseOpenAPISpec(specPath)
	if err != nil {
		return nil, err
	}

	// If no tags specified, return all schemas
	if len(tags) == 0 {
		return allSchemas, nil
	}

	// Filter schemas by tags
	// Note: This would need to be enhanced to actually check which schemas
	// are used by operations with the specified tags
	filtered := make(map[string]*ParsedSchema)
	for name, schema := range allSchemas {
		// For now, use domain inference as a proxy for tags
		for _, tag := range tags {
			if strings.EqualFold(schema.Domain, tag) {
				filtered[name] = schema
				break
			}
		}
	}

	return filtered, nil
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

// WalkAllSchemas walks through all schemas including nested ones.
func (p *Parser) WalkAllSchemas(callback func(name string, schema *ParsedSchema) error) error {
	if p.doc == nil {
		return fmt.Errorf("no document loaded")
	}

	ctx := context.Background()
	for item := range openapi.Walk(ctx, p.doc) {
		err := item.Match(openapi.Matcher{
			Schema: func(schema *oas3.JSONSchema[oas3.Referenceable]) error {
				if schema.IsLeft() {
					// Process the schema
					// Note: We don't have the name here, would need to track context
					return nil
				}
				return nil
			},
		})
		if err != nil {
			return err
		}
	}

	// Also walk our cached schemas
	for name, schema := range p.schemas {
		if err := callback(name, schema); err != nil {
			return err
		}
	}

	return nil
}

// GetDefaultValues extracts all default values from a schema.
func (p *Parser) GetDefaultValues(schemaName string) (map[string]interface{}, error) {
	schema, exists := p.schemas[schemaName]
	if !exists {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}

	defaults := make(map[string]interface{})

	// Extract defaults from properties
	for propName, prop := range schema.Properties {
		if prop.Default != nil {
			defaults[propName] = prop.Default
		}
	}

	return defaults, nil
}

// GetAllConfigDefaults recursively extracts all defaults including from nested schemas.
// This is especially useful for ArchesConfig which references other config schemas.
func (p *Parser) GetAllConfigDefaults(schemaName string) (map[string]interface{}, error) {
	if p.doc == nil || p.doc.Components == nil || p.doc.Components.Schemas == nil {
		return nil, fmt.Errorf("no document loaded")
	}

	schemaRef := p.doc.Components.Schemas.GetOrZero(schemaName)
	if schemaRef == nil || !schemaRef.IsLeft() {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}

	schema := schemaRef.GetLeft()
	return p.extractDefaultsRecursive(schema, schemaName)
}

// extractDefaultsRecursive recursively extracts defaults from a schema and its references.
func (p *Parser) extractDefaultsRecursive(schema *oas3.Schema, path string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if schema.Properties != nil {
		for propName := range schema.Properties.Keys() {
			propRef := schema.Properties.GetOrZero(propName)
			if propRef == nil {
				continue
			}

			if propRef.IsLeft() {
				// Direct schema
				prop := propRef.GetLeft()

				// Check for default value
				if prop.Default != nil {
					var defaultValue interface{}
					if err := prop.Default.Decode(&defaultValue); err == nil {
						result[propName] = defaultValue
					}
				}

				// If it's an object with properties, recurse
				if prop.Type != nil {
					types := prop.GetType()
					if len(types) > 0 && types[0] == "object" && prop.Properties != nil {
						// Recursively get defaults from nested object
						nested, err := p.extractDefaultsRecursive(prop, path+"."+propName)
						if err == nil && len(nested) > 0 {
							// Store nested defaults as a map
							result[propName] = nested
						}
					}
				}
			}
			// TODO: Handle references (IsRight) with Speakeasy's API
			// For now, we're only handling direct schemas
		}
	}

	return result, nil
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
	return len(schema.XCodegen.Events) > 0 || len(schema.Events) > 0
}

// NeedsAdapter checks if a schema needs adapter generation.
func NeedsAdapter(schema *ParsedSchema) bool {
	if schema == nil || schema.XCodegen == nil {
		return false
	}
	// Check if adapter is configured (has mappers or custom mappings)
	return schema.XCodegen.Adapter.GenerateMappers ||
		len(schema.XCodegen.Adapter.CustomMappings) > 0
}
