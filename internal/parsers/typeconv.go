package parsers

import (
	"strconv"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
)

// SchemaToGoType converts a JSON Schema to a Go type
func SchemaToGoType(schema *oas3.Schema) string {
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
				itemSchema := schema.Items.GetLeft()
				if itemSchema != nil {
					itemTypes := itemSchema.GetType()
					// Check if the array item is an object with properties
					if len(itemTypes) > 0 && string(itemTypes[0]) == "object" &&
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
		return "[]interface{}"
	}

	// Check for object
	if schemaType == "object" {
		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.IsLeft() {
				valueType := SchemaToGoType(schema.AdditionalProperties.GetLeft())
				return "map[string]" + valueType
			}
		}
		// For objects with properties, just use map[string]interface{}
		return "map[string]interface{}"
	}

	// Default
	return "interface{}"
}

// SchemaToSQLType converts a JSON Schema to a SQL type
func SchemaToSQLType(schema *oas3.Schema, dialect string) string {
	if schema == nil {
		return "TEXT"
	}

	// Get the types array from the schema
	types := schema.GetType()
	if len(types) == 0 {
		return "TEXT"
	}

	// Use the first type
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

// InferGoType infers the Go type for a field based on its properties
func InferGoType(field FieldDef) string {
	// Check format first
	switch field.Format {
	case "uuid":
		return "uuid.UUID"
	case "date-time":
		return "time.Time"
	case "email":
		return "string"
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
