package generators

import (
	"testing"

	"github.com/archesai/archesai/internal/spec"
)

func TestToHumanReadable(t *testing.T) {
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
		got := toHumanReadable(tt.input)
		if got != tt.expected {
			t.Errorf("toHumanReadable(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestMapPropertyToColumn(t *testing.T) {
	parent := &spec.Schema{Name: "User"}

	tests := []struct {
		name     string
		prop     *spec.Schema
		expected *DataTableColumn
	}{
		{
			name:     "skip ID",
			prop:     &spec.Schema{Name: "ID"},
			expected: nil,
		},
		{
			name: "text field",
			prop: &spec.Schema{Name: "firstName", GoType: "string"},
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
			prop: &spec.Schema{Name: "createdAt", Format: "date-time"},
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
			prop: &spec.Schema{Name: "active", GoType: "bool"},
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
			prop: &spec.Schema{
				Name: "status",
				Enum: []string{"active", "inactive"},
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
			prop: &spec.Schema{Name: "name", GoType: "string"},
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
			got := mapPropertyToColumn(tt.prop, parent)

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

func TestMapPropertyToFormField(t *testing.T) {
	tests := []struct {
		name     string
		prop     *spec.Schema
		expected *FormField
	}{
		{
			name:     "skip ID",
			prop:     &spec.Schema{Name: "ID"},
			expected: nil,
		},
		{
			name:     "skip CreatedAt",
			prop:     &spec.Schema{Name: "CreatedAt"},
			expected: nil,
		},
		{
			name:     "skip UpdatedAt",
			prop:     &spec.Schema{Name: "UpdatedAt"},
			expected: nil,
		},
		{
			name: "text field",
			prop: &spec.Schema{Name: "firstName", GoType: "string"},
			expected: &FormField{
				Name:     "firstName",
				Label:    "First Name",
				Type:     "text",
				Required: true,
			},
		},
		{
			name: "optional field",
			prop: &spec.Schema{Name: "nickname", GoType: "string", Nullable: true},
			expected: &FormField{
				Name:     "nickname",
				Label:    "Nickname",
				Type:     "text",
				Required: false,
			},
		},
		{
			name: "boolean field",
			prop: &spec.Schema{Name: "active", GoType: "bool"},
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
			prop: &spec.Schema{Name: "email", Format: "email"},
			expected: &FormField{
				Name:     "email",
				Label:    "Email",
				Type:     "email",
				Required: true,
			},
		},
		{
			name: "date field",
			prop: &spec.Schema{Name: "birthDate", Format: "date"},
			expected: &FormField{
				Name:     "birthDate",
				Label:    "Birth Date",
				Type:     "date",
				Required: true,
			},
		},
		{
			name: "select field from enum",
			prop: &spec.Schema{
				Name: "role",
				Enum: []string{"admin", "user"},
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
			got := mapPropertyToFormField(tt.prop)

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

func TestFindEntityOperations(t *testing.T) {
	ctx := &GeneratorContext{
		Spec: &spec.Spec{
			Operations: []spec.Operation{
				{ID: "ListUsers", Method: "GET"},
				{ID: "GetUser", Method: "GET"},
				{ID: "CreateUser", Method: "POST"},
				{ID: "UpdateUser", Method: "PUT"},
				{ID: "DeleteUser", Method: "DELETE"},
				{ID: "GetProducts", Method: "GET"}, // Different entity
			},
		},
	}

	ops := findEntityOperations(ctx, "User")

	if !ops.HasList {
		t.Error("expected HasList to be true")
	}
	if !ops.HasGet {
		t.Error("expected HasGet to be true")
	}
	if !ops.HasCreate {
		t.Error("expected HasCreate to be true")
	}
	if !ops.HasUpdate {
		t.Error("expected HasUpdate to be true")
	}
	if !ops.HasDelete {
		t.Error("expected HasDelete to be true")
	}
	if ops.ListOp == nil || ops.ListOp.ID != "ListUsers" {
		t.Error("expected ListOp to be ListUsers")
	}
}

func TestBuildEntityFrontendData(t *testing.T) {
	ctx := &GeneratorContext{
		Spec: &spec.Spec{
			Operations: []spec.Operation{
				{ID: "ListTodos", Method: "GET"},
				{ID: "CreateTodo", Method: "POST"},
			},
		},
		ProjectName: "example/myapp",
	}

	// Test with a simple schema (Properties are complex in the real spec)
	schema := &spec.Schema{
		Name: "Todo",
	}

	data := BuildEntityFrontendData(ctx, schema)

	if data.EntityName != "Todo" {
		t.Errorf("EntityName = %q, want %q", data.EntityName, "Todo")
	}
	if data.EntityNameLower != "todo" {
		t.Errorf("EntityNameLower = %q, want %q", data.EntityNameLower, "todo")
	}
	if data.EntityNameKebab != "todo" {
		t.Errorf("EntityNameKebab = %q, want %q", data.EntityNameKebab, "todo")
	}
	if data.EntityKey != "todos" {
		t.Errorf("EntityKey = %q, want %q", data.EntityKey, "todos")
	}

	// Check operations
	if !data.Operations.HasList {
		t.Error("expected HasList to be true")
	}
	if !data.Operations.HasCreate {
		t.Error("expected HasCreate to be true")
	}
}

func TestBuildEntityFrontendData_FormFieldsFromRequestBodies(t *testing.T) {
	// Create schema for create request body with required fields
	createRequestBodySchema := &spec.Schema{
		Name:   "CreateTodoRequestBody",
		GoType: "CreateTodoRequestBody",
		Properties: map[string]*spec.Ref[spec.Schema]{
			"title": {Value: &spec.Schema{
				Name:   "title",
				GoType: "string",
			}},
			"description": {Value: &spec.Schema{
				Name:     "description",
				GoType:   "string",
				Nullable: true,
			}},
		},
	}

	// Create schema for update request body with optional fields
	updateRequestBodySchema := &spec.Schema{
		Name:   "UpdateTodoRequestBody",
		GoType: "UpdateTodoRequestBody",
		Properties: map[string]*spec.Ref[spec.Schema]{
			"title": {Value: &spec.Schema{
				Name:     "title",
				GoType:   "string",
				Nullable: true,
			}},
			"completed": {Value: &spec.Schema{
				Name:     "completed",
				GoType:   "bool",
				Nullable: true,
			}},
		},
	}

	ctx := &GeneratorContext{
		Spec: &spec.Spec{
			Operations: []spec.Operation{
				{ID: "ListTodos", Method: "GET"},
				{
					ID:     "CreateTodo",
					Method: "POST",
					RequestBody: &spec.RequestBody{
						Schema:   createRequestBodySchema,
						Required: true,
					},
				},
				{
					ID:     "UpdateTodo",
					Method: "PATCH",
					RequestBody: &spec.RequestBody{
						Schema:   updateRequestBodySchema,
						Required: false,
					},
				},
			},
		},
		ProjectName: "example/myapp",
	}

	schema := &spec.Schema{
		Name: "Todo",
	}

	data := BuildEntityFrontendData(ctx, schema)

	// Verify create form fields are from create request body
	if len(data.CreateFormFields) != 2 {
		t.Errorf("CreateFormFields length = %d, want 2", len(data.CreateFormFields))
	}

	// Verify update form fields are from update request body
	if len(data.UpdateFormFields) != 2 {
		t.Errorf("UpdateFormFields length = %d, want 2", len(data.UpdateFormFields))
	}

	// Check that form fields have correct types
	createFieldNames := make(map[string]bool)
	for _, f := range data.CreateFormFields {
		createFieldNames[f.Name] = true
	}
	if !createFieldNames["title"] || !createFieldNames["description"] {
		t.Errorf("CreateFormFields missing expected fields, got: %v", createFieldNames)
	}

	updateFieldNames := make(map[string]bool)
	for _, f := range data.UpdateFormFields {
		updateFieldNames[f.Name] = true
	}
	if !updateFieldNames["title"] || !updateFieldNames["completed"] {
		t.Errorf("UpdateFormFields missing expected fields, got: %v", updateFieldNames)
	}
}

func TestBuildEntityFrontendData_NoRequestBodies(t *testing.T) {
	ctx := &GeneratorContext{
		Spec: &spec.Spec{
			Operations: []spec.Operation{
				{ID: "ListTodos", Method: "GET"},
				{ID: "GetTodo", Method: "GET"},
			},
		},
		ProjectName: "example/myapp",
	}

	schema := &spec.Schema{
		Name: "Todo",
	}

	data := BuildEntityFrontendData(ctx, schema)

	// When there are no create/update operations, form fields should be nil/empty
	if len(data.CreateFormFields) != 0 {
		t.Errorf("CreateFormFields should be empty when no create operation, got %d", len(data.CreateFormFields))
	}
	if len(data.UpdateFormFields) != 0 {
		t.Errorf("UpdateFormFields should be empty when no update operation, got %d", len(data.UpdateFormFields))
	}
}
