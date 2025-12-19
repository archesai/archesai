// Package templates provides file writing utilities for code generation.
package templates

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

// Renderer handles template rendering with proper formatting.
type Renderer struct {
	goTemplates  *template.Template
	tsxTemplates *template.Template
}

// NewRenderer creates a new renderer with default settings.
func NewRenderer(templates *Templates) *Renderer {
	return &Renderer{
		goTemplates:  templates.Go,
		tsxTemplates: templates.TSX,
	}
}

// Render renders a template to the provided writer.
// It selects the correct template collection based on the template path.
// Go templates (.go.tmpl) are formatted with gofmt/goimports.
func (r *Renderer) Render(w io.Writer, templateName string, data any) error {
	var tmpl *template.Template
	var isTSX bool

	// Select template collection based on path prefix
	// Templates are named by their basename, so strip the prefix for lookup
	lookupName := templateName
	if strings.HasPrefix(templateName, "tsx/") {
		isTSX = true
		if r.tsxTemplates == nil {
			return fmt.Errorf("TSX templates not configured")
		}
		// Extract basename for lookup (e.g., "tsx/components/datatable.tsx.tmpl" -> "datatable.tsx.tmpl")
		parts := strings.Split(templateName, "/")
		lookupName = parts[len(parts)-1]
		tmpl = r.tsxTemplates.Lookup(lookupName)
	} else {
		if r.goTemplates == nil {
			return fmt.Errorf("go templates not configured")
		}
		tmpl = r.goTemplates.Lookup(lookupName)
	}

	if tmpl == nil {
		return fmt.Errorf("template %s not found (lookup: %s)", templateName, lookupName)
	}

	// For TSX templates or non-Go templates, just execute directly
	if isTSX || !strings.HasSuffix(templateName, ".go.tmpl") {
		return tmpl.Execute(w, data)
	}

	// For Go templates, we need to format
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	content := buf.Bytes()

	// Format Go code with gofmt
	formatted, err := format.Source(content)
	if err != nil {
		return fmt.Errorf("failed to format Go code from template %s: %w", templateName, err)
	}

	// Apply goimports to fix imports (add missing, remove unused)
	imported, err := imports.Process("", formatted, &imports.Options{
		Fragment:  false,
		AllErrors: false,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	})
	if err != nil {
		// If goimports fails, at least use gofmt result
		content = formatted
	} else {
		content = imported
	}

	// Write the formatted content to the writer
	_, err = w.Write(content)
	return err
}

// RenderToString renders a template to a string.
func (r *Renderer) RenderToString(templateName string, data any) (string, error) {
	var buf bytes.Buffer
	if err := r.Render(&buf, templateName, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
