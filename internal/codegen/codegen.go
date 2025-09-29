// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"bytes"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/speakeasy-api/openapi/jsonschema/oas3"

	"github.com/archesai/archesai/internal/parsers"
)

// Run executes the code generator with the given OpenAPI path (simplified entry point)
func Run(openapiPath string) error {
	_, err := Generate(openapiPath, Configuration{
		SpecPath: openapiPath,
	})
	return err
}

// RunFromJSONSchema generates a Go struct from a standalone JSON Schema file
func RunFromJSONSchema(schemaPath, outputPath string) error {
	// Extract struct name from schema filename
	structName := strings.TrimSuffix(filepath.Base(schemaPath), filepath.Ext(schemaPath))

	// Extract package name from output path
	packageName := filepath.Base(filepath.Dir(outputPath))
	if packageName == "." || packageName == "" {
		return fmt.Errorf("cannot infer package name from output path %s", outputPath)
	}

	// Load templates
	loadedTemplates, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Load schema directly
	schema, xcodegen, err := parsers.ParseJSONSchema(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Set the title if not present (use structName)
	if schema.GetTitle() == "" {
		title := structName
		schema.Title = &title
	}

	// Create a minimal global state for the generation
	state := NewGlobalState()
	state.Templates = loadedTemplates
	state.FileWriter = NewFileWriter().WithOverwrite(true).WithHeader(DefaultHeader())

	// Create a processed schema entry
	processed := &parsers.ProcessedSchema{
		Schema:   schema,
		Fields:   parsers.ExtractFields(schema),
		XCodegen: xcodegen,
	}
	state.ProcessedSchemas = map[string]*parsers.ProcessedSchema{
		structName: processed,
	}

	// Find all referenced schemas recursively (same as generateBatchedValueObjects)
	referenced := findReferencedSchemas(state, schema)

	// Add referenced schemas to the state
	for _, refSchema := range referenced {
		refName := refSchema.GetTitle()
		if refName != "" && state.ProcessedSchemas[refName] == nil {
			state.ProcessedSchemas[refName] = &parsers.ProcessedSchema{
				Schema: refSchema,
				Fields: parsers.ExtractFields(refSchema),
			}
		}
	}

	// Collect all schemas to generate
	var schemasToGenerate []*oas3.Schema
	schemasToGenerate = append(schemasToGenerate, referenced...)
	schemasToGenerate = append(schemasToGenerate, schema)

	// Extract fields for all schemas
	allTypes := []map[string]interface{}{}
	hasTimeFields := false
	hasUUIDFields := false

	for _, s := range schemasToGenerate {
		fields := parsers.ExtractFields(s)
		if len(fields) == 0 {
			continue
		}

		// Sort fields alphabetically
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})

		// Check for time and UUID fields
		for _, field := range fields {
			if strings.Contains(field.GoType, "time.Time") {
				hasTimeFields = true
			}
			if strings.Contains(field.GoType, "uuid.UUID") {
				hasUUIDFields = true
			}
		}

		// Convert fields to template format
		var templateFields []map[string]interface{}
		for _, field := range fields {
			jsonName := field.JSONTag
			if jsonName == "" {
				jsonName = field.Name
			}
			// Remove ,omitempty if present as template will add it
			if strings.Contains(jsonName, ",omitempty") {
				jsonName = strings.Split(jsonName, ",")[0]
			}

			templateFields = append(templateFields, map[string]interface{}{
				"FieldName":    field.FieldName,
				"GoType":       field.GoType,
				"JSONName":     jsonName,
				"YAMLName":     field.YAMLTag,
				"Required":     field.Required,
				"Description":  field.Description,
				"DefaultValue": field.DefaultValue,
			})
		}

		// Get the schema title, ensuring it's set
		title := s.GetTitle()
		if title == "" && s == schema {
			title = structName
		}

		allTypes = append(allTypes, map[string]interface{}{
			"Name":        title,
			"Fields":      templateFields,
			"Description": s.GetDescription(),
		})
	}

	// Get the schema template
	tmpl, ok := loadedTemplates["schema.tmpl"]
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	// Create template data
	data := map[string]interface{}{
		"Package":       packageName,
		"Types":         allTypes,
		"HasTimeFields": hasTimeFields,
		"HasUUIDFields": hasUUIDFields,
		"IsValueObject": true, // Treat standalone schemas as value objects
	}

	// Generate content
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := state.FileWriter.WriteFile(outputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Info("Generated struct from JSON Schema",
		slog.String("schema", schemaPath),
		slog.String("struct", structName),
		slog.String("output", outputPath))

	return nil
}
