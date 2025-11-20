package parsers

import (
	"os"
	"path/filepath"
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
		validateSchema   func(t *testing.T, schema any)
		validateXcodegen func(t *testing.T, xcodegen *any)
	}{
		{
			name:     "valid schema with x-codegen",
			filePath: "../../test/data/parsers/schemas/with-x-codegen.yaml",
			wantErr:  false,
			validateSchema: func(t *testing.T, schema any) {
				assert.NotNil(t, schema)
			},
		},
		{
			name:     "simple schema without x-codegen",
			filePath: "../../test/data/parsers/schemas/simple.yaml",
			wantErr:  false,
			validateSchema: func(t *testing.T, schema any) {
				assert.NotNil(t, schema)
			},
			validateXcodegen: func(t *testing.T, xcodegen *any) {
				assert.Nil(t, xcodegen) // Should be nil for schemas without x-codegen
			},
		},
		{
			name:     "complex schema with nested objects",
			filePath: "../../test/data/parsers/schemas/complex.yaml",
			wantErr:  false,
			validateSchema: func(t *testing.T, schema any) {
				assert.NotNil(t, schema)
			},
			validateXcodegen: func(t *testing.T, xcodegen *any) {
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
			parser := NewJSONSchemaParser(nil)
			schema, err := parser.ParseFile(tt.filePath)
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
		wantErr    bool
		validate   func(t *testing.T, result *SchemaDef)
	}{
		{
			name:       "process simple schema",
			schemaFile: "../../test/data/parsers/schemas/simple.yaml",
			wantErr:    false,
			validate: func(t *testing.T, result *SchemaDef) {
				assert.NotNil(t, result)
				assert.Equal(t, "SimpleSchema", result.Name)
				assert.NotNil(t, result.Schema)

				// Check fields by JSON tag
				assert.NotNil(t, result.GetSortedProperties())
				fieldJSONTags := make(map[string]bool)
				for _, f := range result.GetSortedProperties() {
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
			wantErr:    false,
			validate: func(t *testing.T, result *SchemaDef) {
				assert.NotNil(t, result)
				assert.Equal(t, "ComplexSchema", result.Name)
				assert.NotNil(t, result.Schema)

				// Check fields by JSON tag
				fieldJSONTags := make(map[string]bool)
				for _, f := range result.GetSortedProperties() {
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
			wantErr:    false,
			validate: func(t *testing.T, result *SchemaDef) {
				assert.NotNil(t, result)
				assert.Equal(t, "UserEntity", result.Name)
				assert.Equal(t, XCodegenSchemaType("entity"), result.XCodegenSchemaType)

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load schema from file
			parser := NewJSONSchemaParser(nil)
			result, err := parser.ParseFile(tt.schemaFile)
			require.NoError(t, err)

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
		parser := NewJSONSchemaParser(nil)
		result, err := parser.ParseFile("")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestExtractFields(t *testing.T) {
	// Test with simple schema
	t.Run("simple schema fields", func(t *testing.T) {
		parser := NewJSONSchemaParser(nil)
		schema, err := parser.ParseFile("../../test/data/parsers/schemas/simple.yaml")
		require.NoError(t, err)

		assert.NotEmpty(t, schema.GetSortedProperties())
		assert.Len(t, schema.GetSortedProperties(), 6)

		// Check field details - map by JSON tag since field Names are no
		fieldMap := make(map[string]SchemaDef)
		for _, f := range schema.GetSortedProperties() {
			fieldMap[f.JSONTag] = *f
		}

		// Test id field
		assert.Equal(t, "uuid", fieldMap["id"].Format)
		assert.True(t, schema.IsPropertyRequired("id"))

		// Test name field
		assert.True(t, schema.IsPropertyRequired("name"))

		// Test age field (will have omitempty tag)
		assert.False(t, schema.IsPropertyRequired("age"))

		// Test isActive field (will have omitempty tag)
		assert.False(t, schema.IsPropertyRequired("isActive"))

		// Test email field
		assert.Equal(t, "email", fieldMap["email"].Format)
		assert.True(t, schema.IsPropertyRequired("email"))

		// Test createdAt field (will have omitempty tag)
		assert.Equal(t, "date-time", fieldMap["createdAt,omitempty"].Format)
		assert.False(t, schema.IsPropertyRequired("createdAt"))
	})

	// Test with complex schema
	t.Run("complex schema fields", func(t *testing.T) {
		parser := NewJSONSchemaParser(nil)
		schema, err := parser.ParseFile("../../test/data/parsers/schemas/complex.yaml")
		require.NoError(t, err)
		assert.NotEmpty(t, schema.GetSortedProperties())

		fieldMap := make(map[string]SchemaDef)
		for _, f := range schema.GetSortedProperties() {
			fieldMap[f.Name] = *f
		}

	})

	// Test with nil schema
	t.Run("nil schema", func(t *testing.T) {
		parser := NewJSONSchemaParser(nil)
		jsonSchema, err := parser.ParseFile("")
		assert.Error(t, err)
		assert.Nil(t, jsonSchema)
	})
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
	parser := NewJSONSchemaParser(nil)
	schema, err := parser.ParseFile("../../test/data/parsers/schemas/complex.yaml")
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	// We should have 5 top-level fields
	expectedFieldCount := 5
	if len(schema.GetSortedProperties()) != expectedFieldCount {
		t.Errorf(
			"Expected %d fields, got %d",
			expectedFieldCount,
			len(schema.GetSortedProperties()),
		)
		for _, f := range schema.GetSortedProperties() {
			t.Logf("Field: %s, Type: %s", f.Name, f.GoType)
		}
	}

	// Check for the profile field
	var profileField *SchemaDef
	for i := range schema.GetSortedProperties() {
		if schema.GetSortedProperties()[i].Name == "Profile" {
			profileField = schema.GetSortedProperties()[i]
			break
		}
	}

	if profileField == nil {
		t.Fatal("Profile field not found")
	}

	// The profile field should be a named type (not inline struct)
	// The code generator creates named types for nested objects for better maintainability
	expectedProfileType := "ComplexSchemaProfile"
	if profileField.GoType != expectedProfileType {
		t.Errorf(
			"Expected Profile type to be '%s', got: %s",
			expectedProfileType,
			profileField.GoType,
		)
	}

	// Check the addresses field (array of objects)
	var addressesField *SchemaDef
	for i := range schema.GetSortedProperties() {
		if schema.GetSortedProperties()[i].Name == "Addresses" {
			addressesField = schema.GetSortedProperties()[i]
			break
		}
	}

	if addressesField == nil {
		t.Fatal("Addresses field not found")
	}

	// Addresses should be an array of named type (not inline struct)
	// The code generator creates named types for array item objects for better maintainability
	expectedAddressesType := "[]ComplexSchemaAddressesItem"
	if addressesField.GoType != expectedAddressesType {
		t.Errorf(
			"Expected Addresses type to be '%s', got: %s",
			expectedAddressesType,
			addressesField.GoType,
		)
	}

	// Check the metadata field (additionalProperties)
	var metadataField *SchemaDef
	for i := range schema.GetSortedProperties() {
		if schema.GetSortedProperties()[i].Name == "Metadata" {
			metadataField = schema.GetSortedProperties()[i]
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
