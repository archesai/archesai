package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/archesai/archesai/internal/ref"
)

func TestLoader_LoadSchemaFromBytes_Basic(t *testing.T) {
	schemaYAML := `
title: User
type: object
x-codegen-schema-type: valueobject
required:
  - name
  - email
properties:
  name:
    type: string
  email:
    type: string
    format: email
  age:
    type: integer
`

	loader := NewLoader("")
	schema, err := loader.LoadSchemaFromBytes([]byte(schemaYAML), "User.yaml")
	if err != nil {
		t.Fatalf("failed to load schema: %v", err)
	}

	if schema.Title != "User" {
		t.Errorf("expected title 'User', got %q", schema.Title)
	}

	if schema.GoType != "User" {
		t.Errorf("expected GoType 'User', got %q", schema.GoType)
	}

	if len(schema.Properties) != 3 {
		t.Errorf("expected 3 properties, got %d", len(schema.Properties))
	}

	// Check name property
	nameProp := schema.GetProperty("Name")
	if nameProp == nil {
		t.Fatal("expected 'Name' property to exist")
	}
	if nameProp.GoType != "string" {
		t.Errorf("expected Name.GoType 'string', got %q", nameProp.GoType)
	}
	if nameProp.JSONTag != "name" {
		t.Errorf("expected Name.JSONTag 'name', got %q", nameProp.JSONTag)
	}

	// Check email property
	emailProp := schema.GetProperty("Email")
	if emailProp == nil {
		t.Fatal("expected 'Email' property to exist")
	}
	if emailProp.GoType != "string" {
		t.Errorf("expected Email.GoType 'string', got %q", emailProp.GoType)
	}

	// Check age property (optional)
	ageProp := schema.GetProperty("Age")
	if ageProp == nil {
		t.Fatal("expected 'Age' property to exist")
	}
	if ageProp.GoType != "int" {
		t.Errorf("expected Age.GoType 'int', got %q", ageProp.GoType)
	}
	if ageProp.JSONTag != "age,omitempty" {
		t.Errorf("expected Age.JSONTag 'age,omitempty', got %q", ageProp.JSONTag)
	}
}

func TestLoader_ProcessProperty_Types(t *testing.T) {
	tests := []struct {
		name       string
		schemaType string
		format     string
		wantGoType string
	}{
		{"string", "string", "", "string"},
		{"string-email", "string", "email", "string"},
		{"string-uuid", "string", "uuid", "uuid.UUID"},
		{"string-datetime", "string", "date-time", "time.Time"},
		{"integer", "integer", "", "int"},
		{"integer-int32", "integer", "int32", "int32"},
		{"integer-int64", "integer", "int64", "int64"},
		{"number", "number", "", "float64"},
		{"number-float", "number", "float", "float32"},
		{"boolean", "boolean", "", "bool"},
	}

	loader := NewLoader("")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prop := &Schema{
				Type:   PropertyType{Types: []string{tt.schemaType}},
				Format: tt.format,
			}
			loader.ProcessProperty(prop, "Field", "field", "Parent")

			if prop.GoType != tt.wantGoType {
				t.Errorf("expected GoType %q, got %q", tt.wantGoType, prop.GoType)
			}
		})
	}
}

func TestLoader_ProcessProperty_Array(t *testing.T) {
	loader := NewLoader("")

	// Array of strings
	prop := &Schema{
		Type: PropertyType{Types: []string{"array"}},
		Items: ref.NewInline(&Schema{
			Type: PropertyType{Types: []string{"string"}},
		}),
	}
	loader.ProcessProperty(prop, "Tags", "tags", "User")

	if prop.GoType != "[]string" {
		t.Errorf("expected GoType '[]string', got %q", prop.GoType)
	}
}

func TestLoader_ProcessProperty_NestedObject(t *testing.T) {
	loader := NewLoader("")

	prop := &Schema{
		Type: PropertyType{Types: []string{"object"}},
		Properties: map[string]*ref.Ref[Schema]{
			"street": ref.NewInline(&Schema{
				Type: PropertyType{Types: []string{"string"}},
			}),
			"city": ref.NewInline(&Schema{
				Type: PropertyType{Types: []string{"string"}},
			}),
		},
	}
	loader.ProcessProperty(prop, "Address", "address", "User")

	if prop.GoType != "UserAddress" {
		t.Errorf("expected GoType 'UserAddress', got %q", prop.GoType)
	}

	// Check nested properties were processed
	streetProp := prop.GetProperty("Street")
	if streetProp == nil {
		t.Fatal("expected 'Street' property to exist")
	}
	if streetProp.GoType != "string" {
		t.Errorf("expected Street.GoType 'string', got %q", streetProp.GoType)
	}
}

func TestLoader_LoadSchemaFile_WithRefs(t *testing.T) {
	// Create temp directory with schema files
	tmpDir := t.TempDir()

	// Write base schema
	baseSchema := `
title: Address
type: object
properties:
  street:
    type: string
  city:
    type: string
`
	if err := os.WriteFile(filepath.Join(tmpDir, "Address.yaml"), []byte(baseSchema), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write main schema with ref
	mainSchema := `
title: User
type: object
properties:
  name:
    type: string
  address:
    $ref: ./Address.yaml
`
	mainPath := filepath.Join(tmpDir, "User.yaml")
	if err := os.WriteFile(mainPath, []byte(mainSchema), 0o644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader(tmpDir)
	schema, err := loader.LoadSchemaFile(mainPath)
	if err != nil {
		t.Fatalf("failed to load schema: %v", err)
	}

	if schema.Title != "User" {
		t.Errorf("expected title 'User', got %q", schema.Title)
	}

	// Check that address property was resolved
	addressProp := schema.GetProperty("Address")
	if addressProp == nil {
		t.Fatal("expected 'Address' property to exist")
	}
	if addressProp.GoType != "Address" {
		t.Errorf("expected Address.GoType 'Address', got %q", addressProp.GoType)
	}
}

func TestLoader_EntityBaseFields(t *testing.T) {
	schemaYAML := `
title: User
type: object
x-codegen-schema-type: entity
properties:
  name:
    type: string
`

	loader := NewLoader("")
	schema, err := loader.LoadSchemaFromBytes([]byte(schemaYAML), "User.yaml")
	if err != nil {
		t.Fatalf("failed to load schema: %v", err)
	}

	// Entity schemas should have base fields added
	if schema.GetProperty("ID") == nil {
		t.Error("expected 'ID' base field to exist")
	}
	if schema.GetProperty("CreatedAt") == nil {
		t.Error("expected 'CreatedAt' base field to exist")
	}
	if schema.GetProperty("UpdatedAt") == nil {
		t.Error("expected 'UpdatedAt' base field to exist")
	}
}

func TestLoader_NullableOneOf(t *testing.T) {
	schemaYAML := `
title: User
type: object
properties:
  nickname:
    oneOf:
      - type: string
        minLength: 1
      - type: 'null'
`

	loader := NewLoader("")
	schema, err := loader.LoadSchemaFromBytes([]byte(schemaYAML), "User.yaml")
	if err != nil {
		t.Fatalf("failed to load schema: %v", err)
	}

	nicknameProp := schema.GetProperty("Nickname")
	if nicknameProp == nil {
		t.Fatal("expected 'Nickname' property to exist")
	}
	if !nicknameProp.Nullable {
		t.Error("expected Nickname to be nullable")
	}
	if nicknameProp.GoType != "string" {
		t.Errorf("expected Nickname.GoType 'string', got %q", nicknameProp.GoType)
	}
}
