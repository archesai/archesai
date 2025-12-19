package codegen

import (
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/schema"
	"github.com/archesai/archesai/internal/spec"
)

// GroupSQLite is the generator group for SQLite repositories.
const GroupSQLite = "sqlite"

const (
	genSQLite        = "sqlite_repository"
	genSQLiteDB      = "sqlite_db"
	genSQLiteQueries = "db_queries"
)

// DatabaseTemplateData holds the data for rendering database repository templates.
// Used by database generators to generate concrete implementations.
type DatabaseTemplateData struct {
	Entity          *schema.Schema
	ProjectName     string
	ModelImportPath string
}

// getDatabaseImportPath returns the schema import path for database generators.
func getDatabaseImportPath(s *spec.Spec, sch *schema.Schema) string {
	internalContext := s.InternalContext()
	if sch.IsInternal(internalContext) && sch.XInternal != "" {
		return spec.InternalPackageBase + "/" + sch.XInternal + "/schemas"
	}
	return s.ProjectName + "/schemas"
}

// generateSQLite generates SQLite repository implementations.
func (c *Codegen) generateSQLite(s *spec.Spec) error {
	for _, sch := range s.AllEntitySchemas() {
		path := filepath.Join(
			"database",
			"sqlite",
			"repositories",
			strings.ToLower(sch.Title)+"_repository.gen.go",
		)
		data := &DatabaseTemplateData{
			Entity:          sch,
			ProjectName:     s.ProjectName,
			ModelImportPath: getDatabaseImportPath(s, sch),
		}

		if err := c.RenderToFile(genSQLite+".go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// generateSQLiteDB generates the SQLite database setup file.
func (c *Codegen) generateSQLiteDB(s *spec.Spec) error {
	if len(s.AllEntitySchemas()) == 0 {
		return nil
	}

	path := filepath.Join("database", "sqlite", "repositories", "db.gen.go")
	return c.RenderToFile(genSQLiteDB+".go.tmpl", path, nil)
}

// generateSQLiteQueries generates SQL query files for SQLite.
func (c *Codegen) generateSQLiteQueries(s *spec.Spec) error {
	for _, sch := range s.AllEntitySchemas() {
		path := filepath.Join(
			"database",
			"sqlite",
			"queries",
			strings.ToLower(sch.Title)+"s.gen.sql",
		)
		data := &DatabaseTemplateData{
			Entity:          sch,
			ProjectName:     s.ProjectName,
			ModelImportPath: getDatabaseImportPath(s, sch),
		}

		if err := c.RenderToFile(genSQLiteQueries+".sql.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}
