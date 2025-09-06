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
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// TemplateFuncs returns common template functions used across all generators.
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"title":       Title,
		"lower":       strings.ToLower,
		"upper":       strings.ToUpper,
		"camelCase":   CamelCase,
		"pascalCase":  PascalCase,
		"snakeCase":   SnakeCase,
		"kebabCase":   KebabCase,
		"pluralize":   Pluralize,
		"singularize": Singularize,
		"join":        strings.Join,
		"contains":    Contains,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"trimPrefix":  strings.TrimPrefix,
		"trimSuffix":  strings.TrimSuffix,
		"replace":     strings.ReplaceAll,
		"quote":       Quote,
		"indent":      Indent,
		"comment":     Comment,
	}
}

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

// Title capitalizes the first letter of a string.
func Title(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// CamelCase converts a string to camelCase.
func CamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Handle snake_case, kebab-case, and space-separated
	parts := splitWords(s)
	if len(parts) == 0 {
		return s
	}

	// First part is lowercase
	result := strings.ToLower(parts[0])

	// Remaining parts are title case
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			result += Title(strings.ToLower(parts[i]))
		}
	}

	return result
}

// PascalCase converts a string to PascalCase.
func PascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	parts := splitWords(s)
	result := ""

	for _, part := range parts {
		if part != "" {
			result += Title(strings.ToLower(part))
		}
	}

	return result
}

// SnakeCase converts a string to snake_case.
func SnakeCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var result bytes.Buffer
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// Add underscore before uppercase letters (except at the start)
			if i > 0 && s[i-1] != '_' && s[i-1] != '-' && s[i-1] != ' ' {
				result.WriteRune('_')
			}
		}
		if r == '-' || r == ' ' {
			result.WriteRune('_')
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}

// KebabCase converts a string to kebab-case.
func KebabCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var result bytes.Buffer
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// Add hyphen before uppercase letters (except at the start)
			if i > 0 && s[i-1] != '-' && s[i-1] != '_' && s[i-1] != ' ' {
				result.WriteRune('-')
			}
		}
		if r == '_' || r == ' ' {
			result.WriteRune('-')
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}

// Pluralize converts a singular word to plural using simple rules.
func Pluralize(word string) string {
	if word == "" {
		return word
	}

	// Special cases
	special := map[string]string{
		"person": "people",
		"child":  "children",
		"mouse":  "mice",
		"tooth":  "teeth",
		"foot":   "feet",
		"goose":  "geese",
		"man":    "men",
		"woman":  "women",
	}

	lower := strings.ToLower(word)
	if plural, ok := special[lower]; ok {
		// Preserve original case
		if unicode.IsUpper(rune(word[0])) {
			return Title(plural)
		}
		return plural
	}

	// Words ending in 'y' preceded by a consonant
	if strings.HasSuffix(word, "y") && len(word) > 1 {
		beforeY := word[len(word)-2]
		if !isVowel(beforeY) {
			return word[:len(word)-1] + "ies"
		}
		return word + "s"
	}

	// Words ending in 's', 'x', 'z', 'ch', 'sh'
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
		strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") ||
		strings.HasSuffix(word, "sh") {
		return word + "es"
	}

	// Words ending in 'f' or 'fe'
	if strings.HasSuffix(word, "f") {
		return word[:len(word)-1] + "ves"
	}
	if strings.HasSuffix(word, "fe") {
		return word[:len(word)-2] + "ves"
	}

	// Default: add 's'
	return word + "s"
}

// Singularize converts a plural word to singular using simple rules.
func Singularize(word string) string {
	if word == "" {
		return word
	}

	// Special cases
	special := map[string]string{
		"people":   "person",
		"children": "child",
		"mice":     "mouse",
		"teeth":    "tooth",
		"feet":     "foot",
		"geese":    "goose",
		"men":      "man",
		"women":    "woman",
	}

	lower := strings.ToLower(word)
	if singular, ok := special[lower]; ok {
		// Preserve original case
		if unicode.IsUpper(rune(word[0])) {
			return Title(singular)
		}
		return singular
	}

	// Words ending in 'ies'
	if strings.HasSuffix(word, "ies") {
		return word[:len(word)-3] + "y"
	}

	// Words ending in 'ves'
	if strings.HasSuffix(word, "ves") {
		return word[:len(word)-3] + "f"
	}

	// Words ending in 'es'
	if strings.HasSuffix(word, "es") {
		// Check if it's 'ses', 'xes', 'zes', 'ches', 'shes'
		if len(word) > 3 {
			beforeES := word[len(word)-3]
			if beforeES == 's' || beforeES == 'x' || beforeES == 'z' {
				return word[:len(word)-2]
			}
		}
		if strings.HasSuffix(word, "ches") || strings.HasSuffix(word, "shes") {
			return word[:len(word)-2]
		}
		return word[:len(word)-1] // Remove just 's'
	}

	// Words ending in 's'
	if strings.HasSuffix(word, "s") && !strings.HasSuffix(word, "ss") {
		return word[:len(word)-1]
	}

	// Already singular or unknown pattern
	return word
}

// Contains checks if a slice contains a string.
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Quote wraps a string in quotes.
func Quote(s string) string {
	return `"` + s + `"`
}

// Indent indents each line of a string by the specified number of tabs.
func Indent(n int, s string) string {
	if s == "" {
		return s
	}

	indent := strings.Repeat("\t", n)
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// Comment formats a string as a Go comment.
func Comment(s string) string {
	if s == "" {
		return ""
	}

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "// " + line
		} else {
			lines[i] = "//"
		}
	}
	return strings.Join(lines, "\n")
}

// splitWords splits a string into words based on various delimiters.
func splitWords(s string) []string {
	var words []string
	var current bytes.Buffer

	for i, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		} else if i > 0 && unicode.IsUpper(r) && !unicode.IsUpper(rune(s[i-1])) {
			// Start of new word (camelCase boundary)
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			current.WriteRune(r)
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// isVowel checks if a byte represents a vowel.
func isVowel(b byte) bool {
	return b == 'a' || b == 'e' || b == 'i' || b == 'o' || b == 'u' ||
		b == 'A' || b == 'E' || b == 'I' || b == 'O' || b == 'U'
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

// FileWriter handles writing generated files with proper formatting and headers.
type FileWriter struct {
	// Header to prepend to all generated files
	header string

	// Whether to overwrite existing files
	overwrite bool

	// Whether to format Go code
	formatCode bool
}

// NewFileWriter creates a new file writer with default settings.
func NewFileWriter() *FileWriter {
	return &FileWriter{
		header:     DefaultHeader(),
		overwrite:  true,
		formatCode: true,
	}
}

// WithHeader sets a custom header for generated files.
func (fw *FileWriter) WithHeader(header string) *FileWriter {
	fw.header = header
	return fw
}

// WithOverwrite sets whether to overwrite existing files.
func (fw *FileWriter) WithOverwrite(overwrite bool) *FileWriter {
	fw.overwrite = overwrite
	return fw
}

// WithFormatting sets whether to format Go code.
func (fw *FileWriter) WithFormatting(format bool) *FileWriter {
	fw.formatCode = format
	return fw
}

// WriteFile writes content to a file with proper headers and formatting.
func (fw *FileWriter) WriteFile(path string, content []byte) error {
	// Check if file exists and overwrite is disabled
	if !fw.overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file already exists and overwrite is disabled: %s", path)
		}
	}

	// Ensure the file has .gen.go suffix if it's a Go file
	if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".gen.go") {
		// Only enforce for generated files in internal/
		if strings.Contains(path, "internal/") && !strings.Contains(path, "_manual.go") {
			base := strings.TrimSuffix(path, ".go")
			path = base + ".gen.go"
		}
	}

	// Add header if it's a Go file
	if strings.HasSuffix(path, ".go") && fw.header != "" {
		content = fw.addHeader(content)
	}

	// Format Go code if enabled
	if fw.formatCode && strings.HasSuffix(path, ".go") {
		formatted, err := format.Source(content)
		if err != nil {
			// If formatting fails, write unformatted with a warning comment
			warning := []byte("// WARNING: This file could not be formatted automatically\n")
			content = append(warning, content...)
		} else {
			content = formatted
		}
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// WriteTemplate executes a template and writes the result to a file.
func (fw *FileWriter) WriteTemplate(path string, tmpl *template.Template, data interface{}) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return fw.WriteFile(path, buf.Bytes())
}

// WriteTemplateString parses a template string and writes the result to a file.
func (fw *FileWriter) WriteTemplateString(path, tmplStr string, data interface{}) error {
	tmpl, err := template.New(filepath.Base(path)).Funcs(TemplateFuncs()).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return fw.WriteTemplate(path, tmpl, data)
}

// addHeader prepends the header to content.
func (fw *FileWriter) addHeader(content []byte) []byte {
	// Skip if content already has a header
	contentStr := string(content)
	if strings.HasPrefix(contentStr, "// Code generated") ||
		strings.HasPrefix(contentStr, "// Package") {
		// Check if it already has our header
		if strings.Contains(contentStr, "DO NOT EDIT") {
			return content
		}
	}

	// Find the package declaration
	lines := strings.Split(contentStr, "\n")
	var headerLines []string
	packageIndex := -1

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			packageIndex = i
			break
		}
	}

	// Add header before package declaration
	if packageIndex >= 0 {
		headerLines = append(headerLines, fw.header)
		if !strings.HasSuffix(fw.header, "\n") {
			headerLines = append(headerLines, "")
		}
		headerLines = append(headerLines, lines[packageIndex:]...)
		return []byte(strings.Join(headerLines, "\n"))
	}

	// No package declaration found, prepend header
	return []byte(fw.header + "\n" + contentStr)
}

// DefaultHeader returns the default header for generated files.
func DefaultHeader() string {
	return `// Code generated by archesai-codegen. DO NOT EDIT.`
}

// HeaderWithSource returns a header that includes the source file.
func HeaderWithSource(source string) string {
	return fmt.Sprintf(`// Code generated by archesai-codegen. DO NOT EDIT.
// Source: %s`, source)
}

// HeaderWithDetails returns a detailed header for generated files.
func HeaderWithDetails(generator, source string) string {
	return fmt.Sprintf(`// Code generated by archesai-codegen. DO NOT EDIT.
// Generator: %s
// Source: %s`, generator, source)
}

// EnsureGenSuffix ensures a file path has the .gen.go suffix.
func EnsureGenSuffix(path string) string {
	if !strings.HasSuffix(path, ".go") {
		return path
	}

	if strings.HasSuffix(path, ".gen.go") {
		return path
	}

	// Don't change manual files
	if strings.Contains(path, "_manual") || strings.Contains(path, "manual_") {
		return path
	}

	base := strings.TrimSuffix(path, ".go")
	return base + ".gen.go"
}

// IsGeneratedFile checks if a file is a generated file.
func IsGeneratedFile(path string) bool {
	// Check by suffix
	if strings.HasSuffix(path, ".gen.go") {
		return true
	}

	// Check by reading the file header
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	contentStr := string(content)
	return strings.Contains(contentStr, "Code generated") &&
		strings.Contains(contentStr, "DO NOT EDIT")
}

// BackupFile creates a backup of an existing file before overwriting.
func BackupFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // No file to backup
	}

	// Create backup path
	backupPath := path + ".backup"

	// Read original file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}

// CleanGeneratedFiles removes all generated files in a directory.
func CleanGeneratedFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Remove generated files
		if IsGeneratedFile(path) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", path, err)
			}
		}

		return nil
	})
}

// OutputPath generates the output path for a generated file.
func OutputPath(baseDir, domain, fileName string) string {
	return filepath.Join(baseDir, domain, EnsureGenSuffix(fileName))
}

// DomainPath generates the base path for a domain.
func DomainPath(baseDir, domain string) string {
	return filepath.Join(baseDir, domain)
}

// AdapterPath generates the path for adapter implementations.
func AdapterPath(baseDir, domain, adapter, fileName string) string {
	return filepath.Join(baseDir, domain, "adapters", adapter, EnsureGenSuffix(fileName))
}
