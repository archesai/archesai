package schema

import (
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestPropertyType_PrimaryType(t *testing.T) {
	tests := []struct {
		name  string
		types []string
		want  string
	}{
		{
			name:  "single type",
			types: []string{"string"},
			want:  "string",
		},
		{
			name:  "multiple types",
			types: []string{"string", "integer"},
			want:  "string",
		},
		{
			name:  "empty types",
			types: []string{},
			want:  "",
		},
		{
			name:  "nil types",
			types: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PropertyType{Types: tt.types}
			if got := p.PrimaryType(); got != tt.want {
				t.Errorf("PrimaryType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPropertyType_UnmarshalYAML(t *testing.T) {
	t.Run("single string type", func(t *testing.T) {
		yamlData := `string`
		var p PropertyType
		if err := yaml.Unmarshal([]byte(yamlData), &p); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if len(p.Types) != 1 || p.Types[0] != "string" {
			t.Errorf("Types = %v, want [string]", p.Types)
		}
		if p.Nullable {
			t.Error("Nullable should be false for single string")
		}
	})

	t.Run("array without null", func(t *testing.T) {
		yamlData := `[string, integer]`
		var p PropertyType
		if err := yaml.Unmarshal([]byte(yamlData), &p); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if len(p.Types) != 2 || p.Types[0] != "string" || p.Types[1] != "integer" {
			t.Errorf("Types = %v, want [string, integer]", p.Types)
		}
		if p.Nullable {
			t.Error("Nullable should be false")
		}
	})

	t.Run("array with null", func(t *testing.T) {
		yamlData := `[string, "null"]`
		var p PropertyType
		if err := yaml.Unmarshal([]byte(yamlData), &p); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if len(p.Types) != 1 || p.Types[0] != "string" {
			t.Errorf("Types = %v, want [string]", p.Types)
		}
		if !p.Nullable {
			t.Error("Nullable should be true")
		}
	})

	t.Run("null at beginning", func(t *testing.T) {
		yamlData := `["null", string]`
		var p PropertyType
		if err := yaml.Unmarshal([]byte(yamlData), &p); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if len(p.Types) != 1 || p.Types[0] != "string" {
			t.Errorf("Types = %v, want [string]", p.Types)
		}
		if !p.Nullable {
			t.Error("Nullable should be true")
		}
	})
}

func TestPropertyType_MarshalYAML(t *testing.T) {
	t.Run("single non-nullable type", func(t *testing.T) {
		p := PropertyType{Types: []string{"string"}, Nullable: false}
		data, err := yaml.Marshal(&p)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "string\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})

	t.Run("single nullable type", func(t *testing.T) {
		p := PropertyType{Types: []string{"string"}, Nullable: true}
		data, err := yaml.Marshal(&p)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "- string\n- \"null\"\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})

	t.Run("multiple types non-nullable", func(t *testing.T) {
		p := PropertyType{Types: []string{"string", "integer"}, Nullable: false}
		data, err := yaml.Marshal(&p)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "- string\n- integer\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})

	t.Run("multiple types nullable", func(t *testing.T) {
		p := PropertyType{Types: []string{"string", "integer"}, Nullable: true}
		data, err := yaml.Marshal(&p)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		expected := "- string\n- integer\n- \"null\"\n"
		if string(data) != expected {
			t.Errorf("Marshal = %q, want %q", string(data), expected)
		}
	})
}

func TestType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want Type
	}{
		{"entity", "entity", TypeEntity},
		{"valueobject", "valueobject", TypeValueObject},
		{"custom value", "response", Type("response")},
		{"empty string", `""`, Type("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Type
			if err := yaml.Unmarshal([]byte(tt.yaml), &s); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if s != tt.want {
				t.Errorf("Type = %q, want %q", s, tt.want)
			}
		})
	}
}

func TestType_MarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		s    Type
		want string
	}{
		{"entity", TypeEntity, "entity\n"},
		{"valueobject", TypeValueObject, "valueobject\n"},
		{"empty returns null", Type(""), "null\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(&tt.s)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("Marshal = %q, want %q", string(data), tt.want)
			}
		})
	}
}

func TestTypeConstants(t *testing.T) {
	// Verify type constants are correct
	tests := []struct {
		constant string
		expected string
	}{
		{TypeString, "string"},
		{TypeInteger, "integer"},
		{TypeNumber, "number"},
		{TypeBoolean, "boolean"},
		{TypeArray, "array"},
		{TypeObject, "object"},
		{TypeNull, "null"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("Type constant = %q, want %q", tt.constant, tt.expected)
		}
	}
}

func TestFormatConstants(t *testing.T) {
	// Verify format constants are correct
	tests := []struct {
		constant string
		expected string
	}{
		{FormatDateTime, "date-time"},
		{FormatDate, "date"},
		{FormatUUID, "uuid"},
		{FormatEmail, "email"},
		{FormatURI, "uri"},
		{FormatHostname, "hostname"},
		{FormatPassword, "password"},
		{FormatInt32, "int32"},
		{FormatInt64, "int64"},
		{FormatFloat, "float"},
		{FormatDouble, "double"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("Format constant = %q, want %q", tt.constant, tt.expected)
		}
	}
}
