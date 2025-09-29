// Package parsers provides OpenAPI and JSON Schema parsing functionality
package parsers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
	"github.com/speakeasy-api/openapi/marshaller"
)

const (
	tagOmitempty = ",omitempty"
)

// ParseJSONSchema loads a schema from a YAML file
func ParseJSONSchema(filePath string) (*oas3.Schema, *XCodegenExtension, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read schema file %s: %w", filePath, err)
	}

	ctx := context.Background()

	// Unmarshal directly to a JSONSchema
	var schema oas3.JSONSchema[oas3.Concrete]
	validationErrs, err := marshaller.Unmarshal(ctx, bytes.NewReader(data), &schema)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	// Validate the schema
	additionalErrs := schema.Validate(ctx)
	validationErrs = append(validationErrs, additionalErrs...)

	// Log validation errors if any
	if len(validationErrs) > 0 {
		for _, vErr := range validationErrs {
			fmt.Printf("Schema validation warning: %s\n", vErr.Error())
		}
	}

	// Get the concrete schema
	var concreteSchema *oas3.Schema
	if schema.IsLeft() {
		concreteSchema = schema.GetLeft()
	} else {
		// It's a reference, we need to resolve it
		return nil, nil, fmt.Errorf("schema is a reference, not a concrete schema")
	}

	// Parse x-codegen extension if present
	var xcodegen *XCodegenExtension
	parser := NewXCodegenParser()
	if concreteSchema.Extensions != nil {
		xcodegen = parser.Parse(*concreteSchema.Extensions)
	}

	return concreteSchema, xcodegen, nil
}

// ProcessSchema processes a single schema and extracts all relevant information
func ProcessSchema(schema *oas3.Schema, name string) (*ProcessedSchema, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	result := &ProcessedSchema{
		Name:   name,
		Title:  extractTitle(schema, name),
		Schema: schema,
	}

	// Extract description
	if schema.Description != nil {
		result.Description = *schema.Description
	}

	// Parse x-codegen extension
	if schema.Extensions != nil {
		xcodegenParser := NewXCodegenParser()
		result.XCodegen = xcodegenParser.Parse(*schema.Extensions)
		if result.XCodegen != nil {
			// Set type flags based on x-codegen
			switch result.XCodegen.SchemaType {
			case "entity":
				result.IsEntity = true
				result.HasDomainEvents = true
			case "aggregate":
				result.IsAggregate = true
				result.HasDomainEvents = true
			case "valueobject":
				result.IsValueObject = true
			case "dto":
				result.IsDTO = true
			}
		}
	}

	// Check if it's an enum
	if len(schema.Enum) > 0 {
		result.IsEnum = true
		result.EnumValues = ExtractEnumValues(schema)
		return result, nil
	}

	// Extract fields
	result.Fields = ExtractFields(schema)

	// Sort fields alphabetically for consistency
	sort.Slice(result.Fields, func(i, j int) bool {
		return result.Fields[i].Name < result.Fields[j].Name
	})

	// Separate required and optional fields
	for _, field := range result.Fields {
		if field.Required {
			result.RequiredFields = append(result.RequiredFields, field)
		} else {
			result.OptionalFields = append(result.OptionalFields, field)
		}
	}

	return result, nil
}

// ExtractFields extracts all fields from a schema including allOf compositions
func ExtractFields(schema *oas3.Schema) []FieldDef {
	if schema == nil {
		return nil
	}

	fieldMap := make(map[string]FieldDef)
	requiredFields := ExtractRequiredFields(schema)

	// Process allOf references first
	processAllOfFields(schema, fieldMap, requiredFields)

	// Process direct properties (these override any from allOf)
	processDirectProperties(schema, fieldMap, requiredFields)

	// Convert map to slice
	var fields []FieldDef
	for _, field := range fieldMap {
		field.Required = requiredFields[field.Name]

		// Handle optional fields
		if !field.Required {
			if !strings.Contains(field.JSONTag, tagOmitempty) {
				field.JSONTag += tagOmitempty
			}
			if field.YAMLTag != "" && !strings.Contains(field.YAMLTag, tagOmitempty) {
				field.YAMLTag += tagOmitempty
			}
		}

		// Apply type inference
		field.GoType = InferGoType(field)

		// Make optional fields pointers (after type inference)
		if !field.Required && !strings.HasPrefix(field.GoType, "*") &&
			!strings.HasPrefix(field.GoType, "[]") &&
			!strings.HasPrefix(field.GoType, "map[") {
			field.GoType = "*" + field.GoType
		}

		fields = append(fields, field)
	}

	return fields
}

// processAllOfFields processes fields from allOf compositions
func processAllOfFields(
	schema *oas3.Schema,
	fieldMap map[string]FieldDef,
	requiredFields map[string]bool,
) {
	for _, allOfItem := range schema.AllOf {
		resolvedObject := allOfItem.GetResolvedObject()
		if resolvedObject != nil && resolvedObject.IsLeft() {
			refSchema := resolvedObject.GetLeft()
			if refSchema != nil {
				// Recursively process allOf schemas
				processAllOfFields(refSchema, fieldMap, requiredFields)

				// Process direct properties
				processSchemaProperties(refSchema, fieldMap)

				// Add required fields
				for _, req := range refSchema.Required {
					fieldName := NormalizeFieldName(req)
					requiredFields[fieldName] = true
				}
			}
		}
	}
}

// processDirectProperties processes direct properties of a schema
func processDirectProperties(
	schema *oas3.Schema,
	fieldMap map[string]FieldDef,
	requiredFields map[string]bool,
) {
	processSchemaProperties(schema, fieldMap)

	// Add required fields from this schema
	for _, req := range schema.Required {
		fieldName := NormalizeFieldName(req)
		requiredFields[fieldName] = true
	}
}

// generateInlineStructType generates an inline struct type definition for nested objects
func generateInlineStructType(schema *oas3.Schema) string {
	if schema == nil || schema.Properties == nil || schema.Properties.Len() == 0 {
		return goTypeInterface
	}

	var structFields []string

	// Process each property in the nested object
	for propName := range schema.Properties.Keys() {
		propRef := schema.Properties.GetOrZero(propName)
		if propRef == nil || !propRef.IsLeft() {
			continue
		}

		prop := propRef.GetLeft()
		if prop == nil {
			continue
		}

		fieldDef := generateStructField(schema, propName, prop)
		if fieldDef != "" {
			structFields = append(structFields, fieldDef)
		}
	}

	if len(structFields) == 0 {
		return goTypeInterface
	}

	// Sort fields for consistency
	sort.Strings(structFields)

	// Build the inline struct type
	return "struct {\n\t\t" + strings.Join(structFields, "\n\t\t") + "\n\t}"
}

// generateStructField generates a single struct field definition
func generateStructField(
	schema *oas3.Schema,
	propName string,
	prop *oas3.Schema,
) string {
	fieldName := NormalizeFieldName(propName)
	goType := determineGoType(prop, fieldName)

	if goType == "" {
		goType = "interface{}"
	}

	isRequired := isFieldRequired(schema, propName)

	// Make optional fields pointers
	if !isRequired && !strings.HasPrefix(goType, fieldName+" struct") &&
		!strings.HasPrefix(goType, "*") && !strings.HasPrefix(goType, "[]") &&
		!strings.HasPrefix(goType, "map[") && !strings.HasPrefix(goType, "struct{") {
		goType = "*" + goType
	}

	jsonTag := propName
	if !isRequired {
		jsonTag += tagOmitempty
	}

	// Format the field definition
	if strings.HasPrefix(goType, fieldName+" ") {
		// It's a named inline struct, format it specially
		structDef := strings.TrimPrefix(goType, fieldName+" ")
		return fmt.Sprintf("%s %s `json:\"%s\" yaml:\"%s\"`",
			fieldName, structDef, jsonTag, propName)
	}

	return fmt.Sprintf("%s %s `json:\"%s\" yaml:\"%s\"`",
		fieldName, goType, jsonTag, propName)
}

// determineGoType determines the Go type for a property
func determineGoType(prop *oas3.Schema, fieldName string) string {
	if prop == nil {
		return ""
	}

	// Handle references
	if prop.Ref != nil && prop.Ref.String() != "" {
		return extractTypeFromRef(prop.Ref.String())
	}

	// Handle nested objects
	types := prop.GetType()
	if len(types) > 0 {
		switch string(types[0]) {
		case schemaTypeObject:
			if prop.Properties != nil && prop.Properties.Len() > 0 {
				// Recursively generate inline struct
				return fieldName + " " + generateInlineStructType(prop)
			}
		case "array":
			// Handle array of objects specially
			if prop.Items != nil && prop.Items.IsLeft() {
				itemSchema := prop.Items.GetLeft()
				if itemSchema != nil {
					itemTypes := itemSchema.GetType()
					if len(itemTypes) > 0 && string(itemTypes[0]) == "object" &&
						itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
						// Array of objects with properties
						return "[]" + generateInlineStructType(itemSchema)
					}
				}
			}
		}
	}

	// Use the existing SchemaToGoType for everything else
	return SchemaToGoType(prop)
}

// isFieldRequired checks if a field is required
func isFieldRequired(schema *oas3.Schema, propName string) bool {
	for _, req := range schema.Required {
		if req == propName {
			return true
		}
	}
	return false
}

// processSchemaProperties extracts properties from a schema
func processSchemaProperties(
	schema *oas3.Schema,
	fieldMap map[string]FieldDef,
) {
	if schema.Properties == nil {
		return
	}

	for propName := range schema.Properties.Keys() {
		propRef := schema.Properties.GetOrZero(propName)
		if propRef == nil {
			continue
		}

		// Convert JSON field name to proper Go field name
		fieldName := NormalizeFieldName(propName)

		var goType string
		var prop *oas3.Schema

		if propRef.IsLeft() {
			prop = propRef.GetLeft()
			if prop != nil {
				// Check if it's a reference
				if prop.Ref != nil && prop.Ref.String() != "" {
					goType = extractTypeFromRef(prop.Ref.String())
				} else {
					// Handle nested objects and arrays as inline structs
					types := prop.GetType()
					if len(types) > 0 && string(types[0]) == schemaTypeObject && prop.Properties != nil && prop.Properties.Len() > 0 {
						// Generate inline struct type recursively
						goType = generateInlineStructType(prop)
					} else if len(types) > 0 && string(types[0]) == "array" {
						// Handle arrays specially
						if prop.Items != nil && prop.Items.IsLeft() {
							itemSchema := prop.Items.GetLeft()
							if itemSchema != nil {
								itemTypes := itemSchema.GetType()
								if len(itemTypes) > 0 && string(itemTypes[0]) == "object" &&
									itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
									// Array of objects with properties - generate inline struct
									goType = "[]" + generateInlineStructType(itemSchema)
								} else {
									goType = SchemaToGoType(prop)
								}
							} else {
								goType = SchemaToGoType(prop)
							}
						} else {
							goType = SchemaToGoType(prop)
						}
					} else {
						goType = SchemaToGoType(prop)
					}
				}
			}
		}

		if goType == "" {
			goType = "interface{}"
		}

		field := FieldDef{
			Name:      fieldName,
			FieldName: fieldName,
			GoType:    goType,
			JSONTag:   propName,
			YAMLTag:   propName,
		}

		// Extract additional properties if we have a prop
		if prop != nil {
			extractFieldMetadata(prop, &field)
		}

		fieldMap[field.Name] = field
	}
}

// extractFieldMetadata extracts metadata from a schema property
func extractFieldMetadata(prop *oas3.Schema, field *FieldDef) {
	// Extract format
	if prop.Format != nil {
		field.Format = *prop.Format
	}

	// Extract description
	if prop.Description != nil {
		field.Description = *prop.Description
	}

	// Extract enum values
	if prop.Enum != nil {
		for _, enumVal := range prop.Enum {
			if enumVal != nil {
				var val string
				if err := enumVal.Decode(&val); err == nil {
					field.Enum = append(field.Enum, val)
				}
			}
		}
		if len(field.Enum) > 0 {
			field.IsEnumType = true
		}
	}

	// Check nullable
	if prop.Nullable != nil {
		field.Nullable = *prop.Nullable
	}

	// Extract default value
	if prop.Default != nil {
		field.DefaultValue = ExtractDefaultValue(prop)
	}
}

// ExtractRequiredFields extracts all required fields from a schema including allOf
func ExtractRequiredFields(schema *oas3.Schema) map[string]bool {
	required := make(map[string]bool)

	if schema == nil {
		return required
	}

	// Add required fields from main schema
	for _, field := range schema.Required {
		required[NormalizeFieldName(field)] = true
	}

	// Add required fields from resolved allOf schemas
	for _, allOfItem := range schema.AllOf {
		resolvedObject := allOfItem.GetResolvedObject()
		if resolvedObject != nil && resolvedObject.IsLeft() {
			refSchema := resolvedObject.GetLeft()
			if refSchema != nil {
				for _, field := range refSchema.Required {
					required[NormalizeFieldName(field)] = true
				}
			}
		}
	}

	return required
}

// ExtractEnumValues extracts enum values from a schema
func ExtractEnumValues(schema *oas3.Schema) []string {
	if schema == nil || len(schema.Enum) == 0 {
		return nil
	}

	var values []string
	for _, v := range schema.Enum {
		if v != nil {
			var strVal string
			if err := v.Decode(&strVal); err == nil {
				values = append(values, strVal)
			}
		}
	}
	return values
}

// ExtractDefaultValue extracts the default value from a schema
func ExtractDefaultValue(schema *oas3.Schema) string {
	if schema == nil || schema.Default == nil {
		return ""
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return ""
	}

	// Use the first type
	schemaType := string(types[0])

	// Convert the default value to string representation
	var defaultVal interface{}
	if err := schema.Default.Decode(&defaultVal); err != nil {
		return ""
	}

	// Handle different types
	switch schemaType {
	case "string":
		if str, ok := defaultVal.(string); ok {
			return str
		}
	case "boolean":
		if b, ok := defaultVal.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
	case "integer", "number":
		return fmt.Sprintf("%v", defaultVal)
	}

	return ""
}

// ExtractProperties extracts all properties from a schema including allOf
func ExtractProperties(schema *oas3.Schema) map[string]*oas3.Schema {
	properties := make(map[string]*oas3.Schema)

	// First, add properties from resolved allOf schemas
	for _, allOfItem := range schema.AllOf {
		resolvedObject := allOfItem.GetResolvedObject()
		if resolvedObject != nil && resolvedObject.IsLeft() {
			refSchema := resolvedObject.GetLeft()
			if refSchema != nil && refSchema.Properties != nil {
				for name := range refSchema.Properties.Keys() {
					propRef := refSchema.Properties.GetOrZero(name)
					if propRef != nil && propRef.IsLeft() {
						properties[name] = propRef.GetLeft()
					}
				}
			}
		}
	}

	// Then, add/override with properties from main schema
	if schema.Properties != nil {
		for name := range schema.Properties.Keys() {
			propRef := schema.Properties.GetOrZero(name)
			if propRef != nil && propRef.IsLeft() {
				properties[name] = propRef.GetLeft()
			}
		}
	}

	return properties
}

// extractTitle extracts the title from a schema or uses the default name
func extractTitle(schema *oas3.Schema, defaultName string) string {
	if schema != nil && schema.Title != nil && *schema.Title != "" {
		return *schema.Title
	}
	return defaultName
}

// extractTypeFromRef extracts type name from a reference string
func extractTypeFromRef(ref string) string {
	if strings.HasPrefix(ref, "#/") {
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	} else if strings.HasPrefix(ref, "./") {
		typeName := strings.TrimPrefix(ref, "./")
		typeName = strings.TrimSuffix(typeName, ".yaml")
		return typeName
	}
	return ref
}
