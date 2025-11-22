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

// LoadTemplates loads all templates and returns them as a map
func LoadTemplates() (*template.Template, error) {
	template, err := template.New("root").Funcs(TemplateFuncs()).ParseFS(templatesFS, "tmpl/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return template, nil
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
