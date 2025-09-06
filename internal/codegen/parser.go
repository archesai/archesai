// Package codegen provides schema parsing and type definitions for x-codegen extensions in OpenAPI schemas.
package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Legacy type aliases for backward compatibility - removed as they're now generated in models.gen.go

// PropertyXCodegen represents x-codegen at the property level.
// This remains separate as it's for individual field-level configuration.
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

	// Minimum value (for numbers)
	Minimum *float64 `yaml:"minimum,omitempty" json:"minimum,omitempty"`

	// Maximum value (for numbers)
	Maximum *float64 `yaml:"maximum,omitempty" json:"maximum,omitempty"`

	// Minimum length (for strings)
	MinLength *int `yaml:"minLength,omitempty" json:"minLength,omitempty"`

	// Maximum length (for strings)
	MaxLength *int `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`

	// Pattern (regex for strings)
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty"`

	// x-codegen extension at property level
	XCodegen *PropertyXCodegen `yaml:"x-codegen,omitempty" json:"x-codegen,omitempty"`
}

// Ref represents a reference to another schema.
type Ref struct {
	Ref string `yaml:"$ref,omitempty" json:"$ref,omitempty"`
}

// ParsedSchema represents a fully parsed schema ready for code generation.
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

// GeneratorConfig represents configuration for a specific generator.
type GeneratorConfig struct {
	// Generator name
	Name string

	// Enabled flag
	Enabled bool

	// Output directory
	OutputDir string

	// Template overrides
	Templates map[string]string

	// Custom configuration
	Config map[string]interface{}
}

// ParseConfig represents the overall codegen configuration.
type ParseConfig struct {
	// OpenAPI spec file path
	OpenAPIFile string `yaml:"openapi" json:"openapi"`

	// Domains to generate
	Domains map[string]DomainConfig `yaml:"domains" json:"domains"`

	// Global settings
	Settings GlobalSettings `yaml:"settings" json:"settings"`

	// Generator-specific configurations
	Generators map[string]GeneratorConfig `yaml:"generators,omitempty" json:"generators,omitempty"`
}

// Parser handles parsing of OpenAPI schemas with x-codegen extensions.
type Parser struct {
	// Cache of parsed schemas
	schemas map[string]*ParsedSchema

	// Base directory for resolving references
	baseDir string

	// Warnings accumulated during parsing
	warnings []string
}

// NewParser creates a new schema parser.
func NewParser(baseDir string) *Parser {
	return &Parser{
		schemas:  make(map[string]*ParsedSchema),
		baseDir:  baseDir,
		warnings: []string{},
	}
}

// ParseFile parses a single OpenAPI schema file.
func (p *Parser) ParseFile(filePath string) (*ParsedSchema, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var schema Schema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", filePath, err)
	}

	// Extract schema name from filename
	baseName := filepath.Base(filePath)
	schemaName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	schema.Name = schemaName

	// Infer domain from path or schema name
	domain := p.inferDomain(filePath, schemaName)

	parsed := &ParsedSchema{
		Schema:     schema,
		Domain:     domain,
		SourceFile: filePath,
		ParsedAt:   time.Now(),
		Warnings:   []string{},
		PrimaryKey: "Id",        // Default primary key field
		Type:       schema.Name, // For template compatibility
	}

	// Validate and enhance x-codegen configuration
	p.validateXCodegen(parsed)

	// Cache the parsed schema
	p.schemas[schemaName] = parsed

	return parsed, nil
}

// ParseDirectory parses all schema files in a directory.
func (p *Parser) ParseDirectory(dir string) (map[string]*ParsedSchema, error) {
	schemas := make(map[string]*ParsedSchema)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process YAML files
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		// Skip non-schema files
		if strings.Contains(entry.Name(), "Request") || strings.Contains(entry.Name(), "Response") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		schema, err := p.ParseFile(filePath)
		if err != nil {
			p.warnings = append(p.warnings, fmt.Sprintf("Failed to parse %s: %v", entry.Name(), err))
			continue
		}

		schemas[schema.Name] = schema
	}

	return schemas, nil
}

// ParseOpenAPISpec parses a complete OpenAPI specification file.
func (p *Parser) ParseOpenAPISpec(specPath string) (map[string]*ParsedSchema, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	var spec struct {
		OpenAPI    string `yaml:"openapi"`
		Components struct {
			Schemas map[string]Schema `yaml:"schemas"`
		} `yaml:"components"`
		Paths map[string]interface{} `yaml:"paths"`
	}

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	schemas := make(map[string]*ParsedSchema)

	for name, schema := range spec.Components.Schemas {
		schema.Name = name
		domain := p.inferDomain(specPath, name)

		// Resolve allOf to merge properties
		p.resolveAllOf(&schema, spec.Components.Schemas)

		parsed := &ParsedSchema{
			Schema:     schema,
			Domain:     domain,
			SourceFile: specPath,
			ParsedAt:   time.Now(),
			Warnings:   []string{},
			PrimaryKey: "Id",        // Default primary key field
			Type:       schema.Name, // For template compatibility
		}

		p.validateXCodegen(parsed)
		schemas[name] = parsed
		p.schemas[name] = parsed
	}

	return schemas, nil
}

// ParseWithTags parses schemas filtered by OpenAPI tags.
func (p *Parser) ParseWithTags(specPath string, tags []string) (map[string]*ParsedSchema, error) {
	// First parse all schemas
	allSchemas, err := p.ParseOpenAPISpec(specPath)
	if err != nil {
		return nil, err
	}

	// Filter by tags if provided
	if len(tags) == 0 {
		return allSchemas, nil
	}

	// Create a tag set for quick lookup
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[tag] = true
	}

	// Filter schemas based on tags
	filtered := make(map[string]*ParsedSchema)
	for name, schema := range allSchemas {
		// Check if schema name matches any tag pattern
		if p.matchesTags(name, tagSet) {
			filtered[name] = schema
		}
	}

	return filtered, nil
}

// GetWarnings returns any warnings accumulated during parsing.
func (p *Parser) GetWarnings() []string {
	return p.warnings
}

// ResolveRef resolves a $ref to its schema.
func (p *Parser) ResolveRef(ref string) (*ParsedSchema, error) {
	// Remove the file path prefix if present
	if strings.Contains(ref, "#/") {
		parts := strings.Split(ref, "#/")
		ref = parts[len(parts)-1]
	}

	// Remove components/schemas/ prefix
	ref = strings.TrimPrefix(ref, "components/schemas/")

	// Check cache
	if schema, ok := p.schemas[ref]; ok {
		return schema, nil
	}

	// Try to load from file
	possiblePaths := []string{
		filepath.Join(p.baseDir, "api", "components", "schemas", ref+".yaml"),
		filepath.Join(p.baseDir, ref+".yaml"),
		filepath.Join(p.baseDir, ref+".yml"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return p.ParseFile(path)
		}
	}

	return nil, fmt.Errorf("cannot resolve reference: %s", ref)
}

// validateXCodegen validates and enhances x-codegen configuration.
func (p *Parser) validateXCodegen(schema *ParsedSchema) {
	if schema.XCodegen == nil {
		return
	}

	// Validate repository configuration
	repo := schema.XCodegen.Repository
	// Ensure operations are valid
	validOps := map[string]bool{
		"create":      true,
		"read":        true,
		"update":      true,
		"delete":      true,
		"list":        true,
		"bulk_create": true,
		"bulk_update": true,
		"bulk_delete": true,
	}

	if len(repo.Operations) > 0 {
		for _, op := range repo.Operations {
			if !validOps[string(op)] {
				schema.Warnings = append(schema.Warnings,
					fmt.Sprintf("Invalid repository operation: %s", op))
			}
		}
	}

	// Set default table name if not specified
	// If Postgres config exists but has no table name, set it
	if repo.Postgres.TableName == "" {
		repo.Postgres.TableName = pluralize(strings.ToLower(schema.Name))
	} else if repo.Sqlite.TableName == "" {
		// If SQLite config exists but has no table name, set it
		repo.Sqlite.TableName = pluralize(strings.ToLower(schema.Name))
	} else if repo.Postgres.TableName == "" && repo.Sqlite.TableName == "" {
		// If neither Postgres nor SQLite has a table name, set default Postgres table name
		repo.Postgres.TableName = pluralize(strings.ToLower(schema.Name))
	}

	// Validate cache configuration
	cache := schema.XCodegen.Cache
	if cache.Enabled {
		if cache.Ttl <= 0 {
			cache.Ttl = 300 // Default 5 minutes
		}

		if cache.KeyPattern == "" {
			cache.KeyPattern = fmt.Sprintf("%s:{id}", strings.ToLower(schema.Name))
		}
	}

	// Validate handler configuration
	handler := &schema.XCodegen.Handler
	if handler.PathPrefix == "" {
		prefix := "/" + strings.ToLower(pluralize(schema.Name))
		handler.PathPrefix = prefix
	}

	// Process events
	if len(schema.XCodegen.Events) > 0 {
		for _, eventName := range schema.XCodegen.Events {
			event := Event{
				Type:        CamelCase(schema.Name) + CamelCase(eventName),
				Description: fmt.Sprintf("%s %s event", schema.Name, eventName),
			}
			schema.Events = append(schema.Events, event)
		}
	}

	// Validate property-level x-codegen
	// Sort property keys for consistent processing
	propKeys := GetSortedPropertyKeys(schema.Properties)
	for _, propName := range propKeys {
		prop := schema.Properties[propName]
		if prop.XCodegen != nil {
			// Validate primary key
			if prop.XCodegen.PrimaryKey {
				// Ensure it's indexed
				if prop.XCodegen.Index == nil || !*prop.XCodegen.Index {
					trueVal := true
					prop.XCodegen.Index = &trueVal
				}
				// Ensure it's unique
				if prop.XCodegen.Unique == nil || !*prop.XCodegen.Unique {
					trueVal := true
					prop.XCodegen.Unique = &trueVal
				}
			}

			// Set default column name
			if prop.XCodegen.ColumnName == "" {
				prop.XCodegen.ColumnName = SnakeCase(propName)
			}
			// Update the property back in the map
			schema.Properties[propName] = prop
		}
	}
}

// resolveAllOf resolves allOf references and merges properties into the schema.
func (p *Parser) resolveAllOf(schema *Schema, _ map[string]Schema) {
	if len(schema.AllOf) == 0 {
		return
	}

	// Initialize properties map if nil
	if schema.Properties == nil {
		schema.Properties = make(map[string]Property)
	}

	// Process each item in allOf
	for _, item := range schema.AllOf {
		// Try to convert to a map (inline schema)
		if itemMap, ok := item.(map[string]interface{}); ok {
			// Check if this has properties
			if props, ok := itemMap["properties"].(map[string]interface{}); ok {
				// Merge properties
				for propName, propValue := range props {
					// Convert property value to Property struct
					propData, _ := yaml.Marshal(propValue)
					var prop Property
					if err := yaml.Unmarshal(propData, &prop); err == nil {
						schema.Properties[propName] = prop
					}
				}
			}

			// Also check for required fields
			if required, ok := itemMap["required"].([]interface{}); ok {
				for _, req := range required {
					if reqStr, ok := req.(string); ok {
						schema.Required = append(schema.Required, reqStr)
					}
				}
			}
		}
		// If it's a reference, we skip it for now (would need to resolve BaseEntity etc)
	}
}

// inferDomain infers the domain from file path or schema name.
func (p *Parser) inferDomain(filePath, schemaName string) string {
	// Try to infer from file path first
	if strings.Contains(filePath, "/auth/") || strings.Contains(filePath, "\\auth\\") {
		return "auth"
	}
	if strings.Contains(filePath, "/organizations/") || strings.Contains(filePath, "\\organizations\\") {
		return "organizations"
	}
	if strings.Contains(filePath, "/workflows/") || strings.Contains(filePath, "\\workflows\\") {
		return "workflows"
	}
	if strings.Contains(filePath, "/content/") || strings.Contains(filePath, "\\content\\") {
		return "content"
	}

	// Infer from schema name
	schemaLower := strings.ToLower(schemaName)
	switch {
	case strings.Contains(schemaLower, "user") ||
		strings.Contains(schemaLower, "session") ||
		strings.Contains(schemaLower, "account") ||
		strings.Contains(schemaLower, "auth"):
		return "auth"
	case strings.Contains(schemaLower, "organization") ||
		strings.Contains(schemaLower, "member") ||
		strings.Contains(schemaLower, "team"):
		return "organizations"
	case strings.Contains(schemaLower, "workflow") ||
		strings.Contains(schemaLower, "pipeline") ||
		strings.Contains(schemaLower, "run"):
		return "workflows"
	case strings.Contains(schemaLower, "content") ||
		strings.Contains(schemaLower, "block") ||
		strings.Contains(schemaLower, "page") ||
		strings.Contains(schemaLower, "artifact") ||
		strings.Contains(schemaLower, "label"):
		return "content"
	default:
		return "common"
	}
}

// matchesTags checks if a schema name matches any of the given tags.
func (p *Parser) matchesTags(schemaName string, tags map[string]bool) bool {
	// Direct tag match
	if tags[schemaName] {
		return true
	}

	// Check common patterns
	patterns := map[string][]string{
		"Auth":          {"User", "Session", "Account", "Token"},
		"Users":         {"User", "UserProfile", "UserSettings"},
		"Sessions":      {"Session", "SessionToken"},
		"Organizations": {"Organization", "Member", "Team"},
		"Workflows":     {"Workflow", "Pipeline", "Run", "Tool"},
		"Content":       {"Content", "Block", "Page", "Artifact", "Label"},
	}

	for tag, schemaPatterns := range patterns {
		if tags[tag] {
			for _, pattern := range schemaPatterns {
				if strings.Contains(schemaName, pattern) {
					return true
				}
			}
		}
	}

	return false
}

// pluralize converts a singular word to plural (simple rules).
func pluralize(word string) string {
	if strings.HasSuffix(word, "y") && len(word) > 1 {
		// Check if the letter before 'y' is a vowel
		beforeY := word[len(word)-2]
		if beforeY != 'a' && beforeY != 'e' && beforeY != 'i' && beforeY != 'o' && beforeY != 'u' {
			return word[:len(word)-1] + "ies"
		}
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
		strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}

// GetSchemasByDomain returns all schemas grouped by domain.
func (p *Parser) GetSchemasByDomain() map[string][]*ParsedSchema {
	domains := make(map[string][]*ParsedSchema)

	for _, schema := range p.schemas {
		domains[schema.Domain] = append(domains[schema.Domain], schema)
	}

	return domains
}

// GetSchema returns a specific parsed schema by name.
func (p *Parser) GetSchema(name string) (*ParsedSchema, bool) {
	schema, ok := p.schemas[name]
	return schema, ok
}

// HasXCodegen checks if a schema has any x-codegen configuration.
func HasXCodegen(schema *ParsedSchema) bool {
	if schema.XCodegen != nil {
		return true
	}

	for _, prop := range schema.Properties {
		if prop.XCodegen != nil {
			return true
		}
	}

	return false
}

// NeedsRepository checks if a schema needs repository generation.
func NeedsRepository(schema *ParsedSchema) bool {
	return schema.XCodegen != nil &&
		len(schema.XCodegen.Repository.Operations) > 0
}

// NeedsCache checks if a schema needs cache generation.
func NeedsCache(schema *ParsedSchema) bool {
	return schema.XCodegen != nil &&
		schema.XCodegen.Cache.Enabled
}

// NeedsEvents checks if a schema needs event generation.
func NeedsEvents(schema *ParsedSchema) bool {
	return schema.XCodegen != nil &&
		len(schema.XCodegen.Events) > 0
}

// NeedsHandler checks if a schema needs handler generation.
func NeedsHandler(schema *ParsedSchema) bool {
	return schema.XCodegen != nil && schema.XCodegen.Handler.Generate
}

// NeedsAdapter checks if a schema needs adapter generation.
func NeedsAdapter(schema *ParsedSchema) bool {
	return schema.XCodegen != nil &&
		schema.XCodegen.Adapter.GenerateMappers
}
