// Package codegen provides template management and code generation for the archesai project.
package codegen

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/archesai/archesai/internal/parsers"
)

//go:embed tmpl/*.tmpl
var templatesFS embed.FS

// GetTemplate loads a template by name from the tmpl directory.
func GetTemplate(name string) (string, error) {
	content, err := templatesFS.ReadFile("tmpl/" + name)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", name, err)
	}
	return string(content), nil
}

// LoadTemplates loads all templates and returns them as a map
func LoadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Main templates
	templateFiles := []string{
		"schema.tmpl",
		"controller.tmpl",
		"events.tmpl",
		"repository_postgres.tmpl",
		"repository_sqlite.tmpl",
		"repository.tmpl",
		"command_handler.tmpl",
		"query_handler.tmpl",
		"bootstrap.tmpl",
		"infrastructure.tmpl",
		"schema_hcl.tmpl",
	}

	for _, file := range templateFiles {
		content, err := GetTemplate(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", file, err)
		}

		// Create template with header included
		tmpl := template.New(file).Funcs(TemplateFuncs())

		// Then parse the actual template
		_, err = tmpl.Parse(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", file, err)
		}

		// Store template with its actual name
		templates[file] = tmpl
	}

	return templates, nil
}

// TemplateFuncs returns common template functions used across all generators.
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Case conversion
		"lower":      strings.ToLower,
		"camelCase":  parsers.CamelCase,
		"pascalCase": parsers.PascalCase,
		"snakeCase":  parsers.SnakeCase,
		"kebabCase":  parsers.KebabCase,

		// String utilities
		"pluralize":      parsers.Pluralize,
		"hasPrefix":      strings.HasPrefix,
		"contains":       parsers.Contains, // For slice contains checks
		"stringContains": strings.Contains, // For string contains checks

		// Type checking
		"isPointer": parsers.IsPointer,
		"isSlice":   parsers.IsSlice,

		// Template-specific
		"echoPath": parsers.EchoPath,
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},

		// HCL generation helpers
		"mapToHCLType":           parsers.SchemaToHCLType,
		"mapToSQLiteHCLType":     parsers.SchemaToSQLiteHCLType,
		"formatHCLDefault":       parsers.FormatHCLDefault,
		"formatSQLiteHCLDefault": parsers.FormatSQLiteHCLDefault,
	}
}
