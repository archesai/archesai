package strutil

import (
	"testing"
)

func TestHumanReadable(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "Name"},
		{"firstName", "First Name"},
		{"createdAt", "Created At"},
		{"ID", "ID"},
		{"userID", "User ID"},
		{"HTTPResponse", "HTTP Response"},
		{"APIKey", "API Key"},
	}

	for _, tt := range tests {
		got := HumanReadable(tt.input)
		if got != tt.expected {
			t.Errorf("HumanReadable(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple lowercase", "test", "test"},
		{"simple uppercase", "Test", "test"},
		{"camelCase", "testCase", "test-case"},
		{"PascalCase", "TestCase", "test-case"},
		{"acronym at start", "APIKey", "api-key"},
		{"acronym in middle", "getAPIKey", "get-api-key"},
		{"acronym at end", "myAPI", "my-api"},
		{"multiple words", "PipelineStep", "pipeline-step"},
		{"already kebab", "test-case", "test-case"},
		{"snake_case input", "test_case", "test-case"},
		{"with spaces", "test case", "test-case"},
		{"complex", "CreateAPIKeyRequest", "create-api-key-request"},
		{"all uppercase", "API", "api"},
		{"id suffix", "UserID", "user-id"},
		{"url in name", "APIURLPath", "apiurl-path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("KebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple lowercase", "test", "test"},
		{"simple uppercase", "Test", "test"},
		{"camelCase", "testCase", "test_case"},
		{"PascalCase", "TestCase", "test_case"},
		{"acronym at start", "APIKey", "api_key"},
		{"acronym in middle", "getAPIKey", "get_api_key"},
		{"multiple words", "PipelineStep", "pipeline_step"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple lowercase", "test", "test"},
		{"PascalCase", "TestCase", "testCase"},
		{"kebab-case", "test-case", "testCase"},
		{"snake_case", "test_case", "testCase"},
		{"with space", "test case", "testCase"},
		{"acronym", "APIKey", "apiKey"},
		{"acronym preserve", "createAPIKey", "createAPIKey"},
		{"compound acronym", "APIKeysFilter", "apiKeysFilter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("CamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple lowercase", "test", "Test"},
		{"camelCase", "testCase", "TestCase"},
		{"kebab-case", "test-case", "TestCase"},
		{"snake_case", "test_case", "TestCase"},
		{"already PascalCase", "TestCase", "TestCase"},
		{"acronym", "apiKey", "APIKey"},
		{"acronym preserve", "CreateAPIKey", "CreateAPIKey"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("PascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple word", "user", "users"},
		{"word ending in y", "category", "categories"},
		{"word ending in s", "class", "classes"},
		{"word ending in x", "box", "boxes"},
		{"word ending in ch", "match", "matches"},
		{"irregular", "person", "people"},
		{"uncountable", "health", "health"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pluralize(tt.input)
			if result != tt.expected {
				t.Errorf("Pluralize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
