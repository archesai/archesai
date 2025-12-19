package schema

import "go.yaml.in/yaml/v4"

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
func (p *PropertyType) UnmarshalYAML(unmarshal func(any) error) error {
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
func (p PropertyType) MarshalYAML() (any, error) {
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

// Type represents the x-codegen-schema-type value.
type Type string

// Valid Type values.
const (
	TypeEntity      Type = "entity"
	TypeValueObject Type = "valueobject"
)

// UnmarshalYAML implements custom unmarshaling for Type.
func (t *Type) UnmarshalYAML(node *yaml.Node) error {
	var str string
	if err := node.Decode(&str); err != nil {
		return err
	}
	*t = Type(str)
	return nil
}

// MarshalYAML implements custom marshaling for Type.
func (t Type) MarshalYAML() (any, error) {
	if t == "" {
		return nil, nil
	}
	return string(t), nil
}

// OpenAPI schema type constants.
const (
	TypeString  = "string"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeBoolean = "boolean"
	TypeArray   = "array"
	TypeObject  = "object"
	TypeNull    = "null"
)

// OpenAPI format constants.
const (
	FormatDateTime = "date-time"
	FormatDate     = "date"
	FormatUUID     = "uuid"
	FormatEmail    = "email"
	FormatURI      = "uri"
	FormatHostname = "hostname"
	FormatPassword = "password"
	FormatInt32    = "int32"
	FormatInt64    = "int64"
	FormatFloat    = "float"
	FormatDouble   = "double"
)
