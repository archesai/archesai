// Package codegen provides file writing utilities for code generation.
package codegen

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
	templates *template.Template
}

// NewRenderer creates a new renderer with default settings.
func NewRenderer(templates *template.Template) *Renderer {
	return &Renderer{
		templates: templates,
	}
}

// Render renders a template to the provided writer.
// If the template ends with .go.tmpl, it will format the output.
func (r *Renderer) Render(w io.Writer, templateName string, data any) error {
	if r.templates == nil {
		return fmt.Errorf("templates not configured")
	}

	tmpl := r.templates.Lookup(templateName)
	if tmpl == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	// If writing directly to a non-Go template, just execute it
	if !strings.HasSuffix(templateName, ".go.tmpl") {
		return tmpl.Execute(w, data)
	}

	// For Go templates, we need to format
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	content := buf.Bytes()

	// Format Go code
	// First apply gofmt
	formatted, err := format.Source(content)
	if err != nil {
		return fmt.Errorf("failed to format Go code from template %s: %w", templateName, err)
	}

	// Then apply goimports to fix imports (add missing, remove unused)
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
