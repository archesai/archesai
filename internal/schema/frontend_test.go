package schema

import (
	"testing"

	"github.com/archesai/archesai/internal/ref"
)

func TestSchema_ToDataTableColumn_NilSchema(t *testing.T) {
	var s *Schema
	got := s.ToDataTableColumn(nil)
	if got != nil {
		t.Errorf("ToDataTableColumn on nil schema should return nil, got %+v", got)
	}
}

func TestSchema_ToDataTableColumns(t *testing.T) {
	parent := &Schema{
		Title: "User",
		Properties: map[string]*ref.Ref[Schema]{
			"id": ref.NewInline(&Schema{Title: "ID", GoType: "uuid.UUID"}),
			"name": ref.NewInline(&Schema{
				Title:  "name",
				GoType: "string",
			}),
			"email": ref.NewInline(&Schema{
				Title:  "email",
				GoType: "string",
				Format: "email",
			}),
			"createdAt": ref.NewInline(&Schema{
				Title:  "createdAt",
				Format: "date-time",
			}),
		},
	}

	columns := parent.ToDataTableColumns()

	// ID should be skipped
	if len(columns) != 3 {
		t.Errorf(
			"ToDataTableColumns() returned %d columns, want 3 (ID should be skipped)",
			len(columns),
		)
	}

	// Check that we have the expected columns
	hasName := false
	hasEmail := false
	hasCreatedAt := false
	for _, col := range columns {
		switch col.AccessorKey {
		case "name":
			hasName = true
			if !col.IsLink {
				t.Error("name column should have IsLink=true")
			}
		case "email":
			hasEmail = true
		case "createdAt":
			hasCreatedAt = true
			if col.FilterVariant != "date" {
				t.Errorf("createdAt FilterVariant = %q, want date", col.FilterVariant)
			}
		}
	}

	if !hasName {
		t.Error("columns should include name")
	}
	if !hasEmail {
		t.Error("columns should include email")
	}
	if !hasCreatedAt {
		t.Error("columns should include createdAt")
	}
}

func TestSchema_ToDataTableColumns_NilSchema(t *testing.T) {
	var s *Schema
	got := s.ToDataTableColumns()
	if got != nil {
		t.Errorf("ToDataTableColumns on nil schema should return nil, got %+v", got)
	}
}

func TestSchema_ToFormField_NilSchema(t *testing.T) {
	var s *Schema
	got := s.ToFormField()
	if got != nil {
		t.Errorf("ToFormField on nil schema should return nil, got %+v", got)
	}
}

func TestSchema_ToFormField_ArrayField(t *testing.T) {
	s := &Schema{
		Title: "tags",
		Type:  PropertyType{Types: []string{"array"}},
	}

	got := s.ToFormField()
	if got != nil {
		t.Errorf("ToFormField should return nil for array fields, got %+v", got)
	}
}

func TestSchema_ToFormField_Textarea(t *testing.T) {
	maxLen := 200
	s := &Schema{
		Title:     "description",
		GoType:    "string",
		MaxLength: &maxLen,
	}

	got := s.ToFormField()
	if got == nil {
		t.Fatal("ToFormField should return non-nil for textarea field")
	}

	if got.Type != "textarea" {
		t.Errorf("Type = %q, want textarea", got.Type)
	}
}

func TestSchema_ToFormField_URL(t *testing.T) {
	s := &Schema{
		Title:  "website",
		Format: "uri",
	}

	got := s.ToFormField()
	if got == nil {
		t.Fatal("ToFormField should return non-nil for url field")
	}

	if got.Type != "url" {
		t.Errorf("Type = %q, want url", got.Type)
	}
}

func TestSchema_ToFormField_Password(t *testing.T) {
	s := &Schema{
		Title:  "password",
		Format: "password",
	}

	got := s.ToFormField()
	if got == nil {
		t.Fatal("ToFormField should return non-nil for password field")
	}

	if got.Type != "password" {
		t.Errorf("Type = %q, want password", got.Type)
	}
}

func TestSchema_ToFormFields(t *testing.T) {
	parent := &Schema{
		Title: "CreateUser",
		Properties: map[string]*ref.Ref[Schema]{
			"name": ref.NewInline(&Schema{
				Title:  "name",
				GoType: "string",
			}),
			"email": ref.NewInline(&Schema{
				Title:  "email",
				GoType: "string",
				Format: "email",
			}),
			"ID": ref.NewInline(&Schema{
				Title:  "ID",
				GoType: "uuid.UUID",
			}),
		},
	}

	fields := parent.ToFormFields()

	// ID should be skipped
	if len(fields) != 2 {
		t.Errorf("ToFormFields() returned %d fields, want 2 (ID should be skipped)", len(fields))
	}

	hasName := false
	hasEmail := false
	for _, field := range fields {
		switch field.Name {
		case "name":
			hasName = true
			if field.Type != "text" {
				t.Errorf("name Type = %q, want text", field.Type)
			}
		case "email":
			hasEmail = true
			if field.Type != "email" {
				t.Errorf("email Type = %q, want email", field.Type)
			}
		}
	}

	if !hasName {
		t.Error("fields should include name")
	}
	if !hasEmail {
		t.Error("fields should include email")
	}
}

func TestSchema_ToFormFields_NilSchema(t *testing.T) {
	var s *Schema
	got := s.ToFormFields()
	if got != nil {
		t.Errorf("ToFormFields on nil schema should return nil, got %+v", got)
	}
}

func TestSchema_ToDataTableColumn_JSONTagWithOmitempty(t *testing.T) {
	parent := &Schema{Title: "User"}
	s := &Schema{
		Title:   "firstName",
		GoType:  "string",
		JSONTag: "first_name,omitempty",
	}

	got := s.ToDataTableColumn(parent)
	if got == nil {
		t.Fatal("expected non-nil column")
	}

	if got.AccessorKey != "first_name" {
		t.Errorf(
			"AccessorKey = %q, want first_name (omitempty should be stripped)",
			got.AccessorKey,
		)
	}
}

func TestSchema_ToFormField_JSONTagWithOmitempty(t *testing.T) {
	s := &Schema{
		Title:   "firstName",
		GoType:  "string",
		JSONTag: "first_name,omitempty",
	}

	got := s.ToFormField()
	if got == nil {
		t.Fatal("expected non-nil field")
	}

	if got.Name != "first_name" {
		t.Errorf("Name = %q, want first_name (omitempty should be stripped)", got.Name)
	}
}

func TestSchema_ToDataTableColumn_NilParent(t *testing.T) {
	s := &Schema{
		Title:  "name",
		GoType: "string",
	}

	// Should not panic with nil parent
	got := s.ToDataTableColumn(nil)
	if got == nil {
		t.Fatal("expected non-nil column")
	}

	// IsLink should be true, but LinkParam should be empty with nil parent
	if !got.IsLink {
		t.Error("name field should have IsLink=true")
	}
	if got.LinkParam != "" {
		t.Errorf("LinkParam should be empty with nil parent, got %q", got.LinkParam)
	}
}

func TestSchema_ToDataTableColumn_FieldTypes(t *testing.T) {
	parent := &Schema{Title: "User"}

	tests := []struct {
		name     string
		prop     *Schema
		expected *DataTableColumn
	}{
		{
			name:     "skip ID",
			prop:     &Schema{Title: "ID"},
			expected: nil,
		},
		{
			name: "text field",
			prop: &Schema{Title: "firstName", GoType: "string"},
			expected: &DataTableColumn{
				AccessorKey:   "firstName",
				Label:         "First Name",
				FilterVariant: "text",
				Icon:          "TextIcon",
				EnableFilter:  true,
				EnableSort:    true,
			},
		},
		{
			name: "date field",
			prop: &Schema{Title: "createdAt", Format: "date-time"},
			expected: &DataTableColumn{
				AccessorKey:   "createdAt",
				Label:         "Created At",
				FilterVariant: "date",
				Icon:          "CalendarIcon",
				EnableFilter:  true,
				EnableSort:    true,
			},
		},
		{
			name: "boolean field",
			prop: &Schema{Title: "active", GoType: "bool"},
			expected: &DataTableColumn{
				AccessorKey:   "active",
				Label:         "Active",
				FilterVariant: "boolean",
				Icon:          "CheckIcon",
				EnableFilter:  true,
				EnableSort:    true,
			},
		},
		{
			name: "enum field",
			prop: &Schema{
				Title: "status",
				Enum:  []string{"active", "inactive"},
			},
			expected: &DataTableColumn{
				AccessorKey:   "status",
				Label:         "Status",
				FilterVariant: "multiSelect",
				Icon:          "TextIcon",
				EnableFilter:  true,
				EnableSort:    true,
				Options: []FilterOption{
					{Label: "Active", Value: "active"},
					{Label: "Inactive", Value: "inactive"},
				},
			},
		},
		{
			name: "name field with link",
			prop: &Schema{Title: "name", GoType: "string"},
			expected: &DataTableColumn{
				AccessorKey:   "name",
				Label:         "Name",
				FilterVariant: "text",
				Icon:          "TextIcon",
				EnableFilter:  true,
				EnableSort:    true,
				IsLink:        true,
				LinkParam:     "userID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prop.ToDataTableColumn(parent)

			if tt.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}

			if got == nil {
				t.Fatal("expected non-nil result")
			}

			if got.AccessorKey != tt.expected.AccessorKey {
				t.Errorf("AccessorKey = %q, want %q", got.AccessorKey, tt.expected.AccessorKey)
			}
			if got.Label != tt.expected.Label {
				t.Errorf("Label = %q, want %q", got.Label, tt.expected.Label)
			}
			if got.FilterVariant != tt.expected.FilterVariant {
				t.Errorf(
					"FilterVariant = %q, want %q",
					got.FilterVariant,
					tt.expected.FilterVariant,
				)
			}
			if got.Icon != tt.expected.Icon {
				t.Errorf("Icon = %q, want %q", got.Icon, tt.expected.Icon)
			}
			if got.IsLink != tt.expected.IsLink {
				t.Errorf("IsLink = %v, want %v", got.IsLink, tt.expected.IsLink)
			}
			if got.LinkParam != tt.expected.LinkParam {
				t.Errorf("LinkParam = %q, want %q", got.LinkParam, tt.expected.LinkParam)
			}
			if len(got.Options) != len(tt.expected.Options) {
				t.Errorf("Options length = %d, want %d", len(got.Options), len(tt.expected.Options))
			}
		})
	}
}

func TestSchema_ToFormField_FieldTypes(t *testing.T) {
	tests := []struct {
		name     string
		prop     *Schema
		expected *FormField
	}{
		{
			name:     "skip ID",
			prop:     &Schema{Title: "ID"},
			expected: nil,
		},
		{
			name:     "skip CreatedAt",
			prop:     &Schema{Title: "CreatedAt"},
			expected: nil,
		},
		{
			name:     "skip UpdatedAt",
			prop:     &Schema{Title: "UpdatedAt"},
			expected: nil,
		},
		{
			name: "text field",
			prop: &Schema{Title: "firstName", GoType: "string"},
			expected: &FormField{
				Name:     "firstName",
				Label:    "First Name",
				Type:     "text",
				Required: true,
			},
		},
		{
			name: "optional field",
			prop: &Schema{Title: "nickname", GoType: "string", Nullable: true},
			expected: &FormField{
				Name:     "nickname",
				Label:    "Nickname",
				Type:     "text",
				Required: false,
			},
		},
		{
			name: "boolean field",
			prop: &Schema{Title: "active", GoType: "bool"},
			expected: &FormField{
				Name:         "active",
				Label:        "Active",
				Type:         "checkbox",
				Required:     true,
				DefaultValue: "false",
			},
		},
		{
			name: "email field",
			prop: &Schema{Title: "email", Format: "email"},
			expected: &FormField{
				Name:     "email",
				Label:    "Email",
				Type:     "email",
				Required: true,
			},
		},
		{
			name: "date field",
			prop: &Schema{Title: "birthDate", Format: "date"},
			expected: &FormField{
				Name:     "birthDate",
				Label:    "Birth Date",
				Type:     "date",
				Required: true,
			},
		},
		{
			name: "select field from enum",
			prop: &Schema{
				Title: "role",
				Enum:  []string{"admin", "user"},
			},
			expected: &FormField{
				Name:     "role",
				Label:    "Role",
				Type:     "select",
				Required: true,
				Options: []FormFieldOption{
					{Label: "Admin", Value: "admin"},
					{Label: "User", Value: "user"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prop.ToFormField()

			if tt.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}

			if got == nil {
				t.Fatal("expected non-nil result")
			}

			if got.Name != tt.expected.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.expected.Name)
			}
			if got.Label != tt.expected.Label {
				t.Errorf("Label = %q, want %q", got.Label, tt.expected.Label)
			}
			if got.Type != tt.expected.Type {
				t.Errorf("Type = %q, want %q", got.Type, tt.expected.Type)
			}
			if got.Required != tt.expected.Required {
				t.Errorf("Required = %v, want %v", got.Required, tt.expected.Required)
			}
			if len(got.Options) != len(tt.expected.Options) {
				t.Errorf("Options length = %d, want %d", len(got.Options), len(tt.expected.Options))
			}
		})
	}
}
