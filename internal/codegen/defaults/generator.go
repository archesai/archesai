// Package defaults generates default configuration values from OpenAPI schemas.
package defaults

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const outputFile = "internal/infrastructure/config/defaults.gen.go"

//go:embed templates/config.go.tmpl
var configTemplate string

// Generator handles generation of config defaults.
type Generator struct{}

// NewGenerator creates a new defaults generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Schema represents an OpenAPI schema
type Schema struct {
	Type       string                   `yaml:"type"`
	Properties map[string]Property      `yaml:"properties"`
	Required   []string                 `yaml:"required"`
	Default    interface{}              `yaml:"default"`
	Enum       []interface{}            `yaml:"enum"`
	Ref        string                   `yaml:"$ref"`
	AllOf      []map[string]interface{} `yaml:"allOf"`
}

// Property represents a schema property
type Property struct {
	Type        string              `yaml:"type"`
	Default     interface{}         `yaml:"default"`
	Description string              `yaml:"description"`
	Ref         string              `yaml:"$ref"`
	Enum        []interface{}       `yaml:"enum"`
	Format      string              `yaml:"format"`
	Minimum     *float64            `yaml:"minimum"`
	Maximum     *float64            `yaml:"maximum"`
	MinLength   *int                `yaml:"minLength"`
	MaxLength   *int                `yaml:"maxLength"`
	Items       *Property           `yaml:"items"`
	Properties  map[string]Property `yaml:"properties"`
}

// Generate generates the config defaults code.
func (g *Generator) Generate() error {
	log.Println("Parsing OpenAPI schemas for default values...")

	// Parse all schema files
	schemas, err := g.loadSchemas("api/components/schemas")
	if err != nil {
		return fmt.Errorf("failed to load schemas: %w", err)
	}

	// Generate the Go code from parsed schemas
	code, err := g.generateCode(schemas)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// Format the code
	formatted, err := format.Source([]byte(code))
	if err != nil {
		log.Printf("Warning: Failed to format code: %v", err)
		// Write unformatted for debugging
		if err := os.WriteFile(outputFile+".debug", []byte(code), 0644); err != nil {
			log.Printf("Failed to write debug file: %v", err)
		}
		return fmt.Errorf("code formatting failed: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputFile, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	log.Printf("Successfully generated %s", outputFile)
	return nil
}

// loadSchemas loads all OpenAPI schema files from a directory
func (g *Generator) loadSchemas(dir string) (map[string]Schema, error) {
	schemas := make(map[string]Schema)

	err := filepath.WalkDir(dir, func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		var schema Schema
		if err := yaml.Unmarshal(data, &schema); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		// Use filename without extension as schema name
		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		schemas[name] = schema

		return nil
	})

	return schemas, err
}

// generateCode generates Go code from parsed schemas
func (g *Generator) generateCode(schemas map[string]Schema) (string, error) {
	// Get ArchesConfig schema and generate its fields
	configFields := ""
	if archesConfig, ok := schemas["ArchesConfig"]; ok {
		configFields = g.generateStructFields(archesConfig, schemas, 2)
	}

	// Parse and execute template
	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"ConfigFields": configFields,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateStructFields generates field assignments for a struct
func (g *Generator) generateStructFields(schema Schema, allSchemas map[string]Schema, indent int) string {
	var lines []string
	indentStr := strings.Repeat("\t", indent)

	// Sort properties for consistent output
	var propNames []string
	for name := range schema.Properties {
		propNames = append(propNames, name)
	}
	sort.Strings(propNames)

	for _, propName := range propNames {
		prop := schema.Properties[propName]
		fieldName := toGoFieldName(propName)

		// Handle different property types
		if prop.Ref != "" {
			// Reference to another schema
			refName := filepath.Base(prop.Ref)
			refName = strings.TrimSuffix(refName, ".yaml")

			if refSchema, ok := allSchemas[refName]; ok {
				if refSchema.Type == "object" || refSchema.Type == "" {
					// Object type - generate nested struct
					lines = append(lines, fmt.Sprintf("%s%s: api.%s{", indentStr, fieldName, refName))

					nestedFields := g.generateStructFields(refSchema, allSchemas, indent+1)
					if nestedFields != "" {
						lines = append(lines, nestedFields)
					}
					lines = append(lines, fmt.Sprintf("%s},", indentStr))
				}
			}
		} else if prop.Default != nil {
			// Property with default value
			if prop.Type == "array" {
				// Arrays always need explicit type, even with default
				lines = append(lines, fmt.Sprintf("%s%s: []string{},", indentStr, fieldName))
			} else {
				value := formatValue(prop.Default, prop.Type, propName, isRequired(propName, schema.Required), prop.Enum)
				lines = append(lines, fmt.Sprintf("%s%s: %s,", indentStr, fieldName, value))
			}
		} else if prop.Type == "object" && prop.Properties != nil {
			// Inline object
			lines = append(lines, fmt.Sprintf("%s%s: {", indentStr, fieldName))

			// Create a temporary schema for the inline object
			tempSchema := Schema{
				Type:       "object",
				Properties: prop.Properties,
			}
			nestedFields := g.generateStructFields(tempSchema, allSchemas, indent+1)
			if nestedFields != "" {
				lines = append(lines, nestedFields)
			}
			lines = append(lines, fmt.Sprintf("%s},", indentStr))
		} else if prop.Type == "array" {
			// Array with no default - still need to initialize
			lines = append(lines, fmt.Sprintf("%s%s: []string{},", indentStr, fieldName))
		}
	}

	return strings.Join(lines, "\n")
}

// toGoFieldName converts a JSON field name to Go field name
func toGoFieldName(name string) string {
	// Handle special cases
	switch name {
	case "cors":
		return "Cors"
	case "llm":
		return "Llm"
	case "url":
		return "Url"
	}

	parts := strings.Split(name, "_")
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.Title(part) //nolint:staticcheck
		}
	}
	return strings.Join(parts, "")
}

// isRequired checks if a property is in the required list
func isRequired(name string, required []string) bool {
	for _, r := range required {
		if r == name {
			return true
		}
	}
	return false
}

// formatValue formats a value for Go code
func formatValue(value interface{}, propType string, _ string, _ bool, _ []interface{}) string {
	switch v := value.(type) {
	case string:
		// With omitzero, all strings (including enums) are direct string values
		return fmt.Sprintf("\"%s\"", v)
	case bool:
		// With omitzero, bools are direct values
		return fmt.Sprintf("%v", v)
	case float64:
		switch propType {
		case "integer":
			return fmt.Sprintf("%d", int(v))
		case "number":
			// Handle float32
			return fmt.Sprintf("float32(%v)", v)
		}
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case []interface{}:
		// Empty array default
		return "[]string{}"
	default:
		return fmt.Sprintf("%v", v)
	}
}
