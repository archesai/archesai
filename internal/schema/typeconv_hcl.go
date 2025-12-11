package schema

import (
	"fmt"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
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

// HCL default value constants for empty arrays/objects.
const (
	hclDefaultEmptyObject = `"{}"`
	hclDefaultEmptyArray  = `"[]"`
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
	if field.Type.PrimaryType() == TypeArray && field.Items != nil {
		itemType := field.Items.GetOrNil().HCLType()
		// Remove sql() wrapper from itemType if present
		if strings.HasPrefix(itemType, `sql("`) && strings.HasSuffix(itemType, `")`) {
			itemType = itemType[5 : len(itemType)-2]
		}
		return `sql("` + itemType + `[]")`
	}

	// Handle basic types
	switch field.Type.PrimaryType() {
	case TypeString:
		// Check for specific field names that need special types
		if strings.Contains(strings.ToLower(field.Title), "embedding") {
			return "sql(\"vector(1536)\")"
		}
		// Check for enums
		if len(field.Enum) > 0 {
			return hclTypeText
		}
		// For now, just use text for all strings
		return hclTypeText
	case TypeInteger:
		if field.Format == FormatInt64 {
			return hclTypeBigint
		}
		return hclTypeInteger
	case TypeNumber:
		if field.Format == FormatFloat {
			return hclTypeReal
		}
		return hclTypeNumeric
	case TypeBoolean:
		return hclTypeBoolean
	case TypeObject:
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
	if field.Type.PrimaryType() == TypeArray {
		return hclTypeSQLiteText
	}

	// Handle basic types
	switch field.Type.PrimaryType() {
	case TypeString:
		// Enums, embeddings, and all strings stored as TEXT
		return hclTypeSQLiteText
	case TypeInteger:
		return hclTypeSQLiteInteger
	case TypeNumber:
		return hclTypeSQLiteReal
	case TypeBoolean:
		return hclTypeSQLiteInteger // SQLite uses 0/1 for booleans
	case TypeObject:
		// For nested objects, use TEXT (stored as JSON)
		return hclTypeSQLiteText
	default:
		return hclTypeSQLiteText
	}
}

// hclDefaultFromYAMLNode extracts and formats a default value from a yaml.Node.
func hclDefaultFromYAMLNode(node *yaml.Node, primaryType string) string {
	if node == nil {
		return ""
	}

	switch node.Tag {
	case "!!str":
		return "\"" + node.Value + "\""
	case "!!int":
		return node.Value
	case "!!float":
		return node.Value
	case "!!bool":
		return node.Value
	}

	// Fallback: infer from schema type
	switch primaryType {
	case TypeBoolean:
		if node.Value == "false" || node.Value == "true" {
			return node.Value
		}
	case TypeInteger:
		if node.Value == "0" {
			return "0"
		}
	case TypeArray:
		if node.Kind == yaml.SequenceNode && len(node.Content) == 0 {
			return hclDefaultEmptyObject
		}
	}

	return ""
}

// HCLDefault returns the formatted default value for HCL.
func (s *Schema) HCLDefault() string {
	if s.Default == nil {
		return ""
	}

	// Handle special SQL functions
	if s.Format == FormatUUID {
		return sqlGenRandomUUID
	}
	if s.Format == FormatDateTime {
		return sqlCurrentTimestamp
	}

	return s.formatHCLDefault()
}

// formatHCLDefault formats the default value based on its Go type.
func (s *Schema) formatHCLDefault() string {
	switch v := s.Default.(type) {
	case string:
		return s.formatStringDefault(v)
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
		if len(v) == 0 && s.Type.PrimaryType() == TypeArray {
			return hclDefaultEmptyObject
		}
		return ""
	default:
		return s.formatUnknownDefault()
	}
}

// formatStringDefault formats a string default value.
func (s *Schema) formatStringDefault(v string) string {
	if strings.HasPrefix(v, "sql(") {
		return v
	}
	if s.Type.PrimaryType() == TypeArray && v == "{}" {
		return hclDefaultEmptyObject
	}
	return "\"" + v + "\""
}

// formatUnknownDefault handles unknown types including yaml.Node.
func (s *Schema) formatUnknownDefault() string {
	// Handle yaml.Node directly
	if node, ok := s.Default.(*yaml.Node); ok {
		return hclDefaultFromYAMLNode(node, s.Type.PrimaryType())
	}

	// For other unknown types, try to infer from string representation
	str := fmt.Sprintf("%v", s.Default)
	switch s.Type.PrimaryType() {
	case TypeBoolean:
		if strings.Contains(str, "false") {
			return "false"
		}
		if strings.Contains(str, "true") {
			return "true"
		}
	case TypeInteger:
		if strings.Contains(str, "0") {
			return "0"
		}
	case TypeArray:
		if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
			return hclDefaultEmptyObject
		}
	}

	return ""
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
		if len(v) == 0 && field.Type.PrimaryType() == TypeArray {
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
	if field.Type.PrimaryType() == TypeArray && v == "{}" {
		return hclDefaultEmptyArray
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
	case TypeBoolean:
		if strings.Contains(str, "false") {
			return "0"
		}
		if strings.Contains(str, "true") {
			return "1"
		}
	case TypeInteger:
		if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
			return "0"
		}
	case TypeArray:
		if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
			return hclDefaultEmptyArray
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
