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
		"title":                        Title,
		"lower":                        strings.ToLower,
		"upper":                        strings.ToUpper,
		"camelCase":                    CamelCase,
		"pascalCase":                   PascalCase,
		"snakeCase":                    SnakeCase,
		"kebabCase":                    KebabCase,
		"pluralize":                    Pluralize,
		"singularize":                  Singularize,
		"join":                         strings.Join,
		"contains":                     Contains,
		"hasPrefix":                    strings.HasPrefix,
		"hasSuffix":                    strings.HasSuffix,
		"trimPrefix":                   strings.TrimPrefix,
		"trimSuffix":                   strings.TrimSuffix,
		"replace":                      strings.ReplaceAll,
		"quote":                        Quote,
		"indent":                       Indent,
		"comment":                      Comment,
		"paramType":                    ParamType,
		"isUUIDParam":                  IsUUIDParam,
		"generateTypeConversion":       GenerateTypeConversion,
		"generateCreateTypeConversion": GenerateCreateTypeConversion,
		"generateUpdateTypeConversion": GenerateUpdateTypeConversion,
		"isUpdateExcluded":             IsUpdateExcluded,
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

// ParamType returns the Go type for a parameter name based on naming conventions.
func ParamType(paramName string) string {
	switch paramName {
	case "id", "ID", "userID", "organizationID", "pipelineID", "runId", "runID", "toolID", "invitationId", "invitationID", "artifactID":
		return goTypeUUIDType
	case "name", "email", "token", "provider", "providerAccountId", "slug", "stripeCustomerId", "inviterId":
		return "string" //nolint:goconst // Go type
	default:
		return "string" //nolint:goconst // Go type
	}
}

// IsUUIDParam checks if a parameter should be treated as a UUID type.
func IsUUIDParam(paramName string) bool {
	return paramName == "id" || paramName == "userID"
}

// ToSnakeCase converts camelCase to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// GenerateTypeConversion generates the appropriate type conversion between Go types and SQLC types
//
//nolint:gocyclo // This function needs to handle many type conversion cases
func GenerateTypeConversion(goType, sqlcType, varPrefix, fieldName string) string {
	if goType == "" || sqlcType == "" {
		// Fallback to original logic if type inference failed
		return varPrefix + "." + fieldName
	}

	// Handle special cases based on type combinations
	switch {
	// Email field conversions (including openapi_types.Email)
	case (goType == goTypeEmail || goType == goTypeEmailFull) && sqlcType == goTypePtrString:
		return "types.Email(stringFromPtr(" + varPrefix + "." + fieldName + "))"
	case (goType == goTypeEmail || goType == goTypeEmailFull) && sqlcType == goTypeString:
		return "types.Email(" + varPrefix + "." + fieldName + ")"

	// Map/JSON conversions
	case goType == goTypeMapString && sqlcType == goTypePtrString:
		return "unmarshalJSON(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeString && sqlcType == goTypeString:
		return varPrefix + "." + fieldName

	// String pointer conversions
	case goType == goTypeString && sqlcType == goTypePtrString:
		return "stringFromPtr(" + varPrefix + "." + fieldName + ")"
	case goType == "*string" && sqlcType == goTypeString:
		return "stringPtr(" + varPrefix + "." + fieldName + ")"

	// UUID conversions
	case goType == goTypeUUIDLiteral && sqlcType == goTypePtrUUID:
		return "uuidFromPtr(" + varPrefix + "." + fieldName + ")"
	case goType == "*UUID" && sqlcType == goTypeUUIDType:
		return "&" + varPrefix + "." + fieldName
	case goType == goTypeUUIDLiteral && sqlcType == goTypeUUIDType:
		return varPrefix + "." + fieldName

	// Time conversions
	case goType == goTypeTimeTime && sqlcType == goTypePtrTime:
		return "timeFromPtr(" + varPrefix + "." + fieldName + ")"
	case goType == "*time.Time" && sqlcType == goTypeTimeTime:
		return "&" + varPrefix + "." + fieldName
	case goType == goTypeTimeTime && sqlcType == goTypeTimeTime:
		return varPrefix + "." + fieldName

	// Boolean conversions
	case goType == goTypeBool && sqlcType == goTypeBool:
		return varPrefix + "." + fieldName

	// Number conversions
	case goType == goTypeFloat32 && sqlcType == goTypeInt32:
		return "float32(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeFloat32 && sqlcType == goTypeFloat64:
		return "float32(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeFloat64 && sqlcType == goTypeFloat32:
		return "float64(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeFloat32 && sqlcType == goTypeFloat32:
		return varPrefix + "." + fieldName
	case goType == goTypeFloat64 && sqlcType == goTypeFloat64:
		return varPrefix + "." + fieldName
	case goType == goTypeInt32 && sqlcType == goTypeInt32:
		return varPrefix + "." + fieldName

	// Enum conversions (when Go type is not a basic type but SQLC type is string)
	case sqlcType == goTypeString && goType != goTypeString && goType != "*string":
		return goType + "(" + varPrefix + "." + fieldName + ")"

	// Default: direct assignment
	default:
		return varPrefix + "." + fieldName
	}
}

// GenerateCreateTypeConversion generates type conversions for Create operations
//
//nolint:gocyclo // This function needs to handle many type conversion cases
func GenerateCreateTypeConversion(goType, sqlcType, varPrefix, fieldName string) string {
	if goType == "" || sqlcType == "" {
		// Fallback if type inference failed
		return varPrefix + "." + fieldName
	}

	// Handle special cases based on type combinations
	switch {
	// Email field conversions (including openapi_types.Email)
	case (goType == goTypeEmail || goType == goTypeEmailFull) && sqlcType == goTypePtrString:
		return "stringPtr(string(" + varPrefix + "." + fieldName + "))"
	case (goType == goTypeEmail || goType == goTypeEmailFull) && sqlcType == goTypeString:
		return "string(" + varPrefix + "." + fieldName + ")"

	// Enum conversions (custom type to string)
	case strings.Contains(goType, "ID") && !strings.Contains(goType, "UUID") && sqlcType == goTypeString:
		// This handles types like AccountProviderID
		return "string(" + varPrefix + "." + fieldName + ")"

	// Time to pointer conversions
	case goType == goTypeTimeTime && sqlcType == goTypePtrTime:
		return "&" + varPrefix + "." + fieldName

	// Map/object to JSON string conversions
	case goType == goTypeMapString && sqlcType == goTypePtrString:
		return "marshalJSON(" + varPrefix + "." + fieldName + ")"

	// String to pointer conversions (for optional fields)
	case goType == goTypeString && sqlcType == goTypePtrString:
		return "stringPtr(" + varPrefix + "." + fieldName + ")"

	// UUID handling
	case goType == goTypeUUIDLiteral && sqlcType == goTypeUUIDType:
		return varPrefix + "." + fieldName
	case goType == goTypeUUIDLiteral && sqlcType == goTypePtrUUID:
		return "&" + varPrefix + "." + fieldName
	case goType == "*UUID" && sqlcType == goTypePtrUUID:
		return varPrefix + "." + fieldName

	// Direct assignments
	case goType == goTypeString && sqlcType == goTypeString:
		return varPrefix + "." + fieldName
	case goType == goTypeBool && sqlcType == goTypeBool:
		return varPrefix + "." + fieldName

	// Number conversions
	case goType == goTypeFloat32 && sqlcType == goTypeInt32:
		return "int32(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeFloat32 && sqlcType == goTypeFloat64:
		return "float64(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeInt32 && sqlcType == goTypeInt32:
		return varPrefix + "." + fieldName

	// Enum conversions (when SQL expects string but Go has custom type)
	case sqlcType == goTypeString && !strings.HasPrefix(goType, goTypeString) && !strings.HasPrefix(goType, goTypePtrString):
		return "string(" + varPrefix + "." + fieldName + ")"

	// Default: direct assignment
	default:
		return varPrefix + "." + fieldName
	}
}

// IsUpdateExcluded checks if a field should be excluded from update operations
func IsUpdateExcluded(fieldName string, excludeList []string) bool {
	for _, excluded := range excludeList {
		if strings.EqualFold(excluded, fieldName) {
			return true
		}
	}
	return false
}

// GenerateUpdateTypeConversion generates type conversions for Update operations
// Update operations require pointer types for all fields to indicate which fields to update
func GenerateUpdateTypeConversion(goType, sqlcType, varPrefix, fieldName string) string {
	// For UPDATE operations, everything needs to be a pointer since all fields are optional

	// Handle special type conversions first
	switch {
	// Email types need to be converted to string first, then to pointer
	case strings.Contains(goType, "Email"):
		return "stringPtr(string(" + varPrefix + "." + fieldName + "))"

	// Custom types (enums, type aliases) that are string-based
	case goType != "" && goType != goTypeString && goType != goTypeBool && goType != goTypeInt &&
		goType != goTypeInt32 && goType != goTypeInt64 && goType != goTypeFloat32 && goType != goTypeFloat64 &&
		goType != goTypeTimeTime && !strings.Contains(goType, "UUID") && !strings.Contains(goType, "*"):
		// It's a custom type like MemberRole, convert to string then pointer
		return "stringPtr(string(" + varPrefix + "." + fieldName + "))"

	// UUID types
	case strings.Contains(goType, "UUID"):
		return "&" + varPrefix + "." + fieldName

	// Map/JSON types
	case strings.Contains(goType, "map[string]"):
		return "marshalJSON(" + varPrefix + "." + fieldName + ")"

	// Basic types - convert to pointer
	case goType == goTypeString:
		return "stringPtr(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeBool:
		return "boolPtr(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeInt32:
		return "int32Ptr(" + varPrefix + "." + fieldName + ")"
	case goType == "int64":
		return "int64Ptr(int64(" + varPrefix + "." + fieldName + "))"
	case goType == goTypeFloat32:
		// Check what SQLC expects for float32 fields
		if strings.Contains(sqlcType, "int32") {
			// Credits and similar fields are float32 but stored as int32
			return "int32Ptr(int32(" + varPrefix + "." + fieldName + "))"
		}
		// For Update operations, SQLC often expects *float64 even when the field is float32
		return "float64Ptr(float64(" + varPrefix + "." + fieldName + "))"
	case goType == goTypeFloat64:
		return "float64Ptr(" + varPrefix + "." + fieldName + ")"
	case goType == goTypeTimeTime:
		return "&" + varPrefix + "." + fieldName

	// Already pointer types
	case strings.HasPrefix(goType, "*"):
		return varPrefix + "." + fieldName

	// Default fallback - try to determine based on sqlcType or just take address
	default:
		if strings.Contains(sqlcType, goTypeString) {
			return "stringPtr(" + varPrefix + "." + fieldName + ")"
		} else if strings.Contains(sqlcType, goTypeBool) {
			return "boolPtr(" + varPrefix + "." + fieldName + ")"
		} else if strings.Contains(sqlcType, goTypeInt) {
			return "int32Ptr(" + varPrefix + "." + fieldName + ")"
		} else if strings.Contains(sqlcType, goTypeFloat32) || strings.Contains(sqlcType, goTypeFloat64) {
			return "float64Ptr(" + varPrefix + "." + fieldName + ")"
		}
		// Last resort - take address
		return "&" + varPrefix + "." + fieldName
	}
}
