package typeconv

import (
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	"github.com/archesai/archesai/internal/spec"
)

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
	if schemaType == spec.SchemaTypeString {
		if schema.Format != "" {
			switch schema.Format {
			case spec.FormatDateTime:
				if dialect == "postgresql" {
					return "TIMESTAMPTZ"
				}
				return "DATETIME"
			case spec.FormatDate:
				return "DATE"
			case spec.FormatUUID:
				return "UUID"
			case spec.FormatEmail:
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
	if schemaType == spec.SchemaTypeInteger {
		if schema.Format != "" {
			switch schema.Format {
			case spec.FormatInt32:
				return SQLTypeInteger
			case spec.FormatInt64:
				return "BIGINT"
			default:
				return SQLTypeInteger
			}
		}
		return SQLTypeInteger
	}

	if schemaType == spec.SchemaTypeNumber {
		if schema.Format != "" {
			switch schema.Format {
			case spec.FormatFloat:
				return "REAL"
			case spec.FormatDouble:
				return "DOUBLE PRECISION"
			default:
				return "NUMERIC"
			}
		}
		return "NUMERIC"
	}

	// Check for boolean
	if schemaType == spec.SchemaTypeBoolean {
		return "BOOLEAN"
	}

	// Check for array or object
	if schemaType == spec.SchemaTypeArray || schemaType == spec.SchemaTypeObject {
		if dialect == "postgresql" {
			return "JSONB"
		}
		return SQLTypeText // SQLite stores JSON as TEXT
	}

	// Default
	return SQLTypeText
}
