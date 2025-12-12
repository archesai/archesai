package spec

import (
	"strconv"
	"strings"
)

// SQL dialect constants.
const (
	SQLDialectPostgres = "postgresql"
	SQLDialectSQLite   = "sqlite"
)

// SQL type constants.
const (
	SQLTypeText      = "TEXT"
	SQLTypeInteger   = "INTEGER"
	SQLTypeBigInt    = "BIGINT"
	SQLTypeBoolean   = "BOOLEAN"
	SQLTypeTimestamp = "TIMESTAMPTZ"
	SQLTypeDateTime  = "DATETIME"
	SQLTypeDate      = "DATE"
	SQLTypeUUID      = "UUID"
	SQLTypeJSONB     = "JSONB"
	SQLTypeNumeric   = "NUMERIC"
	SQLTypeReal      = "REAL"
	SQLTypeDouble    = "DOUBLE PRECISION"
)

// SchemaToSQLType converts a schema type and format to a SQL type.
func SchemaToSQLType(schemaType, format string, maxLength int, dialect string) string {
	// Check for string types
	if schemaType == SchemaTypeString {
		if format != "" {
			switch format {
			case FormatDateTime:
				if dialect == "postgresql" {
					return "TIMESTAMPTZ"
				}
				return "DATETIME"
			case FormatDate:
				return "DATE"
			case FormatUUID:
				return "UUID"
			case FormatEmail:
				if maxLength > 0 {
					if strings.ToUpper(dialect) == "POSTGRESQL" {
						return "VARCHAR(" + strconv.Itoa(maxLength) + ")"
					}
					return SQLTypeText
				}
				return "VARCHAR(255)"
			default:
				if maxLength > 0 {
					return "VARCHAR(" + strconv.Itoa(maxLength) + ")"
				}
				return SQLTypeText
			}
		}
		if maxLength > 0 {
			return "VARCHAR(" + strconv.Itoa(maxLength) + ")"
		}
		return SQLTypeText
	}

	// Check for numeric types
	if schemaType == SchemaTypeInteger {
		if format != "" {
			switch format {
			case FormatInt32:
				return SQLTypeInteger
			case FormatInt64:
				return "BIGINT"
			default:
				return SQLTypeInteger
			}
		}
		return SQLTypeInteger
	}

	if schemaType == SchemaTypeNumber {
		if format != "" {
			switch format {
			case FormatFloat:
				return "REAL"
			case FormatDouble:
				return "DOUBLE PRECISION"
			default:
				return "NUMERIC"
			}
		}
		return "NUMERIC"
	}

	// Check for boolean
	if schemaType == SchemaTypeBoolean {
		return "BOOLEAN"
	}

	// Check for array or object
	if schemaType == SchemaTypeArray || schemaType == SchemaTypeObject {
		if dialect == "postgresql" {
			return "JSONB"
		}
		return SQLTypeText // SQLite stores JSON as TEXT
	}

	// Default
	return SQLTypeText
}
