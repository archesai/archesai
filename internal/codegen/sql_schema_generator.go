package codegen

import (
	"fmt"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
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
