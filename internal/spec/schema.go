package spec

import (
	"fmt"
	"sort"
	"strings"
)

// OpenAPI schema type constants.
const (
	SchemaTypeString  = "string"
	SchemaTypeInteger = "integer"
	SchemaTypeNumber  = "number"
	SchemaTypeBoolean = "boolean"
	SchemaTypeArray   = "array"
	SchemaTypeObject  = "object"
	SchemaTypeNull    = "null"
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

// Schema represents a unified OpenAPI schema definition.
// It combines parsing fields with code generation computed fields.
type Schema struct {
	// === Identity ===
	Name        string `yaml:"-"`                     // Schema/property name (computed, not in YAML)
	Title       string `yaml:"title,omitempty"`       // OpenAPI title field
	Description string `yaml:"description,omitempty"` // OpenAPI description

	// === Type Information ===
	Type   PropertyType `yaml:"type,omitempty"`   // Type information (handles nullable via ["string", "null"])
	Format string       `yaml:"format,omitempty"` // OpenAPI format (uuid, date-time, etc.)
	Enum   []string     `yaml:"enum,omitempty"`   // Enum values

	// === Structural (using Ref[Schema] for references) ===
	Properties           map[string]*Ref[Schema] `yaml:"properties,omitempty"`           // Object properties
	Required             []string                `yaml:"required,omitempty"`             // Required property names
	Items                *Ref[Schema]            `yaml:"items,omitempty"`                // Array items schema
	AllOf                []*Ref[Schema]          `yaml:"allOf,omitempty"`                // Composition schemas
	OneOf                []*Ref[Schema]          `yaml:"oneOf,omitempty"`                // Union schemas (e.g., nullable types with constraints)
	AdditionalProperties *bool                   `yaml:"additionalProperties,omitempty"` // Allow extra properties

	// === Validation Constraints ===
	MinLength *int   `yaml:"minLength,omitempty"` // String min length
	MaxLength *int   `yaml:"maxLength,omitempty"` // String max length
	Minimum   *int64 `yaml:"minimum,omitempty"`   // Number minimum
	Maximum   *int64 `yaml:"maximum,omitempty"`   // Number maximum
	Pattern   string `yaml:"pattern,omitempty"`   // Regex pattern
	MaxItems  *int   `yaml:"maxItems,omitempty"`  // Array max items
	Default   any    `yaml:"default,omitempty"`   // Default value
	Example   any    `yaml:"example,omitempty"`   // Example value

	// === OpenAPI Extensions ===
	XCodegenSchemaType SchemaType         `yaml:"x-codegen-schema-type,omitempty"` // x-codegen-schema-type (entity/valueobject)
	XCodegen           *XCodegenExtension `yaml:"x-codegen,omitempty"`             // x-codegen extension
	XInternal          string             `yaml:"x-internal,omitempty"`            // x-internal context

	// === Computed Fields (populated during resolution) ===
	GoType   string `yaml:"-"` // Computed Go type (e.g., "string", "uuid.UUID", "models.User")
	JSONTag  string `yaml:"-"` // JSON tag value (e.g., "name" or "name,omitempty")
	YAMLTag  string `yaml:"-"` // YAML tag value
	Nullable bool   `yaml:"-"` // Whether field is nullable (from type array or explicit)

	// === Resolution State ===
	resolved bool `yaml:"-"` // Whether this schema has been fully resolved
}

// IsResolved returns whether the schema has been resolved.
func (s *Schema) IsResolved() bool {
	if s == nil {
		return false
	}
	return s.resolved
}

// SetResolved marks the schema as resolved.
func (s *Schema) SetResolved(resolved bool) {
	if s != nil {
		s.resolved = resolved
	}
}

// GetProperty returns a property by name, resolving the reference if needed.
// Returns nil if the property doesn't exist.
func (s *Schema) GetProperty(name string) *Schema {
	if s == nil || s.Properties == nil {
		return nil
	}
	ref, ok := s.Properties[name]
	if !ok {
		return nil
	}
	return ref.GetOrNil()
}

// HasProperties returns true if the schema has any properties.
func (s *Schema) HasProperties() bool {
	return s != nil && len(s.Properties) > 0
}

// HasProperty returns true if the schema has a property with the given name.
func (s *Schema) HasProperty(name string) bool {
	if s == nil || s.Properties == nil {
		return false
	}
	_, ok := s.Properties[name]
	return ok
}

// GetItems returns the array items schema, resolving the reference if needed.
// Returns nil if Items is not set.
func (s *Schema) GetItems() *Schema {
	if s == nil || s.Items == nil {
		return nil
	}
	return s.Items.GetOrNil()
}

// IsRef returns true if this schema was loaded from a $ref.
func (s *Schema) IsRef() bool {
	// A schema that has a name but no properties/type is likely a ref placeholder
	return s != nil && s.Name != "" && !s.HasProperties() && s.Type.PrimaryType() == ""
}

// IsEnum returns true if the schema is an enum.
func (s *Schema) IsEnum() bool {
	return s != nil && len(s.Enum) > 0
}

// IsArray returns true if the schema is an array type.
func (s *Schema) IsArray() bool {
	return s != nil && s.Type.PrimaryType() == "array"
}

// IsObject returns true if the schema is an object type.
func (s *Schema) IsObject() bool {
	return s != nil && s.Type.PrimaryType() == "object"
}

// Clone creates a shallow copy of the schema.
func (s *Schema) Clone() *Schema {
	if s == nil {
		return nil
	}
	clone := *s
	return &clone
}

// SchemaType returns the x-codegen-schema-type value as a string.
// This is used by templates that check for "entity" or "valueobject".
func (s *Schema) SchemaType() string {
	if s == nil {
		return ""
	}
	return string(s.XCodegenSchemaType)
}

// IsInternal returns true if this schema should be imported from another package instead of generated.
func (s *Schema) IsInternal(context string) bool {
	if s == nil || s.XInternal == "" {
		return false
	}
	if context == "" {
		return true
	}
	return s.XInternal != context
}

// HasDomainEvents returns true if the schema has domain events.
func (s *Schema) HasDomainEvents() bool {
	return s != nil && s.XCodegenSchemaType == SchemaTypeEntity
}

// GetRepositoryIndices returns the repository indices.
func (s *Schema) GetRepositoryIndices() []string {
	if s == nil || s.XCodegen == nil || s.XCodegen.Repository == nil {
		return nil
	}
	return s.XCodegen.Repository.Indices
}

// GetRepositoryRelations returns the repository relations.
func (s *Schema) GetRepositoryRelations() []XCodegenExtensionRepositoryRelationsItem {
	if s == nil || s.XCodegen == nil || s.XCodegen.Repository == nil {
		return nil
	}
	return s.XCodegen.Repository.Relations
}

// GetRequiredProperties returns only the required properties as resolved schemas.
func (s *Schema) GetRequiredProperties() map[string]*Schema {
	if s == nil {
		return nil
	}
	required := make(map[string]*Schema)
	for _, name := range s.Required {
		if propRef, exists := s.Properties[name]; exists {
			if prop := propRef.GetOrNil(); prop != nil {
				required[name] = prop
			}
		}
	}
	return required
}

// GetOptionalProperties returns only the optional properties as resolved schemas.
func (s *Schema) GetOptionalProperties() map[string]*Schema {
	if s == nil {
		return nil
	}
	optional := make(map[string]*Schema)
	requiredMap := make(map[string]bool)
	for _, name := range s.Required {
		requiredMap[name] = true
	}
	for name, propRef := range s.Properties {
		if !requiredMap[name] {
			if prop := propRef.GetOrNil(); prop != nil {
				optional[name] = prop
			}
		}
	}
	return optional
}

// GetSortedProperties returns properties sorted by key for consistent iteration.
// Priority: ID -> CreatedAt -> UpdatedAt -> alphabetical
func (s *Schema) GetSortedProperties() []*Schema {
	if s == nil || s.Properties == nil {
		return nil
	}

	var names []string
	for name := range s.Properties {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		priority := func(name string) int {
			switch strings.ToLower(name) {
			case "id":
				return 0
			case "createdat", "created_at":
				return 1
			case "updatedat", "updated_at":
				return 2
			default:
				return 999
			}
		}
		priI := priority(names[i])
		priJ := priority(names[j])
		if priI == priJ {
			return names[i] < names[j]
		}
		return priI < priJ
	})

	var sorted []*Schema
	for _, name := range names {
		propRef := s.Properties[name]
		prop := propRef.GetOrNil()
		if prop == nil {
			continue
		}
		if prop.Name == "" {
			prop.Name = name
		}
		sorted = append(sorted, prop)
	}
	return sorted
}

// IsPropertyRequired checks if a property is required.
func (s *Schema) IsPropertyRequired(propName string) bool {
	if s == nil {
		return false
	}
	for _, req := range s.Required {
		if req == propName {
			return true
		}
	}
	return false
}

// CollectNestedTypes walks the schema tree and collects all nested object types.
func (s *Schema) CollectNestedTypes() []*Schema {
	if s == nil {
		return nil
	}
	var nestedTypes []*Schema
	s.collectNestedTypesRecursive(&nestedTypes, make(map[string]bool))
	return nestedTypes
}

func (s *Schema) collectNestedTypesRecursive(
	result *[]*Schema,
	visited map[string]bool,
) {
	if s == nil || s.Properties == nil {
		return
	}
	for _, prop := range s.GetSortedProperties() {
		if prop == nil {
			continue
		}
		if prop.Type.PrimaryType() == SchemaTypeObject && prop.Name != "" && prop.GoType != "" &&
			!strings.Contains(prop.GoType, ".") && len(prop.Properties) > 0 {
			if !visited[prop.GoType] {
				visited[prop.GoType] = true
				*result = append(*result, prop)
				prop.collectNestedTypesRecursive(result, visited)
			}
		}
		if prop.Type.PrimaryType() == "array" && prop.Items != nil {
			items := prop.Items.GetOrNil()
			if items != nil && items.Type.PrimaryType() == SchemaTypeObject && items.Name != "" &&
				items.GoType != "" && !strings.Contains(items.GoType, ".") &&
				len(items.Properties) > 0 {
				if !visited[items.GoType] {
					visited[items.GoType] = true
					*result = append(*result, items)
					items.collectNestedTypesRecursive(result, visited)
				}
			} else if items != nil {
				items.collectNestedTypesRecursive(result, visited)
			}
		}
	}
}

// IsOptional returns true if the property is optional (has omitempty tag).
func (s *Schema) IsOptional() bool {
	if s == nil {
		return false
	}
	return strings.Contains(s.JSONTag, ",omitempty")
}

// IsRequired returns true if the property is required (doesn't have omitempty tag).
func (s *Schema) IsRequired() bool {
	return !s.IsOptional()
}

// NeedsPointer returns true if the field type needs to be a pointer.
func (s *Schema) NeedsPointer() bool {
	if s == nil {
		return false
	}
	if s.IsOptional() {
		return !strings.HasPrefix(s.GoType, "*") &&
			!strings.HasPrefix(s.GoType, "[]") &&
			!strings.HasPrefix(s.GoType, "map")
	}
	return s.Nullable
}

// GetFieldType returns the complete Go type for this field, including pointer if needed.
func (s *Schema) GetFieldType(parentSchema *Schema) string {
	if s == nil {
		return ""
	}
	baseType := s.GoType
	if len(s.Enum) > 0 && parentSchema != nil {
		baseType = parentSchema.Name + s.Name
	}
	if s.NeedsPointer() {
		return "*" + baseType
	}
	return baseType
}

// GetBaseType returns the base Go type without pointer or qualification.
func (s *Schema) GetBaseType(parentSchema *Schema) string {
	if s == nil {
		return ""
	}
	if len(s.Enum) > 0 && parentSchema != nil {
		return parentSchema.Name + s.Name
	}
	return s.GoType
}

// NeedsValidation returns true if this field requires validation in constructor.
func (s *Schema) NeedsValidation(isEntity bool, parentSchema *Schema) bool {
	if s == nil {
		return false
	}
	if isEntity {
		if s.IsSpecialField() {
			return false
		}
		if parentSchema == nil || !parentSchema.IsPropertyRequired(s.JSONTag) || s.IsOptional() {
			return false
		}
	} else {
		if s.IsOptional() {
			return false
		}
	}
	if len(s.Enum) > 0 {
		return true
	}
	if s.GoType == "string" && !s.Nullable {
		return true
	}
	if (s.GoType == "int" || s.GoType == "float64") && !s.Nullable {
		return true
	}
	if s.GoType == "uuid.UUID" && !s.Nullable {
		return true
	}
	return false
}

// IsSpecialField returns true if this is a special entity field.
func (s *Schema) IsSpecialField() bool {
	if s == nil {
		return false
	}
	return s.Name == "ID" || s.Name == "CreatedAt" || s.Name == "UpdatedAt"
}

// ShouldIncludeInConstructor returns true if field should be in entity constructor.
func (s *Schema) ShouldIncludeInConstructor(parentSchema *Schema) bool {
	if s == nil || parentSchema == nil {
		return false
	}
	if parentSchema.XCodegenSchemaType == SchemaTypeEntity {
		return parentSchema.IsPropertyRequired(s.JSONTag) &&
			!s.IsOptional() &&
			!s.IsSpecialField()
	}
	return true
}

// GetValidationError returns the validation error message for this field.
func (s *Schema) GetValidationError(_ *Schema) string {
	if s == nil {
		return ""
	}
	if len(s.Enum) > 0 {
		return fmt.Sprintf("invalid %s: %%s", s.Name)
	}
	if s.GoType == "string" {
		return fmt.Sprintf("%s cannot be empty", s.Name)
	}
	if s.GoType == "int" || s.GoType == "float64" {
		return fmt.Sprintf("%s cannot be negative", s.Name)
	}
	if s.GoType == "uuid.UUID" {
		return fmt.Sprintf("%s cannot be nil UUID", s.Name)
	}
	return fmt.Sprintf("invalid %s", s.Name)
}

// GetCreateParamValue returns the Go expression for mapping entity field to DB param in Create.
func (s *Schema) GetCreateParamValue(entityVar string, _ *Schema) string {
	if s == nil {
		return ""
	}
	fieldAccess := fmt.Sprintf("%s.%s", entityVar, s.Name)
	if s.IsOptional() {
		return fieldAccess
	}
	if s.Nullable && len(s.Enum) > 0 {
		return fmt.Sprintf(
			"func() *string { if %s == nil { return nil }; s := string(*%s); return &s }()",
			fieldAccess, fieldAccess,
		)
	}
	if len(s.Enum) > 0 {
		return fmt.Sprintf("string(%s)", fieldAccess)
	}
	if s.Nullable || strings.HasPrefix(s.GoType, "*") || strings.HasPrefix(s.GoType, "[]") {
		return fieldAccess
	}
	return fieldAccess
}

// GetUpdateParamValue returns the Go expression for mapping entity field to DB param in Update.
func (s *Schema) GetUpdateParamValue(entityVar string, _ *Schema) string {
	if s == nil {
		return ""
	}
	fieldAccess := fmt.Sprintf("%s.%s", entityVar, s.Name)
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Name[0])), s.Name[1:]) + "Str"
	if s.IsOptional() {
		return fieldAccess
	}
	if s.Nullable && len(s.Enum) > 0 {
		return fmt.Sprintf(
			"func() *string { if %s == nil { return nil }; s := string(*%s); return &s }()",
			fieldAccess, fieldAccess,
		)
	}
	if len(s.Enum) > 0 {
		return "&" + varName
	}
	if s.Nullable || strings.HasPrefix(s.GoType, "*") || strings.HasPrefix(s.GoType, "[]") {
		return fieldAccess
	}
	return "&" + fieldAccess
}

// NeedsUpdateVarDeclaration returns true if Update operation needs a variable declaration.
func (s *Schema) NeedsUpdateVarDeclaration() bool {
	if s == nil {
		return false
	}
	return !s.IsOptional() && len(s.Enum) > 0 && !s.Nullable
}

// GetUpdateVarDeclaration returns the variable declaration for Update operation.
func (s *Schema) GetUpdateVarDeclaration(entityVar string) string {
	if s == nil || !s.NeedsUpdateVarDeclaration() {
		return ""
	}
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Name[0])), s.Name[1:]) + "Str"
	return fmt.Sprintf("%s := string(%s.%s)", varName, entityVar, s.Name)
}

// GetDBMapValue returns the Go expression for mapping DB field to entity field.
func (s *Schema) GetDBMapValue(dbVar string, parentSchema *Schema) string {
	if s == nil {
		return ""
	}
	fieldAccess := fmt.Sprintf("%s.%s", dbVar, s.Name)
	enumType := ""
	if parentSchema != nil && len(s.Enum) > 0 {
		enumType = fmt.Sprintf("models.%s%s", parentSchema.Name, s.Name)
	}
	if s.Nullable && len(s.Enum) > 0 {
		return fmt.Sprintf(
			"func() *%s { if %s == nil { return nil }; v := %s(*%s); return &v }()",
			enumType, fieldAccess, enumType, fieldAccess,
		)
	}
	if len(s.Enum) > 0 {
		return fmt.Sprintf("%s(%s)", enumType, fieldAccess)
	}
	return fieldAccess
}
