package typeconv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/archesai/archesai/internal/spec"
)

// SchemaToHCLType converts a Schema to HCL/PostgreSQL type string
func SchemaToHCLType(field *spec.Schema) string {
	// Handle special formats first
	switch field.Format {
	case spec.FormatUUID:
		return hclTypeUUID
	case spec.FormatDateTime:
		return hclTypeTimestamptz
	case spec.FormatDate:
		return hclTypeDate
	case "time":
		return hclTypeTime
	case "ipv4", "ipv6":
		return hclTypeText
	case spec.FormatEmail, spec.FormatURI, spec.FormatHostname:
		return hclTypeText
	case "binary":
		return hclTypeBytea
	case spec.FormatInt32:
		return hclTypeInteger
	case spec.FormatInt64:
		return hclTypeBigint
	case spec.FormatFloat:
		return hclTypeReal
	case spec.FormatDouble:
		return hclTypeDoublePrecision
	}

	// Handle arrays
	if field.Type == spec.SchemaTypeArray && field.Items != nil {
		itemType := SchemaToHCLType(field.Items)
		// Remove sql() wrapper from itemType if present
		if strings.HasPrefix(itemType, `sql("`) && strings.HasSuffix(itemType, `")`) {
			itemType = itemType[5 : len(itemType)-2]
		}
		return `sql("` + itemType + `[]")`
	}

	// Handle basic types
	switch field.Type {
	case spec.SchemaTypeString:
		// Check for specific field names that need special types
		if strings.Contains(strings.ToLower(field.Name), "embedding") {
			return "sql(\"vector(1536)\")"
		}
		// Check for enums
		if len(field.Enum) > 0 {
			return hclTypeText
		}
		// For now, just use text for all strings
		// TODO: Add MaxLength support when available in Schema
		return hclTypeText
	case spec.SchemaTypeInteger:
		if field.Format == spec.FormatInt64 {
			return hclTypeBigint
		}
		return hclTypeInteger
	case spec.SchemaTypeNumber:
		if field.Format == spec.FormatFloat {
			return hclTypeReal
		}
		return hclTypeNumeric
	case spec.SchemaTypeBoolean:
		return hclTypeBoolean
	case spec.SchemaTypeObject:
		// For nested objects, use JSONB
		return hclTypeJSONB
	default:
		return hclTypeText
	}
}

// SchemaToSQLiteHCLType converts a Schema to HCL/SQLite type string
func SchemaToSQLiteHCLType(field *spec.Schema) string {
	// SQLite has a simpler type system: TEXT, INTEGER, REAL, BLOB
	// Handle special formats first
	switch field.Format {
	case spec.FormatUUID:
		return hclTypeSQLiteText // UUIDs stored as TEXT in SQLite
	case spec.FormatDateTime, spec.FormatDate, "time":
		return hclTypeSQLiteText // Timestamps stored as TEXT in SQLite
	case "ipv4", "ipv6":
		return hclTypeSQLiteText
	case spec.FormatEmail, spec.FormatURI, spec.FormatHostname:
		return hclTypeSQLiteText
	case "binary":
		return hclTypeSQLiteBlob
	case spec.FormatInt32, spec.FormatInt64:
		return hclTypeSQLiteInteger
	case spec.FormatFloat, spec.FormatDouble:
		return hclTypeSQLiteReal
	}

	// Handle arrays - SQLite stores arrays as TEXT (JSON)
	if field.Type == spec.SchemaTypeArray {
		return hclTypeSQLiteText
	}

	// Handle basic types
	switch field.Type {
	case spec.SchemaTypeString:
		// Enums, embeddings, and all strings stored as TEXT
		return hclTypeSQLiteText
	case spec.SchemaTypeInteger:
		return hclTypeSQLiteInteger
	case spec.SchemaTypeNumber:
		return hclTypeSQLiteReal
	case spec.SchemaTypeBoolean:
		return hclTypeSQLiteInteger // SQLite uses 0/1 for booleans
	case spec.SchemaTypeObject:
		// For nested objects, use TEXT (stored as JSON)
		return hclTypeSQLiteText
	default:
		return hclTypeSQLiteText
	}
}

// FormatHCLDefault formats default values for HCL
func FormatHCLDefault(field *spec.Schema) string {
	if field.DefaultValue == nil {
		return ""
	}

	// Handle special SQL functions
	switch field.Format {
	case spec.FormatUUID:
		return sqlGenRandomUUID
	case spec.FormatDateTime:
		return sqlCurrentTimestamp
	}

	// Format based on type
	switch v := field.DefaultValue.(type) {
	case string:
		// Check if it's a SQL expression
		if strings.HasPrefix(v, "sql(") {
			return v
		}
		// For arrays represented as strings
		if field.Type == spec.SchemaTypeArray && v == "{}" {
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
		if field.Type == spec.SchemaTypeBoolean {
			if strings.Contains(str, "false") {
				return "false"
			}
			if strings.Contains(str, "true") {
				return "true"
			}
		}

		if field.Type == spec.SchemaTypeInteger {
			// Check for zero value
			if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
				return "0"
			}
		}

		if field.Type == spec.SchemaTypeArray {
			if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
				return "\"{}\""
			}
		}

		// Don't output malformed defaults
		return ""
	}
}

// FormatSQLiteHCLDefault formats default values for HCL/SQLite
func FormatSQLiteHCLDefault(field *spec.Schema) string {
	if field.DefaultValue == nil {
		return ""
	}

	// Handle special SQL functions for SQLite
	switch field.Format {
	case spec.FormatUUID:
		return sqlRandomBlob
	case spec.FormatDateTime:
		return sqlCurrentTimestamp
	}

	// Format based on type
	switch v := field.DefaultValue.(type) {
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
	default:
		return formatSQLiteComplexDefault(field)
	}
}

// formatSQLiteStringDefault formats string defaults for SQLite
func formatSQLiteStringDefault(v string, field *spec.Schema) string {
	if strings.HasPrefix(v, "sql(") {
		return v
	}
	if field.Type == spec.SchemaTypeArray && v == "{}" {
		return "\"[]\""
	}
	return "\"" + v + "\""
}

// formatSQLiteComplexDefault handles complex default value types for SQLite
func formatSQLiteComplexDefault(field *spec.Schema) string {
	str := fmt.Sprintf("%v", field.DefaultValue)

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
	switch field.Type {
	case spec.SchemaTypeBoolean:
		if strings.Contains(str, "false") {
			return "0"
		}
		if strings.Contains(str, "true") {
			return "1"
		}
	case spec.SchemaTypeInteger:
		if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
			return "0"
		}
	case spec.SchemaTypeArray:
		if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
			return "\"[]\""
		}
	}

	return ""
}
