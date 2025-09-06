// Package codegen provides template functions for code generation.
package codegen

import (
	"bytes"
	"strings"
	"text/template"
	"unicode"
)

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
