package generators

import (
	"path/filepath"
	"strings"
)

// GeneratePostgres generates PostgreSQL repository implementations.
func GeneratePostgres(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		path := filepath.Join(
			"database",
			"postgres",
			"repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		data := &DatabaseTemplateData{
			Entity:          schema,
			ProjectName:     ctx.ProjectName,
			ModelImportPath: getDatabaseImportPaths(ctx, schema),
		}

		if err := ctx.RenderToFile("postgres.go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GeneratePostgresQueries generates SQL query files for PostgreSQL.
func GeneratePostgresQueries(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		path := filepath.Join(
			"database",
			"postgres",
			"queries",
			strings.ToLower(schema.Name)+"s.gen.sql",
		)
		data := &DatabaseTemplateData{
			Entity:          schema,
			ProjectName:     ctx.ProjectName,
			ModelImportPath: getDatabaseImportPaths(ctx, schema),
		}

		if err := ctx.RenderToFile("sql_queries.sql.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateSQLite generates SQLite repository implementations.
func GenerateSQLite(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		path := filepath.Join(
			"database",
			"sqlite",
			"repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		data := &DatabaseTemplateData{
			Entity:          schema,
			ProjectName:     ctx.ProjectName,
			ModelImportPath: getDatabaseImportPaths(ctx, schema),
		}

		if err := ctx.RenderToFile("sqlite.go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateSQLiteDB generates the SQLite database setup file.
func GenerateSQLiteDB(ctx *GeneratorContext) error {
	if len(ctx.AllEntitySchemas()) == 0 {
		return nil
	}

	path := filepath.Join("database", "sqlite", "repositories", "db.gen.go")
	return ctx.RenderToFile("sqlite_db.go.tmpl", path, nil)
}

// GenerateSQLiteQueries generates SQL query files for SQLite.
func GenerateSQLiteQueries(ctx *GeneratorContext) error {
	for _, schema := range ctx.AllEntitySchemas() {
		path := filepath.Join(
			"database",
			"sqlite",
			"queries",
			strings.ToLower(schema.Name)+"s.gen.sql",
		)
		data := &DatabaseTemplateData{
			Entity:          schema,
			ProjectName:     ctx.ProjectName,
			ModelImportPath: getDatabaseImportPaths(ctx, schema),
		}

		if err := ctx.RenderToFile("sql_queries.sql.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}
