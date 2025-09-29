package codegen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"

	"github.com/archesai/archesai/internal/parsers"
)

// GenerateModels generates all model types (DTOs, entities, value objects)
func (g *Generator) GenerateModels(schemas map[string]*parsers.ProcessedSchema) error {
	for name, processed := range schemas {
		if processed.Schema == nil || processed.XCodegen == nil {
			continue
		}

		if err := g.generateModel(processed, name, processed.XCodegen.GetSchemaType(), nil); err != nil {
			return fmt.Errorf(
				"failed to generate %s %s: %w",
				processed.XCodegen.GetSchemaType(),
				name,
				err,
			)
		}
	}

	return nil
}

// generateModel is the unified function for generating all model types (DTOs, entities, value objects)
func (g *Generator) generateModel(
	schema *parsers.ProcessedSchema,
	name string,
	modelType string,
	customOutputDir *string,
) error {

	// Determine package and output path based on model type
	var packageName, outputDir string
	var isEntity, isValueObject bool
	var processedSchemas map[string]*parsers.ProcessedSchema

	switch modelType {
	case schemaTypeEntity:
		packageName = "entities"
		outputDir = "internal/core/entities"
		if customOutputDir != nil {
			outputDir = *customOutputDir
			packageName = filepath.Base(*customOutputDir)
		}
		isEntity = true
	case schemaTypeValueObject:
		packageName = "valueobjects"
		outputDir = "internal/core/valueobjects"
		if customOutputDir != nil {
			outputDir = *customOutputDir
			packageName = filepath.Base(*customOutputDir)
		}
		isValueObject = true

		// Handle batched value objects for referenced schemas
		if referencedSchemas := g.findReferencedSchemas(schema.Schema, processedSchemas); len(
			referencedSchemas,
		) > 0 {
			allSchemas := append([]*oas3.Schema{schema.Schema}, referencedSchemas...)
			return g.generateBatchedModels(
				allSchemas,
				strings.ToLower(name),
				packageName,
				outputDir,
			)
		}
	default:
		return fmt.Errorf("unsupported model type: %s", modelType)
	}

	// Extract fields
	fields := parsers.ExtractFields(schema.Schema)
	if len(fields) == 0 {
		return nil // Skip if no fields
	}

	// Sort fields alphabetically
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	// For domain types (Entity, Aggregate, ValueObject)
	var requiredFields, optionalFields []parsers.FieldDef
	hasIDField := false

	for _, field := range fields {
		// Categorize fields
		if field.Required {
			requiredFields = append(requiredFields, field)
		} else {
			optionalFields = append(optionalFields, field)
		}

		// Check field types
		if field.FieldName == "ID" {
			hasIDField = true
		}
	}

	// Only entities with ID fields should have domain events
	hasDomainEvents := hasIDField && !isValueObject

	desc := schema.Description
	data := map[string]interface{}{
		"Package":         packageName,
		"Name":            name,
		"NameLower":       strings.ToLower(name),
		"NamePlural":      Pluralize(name),
		"Description":     desc,
		"Type":            modelType,
		"Fields":          fields,
		"RequiredFields":  requiredFields,
		"OptionalFields":  optionalFields,
		"HasDomainEvents": hasDomainEvents,
		"IsEntity":        isEntity,
		"IsValueObject":   isValueObject,
	}

	outputPath := filepath.Join(outputDir, strings.ToLower(name)+".gen.go")
	tmpl, ok := g.templates["schema.tmpl"]
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	return g.filewriter.WriteTemplate(outputPath, tmpl, data)
}

// generateBatchedModels generates multiple models in a single file (used for value objects with references)
func (g *Generator) generateBatchedModels(
	schemas []*oas3.Schema,
	outputName string,
	packageName string,
	outputDir string,
) error {
	allTypes := []map[string]interface{}{}

	for _, schema := range schemas {
		fields := parsers.ExtractFields(schema)
		if len(fields) == 0 {
			continue
		}

		// Sort fields alphabetically
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})

		// Separate required and optional fields
		var requiredFields, optionalFields []parsers.FieldDef
		for _, field := range fields {
			if field.Required {
				requiredFields = append(requiredFields, field)
			} else {
				optionalFields = append(optionalFields, field)
			}
		}

		title := schema.GetTitle()
		desc := schema.GetDescription()
		typeData := map[string]interface{}{
			"Name":           title,
			"NameLower":      strings.ToLower(title),
			"NamePlural":     Pluralize(title),
			"Description":    desc,
			"Fields":         fields,
			"RequiredFields": requiredFields,
			"OptionalFields": optionalFields,
		}

		allTypes = append(allTypes, typeData)
	}

	// Generate the consolidated file
	data := map[string]interface{}{
		"Package": packageName,
		"Types":   allTypes,
	}

	outputPath := filepath.Join(outputDir, outputName+".gen.go")
	tmpl, ok := g.templates["schema.tmpl"]
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	return g.filewriter.WriteTemplate(outputPath, tmpl, data)
}

// extractSchemaNameFromRef extracts the schema name from a reference string
func (g *Generator) extractSchemaNameFromRef(ref string) string {
	// Handle both local and remote references
	// Local: #/components/schemas/Config
	// Remote: ./ConfigAPI.yaml

	if strings.HasPrefix(ref, "#/") {
		// Local reference
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	} else if strings.HasSuffix(ref, ".yaml") {
		// Remote file reference
		base := filepath.Base(ref)
		return strings.TrimSuffix(base, ".yaml")
	}

	// Default: return the last part after any slash
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

// findReferencedSchemas recursively finds all schemas referenced by a given schema's properties
func (g *Generator) findReferencedSchemas(
	schema *oas3.Schema,
	processedSchemas map[string]*parsers.ProcessedSchema,
) []*oas3.Schema {
	// Track which schemas we've already processed to avoid cycles
	seen := make(map[string]bool)
	var result []*oas3.Schema

	// Helper function to recursively find references
	var findRefs func(*oas3.Schema)
	findRefs = func(s *oas3.Schema) {
		if s == nil || s.Properties == nil {
			return
		}

		// Check each property for references
		for propName := range s.Properties.Keys() {
			propRef := s.Properties.GetOrZero(propName)
			if propRef != nil && propRef.IsLeft() {
				prop := propRef.GetLeft()
				if prop != nil && prop.Ref != nil && prop.Ref.String() != "" {
					// Extract the schema name from the reference
					refString := prop.Ref.String()
					schemaName := g.extractSchemaNameFromRef(refString)

					// Skip if we've already seen this schema
					if seen[schemaName] {
						continue
					}
					seen[schemaName] = true

					// Find the referenced schema in our state
					if processed, exists := processedSchemas[schemaName]; exists {
						// Only include if it doesn't have its own x-codegen
						// (schemas with x-codegen will be generated separately)
						if processed.XCodegen == nil && processed.Schema != nil {
							// Set the title on the schema if it's not already set
							if processed.Schema.Title == nil || *processed.Schema.Title == "" {
								title := schemaName
								processed.Schema.Title = &title
							}
							result = append(result, processed.Schema)

							// Recursively find references in this schema
							findRefs(processed.Schema)
						}
					}
				}
			}
		}
	}

	// Start the recursive search
	findRefs(schema)
	return result
}
