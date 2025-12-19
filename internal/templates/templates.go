// Package templates provides template management and code generation for the archesai project.
package templates

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/archesai/archesai/internal/strutil"
)

// TypeScript type constants for code generation.
const (
	tsTypeString  = "string"
	tsTypeNumber  = "number"
	tsTypeBoolean = "boolean"
	tsTypeUnknown = "unknown"
)

//go:embed tmpl/go/*.tmpl
var goTemplatesFS embed.FS

//go:embed tmpl/tsx/components/*.tmpl tmpl/tsx/routes/*.tmpl tmpl/tsx/config/*.tmpl tmpl/tsx/config/*.css.tmpl tmpl/tsx/lib/*.tmpl
var tsxTemplatesFS embed.FS

// Templates holds both Go and TSX template collections.
type Templates struct {
	Go  *template.Template
	TSX *template.Template
}

// LoadTemplates loads all templates and returns them.
func LoadTemplates() (*Templates, error) {
	goTmpl, err := template.New("go").
		Funcs(TemplateFuncs()).
		ParseFS(goTemplatesFS, "tmpl/go/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go templates: %w", err)
	}

	tsxTmpl, err := template.New("tsx").Delims("[[", "]]").Funcs(TSXTemplateFuncs()).ParseFS(
		tsxTemplatesFS,
		"tmpl/tsx/components/*.tmpl",
		"tmpl/tsx/routes/*.tmpl",
		"tmpl/tsx/config/*.tmpl",
		"tmpl/tsx/config/*.css.tmpl",
		"tmpl/tsx/lib/*.tmpl",
	)
	if err != nil {
		// TSX templates are optional for now - log but don't fail
		tsxTmpl = template.New("tsx").Delims("[[", "]]").Funcs(TSXTemplateFuncs())
	}

	return &Templates{
		Go:  goTmpl,
		TSX: tsxTmpl,
	}, nil
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
		"hasSuffix":      strings.HasSuffix,
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
		"intDeref": func(i *int) int {
			if i == nil {
				return 0
			}
			return *i
		},

		// Type qualification helpers
		"modelType": modelType,

		// String escaping for raw literals
		"rawString": escapeRawString,

		// Arithmetic helpers
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
}

// escapeRawString escapes a string for use in a Go raw string literal (backticks).
// It replaces backticks with a concatenation pattern: ` + "`" + `
func escapeRawString(s string) string {
	return strings.ReplaceAll(s, "`", "` + \"`\" + `")
}

// TSXTemplateFuncs returns template functions for TSX templates.
// It includes all Go template functions plus TSX-specific ones.
func TSXTemplateFuncs() template.FuncMap {
	funcs := TemplateFuncs()

	// Add TSX-specific functions
	funcs["tsType"] = goTypeToTS
	funcs["entityKey"] = func(name string) string {
		return strutil.SnakeCase(name) + "s"
	}
	funcs["upper"] = strings.ToUpper

	return funcs
}

// goTypeToTS converts Go types to TypeScript types.
func goTypeToTS(goType string) string {
	// Handle pointer types
	if strings.HasPrefix(goType, "*") {
		return goTypeToTS(strings.TrimPrefix(goType, "*")) + " | null"
	}

	// Handle slice types
	if strings.HasPrefix(goType, "[]") {
		return goTypeToTS(strings.TrimPrefix(goType, "[]")) + "[]"
	}

	// Map Go types to TypeScript
	switch goType {
	case tsTypeString:
		return tsTypeString
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return tsTypeNumber
	case "float32", "float64":
		return tsTypeNumber
	case "bool":
		return tsTypeBoolean
	case "uuid.UUID":
		return tsTypeString
	case "time.Time":
		return tsTypeString // ISO date string
	case "json.RawMessage":
		return tsTypeUnknown
	default:
		// For schema references, keep as-is (will be generated types)
		return goType
	}
}

// modelType prefixes a type with "schemas." if it's a schema reference (not a primitive or already qualified).
func modelType(goType string) string {
	// Already qualified with package
	if strings.Contains(goType, ".") {
		return goType
	}

	// Handle slice types
	if after, ok := strings.CutPrefix(goType, "[]"); ok {
		inner := after
		return "[]" + modelType(inner)
	}

	// Handle pointer types
	if after, ok := strings.CutPrefix(goType, "*"); ok {
		inner := after
		return "*" + modelType(inner)
	}

	// Primitive types don't need schemas prefix
	switch goType {
	case "string", "int", "int32", "int64", "float32", "float64", "bool", "any", "time.Time":
		return goType
	}

	// uuid.UUID and other qualified types
	if goType == "uuid.UUID" {
		return goType
	}

	// Schema reference - prefix with schemas.
	return "schemas." + goType
}
