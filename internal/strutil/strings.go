package strutil

import (
	"bytes"
	"slices"
	"strings"
	"unicode"
)

// Common acronyms that should be preserved in uppercase
var commonAcronyms = map[string]string{
	"api":     "API",
	"apikey":  "APIKey",
	"apikeys": "APIKeys",
	"url":     "URL",
	"uri":     "URI",
	"uuid":    "UUID",
	"id":      "ID",
	"http":    "HTTP",
	"https":   "HTTPS",
	"sql":     "SQL",
	"json":    "JSON",
	"xml":     "XML",
	"csv":     "CSV",
	"jwt":     "JWT",
	"oauth":   "OAuth",
	"saml":    "SAML",
	"ldap":    "LDAP",
	"dns":     "DNS",
	"tcp":     "TCP",
	"udp":     "UDP",
	"ip":      "IP",
	"vm":      "VM",
	"os":      "OS",
	"cpu":     "CPU",
	"gpu":     "GPU",
	"ram":     "RAM",
	"ssd":     "SSD",
	"hdd":     "HDD",
	"cdn":     "CDN",
	"vpn":     "VPN",
	"ssh":     "SSH",
	"ftp":     "FTP",
	"sftp":    "SFTP",
	"smtp":    "SMTP",
	"imap":    "IMAP",
	"pop":     "POP",
	"aws":     "AWS",
	"gcp":     "GCP",
	"sdk":     "SDK",
	"ci":      "CI",
	"cd":      "CD",
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
				result += PascalCase(lowerPart)
			}
		}
	}

	// Handle Go reserved keywords
	return sanitizeGoKeyword(result)
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
				result += strings.ToUpper(part[:1]) + part[1:]
			}
		}
	}

	return sanitizeGoKeyword(result)
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

	return sanitizeGoKeyword(result.String())
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
		"health": "health", // Health is uncountable, stays the same
		"config": "config", // Config is uncountable, stays the same
		"apikey": "APIKeys",
		"APIKey": "APIKeys",
	}

	lower := strings.ToLower(word)
	if plural, ok := special[lower]; ok {
		// Preserve original case
		if unicode.IsUpper(rune(word[0])) {
			return PascalCase(plural)
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

// Contains checks if a slice contains a string.
func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// IsPointer checks if a Go type string represents a pointer type
func IsPointer(goType string) bool {
	return strings.HasPrefix(goType, "*")
}

// IsSlice checks if a Go type string represents a slice type
func IsSlice(goType string) bool {
	return strings.HasPrefix(goType, "[]")
}

// EchoPath converts OpenAPI path format to Echo router format.
// Converts "/users/{id}" to "/users/:id"
func EchoPath(path string) string {
	// Replace {param} with :param for Echo router format
	result := strings.ReplaceAll(path, "{", ":")
	result = strings.ReplaceAll(result, "}", "")
	return result
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

// SplitByComma splits a string by comma and trims whitespace.
func SplitByComma(s string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			trimmed := TrimSpace(current)
			if trimmed != "" {
				result = append(result, trimmed)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	trimmed := TrimSpace(current)
	if trimmed != "" {
		result = append(result, trimmed)
	}
	return result
}

// TrimSpace trims leading and trailing whitespace.
func TrimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
