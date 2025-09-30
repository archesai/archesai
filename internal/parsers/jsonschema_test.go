package parsers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJSONSchema(t *testing.T) {
	tests := []struct {
		name             string
		filePath         string
		wantErr          bool
		errContains      string
		validateSchema   func(t *testing.T, schema interface{})
		validateXcodegen func(t *testing.T, xcodegen *interface{})
	}{
		{
			name:     "valid schema with x-codegen",
			filePath: "../../test/data/parsers/schemas/with-x-codegen.yaml",
			wantErr:  false,
			validateSchema: func(t *testing.T, schema interface{}) {
				assert.NotNil(t, schema)
			},
		},
		{
			name:     "simple schema without x-codegen",
			filePath: "../../test/data/parsers/schemas/simple.yaml",
			wantErr:  false,
			validateSchema: func(t *testing.T, schema interface{}) {
				assert.NotNil(t, schema)
			},
			validateXcodegen: func(t *testing.T, xcodegen *interface{}) {
				assert.Nil(t, xcodegen) // Should be nil for schemas without x-codegen
			},
		},
		{
			name:     "complex schema with nested objects",
			filePath: "../../test/data/parsers/schemas/complex.yaml",
			wantErr:  false,
			validateSchema: func(t *testing.T, schema interface{}) {
				assert.NotNil(t, schema)
			},
			validateXcodegen: func(t *testing.T, xcodegen *interface{}) {
				assert.Nil(t, xcodegen)
			},
		},
		{
			name:        "non-existent file",
			filePath:    "../../test/data/parsers/schemas/non-existent.yaml",
			wantErr:     true,
			errContains: "failed to read schema file",
		},
		{
			name:        "invalid yaml",
			filePath:    createTempYamlFile(t, "invalid: [\nbad: yaml"),
			wantErr:     true,
			errContains: "failed to unmarshal schema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, _, err := ParseJSONSchema(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.validateSchema != nil {
					tt.validateSchema(t, schema)
				}

			}
		})
	}
}

func TestProcessSchema(t *testing.T) {
	// Test with real schema files
	tests := []struct {
		name       string
		schemaFile string
		schemaName string
		wantErr    bool
		validate   func(t *testing.T, result *ProcessedSchema)
	}{
		{
			name:       "process simple schema",
			schemaFile: "../../test/data/parsers/schemas/simple.yaml",
			schemaName: "SimpleSchema",
			wantErr:    false,
			validate: func(t *testing.T, result *ProcessedSchema) {
				assert.NotNil(t, result)
				assert.Equal(t, "SimpleSchema", result.Name)
				assert.Equal(t, "SimpleSchema", result.Title)
				assert.NotNil(t, result.Schema)
				assert.Len(t, result.RequiredFields, 3)

				// Check required fields by JSON tag
				requiredJSONTags := make(map[string]bool)
				for _, f := range result.RequiredFields {
					requiredJSONTags[f.JSONTag] = true
				}
				assert.True(t, requiredJSONTags["id"])
				assert.True(t, requiredJSONTags["name"])
				assert.True(t, requiredJSONTags["email"])

				// Check fields by JSON tag
				assert.NotNil(t, result.Fields)
				fieldJSONTags := make(map[string]bool)
				for _, f := range result.Fields {
					fieldJSONTags[f.JSONTag] = true
				}
				assert.True(t, fieldJSONTags["id"])
				assert.True(t, fieldJSONTags["name"])
				assert.True(t, fieldJSONTags["email"])
				assert.True(t, fieldJSONTags["age,omitempty"])
				assert.True(t, fieldJSONTags["isActive,omitempty"])
				assert.True(t, fieldJSONTags["createdAt,omitempty"])
			},
		},
		{
			name:       "process complex schema",
			schemaFile: "../../test/data/parsers/schemas/complex.yaml",
			schemaName: "ComplexSchema",
			wantErr:    false,
			validate: func(t *testing.T, result *ProcessedSchema) {
				assert.NotNil(t, result)
				assert.Equal(t, "ComplexSchema", result.Name)
				assert.NotNil(t, result.Schema)

				// Check required fields by JSON tag
				requiredJSONTags := make(map[string]bool)
				for _, f := range result.RequiredFields {
					requiredJSONTags[f.JSONTag] = true
				}
				assert.True(t, requiredJSONTags["id"])
				assert.True(t, requiredJSONTags["profile"])

				// Check fields by JSON tag
				fieldJSONTags := make(map[string]bool)
				for _, f := range result.Fields {
					fieldJSONTags[f.JSONTag] = true
				}
				assert.True(t, fieldJSONTags["id"])
				assert.True(t, fieldJSONTags["profile"])
				assert.True(t, fieldJSONTags["tags,omitempty"])
				assert.True(t, fieldJSONTags["metadata,omitempty"])
				assert.True(t, fieldJSONTags["addresses,omitempty"])
			},
		},
		{
			name:       "process schema with x-codegen",
			schemaFile: "../../test/data/parsers/schemas/with-x-codegen.yaml",
			schemaName: "UserEntity",
			wantErr:    false,
			validate: func(t *testing.T, result *ProcessedSchema) {
				assert.NotNil(t, result)
				assert.Equal(t, "UserEntity", result.Name)
				assert.NotNil(t, result.XCodegen)
				assert.Equal(t, XCodegenExtensionSchemaType("entity"), result.XCodegen.SchemaType)

				// Check required fields by JSON tag
				requiredJSONTags := make(map[string]bool)
				for _, f := range result.RequiredFields {
					requiredJSONTags[f.JSONTag] = true
				}
				assert.True(t, requiredJSONTags["id"])
				assert.True(t, requiredJSONTags["username"])
				assert.True(t, requiredJSONTags["email"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load schema from file
			schema, _, err := ParseJSONSchema(tt.schemaFile)
			require.NoError(t, err)

			result, err := ProcessSchema(schema, tt.schemaName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}

	// Test nil schema
	t.Run("nil schema", func(t *testing.T) {
		result, err := ProcessSchema(nil, "NilSchema")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestExtractTitle(t *testing.T) {
	// Test with actual schema from file
	schema, _, err := ParseJSONSchema("../../test/data/parsers/schemas/simple.yaml")
	require.NoError(t, err)

	// The simple schema doesn't have a title, so it should use the provided name
	title := extractTitle(schema, "ProvidedName")
	assert.Equal(t, "ProvidedName", title)

	// Test with nil schema
	title = extractTitle(nil, "DefaultName")
	assert.Equal(t, "DefaultName", title)
}

func TestExtractFields(t *testing.T) {
	// Test with simple schema
	t.Run("simple schema fields", func(t *testing.T) {
		schema, _, err := ParseJSONSchema("../../test/data/parsers/schemas/simple.yaml")
		require.NoError(t, err)

		fields := ExtractFields(schema)
		assert.NotEmpty(t, fields)
		assert.Len(t, fields, 6)

		// Check field details - map by JSON tag since field Names are now PascalCase
		fieldMap := make(map[string]FieldDef)
		for _, f := range fields {
			fieldMap[f.JSONTag] = f
		}

		// Test id field
		idField := fieldMap["id"]
		assert.Equal(t, "uuid", idField.Format)
		assert.True(t, idField.Required)

		// Test name field
		nameField := fieldMap["name"]
		assert.True(t, nameField.Required)

		// Test age field (will have omitempty tag)
		ageField := fieldMap["age,omitempty"]
		assert.False(t, ageField.Required)

		// Test isActive field (will have omitempty tag)
		isActiveField := fieldMap["isActive,omitempty"]
		assert.False(t, isActiveField.Required)

		// Test email field
		emailField := fieldMap["email"]
		assert.Equal(t, "email", emailField.Format)
		assert.True(t, emailField.Required)

		// Test createdAt field (will have omitempty tag)
		createdAtField := fieldMap["createdAt,omitempty"]
		assert.Equal(t, "date-time", createdAtField.Format)
		assert.False(t, createdAtField.Required)
	})

	// Test with complex schema
	t.Run("complex schema fields", func(t *testing.T) {
		schema, _, err := ParseJSONSchema("../../test/data/parsers/schemas/complex.yaml")
		require.NoError(t, err)

		fields := ExtractFields(schema)
		assert.NotEmpty(t, fields)

		fieldMap := make(map[string]FieldDef)
		for _, f := range fields {
			fieldMap[f.Name] = f
		}

	})

	// Test with nil schema
	t.Run("nil schema", func(t *testing.T) {
		fields := ExtractFields(nil)
		assert.Nil(t, fields)
	})
}

func TestExtractRequiredFields(t *testing.T) {
	// Test with simple schema
	schema, _, err := ParseJSONSchema("../../test/data/parsers/schemas/simple.yaml")
	require.NoError(t, err)

	required := ExtractRequiredFields(schema)
	assert.Len(t, required, 3)
	// ExtractRequiredFields now returns normalized Go field names (PascalCase)
	assert.Contains(t, required, "ID")
	assert.Contains(t, required, "Email")
	assert.Contains(t, required, "Name")

	// Test with nil schema
	required = ExtractRequiredFields(nil)
	assert.Empty(t, required)
}

func TestExtractEnumValues(t *testing.T) {
	// Test with schema that has enum (with-x-codegen has role enum)
	schema, _, err := ParseJSONSchema("../../test/data/parsers/schemas/with-x-codegen.yaml")
	require.NoError(t, err)

	// Get the role field which has enum values
	roleField := schema.Properties.GetOrZero("role")
	assert.NotNil(t, roleField)

	enums := ExtractEnumValues(roleField.Left)
	assert.Len(t, enums, 3)
	assert.Contains(t, enums, "admin")
	assert.Contains(t, enums, "moderator")
	assert.Contains(t, enums, "user")

	// Test with nil schema
	enums = ExtractEnumValues(nil)
	assert.Empty(t, enums)
}

// Helper function to create temporary files for testing
func createTempYamlFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
	return filePath
}

func TestExtractFieldsWithNestedObjects(t *testing.T) {
	// Load the complex schema with nested objects
	schema, _, err := ParseJSONSchema("../../test/data/parsers/schemas/complex.yaml")
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	// Extract fields
	fields := ExtractFields(schema)

	// We should have 5 top-level fields
	expectedFieldCount := 5
	if len(fields) != expectedFieldCount {
		t.Errorf("Expected %d fields, got %d", expectedFieldCount, len(fields))
		for _, f := range fields {
			t.Logf("Field: %s, Type: %s", f.Name, f.GoType)
		}
	}

	// Check for the profile field
	var profileField *FieldDef
	for i := range fields {
		if fields[i].Name == "Profile" {
			profileField = &fields[i]
			break
		}
	}

	if profileField == nil {
		t.Fatal("Profile field not found")
	}

	// The profile field should be an inline struct
	if !strings.Contains(profileField.GoType, "struct {") {
		t.Errorf("Expected Profile to be an inline struct, got: %s", profileField.GoType)
	}

	// Check that the inline struct contains the expected nested fields
	expectedInProfile := []string{
		"FirstName string",
		"LastName string",
		"Avatar *string",
		"Preferences", // This should also be a nested struct
		`json:"firstName"`,
		`json:"lastName"`,
		`json:"avatar,omitempty"`,
	}

	for _, expected := range expectedInProfile {
		if !strings.Contains(profileField.GoType, expected) {
			t.Errorf("Profile struct should contain '%s'\nGot: %s", expected, profileField.GoType)
		}
	}

	// The preferences field within profile should be a nested struct
	if !strings.Contains(profileField.GoType, "Preferences struct {") {
		t.Errorf("Expected Preferences to be a nested inline struct within Profile")
	}

	// Check that preferences has the right fields
	expectedInPreferences := []string{
		"Theme string",
		"Language *string",
		`json:"theme"`,
		`json:"language,omitempty"`,
	}

	for _, expected := range expectedInPreferences {
		if !strings.Contains(profileField.GoType, expected) {
			t.Errorf(
				"Preferences struct should contain '%s'\nGot: %s",
				expected,
				profileField.GoType,
			)
		}
	}

	// Check the addresses field (array of objects)
	var addressesField *FieldDef
	for i := range fields {
		if fields[i].Name == "Addresses" {
			addressesField = &fields[i]
			break
		}
	}

	if addressesField == nil {
		t.Fatal("Addresses field not found")
	}

	// Addresses should be an array of inline structs
	if !strings.HasPrefix(addressesField.GoType, "[]struct {") {
		t.Errorf(
			"Expected Addresses to be an array of inline structs, got: %s",
			addressesField.GoType,
		)
	}

	// Check that the addresses struct contains the expected fields
	expectedInAddresses := []string{
		"Street string",
		"City string",
		"Country string",
		"PostalCode *string",
		`json:"street"`,
		`json:"city"`,
		`json:"country"`,
		`json:"postalCode,omitempty"`,
	}

	for _, expected := range expectedInAddresses {
		if !strings.Contains(addressesField.GoType, expected) {
			t.Errorf(
				"Addresses item struct should contain '%s'\nGot: %s",
				expected,
				addressesField.GoType,
			)
		}
	}

	// Check the metadata field (additionalProperties)
	var metadataField *FieldDef
	for i := range fields {
		if fields[i].Name == "Metadata" {
			metadataField = &fields[i]
			break
		}
	}

	if metadataField == nil {
		t.Fatal("Metadata field not found")
	}

	// Metadata with additionalProperties should be map[string]string
	expectedMetadataType := "map[string]string"
	if metadataField.GoType != expectedMetadataType {
		t.Errorf(
			"Expected Metadata type to be '%s', got: %s",
			expectedMetadataType,
			metadataField.GoType,
		)
	}
}
