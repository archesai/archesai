package schema

import (
	"fmt"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/ref"
)

// Schema represents a unified OpenAPI schema definition.
// It combines parsing fields with code generation computed fields.
type Schema struct {
	// === Core ===
	ID      string                      `yaml:"$id,omitempty"      json:"$id,omitempty"`
	Schema  string                      `yaml:"$schema,omitempty"  json:"$schema,omitempty"`
	Ref     string                      `yaml:"$ref,omitempty"     json:"$ref,omitempty"`
	Comment string                      `yaml:"$comment,omitempty" json:"$comment,omitempty"`
	Defs    map[string]*ref.Ref[Schema] `yaml:"$defs,omitempty"    json:"$defs,omitempty"`

	Anchor        string          `yaml:"$anchor,omitempty"        json:"$anchor,omitempty"`
	DynamicAnchor string          `yaml:"$dynamicAnchor,omitempty" json:"$dynamicAnchor,omitempty"`
	DynamicRef    string          `yaml:"$dynamicRef,omitempty"    json:"$dynamicRef,omitempty"`
	Vocabulary    map[string]bool `yaml:"$vocabulary,omitempty"    json:"$vocabulary,omitempty"`

	// === Metadata ===
	Title       string `yaml:"title,omitempty"       json:"title,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Default     any    `yaml:"default,omitempty"     json:"default,omitempty"`
	Deprecated  bool   `yaml:"deprecated,omitempty"  json:"deprecated,omitempty"`
	ReadOnly    bool   `yaml:"readOnly,omitempty"    json:"readOnly,omitempty"`
	WriteOnly   bool   `yaml:"writeOnly,omitempty"   json:"writeOnly,omitempty"`
	Examples    []any  `yaml:"examples,omitempty"    json:"examples,omitempty"`

	// === Type Information ===
	Type PropertyType `yaml:"type,omitempty" json:"type,omitempty"`
	Enum []string     `yaml:"enum,omitempty" json:"enum,omitempty"`

	// Const is *any because a JSON null (Go nil) is a valid value.
	Const            *any     `yaml:"const,omitempty"            json:"const,omitempty"`
	MultipleOf       *float64 `yaml:"multipleOf,omitempty"       json:"multipleOf,omitempty"`
	Minimum          *float64 `yaml:"minimum,omitempty"          json:"minimum,omitempty"`
	Maximum          *float64 `yaml:"maximum,omitempty"          json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64 `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MinLength        *int     `yaml:"minLength,omitempty"        json:"minLength,omitempty"`
	MaxLength        *int     `yaml:"maxLength,omitempty"        json:"maxLength,omitempty"`
	Pattern          string   `yaml:"pattern,omitempty"          json:"pattern,omitempty"`

	// arrays
	PrefixItems      []*ref.Ref[Schema] `yaml:"prefixItems,omitempty"      json:"prefixItems,omitempty"`
	Items            *ref.Ref[Schema]   `yaml:"items,omitempty"            json:"items,omitempty"`
	MinItems         *int               `yaml:"minItems,omitempty"         json:"minItems,omitempty"`
	MaxItems         *int               `yaml:"maxItems,omitempty"         json:"maxItems,omitempty"`
	AdditionalItems  *ref.Ref[Schema]   `yaml:"additionalItems,omitempty"  json:"additionalItems,omitempty"`
	UniqueItems      bool               `yaml:"uniqueItems,omitempty"      json:"uniqueItems,omitempty"`
	Contains         *ref.Ref[Schema]   `yaml:"contains,omitempty"         json:"contains,omitempty"`
	MinContains      *int               `yaml:"minContains,omitempty"      json:"minContains,omitempty"`
	MaxContains      *int               `yaml:"maxContains,omitempty"      json:"maxContains,omitempty"`
	UnevaluatedItems *ref.Ref[Schema]   `yaml:"unevaluatedItems,omitempty" json:"unevaluatedItems,omitempty"`

	// objects
	MinProperties         *int                        `yaml:"minProperties,omitempty"         json:"minProperties,omitempty"`
	MaxProperties         *int                        `yaml:"maxProperties,omitempty"         json:"maxProperties,omitempty"`
	Required              []string                    `yaml:"required,omitempty"              json:"required,omitempty"`
	DependentRequired     map[string][]string         `yaml:"dependentRequired,omitempty"     json:"dependentRequired,omitempty"`
	Properties            map[string]*ref.Ref[Schema] `yaml:"properties,omitempty"            json:"properties,omitempty"`
	PatternProperties     map[string]*ref.Ref[Schema] `yaml:"patternProperties,omitempty"     json:"patternProperties,omitempty"`
	AdditionalProperties  *bool                       `yaml:"additionalProperties,omitempty"  json:"additionalProperties,omitempty"`
	PropertyNames         *ref.Ref[Schema]            `yaml:"propertyNames,omitempty"         json:"propertyNames,omitempty"`
	UnevaluatedProperties *bool                       `yaml:"unevaluatedProperties,omitempty" json:"unevaluatedProperties,omitempty"`

	// logic
	AllOf []*ref.Ref[Schema] `yaml:"allOf,omitempty" json:"allOf,omitempty"`
	AnyOf []*ref.Ref[Schema] `yaml:"anyOf,omitempty" json:"anyOf,omitempty"`
	OneOf []*ref.Ref[Schema] `yaml:"oneOf,omitempty" json:"oneOf,omitempty"`
	Not   *ref.Ref[Schema]   `yaml:"not,omitempty"   json:"not,omitempty"`

	// conditional
	If               *ref.Ref[Schema]            `yaml:"if,omitempty"               json:"if,omitempty"`
	Then             *ref.Ref[Schema]            `yaml:"then,omitempty"             json:"then,omitempty"`
	Else             *ref.Ref[Schema]            `yaml:"else,omitempty"             json:"else,omitempty"`
	DependentSchemas map[string]*ref.Ref[Schema] `yaml:"dependentSchemas,omitempty" json:"dependentSchemas,omitempty"`

	// other
	ContentEncoding  string           `yaml:"contentEncoding,omitempty"  json:"contentEncoding,omitempty"`
	ContentMediaType string           `yaml:"contentMediaType,omitempty" json:"contentMediaType,omitempty"`
	ContentSchema    *ref.Ref[Schema] `yaml:"contentSchema,omitempty"    json:"contentSchema,omitempty"`

	Format string `yaml:"format,omitempty" json:"format,omitempty"`

	// Extra allows for additional keywords beyond those specified.
	Extra map[string]any `yaml:"-" json:"-"`

	// === OpenAPI Extensions ===
	XCodegenSchemaType Type               `yaml:"x-codegen-schema-type,omitempty" json:"x-codegen-schema-type,omitempty"`
	XCodegen           *XCodegenExtension `yaml:"x-codegen,omitempty"             json:"x-codegen,omitempty"`
	XInternal          string             `yaml:"x-internal,omitempty"            json:"x-internal,omitempty"`

	// === Computed Fields (populated during resolution) ===
	GoType   string `yaml:"-" json:"-"`
	JSONTag  string `yaml:"-" json:"-"`
	YAMLTag  string `yaml:"-" json:"-"`
	Nullable bool   `yaml:"-" json:"-"`

	// === Resolution State ===
	resolved bool `yaml:"-" json:"-"`
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
	r, ok := s.Properties[name]
	if !ok {
		return nil
	}
	return r.GetOrNil()
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
	return s != nil && s.Title != "" && !s.HasProperties() && s.Type.PrimaryType() == ""
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

// SchemaTypeValue returns the x-codegen-schema-type value as a string.
// This is used by templates that check for "entity" or "valueobject".
func (s *Schema) SchemaTypeValue() string {
	if s == nil {
		return ""
	}
	return string(s.XCodegenSchemaType)
}

// SchemaType returns the x-codegen-schema-type value as a string.
// Alias for SchemaTypeValue, used by templates.
func (s *Schema) SchemaType() string {
	return s.SchemaTypeValue()
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
	return s != nil && s.XCodegenSchemaType == TypeEntity
}

// IsEntity returns true if the schema is marked as an entity (x-codegen-schema-type: entity).
func (s *Schema) IsEntity() bool {
	return s != nil && s.XCodegenSchemaType == TypeEntity
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
		if prop.Title == "" {
			prop.Title = name
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
		if prop.Type.PrimaryType() == TypeObject && prop.Title != "" && prop.GoType != "" &&
			!strings.Contains(prop.GoType, ".") && len(prop.Properties) > 0 {
			if !visited[prop.GoType] {
				visited[prop.GoType] = true
				*result = append(*result, prop)
				prop.collectNestedTypesRecursive(result, visited)
			}
		}
		if prop.Type.PrimaryType() == "array" && prop.Items != nil {
			items := prop.Items.GetOrNil()
			if items != nil && items.Type.PrimaryType() == TypeObject && items.Title != "" &&
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
		// Use Name for enum types to match the local type definition in templates.
		// For nested types, GoType may include a package prefix (e.g., authmodels.Account)
		// but enum types are defined locally without the prefix.
		parentName := parentSchema.GoType
		// Strip package prefix if present (e.g., "authmodels.Account" -> "Account")
		if dotIdx := strings.LastIndex(parentName, "."); dotIdx >= 0 {
			parentName = parentName[dotIdx+1:]
		}
		if parentName == "" {
			parentName = parentSchema.Title
		}
		baseType = parentName + s.Title
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
		// Use Name for enum types to match the local type definition in templates.
		parentName := parentSchema.GoType
		// Strip package prefix if present (e.g., "authmodels.Account" -> "Account")
		if dotIdx := strings.LastIndex(parentName, "."); dotIdx >= 0 {
			parentName = parentName[dotIdx+1:]
		}
		if parentName == "" {
			parentName = parentSchema.Title
		}
		return parentName + s.Title
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
	if s.GoType == GoTypeString && !s.Nullable {
		return true
	}
	if (s.GoType == GoTypeInt || s.GoType == GoTypeFloat64) && !s.Nullable {
		return true
	}
	if s.GoType == GoTypeUUID && !s.Nullable {
		return true
	}
	return false
}

// IsSpecialField returns true if this is a special entity field.
func (s *Schema) IsSpecialField() bool {
	if s == nil {
		return false
	}
	return s.Title == "ID" || s.Title == "CreatedAt" || s.Title == "UpdatedAt"
}

// ShouldIncludeInConstructor returns true if field should be in entity constructor.
func (s *Schema) ShouldIncludeInConstructor(parentSchema *Schema) bool {
	if s == nil || parentSchema == nil {
		return false
	}
	if parentSchema.XCodegenSchemaType == TypeEntity {
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
		return fmt.Sprintf("invalid %s: %%s", s.Title)
	}
	if s.GoType == "string" {
		return fmt.Sprintf("%s cannot be empty", s.Title)
	}
	if s.GoType == "int" || s.GoType == "float64" {
		return fmt.Sprintf("%s cannot be negative", s.Title)
	}
	if s.GoType == "uuid.UUID" {
		return fmt.Sprintf("%s cannot be nil UUID", s.Title)
	}
	return fmt.Sprintf("invalid %s", s.Title)
}

// GetCreateParamValue returns the Go expression for mapping entity field to DB param in Create.
func (s *Schema) GetCreateParamValue(entityVar string, _ *Schema) string {
	if s == nil {
		return ""
	}
	fieldAccess := fmt.Sprintf("%s.%s", entityVar, s.Title)
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
	fieldAccess := fmt.Sprintf("%s.%s", entityVar, s.Title)
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Title[0])), s.Title[1:]) + "Str"
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
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Title[0])), s.Title[1:]) + "Str"
	return fmt.Sprintf("%s := string(%s.%s)", varName, entityVar, s.Title)
}

// GetDBMapValue returns the Go expression for mapping DB field to entity field.
func (s *Schema) GetDBMapValue(dbVar string, parentSchema *Schema) string {
	if s == nil {
		return ""
	}
	fieldAccess := fmt.Sprintf("%s.%s", dbVar, s.Title)
	enumType := ""
	if parentSchema != nil && len(s.Enum) > 0 {
		enumType = fmt.Sprintf("schemas.%s%s", parentSchema.Title, s.Title)
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

// ConvertFileRefs recursively converts file $refs to internal refs in the schema.
func (s *Schema) ConvertFileRefs(fileRefToInternalRef func(string, string) string) {
	if s == nil {
		return
	}

	// Convert the schema's own ref if present
	if s.Ref != "" {
		s.Ref = fileRefToInternalRef(s.Ref, "")
	}

	// Convert refs in Properties
	for _, propRef := range s.Properties {
		convertRefFileRefs(propRef, fileRefToInternalRef)
	}

	// Convert refs in Items
	convertRefFileRefs(s.Items, fileRefToInternalRef)

	// Convert refs in AllOf
	for _, r := range s.AllOf {
		convertRefFileRefs(r, fileRefToInternalRef)
	}

	// Convert refs in AnyOf
	for _, r := range s.AnyOf {
		convertRefFileRefs(r, fileRefToInternalRef)
	}

	// Convert refs in OneOf
	for _, r := range s.OneOf {
		convertRefFileRefs(r, fileRefToInternalRef)
	}

	// Convert refs in Not
	convertRefFileRefs(s.Not, fileRefToInternalRef)

	// Convert refs in AdditionalItems
	convertRefFileRefs(s.AdditionalItems, fileRefToInternalRef)

	// Convert refs in Contains
	convertRefFileRefs(s.Contains, fileRefToInternalRef)

	// Convert refs in PrefixItems
	for _, r := range s.PrefixItems {
		convertRefFileRefs(r, fileRefToInternalRef)
	}

	// Convert refs in PropertyNames
	convertRefFileRefs(s.PropertyNames, fileRefToInternalRef)

	// Convert refs in If/Then/Else
	convertRefFileRefs(s.If, fileRefToInternalRef)
	convertRefFileRefs(s.Then, fileRefToInternalRef)
	convertRefFileRefs(s.Else, fileRefToInternalRef)

	// Convert refs in PatternProperties
	for _, r := range s.PatternProperties {
		convertRefFileRefs(r, fileRefToInternalRef)
	}

	// Convert refs in DependentSchemas
	for _, r := range s.DependentSchemas {
		convertRefFileRefs(r, fileRefToInternalRef)
	}

	// Convert refs in ContentSchema
	convertRefFileRefs(s.ContentSchema, fileRefToInternalRef)

	// Convert refs in UnevaluatedItems
	convertRefFileRefs(s.UnevaluatedItems, fileRefToInternalRef)

	// Convert refs in $defs
	for _, r := range s.Defs {
		convertRefFileRefs(r, fileRefToInternalRef)
	}
}

// convertRefFileRefs converts file refs in a Ref[Schema].
func convertRefFileRefs(r *ref.Ref[Schema], fileRefToInternalRef func(string, string) string) {
	if r == nil {
		return
	}
	if r.RefPath != "" {
		r.RefPath = fileRefToInternalRef(r.RefPath, "")
	}
	if r.Value != nil {
		r.Value.ConvertFileRefs(fileRefToInternalRef)
	}
}
