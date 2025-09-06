// Package repository contains OpenAPI parsing logic for x-codegen extensions
package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// XCodegenConfig represents the x-codegen extension in OpenAPI schemas
type XCodegenConfig struct {
	Repository *RepositoryConfig `yaml:"repository,omitempty"`
	Cache      *CacheConfig      `yaml:"cache,omitempty"`
	Events     []string          `yaml:"events,omitempty"`
	Service    *ServiceConfig    `yaml:"service,omitempty"`
}

// RepositoryConfig defines repository generation settings
type RepositoryConfig struct {
	Operations        []string                 `yaml:"operations"`
	Indices           []string                 `yaml:"indices"`
	AdditionalMethods []AdditionalMethodConfig `yaml:"additional_methods,omitempty"`
}

// AdditionalMethodConfig defines custom repository methods
type AdditionalMethodConfig struct {
	Name    string   `yaml:"name"`
	Params  []string `yaml:"params"`
	Returns string   `yaml:"returns"`
}

// CacheConfig defines cache generation settings
type CacheConfig struct {
	Enabled      bool     `yaml:"enabled"`
	TTL          int      `yaml:"ttl"`
	KeyPattern   string   `yaml:"key_pattern"`
	InvalidateOn []string `yaml:"invalidate_on"`
}

// ServiceConfig defines service layer generation settings
type ServiceConfig struct {
	BusinessMethods []string `yaml:"business_methods"`
}

// SchemaDefinition represents an OpenAPI schema with x-codegen
type SchemaDefinition struct {
	Description string                   `yaml:"description"`
	XCodegen    *XCodegenConfig          `yaml:"x-codegen,omitempty"`
	AllOf       []map[string]interface{} `yaml:"allOf,omitempty"`
	Type        string                   `yaml:"type,omitempty"`
	Properties  map[string]PropertyDef   `yaml:"properties,omitempty"`
	Required    []string                 `yaml:"required,omitempty"`
}

// PropertyDef represents a property definition with x-codegen
type PropertyDef struct {
	Description string                 `yaml:"description"`
	Type        string                 `yaml:"type"`
	Format      string                 `yaml:"format,omitempty"`
	XCodegen    *PropertyCodegenConfig `yaml:"x-codegen,omitempty"`
}

// PropertyCodegenConfig defines property-level codegen settings
type PropertyCodegenConfig struct {
	PrimaryKey bool   `yaml:"primary-key,omitempty"`
	Unique     bool   `yaml:"unique,omitempty"`
	Index      bool   `yaml:"index,omitempty"`
	Searchable bool   `yaml:"searchable,omitempty"`
	Validation string `yaml:"validation,omitempty"`
}

// OpenAPISpec represents the full OpenAPI specification
type OpenAPISpec struct {
	OpenAPI    string                 `yaml:"openapi"`
	Components ComponentsSpec         `yaml:"components"`
	Paths      map[string]interface{} `yaml:"paths"`
}

// ComponentsSpec represents the components section
type ComponentsSpec struct {
	Schemas map[string]SchemaDefinition `yaml:"schemas"`
}

// ParseOpenAPIWithXCodegen parses an OpenAPI spec file and extracts x-codegen extensions
func ParseOpenAPIWithXCodegen(specPath string) (*OpenAPISpec, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	return &spec, nil
}

// ParseSchemaFile parses a single schema file with x-codegen extensions
func ParseSchemaFile(schemaPath string) (*SchemaDefinition, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var schema SchemaDefinition
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return &schema, nil
}

// GetEntitiesWithXCodegen extracts all entities with x-codegen from an OpenAPI spec
func GetEntitiesWithXCodegen(spec *OpenAPISpec, tags []string) map[string]*EntityWithCodegen {
	entities := make(map[string]*EntityWithCodegen)

	// Filter schemas by tags if needed
	// For now, we'll look for entities that have x-codegen.repository defined
	for name, schema := range spec.Components.Schemas {
		if schema.XCodegen != nil && schema.XCodegen.Repository != nil {
			entity := &EntityWithCodegen{
				Name:     name,
				Schema:   schema,
				XCodegen: schema.XCodegen,
			}
			entities[name] = entity
		}
	}

	return entities
}

// EntityWithCodegen represents an entity with its x-codegen configuration
type EntityWithCodegen struct {
	Name     string
	Schema   SchemaDefinition
	XCodegen *XCodegenConfig
}

// ParseSchemaDirectory parses all schema files in a directory
func ParseSchemaDirectory(dir string) (map[string]*SchemaDefinition, error) {
	schemas := make(map[string]*SchemaDefinition)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema directory: %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		// Skip non-entity files
		if strings.Contains(file.Name(), "Request") || strings.Contains(file.Name(), "Response") {
			continue
		}

		schemaPath := filepath.Join(dir, file.Name())
		schema, err := ParseSchemaFile(schemaPath)
		if err != nil {
			// Log but continue with other schemas
			fmt.Printf("Warning: failed to parse %s: %v\n", file.Name(), err)
			continue
		}

		// Extract entity name from filename
		entityName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		schemas[entityName] = schema
	}

	return schemas, nil
}

// GenerateRepositoryMethods generates method signatures based on x-codegen config
func GenerateRepositoryMethods(entity *EntityWithCodegen) []MethodSignature {
	var methods []MethodSignature

	if entity.XCodegen.Repository == nil {
		return methods
	}

	repo := entity.XCodegen.Repository

	// Generate standard CRUD operations
	for _, op := range repo.Operations {
		switch op {
		case "create":
			methods = append(methods, MethodSignature{
				Name:    fmt.Sprintf("Create%s", entity.Name),
				Params:  []string{"ctx context.Context", fmt.Sprintf("entity *%s", entity.Name)},
				Returns: []string{fmt.Sprintf("*%s", entity.Name), "error"},
			})
		case "read":
			methods = append(methods, MethodSignature{
				Name:    fmt.Sprintf("Get%sByID", entity.Name),
				Params:  []string{"ctx context.Context", "id uuid.UUID"},
				Returns: []string{fmt.Sprintf("*%s", entity.Name), "error"},
			})
		case "update":
			methods = append(methods, MethodSignature{
				Name:    fmt.Sprintf("Update%s", entity.Name),
				Params:  []string{"ctx context.Context", "id uuid.UUID", fmt.Sprintf("entity *%s", entity.Name)},
				Returns: []string{fmt.Sprintf("*%s", entity.Name), "error"},
			})
		case "delete":
			methods = append(methods, MethodSignature{
				Name:    fmt.Sprintf("Delete%s", entity.Name),
				Params:  []string{"ctx context.Context", "id uuid.UUID"},
				Returns: []string{"error"},
			})
		case "list":
			methods = append(methods, MethodSignature{
				Name:    fmt.Sprintf("List%ss", entity.Name),
				Params:  []string{"ctx context.Context", fmt.Sprintf("params List%ssParams", entity.Name)},
				Returns: []string{fmt.Sprintf("[]*%s", entity.Name), "int64", "error"},
			})
		}
	}

	// Generate additional methods
	for _, method := range repo.AdditionalMethods {
		sig := MethodSignature{
			Name:   method.Name,
			Params: []string{"ctx context.Context"},
		}

		// Add parameters
		for _, param := range method.Params {
			paramType := "string" // Default type
			if param == "id" || strings.HasSuffix(param, "Id") || strings.HasSuffix(param, "ID") {
				paramType = "uuid.UUID"
			}
			sig.Params = append(sig.Params, fmt.Sprintf("%s %s", param, paramType))
		}

		// Add returns
		switch method.Returns {
		case "single":
			sig.Returns = []string{fmt.Sprintf("*%s", entity.Name), "error"}
		case "multiple":
			sig.Returns = []string{fmt.Sprintf("[]*%s", entity.Name), "error"}
		case "void":
			sig.Returns = []string{"error"}
		default:
			sig.Returns = []string{"error"}
		}

		methods = append(methods, sig)
	}

	return methods
}

// MethodSignature represents a Go method signature
type MethodSignature struct {
	Name    string
	Params  []string
	Returns []string
}
