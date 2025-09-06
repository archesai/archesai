// Package codegen provides shared utilities for all code generators.
package codegen

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// GetTemplate loads a template by name from the central templates directory.
func GetTemplate(name string) (string, error) {
	content, err := templatesFS.ReadFile("templates/" + name)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", name, err)
	}
	return string(content), nil
}

// GetTemplateFS returns the embedded template filesystem for direct access.
func GetTemplateFS() embed.FS {
	return templatesFS
}

// ParseTemplate loads and parses a template with common functions.
func ParseTemplate(name string) (*template.Template, error) {
	content, err := GetTemplate(name)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Funcs(TemplateFuncs()).Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return tmpl, nil
}

// GetSortedPropertyKeys returns sorted property names from a map of properties.
// This ensures consistent ordering when iterating over schema properties.
func GetSortedPropertyKeys(properties map[string]Property) []string {
	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AddYamlTags adds YAML tags to Go struct fields that have JSON tags but no YAML tags.
// This functionality was moved from tools/codegen/add_yaml_tags.go to consolidate codegen utilities.
func AddYamlTags(filename string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	// Walk through all struct fields and add yaml tags
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Field:
			if x.Tag != nil {
				// Parse existing tag
				tag := reflect.StructTag(strings.Trim(x.Tag.Value, "`"))
				jsonTag := tag.Get("json")

				if jsonTag != "" && tag.Get("yaml") == "" {
					// Add yaml tag with same value as json tag
					newTag := fmt.Sprintf("`json:\"%s\" yaml:\"%s\"`", jsonTag, jsonTag)
					x.Tag.Value = newTag
				}
			}
		}
		return true
	})

	// Format and write the modified file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("error formatting code: %w", err)
	}

	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
