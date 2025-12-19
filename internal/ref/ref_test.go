package ref

import (
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestNewRef(t *testing.T) {
	ref := NewRef[string]("#/components/schemas/User")

	if ref.RefPath != "#/components/schemas/User" {
		t.Errorf("RefPath = %q, want %q", ref.RefPath, "#/components/schemas/User")
	}
	if ref.Value != nil {
		t.Error("Value should be nil for unresolved ref")
	}
	if ref.state != StateUnresolved {
		t.Errorf("state = %v, want StateUnresolved", ref.state)
	}
}

func TestNewInline(t *testing.T) {
	value := "inline value"
	ref := NewInline(&value)

	if ref.RefPath != "" {
		t.Errorf("RefPath = %q, want empty", ref.RefPath)
	}
	if ref.Value == nil {
		t.Error("Value should not be nil for inline ref")
	}
	if *ref.Value != "inline value" {
		t.Errorf("Value = %q, want %q", *ref.Value, "inline value")
	}
	if ref.state != StateInline {
		t.Errorf("state = %v, want StateInline", ref.state)
	}
}

func TestNewResolved(t *testing.T) {
	value := "resolved value"
	ref := NewResolved("#/components/schemas/User", &value)

	if ref.RefPath != "#/components/schemas/User" {
		t.Errorf("RefPath = %q, want %q", ref.RefPath, "#/components/schemas/User")
	}
	if ref.Value == nil {
		t.Error("Value should not be nil for resolved ref")
	}
	if *ref.Value != "resolved value" {
		t.Errorf("Value = %q, want %q", *ref.Value, "resolved value")
	}
	if ref.state != StateResolved {
		t.Errorf("state = %v, want StateResolved", ref.state)
	}
}

func TestRef_IsRef(t *testing.T) {
	tests := []struct {
		name string
		ref  *Ref[string]
		want bool
	}{
		{
			name: "nil ref",
			ref:  nil,
			want: false,
		},
		{
			name: "unresolved ref",
			ref:  NewRef[string]("#/components/schemas/User"),
			want: true,
		},
		{
			name: "resolved ref",
			ref:  NewResolved("#/components/schemas/User", ptr("value")),
			want: true,
		},
		{
			name: "inline value",
			ref:  NewInline(ptr("inline")),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.IsRef(); got != tt.want {
				t.Errorf("IsRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRef_IsInline(t *testing.T) {
	tests := []struct {
		name string
		ref  *Ref[string]
		want bool
	}{
		{
			name: "nil ref",
			ref:  nil,
			want: false,
		},
		{
			name: "unresolved ref",
			ref:  NewRef[string]("#/components/schemas/User"),
			want: false,
		},
		{
			name: "resolved ref",
			ref:  NewResolved("#/components/schemas/User", ptr("value")),
			want: false,
		},
		{
			name: "inline value",
			ref:  NewInline(ptr("inline")),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.IsInline(); got != tt.want {
				t.Errorf("IsInline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRef_IsResolved(t *testing.T) {
	tests := []struct {
		name string
		ref  *Ref[string]
		want bool
	}{
		{
			name: "nil ref",
			ref:  nil,
			want: false,
		},
		{
			name: "unresolved ref",
			ref:  NewRef[string]("#/components/schemas/User"),
			want: false,
		},
		{
			name: "resolved ref",
			ref:  NewResolved("#/components/schemas/User", ptr("value")),
			want: true,
		},
		{
			name: "inline value",
			ref:  NewInline(ptr("inline")),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.IsResolved(); got != tt.want {
				t.Errorf("IsResolved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRef_Get(t *testing.T) {
	t.Run("returns value when resolved", func(t *testing.T) {
		ref := NewResolved("#/ref", ptr("value"))
		got := ref.Get()
		if got == nil || *got != "value" {
			t.Errorf("Get() = %v, want pointer to 'value'", got)
		}
	})

	t.Run("returns value when inline", func(t *testing.T) {
		ref := NewInline(ptr("inline"))
		got := ref.Get()
		if got == nil || *got != "inline" {
			t.Errorf("Get() = %v, want pointer to 'inline'", got)
		}
	})

	t.Run("panics on nil ref", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Get() should panic on nil ref")
			}
		}()
		var ref *Ref[string]
		ref.Get()
	})

	t.Run("panics on unresolved ref", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Get() should panic on unresolved ref")
			}
		}()
		ref := NewRef[string]("#/ref")
		ref.Get()
	})
}

func TestRef_GetOrNil(t *testing.T) {
	tests := []struct {
		name string
		ref  *Ref[string]
		want *string
	}{
		{
			name: "nil ref",
			ref:  nil,
			want: nil,
		},
		{
			name: "unresolved ref",
			ref:  NewRef[string]("#/ref"),
			want: nil,
		},
		{
			name: "resolved ref",
			ref:  NewResolved("#/ref", ptr("value")),
			want: ptr("value"),
		},
		{
			name: "inline value",
			ref:  NewInline(ptr("inline")),
			want: ptr("inline"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ref.GetOrNil()
			if tt.want == nil {
				if got != nil {
					t.Errorf("GetOrNil() = %v, want nil", got)
				}
			} else {
				if got == nil || *got != *tt.want {
					t.Errorf("GetOrNil() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRef_Resolve(t *testing.T) {
	t.Run("resolves unresolved ref", func(t *testing.T) {
		ref := NewRef[string]("#/ref")
		value := "resolved"
		ref.Resolve(&value)

		if ref.Value == nil || *ref.Value != "resolved" {
			t.Errorf("Value = %v, want pointer to 'resolved'", ref.Value)
		}
		if ref.state != StateResolved {
			t.Errorf("state = %v, want StateResolved", ref.state)
		}
	})

	t.Run("updates inline value", func(t *testing.T) {
		ref := NewInline(ptr("old"))
		value := "new"
		ref.Resolve(&value)

		if ref.Value == nil || *ref.Value != "new" {
			t.Errorf("Value = %v, want pointer to 'new'", ref.Value)
		}
		if ref.state != StateInline {
			t.Errorf("state = %v, want StateInline", ref.state)
		}
	})

	t.Run("no-op on nil ref", func(_ *testing.T) {
		var ref *Ref[string]
		value := "test"
		ref.Resolve(&value) // should not panic
	})
}

func TestRef_State(t *testing.T) {
	tests := []struct {
		name string
		ref  *Ref[string]
		want State
	}{
		{
			name: "nil ref",
			ref:  nil,
			want: StateUnresolved,
		},
		{
			name: "unresolved ref",
			ref:  NewRef[string]("#/ref"),
			want: StateUnresolved,
		},
		{
			name: "resolved ref",
			ref:  NewResolved("#/ref", ptr("value")),
			want: StateResolved,
		},
		{
			name: "inline value",
			ref:  NewInline(ptr("inline")),
			want: StateInline,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.State(); got != tt.want {
				t.Errorf("State() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRef_UnmarshalYAML(t *testing.T) {
	type TestSchema struct {
		Title string `yaml:"title"`
		Type  string `yaml:"type"`
	}

	t.Run("parses $ref", func(t *testing.T) {
		yamlData := `$ref: "#/components/schemas/User"`
		var ref Ref[TestSchema]
		if err := yaml.Unmarshal([]byte(yamlData), &ref); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if ref.RefPath != "#/components/schemas/User" {
			t.Errorf("RefPath = %q, want %q", ref.RefPath, "#/components/schemas/User")
		}
		if ref.state != StateUnresolved {
			t.Errorf("state = %v, want StateUnresolved", ref.state)
		}
	})

	t.Run("parses inline mapping", func(t *testing.T) {
		yamlData := `
title: User
type: object
`
		var ref Ref[TestSchema]
		if err := yaml.Unmarshal([]byte(yamlData), &ref); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if ref.RefPath != "" {
			t.Errorf("RefPath = %q, want empty", ref.RefPath)
		}
		if ref.Value == nil {
			t.Fatal("Value should not be nil")
		}
		if ref.Value.Title != "User" {
			t.Errorf("Value.Title = %q, want %q", ref.Value.Title, "User")
		}
		if ref.Value.Type != "object" {
			t.Errorf("Value.Type = %q, want %q", ref.Value.Type, "object")
		}
		if ref.state != StateInline {
			t.Errorf("state = %v, want StateInline", ref.state)
		}
	})

	t.Run("parses inline scalar", func(t *testing.T) {
		yamlData := `hello`
		var ref Ref[string]
		if err := yaml.Unmarshal([]byte(yamlData), &ref); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if ref.RefPath != "" {
			t.Errorf("RefPath = %q, want empty", ref.RefPath)
		}
		if ref.Value == nil || *ref.Value != "hello" {
			t.Errorf("Value = %v, want pointer to 'hello'", ref.Value)
		}
		if ref.state != StateInline {
			t.Errorf("state = %v, want StateInline", ref.state)
		}
	})
}

func TestRef_MarshalYAML(t *testing.T) {
	type TestSchema struct {
		Title string `yaml:"title"`
	}

	t.Run("marshals $ref", func(t *testing.T) {
		ref := NewRef[TestSchema]("#/components/schemas/User")
		data, err := yaml.Marshal(&ref)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "$ref: '#/components/schemas/User'\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})

	t.Run("marshals resolved ref as $ref", func(t *testing.T) {
		value := TestSchema{Title: "User"}
		ref := NewResolved("#/components/schemas/User", &value)
		data, err := yaml.Marshal(&ref)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "$ref: '#/components/schemas/User'\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})

	t.Run("marshals inline value", func(t *testing.T) {
		value := TestSchema{Title: "User"}
		ref := NewInline(&value)
		data, err := yaml.Marshal(&ref)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "title: User\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})

	t.Run("marshals nil ref", func(t *testing.T) {
		var ref *Ref[TestSchema]
		data, err := yaml.Marshal(ref)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "null\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})
}

// ptr is a helper to create a pointer to a value.
func ptr[T any](v T) *T {
	return &v
}
