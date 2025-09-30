package codegen

import (
	"bytes"
	"slices"
	"strings"
	"text/template"
	"unicode"
)

// Constants for commonly used type strings
const (
	typeString = "string"
	typeUUID   = "uuid.UUID"
	typeNil    = "nil"
)

// TemplateFuncs returns common template functions used across all generators.
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Case conversion
		"title":          Title,
		"camelCase":      CamelCase,
		"pascalCase":     PascalCase,
		"snakeCase":      SnakeCase,
		"kebabCase":      KebabCase,
		"ToConstantCase": ToConstantCase,

		// String utilities
		"lower":       strings.ToLower,
		"upper":       strings.ToUpper,
		"pluralize":   Pluralize,
		"singularize": Singularize,
		"join":        strings.Join,
		"contains":    Contains,
		"strContains": strings.Contains,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"trimPrefix":  strings.TrimPrefix,
		"trimSuffix":  strings.TrimSuffix,
		"replace":     strings.ReplaceAll,
		"quote":       Quote,
		"indent":      Indent,
		"comment":     Comment,

		// Template utilities
		"dict": Dict,

		// Type checking and validation
		"isPointer":       IsPointer,
		"isSlice":         IsSlice,
		"isMap":           IsMap,
		"isStruct":        IsStruct,
		"isPrimitive":     IsPrimitive,
		"isOptional":      IsOptional,
		"needsValidation": NeedsValidation,
		"defaultValue":    DefaultValueForType,

		// Template-specific functions
		"paramType":         ParamType,
		"isUUIDParam":       IsUUIDParam,
		"isUpdateExcluded":  IsUpdateExcluded,
		"default":           DefaultValue,
		"echoPath":          EchoPath,
		"fieldName":         FieldName,
		"exportedFieldName": ExportedFieldName,
	}
}

// Title capitalizes the first letter of a string and fixes common initialisms.
func Title(s string) string {
	if s == "" {
		return s
	}

	// First capitalize the first letter
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	result := string(r)

	// Fix common initialisms
	initialisms := map[string]string{
		"Id":   "ID",
		"Uuid": "UUID",
		"Api":  "API",
		"Url":  "URL",
		"Http": "HTTP",
		"Json": "JSON",
		"Sql":  "SQL",
		"Xml":  "XML",
		"Html": "HTML",
		"Css":  "CSS",
		"Jwt":  "JWT",
	}

	// Check for initialisms at the beginning
	for prefix, replacement := range initialisms {
		if strings.HasPrefix(result, prefix) && len(result) > len(prefix) {
			// Make sure the next character is uppercase (e.g., IdToken -> IDToken)
			nextChar := rune(result[len(prefix)])
			if unicode.IsUpper(nextChar) {
				result = replacement + result[len(prefix):]
				break
			}
		}
	}

	// Check for initialisms at the end of the string
	for suffix, replacement := range initialisms {
		if strings.HasSuffix(result, suffix) {
			result = result[:len(result)-len(suffix)] + replacement
			break
		}
	}

	return result
}

// CamelCase converts a string to camelCase, preserving common acronyms.
func CamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Special case for single acronyms - just lowercase them
	if acronym, ok := commonAcronyms[strings.ToLower(s)]; ok {
		if s == acronym {
			// It's an acronym in uppercase (e.g., "ID", "API"), return lowercase
			return strings.ToLower(s)
		}
	}

	// Check if string already has the right format with acronyms
	// e.g., "createAPIKey" should become "createAPIKey" (already correct)
	// e.g., "CreateAPIKey" should become "createAPIKey" (lowercase first)
	if hasCorrectCasing(s) {
		// Just lowercase the first character if it's uppercase
		if len(s) > 0 && unicode.IsUpper(rune(s[0])) {
			// Check if it starts with an acronym
			for _, acronym := range commonAcronyms {
				if strings.HasPrefix(s, acronym) {
					// Lowercase the entire acronym at the beginning
					return strings.ToLower(acronym) + s[len(acronym):]
				}
			}
			// Otherwise just lowercase the first character
			return strings.ToLower(s[:1]) + s[1:]
		}
		return s
	}

	parts := splitWords(s)
	if len(parts) == 0 {
		return s
	}

	// First part is lowercase (but preserve acronym if it is one)
	firstPart := strings.ToLower(parts[0])
	var result string
	if acronym, ok := commonAcronyms[firstPart]; ok && len(parts) == 1 {
		// If it's just a single acronym, keep it lowercase
		result = strings.ToLower(acronym)
	} else {
		result = firstPart
	}

	// Remaining parts are title case with acronym handling
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			lowerPart := strings.ToLower(parts[i])
			if acronym, ok := commonAcronyms[lowerPart]; ok {
				result += acronym
			} else {
				result += Title(lowerPart)
			}
		}
	}

	// Handle Go reserved keywords
	return sanitizeGoKeyword(result)
}

// sanitizeGoKeyword handles Go reserved keywords by appending an underscore
func sanitizeGoKeyword(name string) string {
	reservedKeywords := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,
	}

	if reservedKeywords[name] {
		return name + "_"
	}
	return name
}

// Common acronyms that should be preserved in uppercase
var commonAcronyms = map[string]string{
	"api":   "API",
	"url":   "URL",
	"uri":   "URI",
	"uuid":  "UUID",
	"id":    "ID",
	"http":  "HTTP",
	"https": "HTTPS",
	"sql":   "SQL",
	"json":  "JSON",
	"xml":   "XML",
	"csv":   "CSV",
	"jwt":   "JWT",
	"oauth": "OAuth",
	"saml":  "SAML",
	"ldap":  "LDAP",
	"dns":   "DNS",
	"tcp":   "TCP",
	"udp":   "UDP",
	"ip":    "IP",
	"vm":    "VM",
	"os":    "OS",
	"cpu":   "CPU",
	"gpu":   "GPU",
	"ram":   "RAM",
	"ssd":   "SSD",
	"hdd":   "HDD",
	"cdn":   "CDN",
	"vpn":   "VPN",
	"ssh":   "SSH",
	"ftp":   "FTP",
	"sftp":  "SFTP",
	"smtp":  "SMTP",
	"imap":  "IMAP",
	"pop":   "POP",
	"aws":   "AWS",
	"gcp":   "GCP",
	"sdk":   "SDK",
	"ci":    "CI",
	"cd":    "CD",
}

// PascalCase converts a string to PascalCase, preserving common acronyms.
func PascalCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// If the string is already in PascalCase with correct acronyms (e.g., "CreateAPIKey"), preserve it
	// But if it starts with lowercase (e.g., "createAPIKey"), we need to capitalize
	if len(s) > 0 && unicode.IsUpper(rune(s[0])) && hasCorrectCasing(s) {
		return s
	}

	// If it starts with lowercase but has correct acronyms (e.g., "createAPIKey"), just capitalize first letter
	if hasCorrectCasing(s) && len(s) > 0 && unicode.IsLower(rune(s[0])) {
		return strings.ToUpper(s[:1]) + s[1:]
	}

	parts := splitWords(s)
	result := ""

	for _, part := range parts {
		if part != "" {
			lowerPart := strings.ToLower(part)
			if acronym, ok := commonAcronyms[lowerPart]; ok {
				result += acronym
			} else {
				result += Title(lowerPart)
			}
		}
	}

	return result
}

// hasCorrectCasing checks if a string already has correct PascalCase with acronyms
func hasCorrectCasing(s string) bool {
	// Check if string contains known acronyms in uppercase
	for _, acronym := range commonAcronyms {
		if strings.Contains(s, acronym) {
			return true
		}
	}
	return false
}

// SnakeCase converts a string to snake_case.
func SnakeCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var result bytes.Buffer
	runes := []rune(s)

	for i, r := range runes {
		if unicode.IsUpper(r) && i > 0 {
			// Look at the previous character
			prevIsLower := unicode.IsLower(runes[i-1])
			prevIsDigit := unicode.IsDigit(runes[i-1])

			// Look at the next character if it exists
			nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])

			// Add underscore before uppercase letter if:
			// 1. Previous char is lowercase or digit
			// 2. OR this is the start of a new word (current is upper, next is lower)
			if prevIsLower || prevIsDigit || nextIsLower {
				// But don't add if previous char is already a separator
				if runes[i-1] != '_' && runes[i-1] != '-' && runes[i-1] != ' ' {
					result.WriteRune('_')
				}
			}
		}

		// Convert separators to underscores
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

// ToConstantCase converts a string to proper constant case with initialisms.
func ToConstantCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// Apply initialisms to the title-cased string
	result := Title(s)

	// Common initialisms that should be uppercase
	initialisms := map[string]string{
		"Id":   "ID",
		"Api":  "API",
		"Url":  "URL",
		"Http": "HTTP",
		"Uuid": "UUID",
		"Json": "JSON",
		"Sql":  "SQL",
		"Xml":  "XML",
		"Html": "HTML",
		"Css":  "CSS",
		"Js":   "JS",
		"Jwt":  "JWT",
	}

	// Replace any word that is an initialism
	for word, replacement := range initialisms {
		// Check at the end of string
		if strings.HasSuffix(result, word) {
			result = result[:len(result)-len(word)] + replacement
		}
		// Check at the beginning
		if strings.HasPrefix(result, word) {
			result = replacement + result[len(word):]
		}
	}

	return result
}

// Capitalize returns a string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
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
		"health": "health", // Health is uncountable, stays the same
		"config": "config", // Config is uncountable, stays the same
		"apikey": "APIKeys",
		"APIKey": "APIKeys",
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
		"apikeys":  "APIKey",
		"Apikeys":  "APIKey",
		"APIKeys":  "APIKey",
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
	return slices.Contains(slice, item)
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

// ParamType returns the Go type for a parameter name.
func ParamType(paramName string) string {
	// UUID types
	if strings.HasSuffix(paramName, "ID") || paramName == "id" {
		return typeUUID
	}
	return typeString
}

// IsUUIDParam checks if a parameter should be treated as a UUID type.
func IsUUIDParam(paramName string) bool {
	return strings.HasSuffix(paramName, "ID") || paramName == "id"
}

// IsUpdateExcluded checks if a field should be excluded from update operations.
func IsUpdateExcluded(fieldName string, excludeList []string) bool {
	for _, excluded := range excludeList {
		if strings.EqualFold(excluded, fieldName) {
			return true
		}
	}
	return false
}

// DefaultValue returns the default value if the first value is empty, otherwise returns the first value.
func DefaultValue(defaultVal string, value string) string {
	if value == "" {
		return defaultVal
	}
	return value
}

// EchoPath converts OpenAPI path format to Echo router format.
// Converts "/users/{id}" to "/users/:id"
func EchoPath(path string) string {
	// Replace {param} with :param for Echo router format
	result := strings.ReplaceAll(path, "{", ":")
	result = strings.ReplaceAll(result, "}", "")
	return result
}

// Dict creates a map from key-value pairs for use in templates.
// Usage: dict "key1" value1 "key2" value2 ...
func Dict(values ...interface{}) map[string]interface{} {
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		if i+1 < len(values) {
			key, ok := values[i].(string)
			if ok {
				dict[key] = values[i+1]
			}
		}
	}
	return dict
}

// IsPointer checks if a Go type string represents a pointer type
func IsPointer(goType string) bool {
	return strings.HasPrefix(goType, "*")
}

// IsSlice checks if a Go type string represents a slice type
func IsSlice(goType string) bool {
	return strings.HasPrefix(goType, "[]")
}

// IsMap checks if a Go type string represents a map type
func IsMap(goType string) bool {
	return strings.HasPrefix(goType, "map[")
}

// IsStruct checks if a Go type string represents a struct type
func IsStruct(goType string) bool {
	// Check if it's not a primitive type and not a pointer/slice/map
	if IsPointer(goType) || IsSlice(goType) || IsMap(goType) {
		return false
	}
	return !IsPrimitive(goType)
}

// IsPrimitive checks if a Go type string represents a primitive type
func IsPrimitive(goType string) bool {
	primitives := map[string]bool{
		"bool":       true,
		"string":     true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"float32":    true,
		"float64":    true,
		"complex64":  true,
		"complex128": true,
		"byte":       true,
		"rune":       true,
		"time.Time":  true,
		"uuid.UUID":  true,
	}

	// Remove pointer prefix if present
	cleanType := strings.TrimPrefix(goType, "*")
	return primitives[cleanType]
}

// IsOptional checks if a field is optional (pointer type)
func IsOptional(goType string) bool {
	return IsPointer(goType)
}

// NeedsValidation checks if a type needs validation
func NeedsValidation(goType string, required bool) bool {
	// Required strings always need validation
	if required && goType == typeString {
		return true
	}
	// UUIDs need validation
	if goType == typeUUID {
		return true
	}
	// Email fields need validation
	if strings.Contains(strings.ToLower(goType), "email") {
		return true
	}
	// URL fields need validation
	if strings.Contains(strings.ToLower(goType), "url") {
		return true
	}
	return false
}

// DefaultValueForType returns the appropriate default value for a Go type
func DefaultValueForType(goType string) string {
	switch {
	case IsPointer(goType):
		return typeNil
	case goType == typeString:
		return `""`
	case goType == "bool":
		return "false"
	case strings.HasPrefix(goType, "int") || strings.HasPrefix(goType, "uint"):
		return "0"
	case strings.HasPrefix(goType, "float"):
		return "0.0"
	case IsSlice(goType):
		return typeNil
	case IsMap(goType):
		return typeNil
	case goType == "time.Time":
		return "time.Time{}"
	case goType == typeUUID:
		return "uuid.Nil"
	default:
		return goType + "{}"
	}
}

// FieldName returns the appropriate field name based on context
func FieldName(name string, unexported bool) string {
	if unexported {
		return CamelCase(name)
	}
	return name
}

// ExportedFieldName always returns an exported field name
func ExportedFieldName(name string) string {
	if name == "" {
		return name
	}
	// Ensure first letter is uppercase
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
