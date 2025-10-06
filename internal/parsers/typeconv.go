package parsers

import (
	"strconv"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
	"github.com/speakeasy-api/openapi/openapi"
)

// Schema type constants
const (
	schemaTypeString  = "string"
	schemaTypeInteger = "integer"
	schemaTypeNumber  = "number"
	schemaTypeBoolean = "boolean"
	schemaTypeArray   = "array"
	schemaTypeObject  = "object"
)

// Format constants
const (
	formatDateTime = "date-time"
	formatDate     = "date"
	formatUUID     = "uuid"
	formatEmail    = "email"
	formatURI      = "uri"
	formatHostname = "hostname"
	formatInt32    = "int32"
	formatInt64    = "int64"
	formatFloat    = "float"
	formatDouble   = "double"
)

// Go type constants
const (
	goTypeInterface = "any"
	goTypeString    = "string"
	goTypeInt       = "int"
	goTypeInt32     = "int32"
	goTypeInt64     = "int64"
	goTypeFloat32   = "float32"
	goTypeFloat64   = "float64"
	goTypeBool      = "bool"
	goTypeTime      = "time.Time"
	goTypeUUID      = "uuid.UUID"
	goTypeMapString = "map[string]any"
	goTypeSliceAny  = "[]any"
)

// SchemaToGoType converts a JSON Schema to a Go type with proper package qualification
func SchemaToGoType(schema *oas3.Schema, doc *openapi.OpenAPI, currentPackage string) string {
	if schema == nil {
		return goTypeInterface
	}

	// Handle references first
	if refType, resolved := resolveSchemaReference(schema, doc, currentPackage); resolved {
		return refType
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return goTypeInterface
	}

	// Use the first type (most schemas have only one type)
	schemaType := string(types[0])

	// Delegate to type-specific handlers
	switch schemaType {
	case schemaTypeString:
		return stringToGoType(schema)
	case schemaTypeInteger:
		return integerToGoType(schema)
	case schemaTypeNumber:
		return numberToGoType(schema)
	case schemaTypeBoolean:
		return goTypeBool
	case schemaTypeArray:
		return arrayToGoType(schema, doc, currentPackage)
	case schemaTypeObject:
		return objectToGoType(schema, doc, currentPackage)
	default:
		return goTypeInterface
	}
}

// resolveSchemaReference resolves a schema reference to a Go type with proper package qualification
func resolveSchemaReference(
	schema *oas3.Schema,
	doc *openapi.OpenAPI,
	currentPackage string,
) (string, bool) {
	if schema.Ref == nil || schema.Ref.String() == "" {
		return "", false
	}

	refString := schema.Ref.String()
	if !strings.HasPrefix(refString, "#/components/schemas/") {
		if doc == nil {
			return "any", true
		}
		return "", false
	}

	schemaName := strings.TrimPrefix(refString, "#/components/schemas/")

	// If doc is nil, return the schema name
	if doc == nil {
		return schemaName, true
	}

	// Look up the schema to get its X-Codegen-Schema-Type
	if doc.Components == nil || doc.Components.Schemas == nil {
		return schemaName, true
	}

	referencedSchema := doc.Components.Schemas.GetOrZero(schemaName)
	if referencedSchema == nil {
		return schemaName, true
	}

	resolvedSchema := referencedSchema.GetResolvedObject()
	if resolvedSchema == nil || resolvedSchema.Left == nil {
		return schemaName, true
	}

	targetPackage := extractTargetPackage(resolvedSchema.Left)
	typeName := resolvedSchema.Left.GetTitle()
	if typeName == "" {
		typeName = schemaName
	}

	return qualifyTypeName(typeName, targetPackage, currentPackage), true
}

// extractTargetPackage extracts the target package from x-codegen extension
func extractTargetPackage(schema *oas3.Schema) string {
	if schema.Extensions == nil {
		return ""
	}

	xExt, found := schema.Extensions.Get("x-codegen")
	if !found {
		return ""
	}

	parser := &XCodegenParser{}
	xcodegen, err := parser.ParseExtension(xExt, schema.GetTitle())
	if err != nil || xcodegen == nil || xcodegen.SchemaType == "" {
		return ""
	}

	return string(xcodegen.GetSchemaType())
}

// qualifyTypeName qualifies a type name with package prefix based on context
func qualifyTypeName(typeName, targetPackage, currentPackage string) string {
	// If we're in a controller context (currentPackage is empty), always add the package prefix
	if currentPackage == "" && targetPackage != "" {
		return addPackagePrefix(typeName, targetPackage)
	}

	// Same package, no prefix needed
	if currentPackage != "" && currentPackage == targetPackage {
		return typeName
	}

	// Different packages, add prefix
	if currentPackage != "" && targetPackage != "" {
		return addPackagePrefix(typeName, targetPackage)
	}

	return typeName
}

// addPackagePrefix adds the appropriate package prefix to a type name
func addPackagePrefix(typeName, targetPackage string) string {
	switch targetPackage {
	case "entity":
		return "entities." + typeName
	case "valueobject":
		return "valueobjects." + typeName
	default:
		return typeName
	}
}

// stringToGoType converts a string schema to a Go type
func stringToGoType(schema *oas3.Schema) string {
	if schema.Format == nil {
		return goTypeString
	}

	switch *schema.Format {
	case formatDateTime, formatDate:
		return goTypeTime
	case formatUUID:
		return goTypeUUID
	case formatEmail, formatURI, formatHostname:
		return goTypeString
	default:
		return goTypeString
	}
}

// integerToGoType converts an integer schema to a Go type
func integerToGoType(schema *oas3.Schema) string {
	if schema.Format == nil {
		return goTypeInt
	}

	switch *schema.Format {
	case formatInt32:
		return goTypeInt32
	case formatInt64:
		return goTypeInt64
	default:
		return goTypeInt
	}
}

// numberToGoType converts a number schema to a Go type
func numberToGoType(schema *oas3.Schema) string {
	if schema.Format == nil {
		return goTypeFloat64
	}

	switch *schema.Format {
	case formatFloat:
		return goTypeFloat32
	case formatDouble:
		return goTypeFloat64
	default:
		return goTypeFloat64
	}
}

// arrayToGoType converts an array schema to a Go type
func arrayToGoType(schema *oas3.Schema, doc *openapi.OpenAPI, currentPackage string) string {
	if schema.Items == nil || !schema.Items.IsLeft() {
		return goTypeSliceAny
	}

	itemSchema := schema.Items.GetLeft()
	if itemSchema == nil {
		return goTypeSliceAny
	}

	itemTypes := itemSchema.GetType()
	// Check if the array item is an object with properties
	if len(itemTypes) > 0 && string(itemTypes[0]) == schemaTypeObject &&
		itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
		// For arrays of objects with properties, we'll need special handling
		// This will be handled by the caller (processSchemaProperties)
		return "[]map[string]any"
	}

	itemType := SchemaToGoType(itemSchema, doc, currentPackage)
	return "[]" + itemType
}

// objectToGoType converts an object schema to a Go type
func objectToGoType(schema *oas3.Schema, doc *openapi.OpenAPI, currentPackage string) string {
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.IsLeft() {
		valueType := SchemaToGoType(schema.AdditionalProperties.GetLeft(), doc, currentPackage)
		return "map[string]" + valueType
	}
	// For objects with properties, just use map[string]any
	return goTypeMapString
}

// SchemaToSQLType converts a JSON Schema to a SQL type
func SchemaToSQLType(schema *oas3.Schema, dialect string) string {
	if schema == nil {
		return SQLTypeText
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return SQLTypeText
	}

	// Use the first type
	schemaType := string(types[0])

	// Check for string types
	if schemaType == schemaTypeString {
		if schema.Format != nil {
			switch *schema.Format {
			case formatDateTime:
				if dialect == "postgresql" {
					return "TIMESTAMPTZ"
				}
				return "DATETIME"
			case formatDate:
				return "DATE"
			case formatUUID:
				return "UUID"
			case formatEmail:
				if schema.MaxLength != nil && *schema.MaxLength > 0 {
					if strings.ToUpper(dialect) == "POSTGRESQL" {
						return "VARCHAR(" + strconv.Itoa(int(*schema.MaxLength)) + ")"
					}
					return SQLTypeText
				}
				return "VARCHAR(255)"
			default:
				if schema.MaxLength != nil && *schema.MaxLength > 0 {
					return "VARCHAR(" + strconv.Itoa(int(*schema.MaxLength)) + ")"
				}
				return SQLTypeText
			}
		}
		if schema.MaxLength != nil && *schema.MaxLength > 0 {
			return "VARCHAR(" + strconv.Itoa(int(*schema.MaxLength)) + ")"
		}
		return SQLTypeText
	}

	// Check for numeric types
	if schemaType == schemaTypeInteger {
		if schema.Format != nil {
			switch *schema.Format {
			case formatInt32:
				return SQLTypeInteger
			case formatInt64:
				return "BIGINT"
			default:
				return SQLTypeInteger
			}
		}
		return SQLTypeInteger
	}

	if schemaType == schemaTypeNumber {
		if schema.Format != nil {
			switch *schema.Format {
			case formatFloat:
				return "REAL"
			case formatDouble:
				return "DOUBLE PRECISION"
			default:
				return "NUMERIC"
			}
		}
		return "NUMERIC"
	}

	// Check for boolean
	if schemaType == schemaTypeBoolean {
		return "BOOLEAN"
	}

	// Check for array or object
	if schemaType == schemaTypeArray || schemaType == schemaTypeObject {
		if dialect == "postgresql" {
			return "JSONB"
		}
		return SQLTypeText // SQLite stores JSON as TEXT
	}

	// Default
	return SQLTypeText
}
