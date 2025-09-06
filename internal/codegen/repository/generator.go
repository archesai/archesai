// Package repository generates repository interfaces and implementations from OpenAPI specs.
package repository

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// Generator handles generation of repository code.
type Generator struct{}

// NewGenerator creates a new repository generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Config represents repository generation configuration.
type Config struct {
	Domains map[string]DomainConfig `yaml:"domains"`
}

// DomainConfig represents configuration for a single domain.
type DomainConfig struct {
	OpenAPI  string                  `yaml:"openapi"`
	Tags     []string                `yaml:"tags"`
	Storage  StorageConfig           `yaml:"storage"`
	Entities map[string]EntityConfig `yaml:"entities,omitempty"`
}

// StorageConfig represents storage adapter configuration.
type StorageConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	SQLite   SQLiteConfig   `yaml:"sqlite"`
}

// PostgresConfig represents PostgreSQL configuration.
type PostgresConfig struct {
	Enabled bool   `yaml:"enabled"`
	Package string `yaml:"package"`
}

// SQLiteConfig represents SQLite configuration.
type SQLiteConfig struct {
	Enabled bool   `yaml:"enabled"`
	Package string `yaml:"package"`
}

// EntityConfig represents configuration for a single entity.
type EntityConfig struct {
	Table      string                 `yaml:"table"`
	Operations []string               `yaml:"operations"` // create, read, update, delete, list
	Indices    []string               `yaml:"indices"`
	Unique     []string               `yaml:"unique"`
	SoftDelete bool                   `yaml:"soft_delete"`
	Fields     map[string]FieldConfig `yaml:"fields,omitempty"`
}

// FieldConfig represents field-specific configuration.
type FieldConfig struct {
	PrimaryKey bool   `yaml:"primary_key"`
	Index      bool   `yaml:"index"`
	Unique     bool   `yaml:"unique"`
	Nullable   bool   `yaml:"nullable"`
	Type       string `yaml:"type,omitempty"`
}

// EntityInfo represents parsed entity information.
type EntityInfo struct {
	Name              string
	Type              string
	Table             string
	Operations        []string
	Fields            []FieldInfo
	PrimaryKey        string
	AdditionalMethods []MethodSignature // From x-codegen
	Indices           []string          // Fields to index from x-codegen
}

// FieldInfo represents field information.
type FieldInfo struct {
	Name       string
	Type       string
	GoType     string
	PrimaryKey bool
	Index      bool
	Unique     bool
	Nullable   bool
}

// TemplateData represents data passed to templates.
type TemplateData struct {
	Domain   string
	Package  string
	Entities []EntityInfo
	Imports  []string
	Config   DomainConfig
}

// Generate generates repository code from OpenAPI schemas with x-codegen extensions.
func (g *Generator) Generate(schemaDir string) error {
	// Find all OpenAPI schema files
	schemaFiles, err := g.findSchemaFiles(schemaDir)
	if err != nil {
		return fmt.Errorf("failed to find schema files: %w", err)
	}

	// Group schemas by domain (derived from directory structure)
	domainSchemas := make(map[string][]*EntityWithCodegen)

	for _, schemaFile := range schemaFiles {
		// Parse schema with x-codegen
		schema, err := ParseSchemaFile(schemaFile)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", schemaFile, err)
			continue
		}

		// Skip if no x-codegen.repository config
		if schema.XCodegen == nil || schema.XCodegen.Repository == nil {
			continue
		}

		// Extract entity name from filename
		baseName := filepath.Base(schemaFile)
		entityName := strings.TrimSuffix(baseName, ".yaml")
		entityName = strings.TrimSuffix(entityName, "Entity") // Remove Entity suffix if present

		// Determine domain from the entity name or tags
		domain := g.inferDomain(entityName)

		entity := &EntityWithCodegen{
			Name:     entityName,
			Schema:   *schema,
			XCodegen: schema.XCodegen,
		}

		domainSchemas[domain] = append(domainSchemas[domain], entity)
	}

	// Generate repository for each domain
	for domain, entities := range domainSchemas {
		if err := g.generateDomainRepository(domain, entities); err != nil {
			return fmt.Errorf("failed to generate repository for %s: %w", domain, err)
		}
	}

	return nil
}

// generateDomainRepository generates repository code for a single domain.
func (g *Generator) generateDomainRepository(domain string, entities []*EntityWithCodegen) error {
	if len(entities) == 0 {
		return nil // No entities to generate
	}

	// Convert entities to EntityInfo for template
	var entityInfos []EntityInfo
	for _, entity := range entities {
		entityInfo := g.convertToEntityInfo(entity)
		entityInfos = append(entityInfos, entityInfo)
	}

	// Template data
	data := TemplateData{
		Domain:   domain,
		Package:  domain,
		Entities: entityInfos,
		Imports: []string{
			"context",
			"github.com/google/uuid",
			"time",
		},
	}

	// Generate repository interface in flat structure
	if err := g.generateFile("repository.go.tmpl", filepath.Join("internal", domain, "repository.gen.go"), data); err != nil {
		return fmt.Errorf("failed to generate repository interface: %w", err)
	}

	// Always generate both PostgreSQL and SQLite implementations
	// Users can choose which to use at runtime

	// PostgreSQL implementation
	postgresData := data
	postgresData.Imports = append([]string{
		"database/sql",
		"fmt",
		"strings",
	}, data.Imports...)
	if err := g.generateFile("repository_postgres.go.tmpl", filepath.Join("internal", domain, "repository_postgres.gen.go"), postgresData); err != nil {
		return fmt.Errorf("failed to generate postgres repository: %w", err)
	}

	// SQLite implementation
	sqliteData := data
	sqliteData.Imports = append([]string{
		"database/sql",
		"fmt",
		"strings",
	}, data.Imports...)
	if err := g.generateFile("repository_sqlite.go.tmpl", filepath.Join("internal", domain, "repository_sqlite.gen.go"), sqliteData); err != nil {
		return fmt.Errorf("failed to generate sqlite repository: %w", err)
	}

	return nil
}

// findSchemaFiles finds all YAML schema files in the given directory.
func (g *Generator) findSchemaFiles(dir string) ([]string, error) {
	var schemaFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for Entity YAML files
		if !info.IsDir() && strings.HasSuffix(path, "Entity.yaml") {
			schemaFiles = append(schemaFiles, path)
		}

		return nil
	})

	return schemaFiles, err
}

// inferDomain infers the domain from an entity name.
func (g *Generator) inferDomain(entityName string) string {
	// Map entity names to domains
	domainMap := map[string]string{
		"User":         "auth",
		"Session":      "auth",
		"Account":      "auth",
		"Organization": "organizations",
		"Member":       "organizations",
		"Workflow":     "workflows",
		"Run":          "workflows",
		"Tool":         "workflows",
		"Pipeline":     "workflows",
		"Content":      "content",
		"Block":        "content",
	}

	if domain, ok := domainMap[entityName]; ok {
		return domain
	}

	// Default to lowercase entity name
	return strings.ToLower(entityName)
}

// convertToEntityInfo converts EntityWithCodegen to EntityInfo for templates.
func (g *Generator) convertToEntityInfo(entity *EntityWithCodegen) EntityInfo {
	info := EntityInfo{
		Name:       entity.Name,
		Type:       entity.Name + "Entity",
		Table:      g.toTableName(entity.Name),
		Operations: entity.XCodegen.Repository.Operations,
		PrimaryKey: "id", // Default
		Indices:    entity.XCodegen.Repository.Indices,
	}

	// Generate additional methods from x-codegen
	methods := GenerateRepositoryMethods(entity)
	// Filter out standard CRUD methods, keep only additional ones
	for _, method := range methods {
		isStandard := false
		for _, op := range entity.XCodegen.Repository.Operations {
			switch op {
			case "create":
				if method.Name == fmt.Sprintf("Create%s", entity.Name) {
					isStandard = true
				}
			case "read":
				if method.Name == fmt.Sprintf("Get%sByID", entity.Name) {
					isStandard = true
				}
			case "update":
				if method.Name == fmt.Sprintf("Update%s", entity.Name) {
					isStandard = true
				}
			case "delete":
				if method.Name == fmt.Sprintf("Delete%s", entity.Name) {
					isStandard = true
				}
			case "list":
				if method.Name == fmt.Sprintf("List%ss", entity.Name) {
					isStandard = true
				}
			}
		}
		if !isStandard {
			info.AdditionalMethods = append(info.AdditionalMethods, method)
		}
	}

	// Extract fields from schema properties
	for fieldName, prop := range entity.Schema.Properties {
		fieldInfo := FieldInfo{
			Name:   fieldName,
			Type:   prop.Type,
			GoType: g.mapGoType(prop),
		}

		// Check for primary key
		if fieldName == "id" {
			fieldInfo.PrimaryKey = true
			info.PrimaryKey = "id"
		}

		// Check for indices in x-codegen
		for _, index := range entity.XCodegen.Repository.Indices {
			if index == fieldName {
				fieldInfo.Index = true
			}
		}

		// Check property-level x-codegen
		if prop.XCodegen != nil {
			fieldInfo.Index = prop.XCodegen.Index
			fieldInfo.Unique = prop.XCodegen.Unique
			if prop.XCodegen.PrimaryKey {
				fieldInfo.PrimaryKey = true
				info.PrimaryKey = fieldName
			}
		}

		// Check if nullable (not in required list)
		fieldInfo.Nullable = !g.isRequired(fieldName, entity.Schema.Required)

		info.Fields = append(info.Fields, fieldInfo)
	}

	return info
}

// toTableName converts an entity name to a table name.
func (g *Generator) toTableName(entityName string) string {
	// Simple pluralization for now
	if strings.HasSuffix(entityName, "y") {
		return strings.ToLower(entityName[:len(entityName)-1] + "ies")
	}
	return strings.ToLower(entityName + "s")
}

// mapGoType maps OpenAPI types to Go types.
func (g *Generator) mapGoType(prop PropertyDef) string {
	switch prop.Type {
	case "string":
		if prop.Format == "uuid" {
			return "uuid.UUID"
		}
		if prop.Format == "date-time" {
			return "time.Time"
		}
		if prop.Format == "email" {
			return "string" // Could be openapi_types.Email
		}
		return "string"
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// isRequired checks if a field is in the required list.
func (g *Generator) isRequired(fieldName string, required []string) bool {
	for _, r := range required {
		if r == fieldName {
			return true
		}
	}
	return false
}

// generateFile generates a file from template.
func (g *Generator) generateFile(templateName, outputPath string, data TemplateData) error {
	// Read template
	tmplContent, err := templatesFS.ReadFile(filepath.Join("templates", templateName))
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	// Create template with helper functions
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return s
			}
			return string(unicode.ToUpper(rune(s[0]))) + s[1:]
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"join":  strings.Join,
		"contains": func(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		"camelCase": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	// Format Go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code for debugging
		formatted = buf.Bytes()
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}
