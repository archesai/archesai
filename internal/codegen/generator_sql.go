package codegen

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/templates"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
)

// Common constants for SQL generation.
const (
	fieldUpdatedAt = "updatedAt"
	fieldCreatedAt = "createdAt"
	fieldID        = "id"
)

const (
	// PostgreSQLDialect generates PostgreSQL-compatible SQL.
	PostgreSQLDialect Dialect = "postgresql"
	// SQLiteDialect generates SQLite-compatible SQL.
	SQLiteDialect Dialect = "sqlite"
)

// Dialect represents the SQL dialect to use for generation.
type Dialect string

// SQLGenerator handles generation of SQL schemas and queries.
type SQLGenerator struct {
	Parser *parsers.Parser
	logger *slog.Logger
}

// NewSQLGenerator creates a new SQL generator.
func NewSQLGenerator(parser *parsers.Parser, logger *slog.Logger) Generator {
	return &SQLGenerator{
		Parser: parser,
		logger: logger,
	}
}

// Generate implements the Generator interface.
func (g *SQLGenerator) Generate(ctx *GeneratorContext) error {
	g.logger.Debug("Running SQL generator")

	if ctx.Config.Generators.SQL == nil ||
		(ctx.Config.Generators.SQL.SchemaDir == nil && ctx.Config.Generators.SQL.QueryDir == nil) {
		return nil
	}

	sqlConfig := ctx.Config.Generators.SQL
	dialect := *sqlConfig.Dialect
	if dialect == "" {
		dialect = "postgresql" // Default to PostgreSQL
	}

	// Create sub-generators
	schemaGen := &SQLSchemaGenerator{
		dialect: Dialect(dialect),
		parser:  g.Parser,
	}
	queryGen := NewSQLQueryGenerator(Dialect(dialect))

	// Process each schema that has database configuration
	for _, parsedSchema := range ctx.Schemas {
		ext, _ := ExtractCodegenExtension(parsedSchema)
		if ext == nil || ext.Database == nil ||
			ext.Database.Table == nil || *ext.Database.Table == "" {
			continue // Skip schemas without database configuration
		}

		g.logger.Debug("Generating SQL for schema", slog.String("schema", parsedSchema.Name))

		// Get table name for file naming
		var tableName string
		if ext.Database.Table != nil && *ext.Database.Table != "" {
			tableName = *ext.Database.Table
		} else {
			tableName = templates.SnakeCase(parsedSchema.Name)
		}

		// Generate CREATE TABLE statement if schema directory is configured
		if sqlConfig.SchemaDir != nil {
			schema, err := schemaGen.GenerateCreateTable(parsedSchema)
			if err != nil {
				g.logger.Warn(
					"Failed to generate schema",
					slog.String("schema", parsedSchema.Name),
					slog.String("error", err.Error()),
				)
				continue
			}

			schemaPath := fmt.Sprintf("%s/%s.sql", sqlConfig.SchemaDir, tableName)
			if err := ctx.FileWriter.WriteFile(schemaPath, []byte(schema)); err != nil {
				return fmt.Errorf("failed to write schema for %s: %w", parsedSchema.Name, err)
			}
		}

		// Generate SQLC queries if query directory is configured
		if sqlConfig.QueryDir != nil {
			queries, err := queryGen.GenerateQueries(parsedSchema)
			if err != nil {
				g.logger.Warn(
					"Failed to generate queries",
					slog.String("schema", parsedSchema.Name),
					slog.String("error", err.Error()),
				)
				continue
			}

			queryPath := fmt.Sprintf("%s/%s.sql", sqlConfig.QueryDir, tableName)
			if err := ctx.FileWriter.WriteFile(queryPath, []byte(queries)); err != nil {
				return fmt.Errorf("failed to write queries for %s: %w", parsedSchema.Name, err)
			}
		}
	}

	return nil
}

// SQLSchemaGenerator generates SQL schema (CREATE TABLE) statements from OpenAPI schemas.
type SQLSchemaGenerator struct {
	dialect Dialect
	parser  *parsers.Parser
}

// GenerateCreateTable generates a CREATE TABLE statement for a schema.
func (g *SQLSchemaGenerator) GenerateCreateTable(schema *parsers.JSONSchema) (string, error) {
	ext, _ := ExtractCodegenExtension(schema)
	if ext == nil || ext.Database == nil || ext.Database.Table == nil || *ext.Database.Table == "" {
		return "", fmt.Errorf("schema %s has no database configuration", schema.Name)
	}

	var tableName string
	if ext.Database.Table != nil && *ext.Database.Table != "" {
		tableName = *ext.Database.Table
	} else {
		tableName = templates.SnakeCase(schema.Name)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

	columns := []string{}

	// Add id column (from Base.yaml)
	columns = append(
		columns,
		g.generateColumn("id", "UUID", false, "PRIMARY KEY DEFAULT gen_random_uuid()"),
	)

	// Add timestamp columns (from Base.yaml)
	columns = append(
		columns,
		g.generateColumn("created_at", "TIMESTAMPTZ", false, "DEFAULT CURRENT_TIMESTAMP"),
	)
	columns = append(
		columns,
		g.generateColumn("updated_at", "TIMESTAMPTZ", false, "DEFAULT CURRENT_TIMESTAMP"),
	)

	// Add columns from properties
	if schema.Properties != nil {
		for name := range schema.Properties.Keys() {
			propRef := schema.Properties.GetOrZero(name)
			if name == "id" || name == "createdAt" ||
				name == "updatedAt" { //nolint:goconst // field names
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
	if ext.Database.Relations != nil {
		for _, rel := range ext.Database.Relations {
			constraint := g.generateForeignKeyConstraint(rel)
			if constraint != "" {
				columns = append(columns, constraint)
			}
		}
	}

	b.WriteString("    " + strings.Join(columns, ",\n    "))
	b.WriteString("\n);\n")

	return b.String(), nil
}

// generateColumn generates a column definition.
func (g *SQLSchemaGenerator) generateColumn(
	name, sqlType string,
	nullable bool,
	extra string,
) string {
	col := fmt.Sprintf("%s %s", name, sqlType)
	if !nullable {
		col += " NOT NULL"
	}
	if extra != "" {
		col += " " + extra
	}
	return col
}

// generateColumnFromProperty generates a column definition from an OpenAPI property.
func (g *SQLSchemaGenerator) generateColumnFromProperty(
	name string,
	prop *oas3.Schema,
	required []string,
) string {
	// Infer from OpenAPI type
	columnName := templates.SnakeCase(name)
	schema := parsers.NewJSONSchema(prop)
	sqlType := schema.SchemaToSQLType(prop, string(g.dialect))
	nullable := !templates.Contains(required, name)
	extra := ""

	return g.generateColumn(columnName, sqlType, nullable, extra)
}

// generateForeignKeyConstraint generates a foreign key constraint.
func (g *SQLSchemaGenerator) generateForeignKeyConstraint(
	rel struct {
		Field      string  `json:"field" yaml:"field"`
		OnDelete   *string `json:"onDelete,omitempty" yaml:"onDelete,omitempty"`
		OnUpdate   *string `json:"onUpdate,omitempty" yaml:"onUpdate,omitempty"`
		References string  `json:"references" yaml:"references"`
	},
) string {
	columnName := templates.SnakeCase(rel.Field)
	parts := strings.Split(rel.References, ".")
	if len(parts) != 2 {
		return ""
	}
	refTable := parts[0]
	refColumn := parts[1]

	constraint := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", columnName, refTable, refColumn)
	if rel.OnDelete != nil && *rel.OnDelete != "" {
		constraint += " ON DELETE " + *rel.OnDelete
	}
	if rel.OnUpdate != nil && *rel.OnUpdate != "" {
		constraint += " ON UPDATE " + *rel.OnUpdate
	}
	return constraint
}

// GenerateIndices generates CREATE INDEX statements.
func (g *SQLSchemaGenerator) GenerateIndices(schema *parsers.JSONSchema) ([]string, error) {
	ext, _ := ExtractCodegenExtension(schema)
	if ext == nil || ext.Database == nil || ext.Database.Table == nil || *ext.Database.Table == "" {
		return nil, nil
	}

	var tableName string
	if ext.Database.Table != nil && *ext.Database.Table != "" {
		tableName = *ext.Database.Table
	} else {
		tableName = templates.SnakeCase(schema.Name)
	}

	var indices []string
	if ext.Database.Indices != nil {
		for _, idx := range ext.Database.Indices {
			var indexName string
			if idx.Name != nil && *idx.Name != "" {
				indexName = *idx.Name
			} else {
				indexName = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(idx.Fields, "_"))
			}

			stmt := "CREATE "
			if idx.Unique != nil && *idx.Unique {
				stmt += "UNIQUE "
			}
			stmt += fmt.Sprintf(
				"INDEX %s ON %s (%s)",
				indexName,
				tableName,
				strings.Join(idx.Fields, ", "),
			)

			if idx.Where != nil && *idx.Where != "" {
				stmt += " WHERE " + *idx.Where
			}
			stmt += ";"

			indices = append(indices, stmt)
		}
	}

	return indices, nil
}

// SQLQueryGenerator generates SQLC query files from OpenAPI schemas.
type SQLQueryGenerator struct {
	dialect Dialect // "postgresql" or "sqlite"
}

// NewSQLQueryGenerator creates a new SQL query generator.
func NewSQLQueryGenerator(dialect Dialect) *SQLQueryGenerator {
	return &SQLQueryGenerator{
		dialect: dialect,
	}
}

// GenerateQueries generates all SQL queries for a schema.
func (g *SQLQueryGenerator) GenerateQueries(schema *parsers.JSONSchema) (string, error) {
	ext, err := ExtractCodegenExtension(schema)
	if err != nil {
		return "", fmt.Errorf("failed to extract x-codegen for %s: %w", schema.Name, err)
	}

	var tableName string
	if ext.Database.Table != nil && *ext.Database.Table != "" {
		tableName = *ext.Database.Table
	} else {
		tableName = templates.SnakeCase(schema.Name)
	}

	var b strings.Builder

	// Generate standard CRUD queries based on repository operations
	for _, op := range ext.Repository.Operations {
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
	if len(ext.Repository.AdditionalMethods) > 0 {
		for _, method := range ext.Repository.AdditionalMethods {
			query := g.generateAdditionalMethodQuery(schema, tableName, method)
			if query != "" {
				b.WriteString(query)
				b.WriteString("\n\n")
			}
		}
	}

	// Generate custom queries from database configuration
	if ext.Database.Queries != nil {
		for _, query := range ext.Database.Queries {
			b.WriteString(g.formatCustomQuery(query))
			b.WriteString("\n\n")
		}
	}

	return b.String(), nil
}

// generateCreateQuery generates an INSERT query.
func (g *SQLQueryGenerator) generateCreateQuery(
	schema *parsers.JSONSchema,
	tableName string,
) string {
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
			if columnName == "account_id" || columnName == "user_id" ||
				columnName == "provider_id" {
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

// generateGetQuery generates a SELECT by ID query.
func (g *SQLQueryGenerator) generateGetQuery(schema *parsers.JSONSchema, tableName string) string {
	return fmt.Sprintf(`-- name: Get%s :one
SELECT * FROM %s
WHERE id = $1 LIMIT 1`, schema.Name, tableName)
}

// generateUpdateQuery generates an UPDATE query.
func (g *SQLQueryGenerator) generateUpdateQuery(
	schema *parsers.JSONSchema,
	tableName string,
) string {
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
			if !templates.Contains(schema.Required, name) {
				updates = append(
					updates,
					fmt.Sprintf(
						"    %s = COALESCE(sqlc.narg(%s), %s)",
						columnName,
						columnName,
						columnName,
					),
				)
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

// generateDeleteQuery generates a DELETE query.
func (g *SQLQueryGenerator) generateDeleteQuery(
	schema *parsers.JSONSchema,
	tableName string,
) string {
	return fmt.Sprintf(`-- name: Delete%s :exec
DELETE FROM %s
WHERE id = $1;`, schema.Name, tableName)
}

// generateListQuery generates a SELECT list query.
func (g *SQLQueryGenerator) generateListQuery(schema *parsers.JSONSchema, tableName string) string {
	return fmt.Sprintf(`-- name: List%ss :many
SELECT * FROM %s
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;`, schema.Name, tableName)
}

// generateAdditionalMethodQuery generates query for additional repository methods.
func (g *SQLQueryGenerator) generateAdditionalMethodQuery(
	_ *parsers.JSONSchema,
	tableName string,
	method interface{},
) string {
	// Type assert to map to access fields
	methodMap, ok := method.(map[string]interface{})
	if !ok {
		return ""
	}

	methodName, _ := methodMap["name"].(string)
	methodReturns, _ := methodMap["returns"].(string)
	methodQuery, _ := methodMap["query"].(string)

	// If custom query is provided, use it
	if methodQuery != "" {
		sqlcType := ":one"
		switch methodReturns {
		case "multiple":
			sqlcType = ":many"
		case "void":
			sqlcType = ":exec"
		case "exists", "count":
			sqlcType = ":one"
		}
		// Add semicolon if not already present
		query := methodQuery
		if !strings.HasSuffix(strings.TrimSpace(query), ";") {
			query = strings.TrimSpace(query) + ";"
		}
		return fmt.Sprintf("-- name: %s %s\n%s", methodName, sqlcType, query)
	}

	// Otherwise, try to infer the query from method name and params
	if strings.HasPrefix(methodName, "GetBy") {
		field := strings.TrimPrefix(methodName, "GetBy")
		columnName := templates.SnakeCase(field)
		return fmt.Sprintf(`-- name: %s :one
SELECT * FROM %s
WHERE %s = $1
LIMIT 1;`, methodName, tableName, columnName)
	}

	if strings.HasPrefix(methodName, "ListBy") {
		field := strings.TrimPrefix(methodName, "ListBy")
		columnName := templates.SnakeCase(field)
		return fmt.Sprintf(`-- name: %s :many
SELECT * FROM %s
WHERE %s = $1
ORDER BY created_at DESC;`, methodName, tableName, columnName)
	}

	if strings.HasPrefix(methodName, "DeleteBy") {
		field := strings.TrimPrefix(methodName, "DeleteBy")
		columnName := templates.SnakeCase(field)
		return fmt.Sprintf(`-- name: %s :exec
DELETE FROM %s
WHERE %s = $1;`, methodName, tableName, columnName)
	}

	return ""
}

// formatCustomQuery formats a custom query from database configuration.
func (g *SQLQueryGenerator) formatCustomQuery(query interface{}) string {
	// Type assert to map to access fields
	queryMap, ok := query.(map[string]interface{})
	if !ok {
		// Try to assert to the struct type directly
		queryStruct, ok := query.(struct {
			Description *string `json:"description,omitempty" yaml:"description,omitempty"`
			Name        string  `json:"name" yaml:"name"`
			SQL         string  `json:"sql" yaml:"sql"`
			Type        string  `json:"type" yaml:"type"`
		})
		if !ok {
			return ""
		}

		var b strings.Builder
		if queryStruct.Description != nil && *queryStruct.Description != "" {
			b.WriteString(fmt.Sprintf("-- %s\n", *queryStruct.Description))
		}

		sqlcType := ":" + queryStruct.Type
		b.WriteString(fmt.Sprintf("-- name: %s %s\n", queryStruct.Name, sqlcType))
		sql := strings.TrimSpace(queryStruct.SQL)
		// Add semicolon if not already present
		if !strings.HasSuffix(sql, ";") {
			sql += ";"
		}
		b.WriteString(sql)

		return b.String()
	}

	var b strings.Builder

	description, _ := queryMap["description"].(*string)
	if description != nil && *description != "" {
		b.WriteString(fmt.Sprintf("-- %s\n", *description))
	}

	name, _ := queryMap["name"].(string)
	queryType, _ := queryMap["type"].(string)
	querySQL, _ := queryMap["sql"].(string)

	sqlcType := ":" + queryType
	b.WriteString(fmt.Sprintf("-- name: %s %s\n", name, sqlcType))
	sql := strings.TrimSpace(querySQL)
	// Add semicolon if not already present
	if !strings.HasSuffix(sql, ";") {
		sql += ";"
	}
	b.WriteString(sql)

	return b.String()
}

// getColumnName gets the database column name for a property.
func (g *SQLQueryGenerator) getColumnName(name string) string {
	return templates.SnakeCase(name)
}
