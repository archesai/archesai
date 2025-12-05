// Package templates provides template management and code generation for the archesai project.
package templates

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/archesai/archesai/internal/strutil"
	"github.com/archesai/archesai/internal/typeconv"
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
		"camelCase":  strutil.CamelCase,
		"pascalCase": strutil.PascalCase,
		"snakeCase":  strutil.SnakeCase,
		"kebabCase":  strutil.KebabCase,

		// String utilities
		"pluralize":      strutil.Pluralize,
		"hasPrefix":      strings.HasPrefix,
		"contains":       strutil.Contains, // For slice contains checks
		"stringContains": strings.Contains, // For string contains checks

		// Type checking
		"isPointer": strutil.IsPointer,
		"isSlice":   strutil.IsSlice,

		// Template-specific
		"echoPath": strutil.EchoPath,
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},

		// HCL generation helpers
		"mapToHCLType":           typeconv.SchemaToHCLType,
		"mapToSQLiteHCLType":     typeconv.SchemaToSQLiteHCLType,
		"formatHCLDefault":       typeconv.FormatHCLDefault,
		"formatSQLiteHCLDefault": typeconv.FormatSQLiteHCLDefault,
	}
}
