package parsers

import (
	"fmt"
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

// SQL function constants
const (
	sqlCurrentTimestamp = "sql(\"CURRENT_TIMESTAMP\")"
	sqlGenRandomUUID    = "sql(\"gen_random_uuid()\")"
	sqlRandomBlob       = "sql(\"lower(hex(randomblob(16)))\")"
)

// HCL/PostgreSQL type constants
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

// HCL/SQLite type constants
const (
	hclTypeSQLiteText    = `sql("TEXT")`
	hclTypeSQLiteInteger = `sql("INTEGER")`
	hclTypeSQLiteReal    = `sql("REAL")`
	hclTypeSQLiteBlob    = `sql("BLOB")`
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

// Constants for various type mappings and validations
const (
	// SQL Dialects
	SQLDialectPostgres = "postgresql"
	SQLDialectSQLite   = "sqlite"

	// Default SQL Types
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

// SchemaToHCLType converts a SchemaDef to HCL/PostgreSQL type string
func SchemaToHCLType(field *SchemaDef) string {
	// Handle special formats first
	switch field.Format {
	case formatUUID:
		return hclTypeUUID
	case formatDateTime:
		return hclTypeTimestamptz
	case formatDate:
		return hclTypeDate
	case "time":
		return hclTypeTime
	case "ipv4", "ipv6":
		return hclTypeText
	case formatEmail, formatURI, formatHostname:
		return hclTypeText
	case "binary":
		return hclTypeBytea
	case formatInt32:
		return hclTypeInteger
	case formatInt64:
		return hclTypeBigint
	case formatFloat:
		return hclTypeReal
	case formatDouble:
		return hclTypeDoublePrecision
	}

	// Handle arrays
	if field.Type == schemaTypeArray && field.Items != nil {
		itemType := SchemaToHCLType(field.Items)
		// Remove sql() wrapper from itemType if present
		if strings.HasPrefix(itemType, `sql("`) && strings.HasSuffix(itemType, `")`) {
			itemType = itemType[5 : len(itemType)-2]
		}
		return `sql("` + itemType + `[]")`
	}

	// Handle basic types
	switch field.Type {
	case schemaTypeString:
		// Check for specific field names that need special types
		if strings.Contains(strings.ToLower(field.Name), "embedding") {
			return "sql(\"vector(1536)\")"
		}
		// Check for enums
		if len(field.Enum) > 0 {
			return hclTypeText
		}
		// For now, just use text for all strings
		// TODO: Add MaxLength support when available in SchemaDef
		return hclTypeText
	case schemaTypeInteger:
		if field.Format == formatInt64 {
			return hclTypeBigint
		}
		return hclTypeInteger
	case schemaTypeNumber:
		if field.Format == formatFloat {
			return hclTypeReal
		}
		return hclTypeNumeric
	case schemaTypeBoolean:
		return hclTypeBoolean
	case schemaTypeObject:
		// For nested objects, use JSONB
		return hclTypeJSONB
	default:
		return hclTypeText
	}
}

// SchemaToSQLiteHCLType converts a SchemaDef to HCL/SQLite type string
func SchemaToSQLiteHCLType(field *SchemaDef) string {
	// SQLite has a simpler type system: TEXT, INTEGER, REAL, BLOB
	// Handle special formats first
	switch field.Format {
	case formatUUID:
		return hclTypeSQLiteText // UUIDs stored as TEXT in SQLite
	case formatDateTime, formatDate, "time":
		return hclTypeSQLiteText // Timestamps stored as TEXT in SQLite
	case "ipv4", "ipv6":
		return hclTypeSQLiteText
	case formatEmail, formatURI, formatHostname:
		return hclTypeSQLiteText
	case "binary":
		return hclTypeSQLiteBlob
	case formatInt32, formatInt64:
		return hclTypeSQLiteInteger
	case formatFloat, formatDouble:
		return hclTypeSQLiteReal
	}

	// Handle arrays - SQLite stores arrays as TEXT (JSON)
	if field.Type == schemaTypeArray {
		return hclTypeSQLiteText
	}

	// Handle basic types
	switch field.Type {
	case schemaTypeString:
		// Enums, embeddings, and all strings stored as TEXT
		return hclTypeSQLiteText
	case schemaTypeInteger:
		return hclTypeSQLiteInteger
	case schemaTypeNumber:
		return hclTypeSQLiteReal
	case schemaTypeBoolean:
		return hclTypeSQLiteInteger // SQLite uses 0/1 for booleans
	case schemaTypeObject:
		// For nested objects, use TEXT (stored as JSON)
		return hclTypeSQLiteText
	default:
		return hclTypeSQLiteText
	}
}

// FormatHCLDefault formats default values for HCL
func FormatHCLDefault(field *SchemaDef) string {
	if field.DefaultValue == nil {
		return ""
	}

	// Handle special SQL functions
	switch field.Format {
	case formatUUID:
		return sqlGenRandomUUID
	case formatDateTime:
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
		if field.Type == schemaTypeArray && v == "{}" {
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
		if field.Type == schemaTypeBoolean {
			if strings.Contains(str, "false") {
				return "false"
			}
			if strings.Contains(str, "true") {
				return "true"
			}
		}

		if field.Type == schemaTypeInteger {
			// Check for zero value
			if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
				return "0"
			}
		}

		if field.Type == schemaTypeArray {
			if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
				return "\"{}\""
			}
		}

		// Don't output malformed defaults
		return ""
	}
}

// FormatSQLiteHCLDefault formats default values for HCL/SQLite
func FormatSQLiteHCLDefault(field *SchemaDef) string {
	if field.DefaultValue == nil {
		return ""
	}

	// Handle special SQL functions for SQLite
	switch field.Format {
	case formatUUID:
		return sqlRandomBlob
	case formatDateTime:
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
func formatSQLiteStringDefault(v string, field *SchemaDef) string {
	if strings.HasPrefix(v, "sql(") {
		return v
	}
	if field.Type == schemaTypeArray && v == "{}" {
		return "\"[]\""
	}
	return "\"" + v + "\""
}

// formatSQLiteComplexDefault handles complex default value types for SQLite
func formatSQLiteComplexDefault(field *SchemaDef) string {
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
	case schemaTypeBoolean:
		if strings.Contains(str, "false") {
			return "0"
		}
		if strings.Contains(str, "true") {
			return "1"
		}
	case schemaTypeInteger:
		if strings.Contains(str, " 0 ") || strings.Contains(str, " 0}") {
			return "0"
		}
	case schemaTypeArray:
		if strings.Contains(str, "[]") || strings.Contains(str, "{}") {
			return "\"[]\""
		}
	}

	return ""
}

// extractYAMLValue extracts a value from YAML-formatted string
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
