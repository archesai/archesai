package spec

import (
	"testing"

	"github.com/archesai/archesai/internal/ref"
	"github.com/archesai/archesai/internal/schema"
)

func TestFindEntityOperations(t *testing.T) {
	s := &Spec{
		Operations: []Operation{
			{ID: "ListUsers", Method: "GET"},
			{ID: "GetUser", Method: "GET"},
			{ID: "CreateUser", Method: "POST"},
			{ID: "UpdateUser", Method: "PUT"},
			{ID: "DeleteUser", Method: "DELETE"},
			{ID: "GetProducts", Method: "GET"}, // Different entity
		},
	}

	ops := s.FindEntityOperations("User")

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
	s := &Spec{
		Operations: []Operation{
			{ID: "ListTodos", Method: "GET"},
			{ID: "CreateTodo", Method: "POST"},
		},
		ProjectName: "example/myapp",
	}

	// Test with a simple schema (Properties are complex in the real spec)
	sch := &schema.Schema{
		Title: "Todo",
	}

	data := s.BuildEntityFrontendData(sch)

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
	createRequestBodySchema := &schema.Schema{
		Title:  "CreateTodoRequestBody",
		GoType: "CreateTodoRequestBody",
		Properties: map[string]*ref.Ref[schema.Schema]{
			"title": {Value: &schema.Schema{
				Title:  "title",
				GoType: "string",
			}},
			"description": {Value: &schema.Schema{
				Title:    "description",
				GoType:   "string",
				Nullable: true,
			}},
		},
	}

	// Create schema for update request body with optional fields
	updateRequestBodySchema := &schema.Schema{
		Title:  "UpdateTodoRequestBody",
		GoType: "UpdateTodoRequestBody",
		Properties: map[string]*ref.Ref[schema.Schema]{
			"title": {Value: &schema.Schema{
				Title:    "title",
				GoType:   "string",
				Nullable: true,
			}},
			"completed": {Value: &schema.Schema{
				Title:    "completed",
				GoType:   "bool",
				Nullable: true,
			}},
		},
	}

	s := &Spec{
		Operations: []Operation{
			{ID: "ListTodos", Method: "GET"},
			{
				ID:     "CreateTodo",
				Method: "POST",
				RequestBody: &RequestBody{
					Schema:   createRequestBodySchema,
					Required: true,
				},
			},
			{
				ID:     "UpdateTodo",
				Method: "PATCH",
				RequestBody: &RequestBody{
					Schema:   updateRequestBodySchema,
					Required: false,
				},
			},
		},
		ProjectName: "example/myapp",
	}

	sch := &schema.Schema{
		Title: "Todo",
	}

	data := s.BuildEntityFrontendData(sch)

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
	s := &Spec{
		Operations: []Operation{
			{ID: "ListTodos", Method: "GET"},
			{ID: "GetTodo", Method: "GET"},
		},
		ProjectName: "example/myapp",
	}

	sch := &schema.Schema{
		Title: "Todo",
	}

	data := s.BuildEntityFrontendData(sch)

	// When there are no create/update operations, form fields should be nil/empty
	if len(data.CreateFormFields) != 0 {
		t.Errorf(
			"CreateFormFields should be empty when no create operation, got %d",
			len(data.CreateFormFields),
		)
	}
	if len(data.UpdateFormFields) != 0 {
		t.Errorf(
			"UpdateFormFields should be empty when no update operation, got %d",
			len(data.UpdateFormFields),
		)
	}
}

func TestEntityOperations_IsNested(t *testing.T) {
	tests := []struct {
		name     string
		ops      EntityOperations
		expected bool
	}{
		{
			name: "not nested - no path params",
			ops: EntityOperations{
				HasList: true,
				ListOp:  &Operation{Parameters: []Param{}},
			},
			expected: false,
		},
		{
			name: "nested - list has path params",
			ops: EntityOperations{
				HasList: true,
				ListOp: &Operation{
					Parameters: []Param{
						{Schema: &schema.Schema{Title: "parentID"}, In: "path"},
					},
				},
			},
			expected: true,
		},
		{
			name: "no list operation",
			ops: EntityOperations{
				HasList: false,
				ListOp:  nil,
			},
			expected: false,
		},
		{
			name: "not nested - only query params",
			ops: EntityOperations{
				HasList: true,
				ListOp: &Operation{
					Parameters: []Param{
						{Schema: &schema.Schema{Title: "filter"}, In: "query"},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ops.IsNested(); got != tt.expected {
				t.Errorf("IsNested() = %v, want %v", got, tt.expected)
			}
		})
	}
}
