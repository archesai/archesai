package spec

import "gopkg.in/yaml.v3"

// Tag represents a tag in the API specification.
type Tag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// SecScheme represents a security scheme.
type SecScheme struct {
	Type        string `yaml:"type"`
	Scheme      string `yaml:"scheme"`
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Description string `yaml:"description"`
}

// Security represents a security requirement.
type Security struct {
	Name   string   // Security scheme name
	Type   string   // Security type (http, apiKey, oauth2)
	Scheme string   // Security scheme (bearer for http, cookie for apiKey)
	Scopes []string // Required scopes
}

// PropertyType represents a schema type that can be a single type or multiple types (for nullable).
// In OpenAPI 3.1, type can be a string ("string") or array (["string", "null"]).
type PropertyType struct {
	Types    []string
	Nullable bool
}

// PrimaryType returns the primary (non-null) type.
func (p *PropertyType) PrimaryType() string {
	if len(p.Types) > 0 {
		return p.Types[0]
	}
	return ""
}

// UnmarshalYAML implements custom unmarshaling for PropertyType.
// Handles both single string and array of strings for type field.
func (p *PropertyType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try as string first
	var single string
	if err := unmarshal(&single); err == nil {
		p.Types = []string{single}
		p.Nullable = false
		return nil
	}

	// Try as array of strings
	var arr []string
	if err := unmarshal(&arr); err != nil {
		return err
	}

	p.Types = make([]string, 0, len(arr))
	for _, t := range arr {
		if t == "null" {
			p.Nullable = true
		} else {
			p.Types = append(p.Types, t)
		}
	}

	return nil
}

// MarshalYAML implements custom marshaling for PropertyType.
// Outputs single string if not nullable, array if nullable.
func (p PropertyType) MarshalYAML() (interface{}, error) {
	if p.Nullable {
		types := make([]string, 0, len(p.Types)+1)
		types = append(types, p.Types...)
		types = append(types, "null")
		return types, nil
	}
	if len(p.Types) == 1 {
		return p.Types[0], nil
	}
	return p.Types, nil
}

// SchemaType represents the x-codegen-schema-type value.
type SchemaType string

// Valid SchemaType values.
const (
	SchemaTypeEntity      SchemaType = "entity"
	SchemaTypeValueObject SchemaType = "valueobject"
)

// UnmarshalYAML implements custom unmarshaling for SchemaType.
func (s *SchemaType) UnmarshalYAML(node *yaml.Node) error {
	var str string
	if err := node.Decode(&str); err != nil {
		return err
	}
	*s = SchemaType(str)
	return nil
}

// MarshalYAML implements custom marshaling for SchemaType.
func (s SchemaType) MarshalYAML() (interface{}, error) {
	if s == "" {
		return nil, nil
	}
	return string(s), nil
}
