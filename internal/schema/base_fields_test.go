package schema

import (
	"testing"

	"github.com/archesai/archesai/internal/ref"
)

func TestBaseFieldDefs(t *testing.T) {
	defs := BaseFieldDefs()

	t.Run("contains ID field", func(t *testing.T) {
		id, ok := defs["ID"]
		if !ok {
			t.Fatal("ID field not found")
		}
		if id.Title != "ID" {
			t.Errorf("ID.Title = %q, want %q", id.Title, "ID")
		}
		if id.Format != FormatUUID {
			t.Errorf("ID.Format = %q, want %q", id.Format, FormatUUID)
		}
		if id.GoType != GoTypeUUID {
			t.Errorf("ID.GoType = %q, want %q", id.GoType, GoTypeUUID)
		}
		if id.JSONTag != "id" {
			t.Errorf("ID.JSONTag = %q, want %q", id.JSONTag, "id")
		}
	})

	t.Run("contains CreatedAt field", func(t *testing.T) {
		createdAt, ok := defs["CreatedAt"]
		if !ok {
			t.Fatal("CreatedAt field not found")
		}
		if createdAt.Format != FormatDateTime {
			t.Errorf("CreatedAt.Format = %q, want %q", createdAt.Format, FormatDateTime)
		}
		if createdAt.GoType != GoTypeTime {
			t.Errorf("CreatedAt.GoType = %q, want %q", createdAt.GoType, GoTypeTime)
		}
		if createdAt.JSONTag != "createdAt" {
			t.Errorf("CreatedAt.JSONTag = %q, want %q", createdAt.JSONTag, "createdAt")
		}
	})

	t.Run("contains UpdatedAt field", func(t *testing.T) {
		updatedAt, ok := defs["UpdatedAt"]
		if !ok {
			t.Fatal("UpdatedAt field not found")
		}
		if updatedAt.Format != FormatDateTime {
			t.Errorf("UpdatedAt.Format = %q, want %q", updatedAt.Format, FormatDateTime)
		}
		if updatedAt.GoType != GoTypeTime {
			t.Errorf("UpdatedAt.GoType = %q, want %q", updatedAt.GoType, GoTypeTime)
		}
		if updatedAt.JSONTag != "updatedAt" {
			t.Errorf("UpdatedAt.JSONTag = %q, want %q", updatedAt.JSONTag, "updatedAt")
		}
	})

	t.Run("has exactly 3 fields", func(t *testing.T) {
		if len(defs) != 3 {
			t.Errorf("len(BaseFieldDefs()) = %d, want 3", len(defs))
		}
	})
}

func TestAddBaseFields(t *testing.T) {
	t.Run("adds fields to empty schema", func(t *testing.T) {
		schema := &Schema{}
		AddBaseFields(schema)

		// Check properties were added
		if schema.Properties == nil {
			t.Fatal("Properties should not be nil")
		}
		if len(schema.Properties) != 3 {
			t.Errorf("len(Properties) = %d, want 3", len(schema.Properties))
		}

		// Check ID was added
		if _, ok := schema.Properties["ID"]; !ok {
			t.Error("ID property not added")
		}
		if _, ok := schema.Properties["CreatedAt"]; !ok {
			t.Error("CreatedAt property not added")
		}
		if _, ok := schema.Properties["UpdatedAt"]; !ok {
			t.Error("UpdatedAt property not added")
		}
	})

	t.Run("adds required fields", func(t *testing.T) {
		schema := &Schema{}
		AddBaseFields(schema)

		requiredMap := make(map[string]bool)
		for _, r := range schema.Required {
			requiredMap[r] = true
		}

		if !requiredMap["id"] {
			t.Error("id not in required list")
		}
		if !requiredMap["createdAt"] {
			t.Error("createdAt not in required list")
		}
		if !requiredMap["updatedAt"] {
			t.Error("updatedAt not in required list")
		}
	})

	t.Run("does not duplicate existing properties", func(t *testing.T) {
		existingID := &Schema{Title: "ID", Description: "Custom ID"}
		schema := &Schema{
			Properties: map[string]*ref.Ref[Schema]{
				"ID": ref.NewInline(existingID),
			},
		}
		AddBaseFields(schema)

		// Should still have 3 properties
		if len(schema.Properties) != 3 {
			t.Errorf("len(Properties) = %d, want 3", len(schema.Properties))
		}

		// Existing ID should keep its description
		id := schema.Properties["ID"].GetOrNil()
		if id == nil {
			t.Fatal("ID property is nil")
		}
		if id.Description != "Custom ID" {
			t.Errorf("ID.Description = %q, want %q", id.Description, "Custom ID")
		}
	})

	t.Run("adds description to existing field without one", func(t *testing.T) {
		existingID := &Schema{Title: "ID"}
		schema := &Schema{
			Properties: map[string]*ref.Ref[Schema]{
				"ID": ref.NewInline(existingID),
			},
		}
		AddBaseFields(schema)

		id := schema.Properties["ID"].GetOrNil()
		if id.Description == "" {
			t.Error("Description should be added to existing field without one")
		}
	})

	t.Run("does not duplicate required fields", func(t *testing.T) {
		schema := &Schema{
			Required: []string{"id", "name"},
		}
		AddBaseFields(schema)

		// Count occurrences of "id"
		count := 0
		for _, r := range schema.Required {
			if r == "id" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("id appears %d times in required, want 1", count)
		}
	})
}

func TestGetSchemasAlias(t *testing.T) {
	tests := []struct {
		include string
		want    string
	}{
		{"server", "serverschemas"},
		{"auth", "authschemas"},
		{"config", "configschemas"},
		{"pipelines", "pipelinesschemas"},
		{"executor", "executorschemas"},
		{"storage", "storageschemas"},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.include, func(t *testing.T) {
			got := GetSchemasAlias(tt.include)
			if got != tt.want {
				t.Errorf("GetSchemasAlias(%q) = %q, want %q", tt.include, got, tt.want)
			}
		})
	}
}
