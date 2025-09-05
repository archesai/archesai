// Package domain generates domain scaffolding for new domains.
package domain

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// Generator handles generation of domain scaffolding.
type Generator struct{}

// NewGenerator creates a new domain generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Config represents domain generation configuration.
type Config struct {
	Name        string // e.g., "billing"
	Package     string // e.g., "billing"
	Description string // e.g., "Billing and subscription management"
	Tables      string // e.g., "subscription,invoice,payment" (comma-separated)
	HasAuth     bool   // Whether domain needs auth middleware
	HasEvents   bool   // Whether domain uses events
}

// Generate generates domain scaffolding based on configuration.
func (g *Generator) Generate(config Config) error {
	// Normalize config
	config.Package = strings.ToLower(config.Name)
	if config.Description == "" {
		config.Description = fmt.Sprintf("%s domain functionality", config.Name)
	}

	// Parse tables from comma-separated string
	tables := ParseTables(config.Tables)

	domainPath := filepath.Join("internal", "domains", config.Package)
	adaptersPath := filepath.Join(domainPath, "adapters")

	// Create directories
	if err := os.MkdirAll(adaptersPath, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files from templates
	files := map[string]string{
		"domain.go.tmpl":     filepath.Join(domainPath, fmt.Sprintf("%s.go", config.Package)),
		"entities.go.tmpl":   filepath.Join(domainPath, "entities.go"),
		"service.go.tmpl":    filepath.Join(domainPath, "service.go"),
		"repository.go.tmpl": filepath.Join(domainPath, "repository.go"),
		"handler.go.tmpl":    filepath.Join(domainPath, "handler.go"),
	}

	if config.HasAuth {
		files["middleware.go.tmpl"] = filepath.Join(domainPath, "middleware.go")
	}

	if config.HasEvents {
		files["events.go.tmpl"] = filepath.Join(domainPath, "events.go")
	}

	// Create function map for templates
	funcMap := template.FuncMap{
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}

	for tmplFile, outputPath := range files {
		tmplPath := filepath.Join("templates", tmplFile)
		tmplContent, err := templatesFS.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", tmplFile, err)
		}

		tmpl, err := template.New(tmplFile).Funcs(funcMap).Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tmplFile, err)
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}(file)

		// Create template data with parsed tables
		data := struct {
			*Config
			Tables []string
		}{
			Config: &config,
			Tables: tables,
		}

		if err := tmpl.Execute(file, data); err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", outputPath, err)
		}
	}

	// Create empty .gitkeep in adapters directory
	gitkeepPath := filepath.Join(adaptersPath, ".gitkeep")
	if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create .gitkeep: %w", err)
	}

	return nil
}

// ParseTables parses a comma-separated string of tables into a slice.
func ParseTables(tablesStr string) []string {
	if tablesStr == "" {
		return []string{}
	}

	tables := strings.Split(tablesStr, ",")
	result := make([]string, 0, len(tables))
	for _, t := range tables {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
