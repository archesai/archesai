package typeconv

import (
	"strings"
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
