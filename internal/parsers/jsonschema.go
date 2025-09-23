package parsers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"

	"github.com/archesai/archesai/internal/templates"
)

// JSONSchema handles JSON Schema operations and transformations.
type JSONSchema struct {
	*oas3.Schema
	Name       string
	Tags       []string
	Extensions map[string]any // Generic extension storage
}

// NewJSONSchema creates a new JSONSchema instance.
func NewJSONSchema(
	schema *oas3.Schema,
) *JSONSchema {
	return &JSONSchema{
		Schema:     schema,
		Extensions: make(map[string]any),
	}
}

// GetExtension retrieves an extension by name.
func (j *JSONSchema) GetExtension(name string) (any, bool) {
	if j.Extensions == nil {
		return nil, false
	}
	ext, ok := j.Extensions[name]
	return ext, ok
}

// ExtractRequiredFields extracts all required fields from a schema including allOf.
func (j *JSONSchema) ExtractRequiredFields() map[string]bool {
	required := make(map[string]bool)

	// Add required fields from main schema
	for _, field := range j.Required {
		required[field] = true
	}

	// Add required fields from resolved allOf schemas
	for _, allOfItem := range j.AllOf {
		resolvedObject := allOfItem.GetResolvedObject()
		if resolvedObject != nil && resolvedObject.IsLeft() {
			refSchema := resolvedObject.GetLeft()
			if refSchema != nil {
				for _, field := range refSchema.Required {
					required[field] = true
				}
			}
		}
	}

	return required
}

// ExtractProperties extracts all properties from a schema including allOf.
func (j *JSONSchema) ExtractProperties() map[string]*oas3.Schema {
	properties := make(map[string]*oas3.Schema)

	// First, add properties from resolved allOf schemas (base schemas)
	for _, allOfItem := range j.AllOf {
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
	if j.Properties != nil {
		for name := range j.Properties.Keys() {
			propRef := j.Properties.GetOrZero(name)
			if propRef != nil && propRef.IsLeft() {
				properties[name] = propRef.GetLeft()
			}
		}
	}

	return properties
}

// =============================================================================
// Type Conversion Methods
// =============================================================================

// SchemaToGoType converts a JSON Schema to a Go type.
func (j *JSONSchema) SchemaToGoType(schema *oas3.Schema) string {
	if schema == nil {
		return "interface{}"
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return "interface{}"
	}

	// Use the first type (most schemas have only one type)
	schemaType := string(types[0])

	// Check for string types with format
	if schemaType == "string" {
		if schema.Format != nil {
			switch *schema.Format {
			case "date-time":
				return "time.Time"
			case "date":
				return "time.Time"
			case "uuid":
				return "uuid.UUID"
			case "email", "uri", "hostname":
				return "string"
			default:
				return "string"
			}
		}
		return "string"
	}

	// Check for numeric types
	if schemaType == "integer" {
		if schema.Format != nil {
			switch *schema.Format {
			case "int32":
				return "int32"
			case "int64":
				return "int64"
			default:
				return "int"
			}
		}
		return "int"
	}

	if schemaType == "number" {
		if schema.Format != nil {
			switch *schema.Format {
			case "float":
				return "float32"
			case "double":
				return "float64"
			default:
				return "float64"
			}
		}
		return "float64"
	}

	// Check for boolean
	if schemaType == "boolean" {
		return "bool"
	}

	// Check for array
	if schemaType == "array" {
		if schema.Items != nil {
			if schema.Items.IsLeft() {
				itemType := j.SchemaToGoType(schema.Items.GetLeft())
				return "[]" + itemType
			}
		}
		return "[]interface{}"
	}

	// Check for object
	if schemaType == "object" {
		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.IsLeft() {
				valueType := j.SchemaToGoType(schema.AdditionalProperties.GetLeft())
				return "map[string]" + valueType
			}
		}
		// If it has properties, generate an inline struct with fields
		if schema.Properties != nil && schema.Properties.Len() > 0 {
			return j.generateInlineStruct(schema)
		}
		return "map[string]interface{}"
	}

	// Default
	return "interface{}"
}

// generateInlineStruct generates an inline struct definition with fields.
func (j *JSONSchema) generateInlineStruct(schema *oas3.Schema) string {
	if schema == nil || schema.Properties == nil {
		return "struct{}"
	}

	var fields []string
	for propName := range schema.Properties.Keys() {
		propRef := schema.Properties.GetOrZero(propName)
		if propRef == nil || !propRef.IsLeft() {
			continue
		}

		prop := propRef.GetLeft()
		if prop == nil {
			continue
		}

		// Convert JSON field name to proper Go field name
		fieldName := templates.PascalCase(propName)
		// Apply initialism fixes
		fieldName = templates.Title(fieldName)

		// Get the Go type for this property
		goType := j.SchemaToGoType(prop)

		// Check if field is required
		isRequired := false
		if schema.Required != nil {
			for _, req := range schema.Required {
				if req == propName {
					isRequired = true
					break
				}
			}
		}

		// Make optional fields pointers (unless they're already pointers, slices, or maps)
		if !isRequired && !strings.HasPrefix(goType, "*") &&
			!strings.HasPrefix(goType, "[]") &&
			!strings.HasPrefix(goType, "map[") {
			goType = "*" + goType
		}

		// Build JSON and YAML tags
		jsonTag := propName
		yamlTag := propName
		if !isRequired {
			jsonTag += ",omitempty"
			yamlTag += ",omitempty"
		}

		// Add description if available
		description := ""
		if prop.Description != nil {
			description = " // " + *prop.Description
		}

		// Build the field definition
		field := fmt.Sprintf("\t%s %s `json:\"%s\" yaml:\"%s\"`%s",
			fieldName, goType, jsonTag, yamlTag, description)
		fields = append(fields, field)
	}

	// Sort fields for consistent output
	sort.Strings(fields)

	if len(fields) == 0 {
		return "struct{}"
	}

	// Return the inline struct definition
	return "struct {\n" + strings.Join(fields, "\n") + "\n}"
}

// SchemaToSQLType converts a JSON Schema to a SQL type.
func (j *JSONSchema) SchemaToSQLType(schema *oas3.Schema, dialect string) string {
	if schema == nil {
		return "TEXT"
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return "TEXT"
	}

	// Use the first type (most schemas have only one type)
	schemaType := string(types[0])

	// Check for string types
	if schemaType == "string" {
		if schema.Format != nil {
			switch *schema.Format {
			case "date-time":
				if dialect == "postgresql" {
					return "TIMESTAMPTZ"
				}
				return "DATETIME"
			case "date":
				return "DATE"
			case "uuid":
				return "UUID"
			case "email":
				if schema.MaxLength != nil && *schema.MaxLength > 0 {
					if strings.ToUpper(dialect) == "POSTGRESQL" {
						return "VARCHAR(" + strconv.Itoa(int(*schema.MaxLength)) + ")"
					}
					return "TEXT"
				}
				return "VARCHAR(255)"
			default:
				if schema.MaxLength != nil && *schema.MaxLength > 0 {
					return "VARCHAR(" + strconv.Itoa(int(*schema.MaxLength)) + ")"
				}
				return "TEXT"
			}
		}
		if schema.MaxLength != nil && *schema.MaxLength > 0 {
			return "VARCHAR(" + strconv.Itoa(int(*schema.MaxLength)) + ")"
		}
		return "TEXT"
	}

	// Check for numeric types
	if schemaType == "integer" {
		if schema.Format != nil {
			switch *schema.Format {
			case "int32":
				return "INTEGER"
			case "int64":
				return "BIGINT"
			default:
				return "INTEGER"
			}
		}
		return "INTEGER"
	}

	if schemaType == "number" {
		if schema.Format != nil {
			switch *schema.Format {
			case "float":
				return "REAL"
			case "double":
				return "DOUBLE PRECISION"
			default:
				return "NUMERIC"
			}
		}
		return "NUMERIC"
	}

	// Check for boolean
	if schemaType == "boolean" {
		return "BOOLEAN"
	}

	// Check for array or object
	if schemaType == "array" || schemaType == "object" {
		if dialect == "postgresql" {
			return "JSONB"
		}
		return "TEXT" // SQLite stores JSON as TEXT
	}

	// Default
	return "TEXT"
}

// ExtractEnumValues extracts enum values from a schema.
func (j *JSONSchema) ExtractEnumValues(schema *oas3.Schema) []string {
	if schema == nil || len(schema.Enum) == 0 {
		return nil
	}

	var values []string
	for _, v := range schema.Enum {
		if v != nil {
			// Convert the value to string representation
			var strVal string
			if err := v.Decode(&strVal); err == nil {
				values = append(values, strVal)
			}
		}
	}
	return values
}

// IsRequired checks if a field is required.
func (j *JSONSchema) IsRequired(fieldName string, required []string) bool {
	for _, r := range required {
		if r == fieldName {
			return true
		}
	}
	return false
}

// ExtractDefaultValue extracts the default value from a schema.
func (j *JSONSchema) ExtractDefaultValue(schema *oas3.Schema) string {
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
			return `"` + str + `"`
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

// InferGoType infers the Go type for a field based on its properties.
func (j *JSONSchema) InferGoType(field templates.FieldData) string {
	// Check format first
	switch field.Format {
	case "uuid":
		return "uuid.UUID"
	case "date-time":
		return "time.Time"
	case "email":
		return "string" // Keep as string for now, not openapi_types.Email
	case "int32":
		return "int32"
	case "int64":
		return "int64"
	case "float":
		return "float32"
	case "double":
		return "float64"
	}

	// Check enum
	if len(field.Enum) > 0 {
		return "string" // Enums are typically strings
	}

	// Use the Type field
	switch field.GoType {
	case "string", "*string":
		return "string"
	case "integer", "*integer":
		return "int"
	case "number", "*number":
		return "float64"
	case "boolean", "*boolean":
		return "bool"
	case "array":
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		// If type starts with *, it's a pointer - extract the base type
		if strings.HasPrefix(field.GoType, "*") {
			return field.GoType[1:]
		}
		// If we have a type, use it
		if field.GoType != "" && field.GoType != "interface{}" {
			return field.GoType
		}
		return "interface{}"
	}

}

// InferSQLCType infers the SQLC type for a field.
func (j *JSONSchema) InferSQLCType(field templates.FieldData) string {
	goType := field.GoType
	if goType == "" {
		goType = j.InferGoType(field)
	}

	// Map Go types to SQLC types
	switch goType {
	case "uuid.UUID":
		return "uuid.UUID"
	case "time.Time":
		return "time.Time"
	case "string":
		return "string"
	case "int", "int32":
		return "int32"
	case "int64":
		return "int64"
	case "float32", "float64":
		return "float64"
	case "bool":
		return "bool"
	case "[]byte":
		return "[]byte"
	case "json.RawMessage", "map[string]interface{}":
		return "json.RawMessage"
	default:
		// For complex types, use json.RawMessage
		if strings.HasPrefix(goType, "[]") || strings.HasPrefix(goType, "map[") {
			return "json.RawMessage"
		}
		return goType
	}
}

// IsPointerType checks if a type should be a pointer based on field properties.
func (j *JSONSchema) IsPointerType(field templates.FieldData) bool {
	// Required fields are not pointers (unless explicitly nullable)
	if field.Required && !field.Nullable {
		return false
	}

	// Optional fields are pointers
	if !field.Required {
		return true
	}

	// Explicitly nullable fields are pointers
	return field.Nullable
}

// WrapPointer wraps a type in a pointer if needed.
func (j *JSONSchema) WrapPointer(goType string, isPointer bool) string {
	if !isPointer || goType == "interface{}" {
		return goType
	}

	if !strings.HasPrefix(goType, "*") {
		return "*" + goType
	}

	return goType
}

// =============================================================================
// Entity Data Preparation Methods
// =============================================================================

// ExtractSchemaFields extracts fields from a schema with allOf support.
func (j *JSONSchema) ExtractSchemaFields() []templates.FieldData {
	if j.Schema == nil {
		return nil
	}

	fieldMap := make(map[string]templates.FieldData)
	requiredFields := make(map[string]bool)

	// Process allOf references first
	j.processAllOfFields(j.Schema, fieldMap, requiredFields)

	// Process direct properties (these override any from allOf)
	j.processDirectProperties(j.Schema, fieldMap, requiredFields)

	// Convert map to slice and update required status
	var fields []templates.FieldData
	for _, field := range fieldMap {
		field.Required = requiredFields[field.Name]

		// Handle optional fields
		if !field.Required {
			if !strings.Contains(field.JSONTag, ",omitempty") {
				field.JSONTag += ",omitempty"
			}
			// For YAML, only add omitempty if not already present
			if field.YAMLTag != "" && !strings.Contains(field.YAMLTag, ",omitempty") {
				field.YAMLTag += ",omitempty"
			}
		}

		// Apply type inference
		field.GoType = j.InferGoType(field)
		field.SQLCType = j.InferSQLCType(field)

		// Make optional fields pointers (after type inference)
		if !field.Required && !strings.HasPrefix(field.GoType, "*") &&
			!strings.HasPrefix(field.GoType, "[]") &&
			!strings.HasPrefix(field.GoType, "map[") {
			field.GoType = "*" + field.GoType
		}

		field.FieldName = field.Name     // Already processed above
		field.SQLCFieldName = field.Name // SQLC uses same naming as Go

		fields = append(fields, field)
	}

	// Sort fields for consistent output
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return fields
}

// processAllOfFields processes allOf references in a schema.
func (j *JSONSchema) processAllOfFields(
	schema *oas3.Schema,
	fieldMap map[string]templates.FieldData,
	requiredFields map[string]bool,
) {
	if schema == nil {
		return
	}

	// Resolve allOf references and extract fields from each
	for _, allOfItem := range j.AllOf {
		resolvedObject := allOfItem.GetResolvedObject()
		if resolvedObject != nil && resolvedObject.IsLeft() {
			refSchema := resolvedObject.GetLeft()
			if refSchema != nil {
				refJsonSchema := NewJSONSchema(refSchema)
				refJsonSchema.extractFieldsFromSchema(fieldMap)

				// Add required fields from this schema
				if refSchema.Required != nil {
					for _, req := range refSchema.Required {
						fieldName := templates.Title(templates.CamelCase(req))
						requiredFields[fieldName] = true
					}
				}
			}
		}
	}
}

// processDirectProperties processes direct schema properties.
func (j *JSONSchema) processDirectProperties(
	schema *oas3.Schema,
	fieldMap map[string]templates.FieldData,
	requiredFields map[string]bool,
) {
	if schema.Properties == nil {
		return
	}

	j.extractFieldsFromSchema(fieldMap)

	// Add required fields from this schema
	if schema.Required != nil {
		for _, req := range schema.Required {
			fieldName := templates.Title(templates.CamelCase(req))
			requiredFields[fieldName] = true
		}
	}
}

// extractFieldsFromSchema extracts fields from a schema into the field map.
func (j *JSONSchema) extractFieldsFromSchema(
	fieldMap map[string]templates.FieldData,
) {
	if j.Properties == nil {
		return
	}

	for propName := range j.Properties.Keys() {
		propRef := j.Properties.GetOrZero(propName)
		if propRef == nil {
			continue
		}

		// Convert JSON field name to proper Go field name
		fieldName := templates.PascalCase(propName)
		// Apply initialism fixes
		fieldName = templates.Title(fieldName)

		var goType string
		var prop *oas3.Schema

		// Properties are always schemas, but they might have a $ref
		if propRef.IsLeft() {
			prop = propRef.GetLeft()
			if prop != nil {
				// Check if this schema has a reference
				if prop.Ref != nil && prop.Ref.String() != "" {
					// It's a reference - extract the type name
					refStr := prop.Ref.String()
					// Handle both local refs (#/components/schemas/ConfigAPI)
					// and file refs (./ConfigAPI.yaml)
					if strings.HasPrefix(refStr, "#/") {
						// Local reference
						parts := strings.Split(refStr, "/")
						if len(parts) > 0 {
							goType = parts[len(parts)-1]
						}
					} else if strings.HasPrefix(refStr, "./") {
						// File reference
						typeName := strings.TrimPrefix(refStr, "./")
						typeName = strings.TrimSuffix(typeName, ".yaml")
						goType = typeName
					} else {
						// Just use the ref as-is
						goType = refStr
					}
				} else {
					// It's an inline schema
					goType = j.SchemaToGoType(prop)
				}
			}
		}

		if goType == "" {
			goType = "interface{}"
		}

		field := templates.FieldData{
			Name:    fieldName,
			GoType:  goType,
			JSONTag: propName,
			YAMLTag: propName,
		}

		// Only extract additional properties if we have a prop (inline schema)
		if prop != nil {
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
						field.Enum = append(field.Enum, fmt.Sprintf("%v", enumVal.Value))
					}
				}
				// Set IsEnumType flag when enum values exist
				if len(field.Enum) > 0 {
					// Generate enum type name (schema name + field name)
					// This will be updated later if we have access to the schema name
					field.IsEnumType = true
				}
			}

			// Check nullable
			if prop.Nullable != nil {
				field.Nullable = *prop.Nullable
			}

			// Extract default value
			if prop.Default != nil {
				var defaultValue interface{}
				if err := prop.Default.Decode(&defaultValue); err == nil {
					field.DefaultValue = fmt.Sprint(defaultValue)
				}
			}
		}

		fieldMap[field.Name] = field
	}
}
