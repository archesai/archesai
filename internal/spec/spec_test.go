package spec

import (
	"testing"

	"github.com/archesai/archesai/internal/schema"
)

func TestSpec_GetOperationsByTag(t *testing.T) {
	tests := []struct {
		name         string
		spec         *Spec
		wantTagCount int
		wantTagNames []string
		wantOpCounts map[string]int
	}{
		{
			name:         "nil spec returns nil",
			spec:         nil,
			wantTagCount: 0,
		},
		{
			name: "empty operations returns nil",
			spec: &Spec{
				Operations: []Operation{},
			},
			wantTagCount: 0,
		},
		{
			name: "single tag with single operation",
			spec: &Spec{
				Operations: []Operation{
					{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
				},
				Schemas: map[string]*schema.Schema{},
			},
			wantTagCount: 1,
			wantTagNames: []string{"Pipeline"},
			wantOpCounts: map[string]int{"Pipeline": 1},
		},
		{
			name: "single tag with multiple operations",
			spec: &Spec{
				Operations: []Operation{
					{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
					{ID: "GetPipeline", Tag: "Pipeline", Method: "GET", Path: "/pipelines/{id}"},
					{
						ID:     "DeletePipeline",
						Tag:    "Pipeline",
						Method: "DELETE",
						Path:   "/pipelines/{id}",
					},
				},
				Schemas: map[string]*schema.Schema{},
			},
			wantTagCount: 1,
			wantTagNames: []string{"Pipeline"},
			wantOpCounts: map[string]int{"Pipeline": 3},
		},
		{
			name: "multiple tags",
			spec: &Spec{
				Operations: []Operation{
					{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
					{ID: "GetPipeline", Tag: "Pipeline", Method: "GET", Path: "/pipelines/{id}"},
					{ID: "CreateUser", Tag: "User", Method: "POST", Path: "/users"},
					{ID: "GetUser", Tag: "User", Method: "GET", Path: "/users/{id}"},
					{ID: "DeleteUser", Tag: "User", Method: "DELETE", Path: "/users/{id}"},
				},
				Schemas: map[string]*schema.Schema{},
			},
			wantTagCount: 2,
			wantTagNames: []string{"Pipeline", "User"},
			wantOpCounts: map[string]int{"Pipeline": 2, "User": 3},
		},
		{
			name: "operations without tag go to Default",
			spec: &Spec{
				Operations: []Operation{
					{ID: "GetHealth", Tag: "", Method: "GET", Path: "/health"},
					{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
				},
				Schemas: map[string]*schema.Schema{},
			},
			wantTagCount: 2,
			wantTagNames: []string{"Default", "Pipeline"},
			wantOpCounts: map[string]int{"Default": 1, "Pipeline": 1},
		},
		{
			name: "operations sorted by path then method within tag",
			spec: &Spec{
				Operations: []Operation{
					{
						ID:     "DeletePipeline",
						Tag:    "Pipeline",
						Method: "DELETE",
						Path:   "/pipelines/{id}",
					},
					{ID: "GetPipeline", Tag: "Pipeline", Method: "GET", Path: "/pipelines/{id}"},
					{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
					{ID: "ListPipelines", Tag: "Pipeline", Method: "GET", Path: "/pipelines"},
				},
				Schemas: map[string]*schema.Schema{},
			},
			wantTagCount: 1,
			wantTagNames: []string{"Pipeline"},
			wantOpCounts: map[string]int{"Pipeline": 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.GetOperationsByTag()

			if len(got) != tt.wantTagCount {
				t.Errorf("GetOperationsByTag() got %d tags, want %d", len(got), tt.wantTagCount)
				return
			}

			if tt.wantTagCount == 0 {
				return
			}

			// Check tag names
			gotNames := make([]string, len(got))
			for i, g := range got {
				gotNames[i] = g.Name
			}
			for i, wantName := range tt.wantTagNames {
				if gotNames[i] != wantName {
					t.Errorf(
						"GetOperationsByTag() tag[%d].Name = %q, want %q",
						i,
						gotNames[i],
						wantName,
					)
				}
			}

			// Check operation counts
			for _, g := range got {
				if wantCount, ok := tt.wantOpCounts[g.Name]; ok {
					if len(g.Operations) != wantCount {
						t.Errorf("GetOperationsByTag() tag[%s] has %d operations, want %d",
							g.Name, len(g.Operations), wantCount)
					}
				}
			}
		})
	}
}

func TestTagGroup_Package(t *testing.T) {
	spec := &Spec{
		Operations: []Operation{
			{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
			{ID: "CreateUser", Tag: "User", Method: "POST", Path: "/users"},
			{ID: "CreateAPIKey", Tag: "APIKey", Method: "POST", Path: "/apikeys"},
		},
		Schemas: map[string]*schema.Schema{},
	}

	groups := spec.GetOperationsByTag()

	wantPackages := map[string]string{
		"APIKey":   "apikey",
		"Pipeline": "pipeline",
		"User":     "user",
	}

	for _, g := range groups {
		if want, ok := wantPackages[g.Name]; ok {
			if g.Package != want {
				t.Errorf(
					"GetOperationsByTag() tag[%s].Package = %q, want %q",
					g.Name,
					g.Package,
					want,
				)
			}
		}
	}
}

func TestTagGroup_Entities(t *testing.T) {
	pipelineSchema := &schema.Schema{
		Title:              "Pipeline",
		XCodegenSchemaType: schema.TypeEntity,
	}
	userSchema := &schema.Schema{
		Title:              "User",
		XCodegenSchemaType: schema.TypeEntity,
	}
	configSchema := &schema.Schema{
		Title:              "Config",
		XCodegenSchemaType: "", // Not an entity
	}

	spec := &Spec{
		Operations: []Operation{
			{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
			{ID: "CreateUser", Tag: "User", Method: "POST", Path: "/users"},
		},
		Schemas: map[string]*schema.Schema{
			"Pipeline": pipelineSchema,
			"User":     userSchema,
			"Config":   configSchema,
		},
	}

	groups := spec.GetOperationsByTag()

	for _, g := range groups {
		switch g.Name {
		case "Pipeline":
			if len(g.Entities) != 1 {
				t.Errorf("Pipeline tag should have 1 entity, got %d", len(g.Entities))
			} else if g.Entities[0].Title != "Pipeline" {
				t.Errorf("Pipeline entity title = %q, want %q", g.Entities[0].Title, "Pipeline")
			}
		case "User":
			if len(g.Entities) != 1 {
				t.Errorf("User tag should have 1 entity, got %d", len(g.Entities))
			} else if g.Entities[0].Title != "User" {
				t.Errorf("User entity title = %q, want %q", g.Entities[0].Title, "User")
			}
		}
	}
}

func TestTagGroup_OperationsSortedWithinTag(t *testing.T) {
	spec := &Spec{
		Operations: []Operation{
			{ID: "DeletePipeline", Tag: "Pipeline", Method: "DELETE", Path: "/pipelines/{id}"},
			{ID: "UpdatePipeline", Tag: "Pipeline", Method: "PATCH", Path: "/pipelines/{id}"},
			{ID: "GetPipeline", Tag: "Pipeline", Method: "GET", Path: "/pipelines/{id}"},
			{ID: "CreatePipeline", Tag: "Pipeline", Method: "POST", Path: "/pipelines"},
			{ID: "ListPipelines", Tag: "Pipeline", Method: "GET", Path: "/pipelines"},
		},
		Schemas: map[string]*schema.Schema{},
	}

	groups := spec.GetOperationsByTag()

	if len(groups) != 1 {
		t.Fatalf("Expected 1 tag group, got %d", len(groups))
	}

	ops := groups[0].Operations
	wantOrder := []string{
		"ListPipelines",
		"CreatePipeline",
		"GetPipeline",
		"UpdatePipeline",
		"DeletePipeline",
	}

	if len(ops) != len(wantOrder) {
		t.Fatalf("Expected %d operations, got %d", len(wantOrder), len(ops))
	}

	for i, want := range wantOrder {
		if ops[i].ID != want {
			t.Errorf("Operation[%d].ID = %q, want %q", i, ops[i].ID, want)
		}
	}
}
