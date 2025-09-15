package codegen

import (
	"fmt"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
)

// Common constants for SQL generation
const (
	fieldUpdatedAt = "updatedAt"
	fieldCreatedAt = "createdAt"
	fieldID        = "id"
)

// Dialect represents the SQL dialect to use for generation
type Dialect string

const (
	// PostgreSQLDialect generates PostgreSQL-compatible SQL
	PostgreSQLDialect Dialect = "postgresql"
	// SQLiteDialect generates SQLite-compatible SQL
	SQLiteDialect Dialect = "sqlite"
)

// SQLSchemaGenerator generates SQL schema (CREATE TABLE) statements from OpenAPI schemas
type SQLSchemaGenerator struct {
	dialect Dialect
}

// NewSQLSchemaGenerator creates a new SQL schema generator
func NewSQLSchemaGenerator(dialect string) *SQLSchemaGenerator {
	return &SQLSchemaGenerator{
		dialect: Dialect(dialect),
	}
}

// GenerateCreateTable generates a CREATE TABLE statement for a schema
func (g *SQLSchemaGenerator) GenerateCreateTable(schema *ParsedSchema) (string, error) {
	if schema.XCodegen == nil || schema.XCodegen.Database.Table == "" {
		return "", fmt.Errorf("schema %s has no database configuration", schema.Name)
	}

	tableName := schema.XCodegen.Database.Table
	if tableName == "" {
		tableName = ToSnakeCase(schema.Name)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

	columns := []string{}

	// Add id column (from Base.yaml)
	columns = append(columns, g.generateColumn("id", "UUID", false, "PRIMARY KEY DEFAULT gen_random_uuid()"))

	// Add timestamp columns (from Base.yaml)
	columns = append(columns, g.generateColumn("created_at", "TIMESTAMPTZ", false, "DEFAULT CURRENT_TIMESTAMP"))
	columns = append(columns, g.generateColumn("updated_at", "TIMESTAMPTZ", false, "DEFAULT CURRENT_TIMESTAMP"))

	// Add columns from properties
	if schema.Properties != nil {
		for name := range schema.Properties.Keys() {
			propRef := schema.Properties.GetOrZero(name)
			if name == "id" || name == "createdAt" || name == "updatedAt" { //nolint:goconst // field names
				continue // Already added from Base
			}

			if propRef != nil && propRef.IsLeft() {
				prop := propRef.GetLeft()
				columnDef := g.generateColumnFromProperty(name, prop, schema.Required)
				if columnDef != "" {
					columns = append(columns, columnDef)
				}
			}
		}
	}

	// Add foreign key constraints
	for _, rel := range schema.XCodegen.Database.Relations {
		constraint := g.generateForeignKeyConstraint(rel)
		if constraint != "" {
			columns = append(columns, constraint)
		}
	}

	b.WriteString("    " + strings.Join(columns, ",\n    "))
	b.WriteString("\n);\n")

	return b.String(), nil
}

// generateColumn generates a column definition
func (g *SQLSchemaGenerator) generateColumn(name, sqlType string, nullable bool, extra string) string {
	col := fmt.Sprintf("%s %s", name, sqlType)
	if !nullable {
		col += " NOT NULL"
	}
	if extra != "" {
		col += " " + extra
	}
	return col
}

// generateColumnFromProperty generates a column definition from an OpenAPI property
func (g *SQLSchemaGenerator) generateColumnFromProperty(name string, prop *oas3.Schema, required []string) string {
	// Infer from OpenAPI type
	columnName := ToSnakeCase(name)
	sqlType := g.inferSQLType(prop)
	nullable := !Contains(required, name)
	extra := ""

	return g.generateColumn(columnName, sqlType, nullable, extra)
}

// inferSQLType infers SQL type from OpenAPI property
func (g *SQLSchemaGenerator) inferSQLType(prop *oas3.Schema) string {
	// Get the types array from the schema
	types := prop.GetType()
	if len(types) == 0 {
		return "TEXT" //nolint:goconst // default SQL type
	}

	// Use the first type (most schemas have only one type)
	schemaType := string(types[0])

	if schemaType == "string" { //nolint:goconst // OpenAPI type
		if prop.Format != nil && *prop.Format == "date-time" {
			if g.dialect == PostgreSQLDialect {
				return "TIMESTAMPTZ"
			}
			return "TIMESTAMP"
		}
		if prop.Format != nil && *prop.Format == "uuid" {
			return "UUID"
		}
		if prop.MaxLength != nil && *prop.MaxLength > 0 && *prop.MaxLength <= 255 {
			return fmt.Sprintf("VARCHAR(%d)", *prop.MaxLength)
		}
		return "TEXT" //nolint:goconst // SQL type
	}
	if schemaType == "integer" {
		if prop.Format != nil && *prop.Format == "int64" {
			return "BIGINT"
		}
		return "INTEGER"
	}
	if schemaType == "number" {
		return "DECIMAL"
	}
	if schemaType == "boolean" {
		return "BOOLEAN"
	}
	if schemaType == "array" {
		// PostgreSQL array type
		if g.dialect == PostgreSQLDialect {
			return "TEXT[]"
		}
		return "TEXT" // SQLite doesn't have native arrays //nolint:goconst
	}
	return "TEXT" //nolint:goconst // default SQL type
}

// generateForeignKeyConstraint generates a foreign key constraint
func (g *SQLSchemaGenerator) generateForeignKeyConstraint(rel struct {
	Field      string                            `json:"field" yaml:"field"`
	OnDelete   XCodegenDatabaseRelationsOnDelete `json:"onDelete,omitempty,omitzero" yaml:"onDelete,omitempty"`
	OnUpdate   XCodegenDatabaseRelationsOnUpdate `json:"onUpdate,omitempty,omitzero" yaml:"onUpdate,omitempty"`
	References string                            `json:"references" yaml:"references"`
}) string {
	columnName := ToSnakeCase(rel.Field)
	parts := strings.Split(rel.References, ".")
	if len(parts) != 2 {
		return ""
	}
	refTable := parts[0]
	refColumn := parts[1]

	constraint := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", columnName, refTable, refColumn)
	if rel.OnDelete != "" {
		constraint += " ON DELETE " + string(rel.OnDelete)
	}
	if rel.OnUpdate != "" {
		constraint += " ON UPDATE " + string(rel.OnUpdate)
	}
	return constraint
}

// GenerateIndices generates CREATE INDEX statements
func (g *SQLSchemaGenerator) GenerateIndices(schema *ParsedSchema) ([]string, error) {
	if schema.XCodegen == nil || schema.XCodegen.Database.Table == "" {
		return nil, nil
	}

	tableName := schema.XCodegen.Database.Table
	if tableName == "" {
		tableName = ToSnakeCase(schema.Name)
	}

	var indices []string
	for _, idx := range schema.XCodegen.Database.Indices {
		indexName := idx.Name
		if indexName == "" {
			indexName = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(idx.Fields, "_"))
		}

		stmt := "CREATE "
		if idx.Unique {
			stmt += "UNIQUE "
		}
		stmt += fmt.Sprintf("INDEX %s ON %s (%s)", indexName, tableName, strings.Join(idx.Fields, ", "))

		if idx.Where != "" {
			stmt += " WHERE " + idx.Where
		}
		stmt += ";"

		indices = append(indices, stmt)
	}

	return indices, nil
}

// SQLQueryGenerator generates SQLC query files from OpenAPI schemas
type SQLQueryGenerator struct {
	dialect Dialect // "postgresql" or "sqlite"
}

// NewSQLQueryGenerator creates a new SQL query generator
func NewSQLQueryGenerator(dialect Dialect) *SQLQueryGenerator {
	return &SQLQueryGenerator{
		dialect: dialect,
	}
}

// GenerateQueries generates all SQL queries for a schema
func (g *SQLQueryGenerator) GenerateQueries(schema *ParsedSchema) (string, error) {
	if schema.XCodegen == nil {
		return "", fmt.Errorf("schema %s has no x-codegen configuration", schema.Name)
	}

	tableName := schema.XCodegen.Database.Table
	if tableName == "" {
		tableName = ToSnakeCase(schema.Name)
	}

	var b strings.Builder

	// Generate standard CRUD queries based on repository operations
	for _, op := range schema.XCodegen.Repository.Operations {
		switch op {
		case "create":
			b.WriteString(g.generateCreateQuery(schema, tableName))
			b.WriteString("\n\n")
		case "read":
			b.WriteString(g.generateGetQuery(schema, tableName))
			b.WriteString("\n\n")
		case "update":
			b.WriteString(g.generateUpdateQuery(schema, tableName))
			b.WriteString("\n\n")
		case "delete":
			b.WriteString(g.generateDeleteQuery(schema, tableName))
			b.WriteString("\n\n")
		case "list":
			b.WriteString(g.generateListQuery(schema, tableName))
			b.WriteString("\n\n")
		}
	}

	// Generate additional method queries
	if len(schema.XCodegen.Repository.AdditionalMethods) > 0 {
		for _, method := range schema.XCodegen.Repository.AdditionalMethods {
			query := g.generateAdditionalMethodQuery(schema, tableName, method)
			if query != "" {
				b.WriteString(query)
				b.WriteString("\n\n")
			}
		}
	}

	// Generate custom queries from database configuration
	if schema.XCodegen.Database.Queries != nil {
		for _, query := range schema.XCodegen.Database.Queries {
			b.WriteString(g.formatCustomQuery(query))
			b.WriteString("\n\n")
		}
	}

	return b.String(), nil
}

// generateCreateQuery generates an INSERT query
func (g *SQLQueryGenerator) generateCreateQuery(schema *ParsedSchema, tableName string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("-- name: Create%s :one\n", schema.Name))
	b.WriteString(fmt.Sprintf("INSERT INTO %s (\n", tableName))

	// Collect column names
	var columns []string
	var values []string
	paramNum := 1

	// Add id column
	columns = append(columns, "id")
	values = append(values, fmt.Sprintf("$%d", paramNum))
	paramNum++

	// Add property columns
	if schema.Properties != nil {
		for name := range schema.Properties.Keys() {
			if name == fieldID || name == fieldCreatedAt || name == fieldUpdatedAt {
				continue // Skip auto-generated fields
			}

			columnName := g.getColumnName(name)
			columns = append(columns, columnName)

			// For INSERT, use positional parameters for required fields
			// Required fields: accountId, userID, providerId from YAML map to
			// account_id, user_id, provider_id in database
			isRequired := false
			if columnName == "account_id" || columnName == "user_id" || columnName == "provider_id" {
				isRequired = true
			}

			if isRequired {
				values = append(values, fmt.Sprintf("$%d", paramNum))
				paramNum++
			} else {
				values = append(values, fmt.Sprintf("sqlc.narg('%s')", columnName))
			}
		}
	}

	b.WriteString("    " + strings.Join(columns, ",\n    "))
	b.WriteString("\n) VALUES (\n    ")
	b.WriteString(strings.Join(values, ", "))
	b.WriteString("\n)\nRETURNING *;")

	return b.String()
}

// generateGetQuery generates a SELECT by ID query
func (g *SQLQueryGenerator) generateGetQuery(schema *ParsedSchema, tableName string) string {
	return fmt.Sprintf(`-- name: Get%s :one
SELECT * FROM %s
WHERE id = $1 LIMIT 1`, schema.Name, tableName)
}

// generateUpdateQuery generates an UPDATE query
func (g *SQLQueryGenerator) generateUpdateQuery(schema *ParsedSchema, tableName string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("-- name: Update%s :one\n", schema.Name))
	b.WriteString(fmt.Sprintf("UPDATE %s\n", tableName))
	b.WriteString("SET \n")

	var updates []string
	if schema.Properties != nil {
		for name := range schema.Properties.Keys() {
			if name == fieldID || name == fieldCreatedAt {
				continue // Skip immutable fields
			}

			if name == fieldUpdatedAt {
				updates = append(updates, "    updated_at = NOW()")
				continue
			}

			columnName := g.getColumnName(name)
			// Use COALESCE for optional fields
			if !Contains(schema.Required, name) {
				updates = append(updates, fmt.Sprintf("    %s = COALESCE(sqlc.narg(%s), %s)", columnName, columnName, columnName))
			} else {
				updates = append(updates, fmt.Sprintf("    %s = sqlc.arg(%s)", columnName, columnName))
			}
		}
	}

	// If no updates were generated (e.g., for schemas with allOf), add a default updated_at
	if len(updates) == 0 {
		updates = append(updates, "    updated_at = NOW()")
	}

	b.WriteString(strings.Join(updates, ",\n"))
	b.WriteString("\nWHERE id = $1")
	b.WriteString("\nRETURNING *;")

	return b.String()
}

// generateDeleteQuery generates a DELETE query
func (g *SQLQueryGenerator) generateDeleteQuery(schema *ParsedSchema, tableName string) string {
	return fmt.Sprintf(`-- name: Delete%s :exec
DELETE FROM %s
WHERE id = $1;`, schema.Name, tableName)
}

// generateListQuery generates a SELECT list query
func (g *SQLQueryGenerator) generateListQuery(schema *ParsedSchema, tableName string) string {
	return fmt.Sprintf(`-- name: List%ss :many
SELECT * FROM %s
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;`, schema.Name, tableName)
}

// generateAdditionalMethodQuery generates query for additional repository methods
func (g *SQLQueryGenerator) generateAdditionalMethodQuery(_ *ParsedSchema, tableName string, method struct {
	Cache   bool                                       `json:"cache,omitempty,omitzero" yaml:"cache,omitempty"`
	Name    string                                     `json:"name" yaml:"name"`
	Params  []string                                   `json:"params" yaml:"params"`
	Query   string                                     `json:"query,omitempty,omitzero" yaml:"query,omitempty"`
	Returns XCodegenRepositoryAdditionalMethodsReturns `json:"returns" yaml:"returns"`
}) string {
	// If custom query is provided, use it
	if method.Query != "" {
		sqlcType := ":one"
		switch method.Returns {
		case "multiple":
			sqlcType = ":many"
		case "void":
			sqlcType = ":exec"
		case "exists", "count":
			sqlcType = ":one"
		}
		// Add semicolon if not already present
		query := method.Query
		if !strings.HasSuffix(strings.TrimSpace(query), ";") {
			query = strings.TrimSpace(query) + ";"
		}
		return fmt.Sprintf("-- name: %s %s\n%s", method.Name, sqlcType, query)
	}

	// Otherwise, try to infer the query from method name and params
	if strings.HasPrefix(method.Name, "GetBy") {
		field := strings.TrimPrefix(method.Name, "GetBy")
		columnName := ToSnakeCase(field)
		return fmt.Sprintf(`-- name: %s :one
SELECT * FROM %s
WHERE %s = $1
LIMIT 1;`, method.Name, tableName, columnName)
	}

	if strings.HasPrefix(method.Name, "ListBy") {
		field := strings.TrimPrefix(method.Name, "ListBy")
		columnName := ToSnakeCase(field)
		return fmt.Sprintf(`-- name: %s :many
SELECT * FROM %s
WHERE %s = $1
ORDER BY created_at DESC;`, method.Name, tableName, columnName)
	}

	if strings.HasPrefix(method.Name, "DeleteBy") {
		field := strings.TrimPrefix(method.Name, "DeleteBy")
		columnName := ToSnakeCase(field)
		return fmt.Sprintf(`-- name: %s :exec
DELETE FROM %s
WHERE %s = $1;`, method.Name, tableName, columnName)
	}

	return ""
}

// formatCustomQuery formats a custom query from database configuration
func (g *SQLQueryGenerator) formatCustomQuery(query struct {
	Description string                      `json:"description,omitempty,omitzero" yaml:"description,omitempty"`
	Name        string                      `json:"name" yaml:"name"`
	SQL         string                      `json:"sql" yaml:"sql"` //nolint:revive // matches generated type
	Type        XCodegenDatabaseQueriesType `json:"type" yaml:"type"`
}) string {
	var b strings.Builder

	if query.Description != "" {
		b.WriteString(fmt.Sprintf("-- %s\n", query.Description))
	}

	sqlcType := ":" + string(query.Type)
	b.WriteString(fmt.Sprintf("-- name: %s %s\n", query.Name, sqlcType))
	sql := strings.TrimSpace(query.SQL)
	// Add semicolon if not already present
	if !strings.HasSuffix(sql, ";") {
		sql += ";"
	}
	b.WriteString(sql)

	return b.String()
}

// getColumnName gets the database column name for a property
func (g *SQLQueryGenerator) getColumnName(name string) string {
	return ToSnakeCase(name)
}
