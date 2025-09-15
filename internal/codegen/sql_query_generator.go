package codegen

import (
	"fmt"
	"strings"
)

const (
	fieldUpdatedAt = "updatedAt"
	fieldCreatedAt = "createdAt"
	fieldID        = "id"
)

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
			// Required fields: accountId, userId, providerId from YAML map to
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
	Sql         string                      `json:"sql" yaml:"sql"` //nolint:revive // matches generated type
	Type        XCodegenDatabaseQueriesType `json:"type" yaml:"type"`
}) string {
	var b strings.Builder

	if query.Description != "" {
		b.WriteString(fmt.Sprintf("-- %s\n", query.Description))
	}

	sqlcType := ":" + string(query.Type)
	b.WriteString(fmt.Sprintf("-- name: %s %s\n", query.Name, sqlcType))
	sql := strings.TrimSpace(query.Sql)
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
