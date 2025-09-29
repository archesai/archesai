package parsers

import (
	"strconv"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
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
	goTypeInterface = "interface{}"
	goTypeString    = "string"
	goTypeInt       = "int"
	goTypeInt32     = "int32"
	goTypeInt64     = "int64"
	goTypeFloat32   = "float32"
	goTypeFloat64   = "float64"
	goTypeBool      = "bool"
	goTypeTime      = "time.Time"
	goTypeUUID      = "uuid.UUID"
	goTypeMapString = "map[string]interface{}"
	goTypeSliceAny  = "[]interface{}"
)

// SchemaToGoType converts a JSON Schema to a Go type
func SchemaToGoType(schema *oas3.Schema) string {
	if schema == nil {
		return goTypeInterface
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return goTypeInterface
	}

	// Use the first type (most schemas have only one type)
	schemaType := string(types[0])

	// Check for string types with format
	if schemaType == schemaTypeString {
		if schema.Format != nil {
			switch *schema.Format {
			case formatDateTime:
				return goTypeTime
			case formatDate:
				return goTypeTime
			case formatUUID:
				return goTypeUUID
			case formatEmail, formatURI, formatHostname:
				return goTypeString
			default:
				return goTypeString
			}
		}
		return goTypeString
	}

	// Check for numeric types
	if schemaType == schemaTypeInteger {
		if schema.Format != nil {
			switch *schema.Format {
			case formatInt32:
				return goTypeInt32
			case formatInt64:
				return goTypeInt64
			default:
				return goTypeInt
			}
		}
		return goTypeInt
	}

	if schemaType == schemaTypeNumber {
		if schema.Format != nil {
			switch *schema.Format {
			case formatFloat:
				return goTypeFloat32
			case formatDouble:
				return goTypeFloat64
			default:
				return goTypeFloat64
			}
		}
		return goTypeFloat64
	}

	// Check for boolean
	if schemaType == schemaTypeBoolean {
		return goTypeBool
	}

	// Check for array
	if schemaType == schemaTypeArray {
		if schema.Items != nil {
			if schema.Items.IsLeft() {
				itemSchema := schema.Items.GetLeft()
				if itemSchema != nil {
					itemTypes := itemSchema.GetType()
					// Check if the array item is an object with properties
					if len(itemTypes) > 0 && string(itemTypes[0]) == schemaTypeObject &&
						itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
						// For arrays of objects with properties, we'll need special handling
						// This will be handled by the caller (processSchemaProperties)
						return "[]map[string]interface{}"
					}
					itemType := SchemaToGoType(itemSchema)
					return "[]" + itemType
				}
			}
		}
		return goTypeSliceAny
	}

	// Check for object
	if schemaType == schemaTypeObject {
		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.IsLeft() {
				valueType := SchemaToGoType(schema.AdditionalProperties.GetLeft())
				return "map[string]" + valueType
			}
		}
		// For objects with properties, just use map[string]interface{}
		return goTypeMapString
	}

	// Default
	return goTypeInterface
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

// InferGoType infers the Go type for a field based on its properties
func InferGoType(field FieldDef) string {
	// Check format first
	switch field.Format {
	case formatUUID:
		return goTypeUUID
	case formatDateTime:
		return goTypeTime
	case formatEmail:
		return goTypeString
	case formatInt32:
		return goTypeInt32
	case formatInt64:
		return goTypeInt64
	case formatFloat:
		return goTypeFloat32
	case formatDouble:
		return goTypeFloat64
	}

	// Check enum
	if len(field.Enum) > 0 {
		return goTypeString // Enums are typically strings
	}

	// Use the Type field
	switch field.GoType {
	case schemaTypeString, "*" + schemaTypeString:
		return goTypeString
	case schemaTypeInteger, "*" + schemaTypeInteger:
		return goTypeInt
	case schemaTypeNumber, "*" + schemaTypeNumber:
		return goTypeFloat64
	case schemaTypeBoolean, "*" + schemaTypeBoolean:
		return goTypeBool
	case schemaTypeArray:
		return goTypeSliceAny
	case schemaTypeObject:
		return goTypeMapString
	default:
		// If type starts with *, it's a pointer - extract the base type
		if strings.HasPrefix(field.GoType, "*") {
			return field.GoType[1:]
		}
		// If we have a type, use it
		if field.GoType != "" && field.GoType != goTypeInterface {
			return field.GoType
		}
		return goTypeInterface
	}
}

// NormalizeFieldName converts a field name to Go conventions
func NormalizeFieldName(name string) string {
	if name == "" {
		return ""
	}

	// Special cases
	switch name {
	case "id", "ID":
		return "ID"
	case "url", "URL":
		return "URL"
	case "api", "API":
		return "API"
	case "jwtSecret":
		return "JWTSecret"
	case "accessTokenTtl":
		return "AccessTokenTTL"
	case "refreshTokenTtl":
		return "RefreshTokenTTL"
	}

	// Convert snake_case to PascalCase
	parts := strings.Split(name, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
