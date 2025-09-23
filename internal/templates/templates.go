// Package templates provides template management and code generation for the archesai project.
package templates

import (
	"embed"
	"fmt"
	"text/template"
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
	templateFiles := []string{
		"echo_server.tmpl",
		"events_nats.tmpl",
		"events_redis.tmpl",
		"events.tmpl",
		"repository_postgres.tmpl",
		"repository_sqlite.tmpl",
		"repository.tmpl",
		"service.tmpl",
		"types.tmpl",
	}

	for _, file := range templateFiles {
		content, err := GetTemplate(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", file, err)
		}

		tmpl, err := template.New(file).Funcs(TemplateFuncs()).Parse(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", file, err)
		}

		// Store template with its actual name
		templates[file] = tmpl

	}

	return templates, nil
}
