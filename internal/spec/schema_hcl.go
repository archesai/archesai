package spec

import (
	"fmt"
	"strconv"
	"strings"
)

// SQL function constants for HCL.
const (
	sqlCurrentTimestamp = "sql(\"CURRENT_TIMESTAMP\")"
	sqlGenRandomUUID    = "sql(\"gen_random_uuid()\")"
	sqlRandomBlob       = "sql(\"lower(hex(randomblob(16)))\")"
)

// HCL/PostgreSQL type constants.
const (
	hclTypeText            = `sql("text")`
	hclTypeInteger         = `sql("integer")`
	hclTypeBigint          = `sql("bigint")`
	hclTypeBoolean         = `sql("boolean")`
	hclTypeUUID            = `sql("uuid")`
	hclTypeTimestamptz     = `sql("timestamptz")`
	hclTypeDate            = `sql("date")`
	hclTypeTime            = `sql("time")`
	hclTypeBytea           = `sql("bytea")`
	hclTypeReal            = `sql("real")`
	hclTypeDoublePrecision = `sql("double precision")`
	hclTypeNumeric         = `sql("numeric")`
	hclTypeJSONB           = `sql("jsonb")`
)

// HCL/SQLite type constants.
const (
	hclTypeSQLiteText    = `sql("TEXT")`
	hclTypeSQLiteInteger = `sql("INTEGER")`
	hclTypeSQLiteReal    = `sql("REAL")`
	hclTypeSQLiteBlob    = `sql("BLOB")`
)

// HCLType returns the HCL/PostgreSQL type string for this schema.
func (s *Schema) HCLType() string {
	field := s
	// Handle special formats first
	switch field.Format {
	case FormatUUID:
		return hclTypeUUID
	case FormatDateTime:
		return hclTypeTimestamptz
	case FormatDate:
		return hclTypeDate
	case "time":
		return hclTypeTime
	case "ipv4", "ipv6":
		return hclTypeText
	case FormatEmail, FormatURI, FormatHostname:
		return hclTypeText
	case "binary":
		return hclTypeBytea
	case FormatInt32:
		return hclTypeInteger
	case FormatInt64:
		return hclTypeBigint
	case FormatFloat:
		return hclTypeReal
	case FormatDouble:
		return hclTypeDoublePrecision
	}

	// Handle arrays
	if field.Type.PrimaryType() == SchemaTypeArray && field.Items != nil {
		itemType := field.Items.GetOrNil().HCLType()
		// Remove sql() wrapper from itemType if present
		if strings.HasPrefix(itemType, `sql("`) && strings.HasSuffix(itemType, `")`) {
			itemType = itemType[5 : len(itemType)-2]
		}
		return `sql("` + itemType + `[]")`
	}

	// Handle basic types
	switch field.Type.PrimaryType() {
	case SchemaTypeString:
		// Check for specific field names that need special types
		if strings.Contains(strings.ToLower(field.Name), "embedding") {
			return "sql(\"vector(1536)\")"
		}
		// Check for enums
		if len(field.Enum) > 0 {
			return hclTypeText
		}
		// For now, just use text for all strings
		return hclTypeText
	case SchemaTypeInteger:
		if field.Format == FormatInt64 {
			return hclTypeBigint
		}
		return hclTypeInteger
	case SchemaTypeNumber:
		if field.Format == FormatFloat {
			return hclTypeReal
		}
		return hclTypeNumeric
	case SchemaTypeBoolean:
		return hclTypeBoolean
	case SchemaTypeObject:
		// For nested objects, use JSONB
		return hclTypeJSONB
	default:
		return hclTypeText
	}
}

// SQLiteHCLType returns the HCL/SQLite type string for this schema.
func (s *Schema) SQLiteHCLType() string {
	field := s
	// SQLite has a simpler type system: TEXT, INTEGER, REAL, BLOB
	// Handle special formats first
	switch field.Format {
	case FormatUUID:
		return hclTypeSQLiteText // UUIDs stored as TEXT in SQLite
	case FormatDateTime, FormatDate, "time":
		return hclTypeSQLiteText // Timestamps stored as TEXT in SQLite
	case "ipv4", "ipv6":
		return hclTypeSQLiteText
	case FormatEmail, FormatURI, FormatHostname:
		return hclTypeSQLiteText
	case "binary":
		return hclTypeSQLiteBlob
	case FormatInt32, FormatInt64:
		return hclTypeSQLiteInteger
	case FormatFloat, FormatDouble:
		return hclTypeSQLiteReal
	}

	// Handle arrays - SQLite stores arrays as TEXT (JSON)
	if field.Type.PrimaryType() == SchemaTypeArray {
		return hclTypeSQLiteText
	}

	// Handle basic types
	switch field.Type.PrimaryType() {
	case SchemaTypeString:
		// Enums, embeddings, and all strings stored as TEXT
		return hclTypeSQLiteText
	case SchemaTypeInteger:
		return hclTypeSQLiteInteger
	case SchemaTypeNumber:
		return hclTypeSQLiteReal
	case SchemaTypeBoolean:
		return hclTypeSQLiteInteger // SQLite uses 0/1 for booleans
	case SchemaTypeObject:
		// For nested objects, use TEXT (stored as JSON)
		return hclTypeSQLiteText
	default:
		return hclTypeSQLiteText
	}
}

// HCLDefault returns the formatted default value for HCL.
func (s *Schema) HCLDefault() string {
	field := s
	if field.Default == nil {
		return ""
	}

	// Handle special SQL functions
	switch field.Format {
	case FormatUUID:
		return sqlGenRandomUUID
	case FormatDateTime:
		return sqlCurrentTimestamp
	}

	// Format based on type
	switch v := field.Default.(type) {
	case string:
		// Check if it's a SQL expression
		if strings.HasPrefix(v, "sql(") {
			return v
		}
		// For arrays represented as strings
		if field.Type.PrimaryType() == SchemaTypeArray && v == "{}" {
			return "\"{}\""
		}
		// Regular string default
		return "\"" + v + "\""
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case []any:
		// Empty array default: [] -> "{}" for PostgreSQL
		if len(v) == 0 && field.Type.PrimaryType() == SchemaTypeArray {
			return "\"{}\""
		}
		return ""
	default:
		// Convert to string to analyze the value
		str := fmt.Sprintf("%v", v)

		// Extract value from YAML struct representation
		// Pattern: "&{... !!str VALUE ...}" or "&{... !!int VALUE ...}"
		if strings.Contains(str, "!!str ") {
			// Find the value after !!str
			parts := strings.Split(str, "!!str ")
			if len(parts) > 1 {
				// Extract the value (it's the next word after !!str)
				valuePart := strings.TrimSpace(parts[1])
				words := strings.Fields(valuePart)
				if len(words) > 0 {
					value := words[0]
					// For string values, wrap in quotes
					return "\"" + value + "\""
				}
			}
		}

		if strings.Contains(str, "!!int ") {
			// Find the value after !!int
			parts := strings.Split(str, "!!int ")
			if len(parts) > 1 {
				// Extract the value (it's the next word after !!int)
				valuePart := strings.TrimSpace(parts[1])
				words := strings.Fields(valuePart)
				if len(words) > 0 {
					return words[0]
				}
			}
		}

		// Handle simple type detection from the formatted string
		if field.Type.PrimaryType() == SchemaTypeBoolean {
			if strings.Contains(str, "false") {
				return "false"
			}
			if strings.Contains(str, "true") {
				return "true"
			}
		}

		if field.Type.PrimaryType() == SchemaTypeInteger {
			// Check for zero value
			if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
				return "0"
			}
		}

		if field.Type.PrimaryType() == SchemaTypeArray {
			if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
				return "\"{}\""
			}
		}

		// Don't output malformed defaults
		return ""
	}
}

// SQLiteHCLDefault returns the formatted default value for HCL/SQLite.
func (s *Schema) SQLiteHCLDefault() string {
	field := s
	if field.Default == nil {
		return ""
	}

	// Handle special SQL functions for SQLite
	switch field.Format {
	case FormatUUID:
		return sqlRandomBlob
	case FormatDateTime:
		return sqlCurrentTimestamp
	}

	// Format based on type
	switch v := field.Default.(type) {
	case string:
		return formatSQLiteStringDefault(v, field)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case []any:
		// Empty array default: [] -> "[]" for SQLite (stored as JSON)
		if len(v) == 0 && field.Type.PrimaryType() == SchemaTypeArray {
			return "\"[]\""
		}
		return ""
	default:
		return formatSQLiteComplexDefault(field)
	}
}

// formatSQLiteStringDefault formats string defaults for SQLite.
func formatSQLiteStringDefault(v string, field *Schema) string {
	if strings.HasPrefix(v, "sql(") {
		return v
	}
	if field.Type.PrimaryType() == SchemaTypeArray && v == "{}" {
		return "\"[]\""
	}
	return "\"" + v + "\""
}

// formatSQLiteComplexDefault handles complex default value types for SQLite.
func formatSQLiteComplexDefault(field *Schema) string {
	str := fmt.Sprintf("%v", field.Default)

	// Extract value from YAML struct representation
	if strings.Contains(str, "!!str ") {
		if val := extractYAMLValue(str, "!!str "); val != "" {
			return "\"" + val + "\""
		}
	}

	if strings.Contains(str, "!!int ") {
		if val := extractYAMLValue(str, "!!int "); val != "" {
			return val
		}
	}

	// Handle simple type detection
	switch field.Type.PrimaryType() {
	case SchemaTypeBoolean:
		if strings.Contains(str, "false") {
			return "0"
		}
		if strings.Contains(str, "true") {
			return "1"
		}
	case SchemaTypeInteger:
		if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
			return "0"
		}
	case SchemaTypeArray:
		if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
			return "\"[]\""
		}
	}

	return ""
}

// extractYAMLValue extracts a value from YAML-formatted string.
func extractYAMLValue(str, marker string) string {
	parts := strings.Split(str, marker)
	if len(parts) > 1 {
		valuePart := strings.TrimSpace(parts[1])
		words := strings.Fields(valuePart)
		if len(words) > 0 {
			return words[0]
		}
	}
	return ""
}
