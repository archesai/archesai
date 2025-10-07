package parsers

import (
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
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
func SchemaToGoType(schema *base.Schema, doc *v3.Document, currentPackage string) string {
	if schema == nil {
		return goTypeInterface
	}

	// Get the types array from the schema
	if len(schema.Type) == 0 {
		return goTypeInterface
	}

	// Use the first type (most schemas have only one type)
	schemaType := schema.Type[0]

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

// SchemaToSQLType converts a JSON Schema to a SQL type
func SchemaToSQLType(schema *base.Schema, dialect string) string {
	if schema == nil {
		return SQLTypeText
	}

	// Get the types array from the schema
	if len(schema.Type) == 0 {
		return SQLTypeText
	}

	// Use the first type
	schemaType := schema.Type[0]

	// Check for string types
	if schemaType == schemaTypeString {
		if schema.Format != "" {
			switch schema.Format {
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
		if schema.Format != "" {
			switch schema.Format {
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
		if schema.Format != "" {
			switch schema.Format {
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

// stringToGoType converts a string schema to a Go type
func stringToGoType(schema *base.Schema) string {
	if schema.Format == "" {
		return goTypeString
	}

	switch schema.Format {
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
func integerToGoType(schema *base.Schema) string {
	if schema.Format == "" {
		return goTypeInt
	}

	switch schema.Format {
	case formatInt32:
		return goTypeInt32
	case formatInt64:
		return goTypeInt64
	default:
		return goTypeInt
	}
}

// numberToGoType converts a number schema to a Go type
func numberToGoType(schema *base.Schema) string {
	if schema.Format == "" {
		return goTypeFloat64
	}

	switch schema.Format {
	case formatFloat:
		return goTypeFloat32
	case formatDouble:
		return goTypeFloat64
	default:
		return goTypeFloat64
	}
}

// arrayToGoType converts an array schema to a Go type
func arrayToGoType(schema *base.Schema, doc *v3.Document, currentPackage string) string {
	if schema.Items == nil || schema.Items.A == nil {
		return goTypeSliceAny
	}

	itemSchema := schema.Items.A.Schema()
	if itemSchema == nil {
		return goTypeSliceAny
	}

	itemTypes := itemSchema.Type
	// Check if the array item is an object with properties
	if len(itemTypes) > 0 && itemTypes[0] == schemaTypeObject &&
		itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
		// For arrays of objects with properties, we'll need special handling
		// This will be handled by the caller (processSchemaProperties)
		return "[]map[string]any"
	}

	itemType := SchemaToGoType(itemSchema, doc, currentPackage)
	return "[]" + itemType
}

// objectToGoType converts an object schema to a Go type
func objectToGoType(schema *base.Schema, doc *v3.Document, currentPackage string) string {
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.IsA() {
		valueSchema := schema.AdditionalProperties.A.Schema()
		if valueSchema != nil {
			valueType := SchemaToGoType(valueSchema, doc, currentPackage)
			return "map[string]" + valueType
		}
	}
	// For objects with properties, just use map[string]any
	return goTypeMapString
}
