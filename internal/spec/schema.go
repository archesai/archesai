package spec

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// XCodegenSchemaType represents the enumeration of valid values for SchemaType
type XCodegenSchemaType string

// Valid SchemaType values
const (
	XCodegenSchemaTypeEntity      XCodegenSchemaType = "entity"
	XCodegenSchemaTypeValueobject XCodegenSchemaType = "valueobject"
)

// Schema represents any schema - from a full object to a single field
type Schema struct {
	// Identity
	Name        string
	Description string

	// Type information
	Type   string // "object", "string", "integer", "array", "boolean", "number"
	Format string // "uuid", "email", "date-time", etc.

	// For objects: nested properties
	Properties map[string]*Schema
	Required   []string // List of required property names

	// For arrays: item schema
	Items *Schema

	// For all types
	Enum         []string
	DefaultValue any
	Nullable     bool

	// Code generation
	GoType  string // Computed Go type name
	JSONTag string // JSON tag value
	YAMLTag string // YAML tag value

	// Extensions
	XCodegen           *XCodegenExtension
	XCodegenSchemaType XCodegenSchemaType
	XInternal          string // When set (e.g., "server", "config"), this schema should be imported not generated

	// Original OpenAPI schema reference
	Schema *base.Schema
}

// IsInternal returns true if this schema should be imported from another package instead of generated.
// If context is empty, it returns true whenever XInternal is set.
// If context is provided, it returns true only when XInternal is set AND doesn't match the context.
func (s *Schema) IsInternal(context string) bool {
	if s.XInternal == "" {
		return false
	}
	// If no context provided, treat all internal schemas as internal
	if context == "" {
		return true
	}
	// Only internal if the internal tag doesn't match the current context
	return s.XInternal != context
}

// IsEnum returns true if the schema is an enum
func (s *Schema) IsEnum() bool {
	return len(s.Enum) > 0
}

// HasDomainEvents returns true if the schema has domain events
func (s *Schema) HasDomainEvents() bool {
	return s.XCodegenSchemaType == XCodegenSchemaTypeEntity
}

// HasProperty returns true if the schema has a property with the given name
func (s *Schema) HasProperty(name string) bool {
	if s.Properties == nil {
		return false
	}
	_, ok := s.Properties[name]
	return ok
}

// GetRepositoryIndices returns the repository indices
func (s *Schema) GetRepositoryIndices() []string {
	if s.XCodegen == nil || s.XCodegen.Repository == nil {
		return nil
	}
	return s.XCodegen.Repository.Indices
}

// GetRepositoryRelations returns the repository relations
func (s *Schema) GetRepositoryRelations() []XCodegenExtensionRepositoryRelationsItem {
	if s.XCodegen == nil || s.XCodegen.Repository == nil {
		return nil
	}
	return s.XCodegen.Repository.Relations
}

// GetRequiredProperties returns only the required properties
func (s *Schema) GetRequiredProperties() map[string]*Schema {
	required := make(map[string]*Schema)
	for _, name := range s.Required {
		if prop, exists := s.Properties[name]; exists {
			required[name] = prop
		}
	}
	return required
}

// GetOptionalProperties returns only the optional properties
func (s *Schema) GetOptionalProperties() map[string]*Schema {
	optional := make(map[string]*Schema)
	requiredMap := make(map[string]bool)
	for _, name := range s.Required {
		requiredMap[name] = true
	}
	for name, prop := range s.Properties {
		if !requiredMap[name] {
			optional[name] = prop
		}
	}
	return optional
}

// GetSortedProperties returns properties sorted by key for consistent iteration
func (s *Schema) GetSortedProperties() []*Schema {
	if s.Properties == nil {
		return nil
	}

	// Create a slice of property names and sort them
	var names []string
	for name := range s.Properties {
		names = append(names, name)
	}
	// FIXME: Ensure id, created_at, and updated_at stay at the top
	// Sort alphabetically but put ID, CreatedAt, and UpdatedAt first
	sort.Slice(names, func(i, j int) bool {
		// Define priority order for special fields
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

		// If both have the same priority, sort alphabetically
		if priI == priJ {
			return names[i] < names[j]
		}

		// Otherwise, sort by priority
		return priI < priJ
	})

	// Build the sorted slice
	var sorted []*Schema
	for _, name := range names {
		prop := s.Properties[name]
		// Ensure the property knows its own name
		if prop.Name == "" {
			prop.Name = name
		}
		sorted = append(sorted, prop)
	}
	return sorted
}

// IsPropertyRequired checks if a property is required
func (s *Schema) IsPropertyRequired(propName string) bool {
	for _, req := range s.Required {
		if req == propName {
			return true
		}
	}
	return false
}

// CollectNestedTypes walks the schema tree and collects all nested object types that should be generated as separate structs.
// Returns nested types in the order they appear in the root struct (depth-first traversal).
func (s *Schema) CollectNestedTypes() []*Schema {
	var nestedTypes []*Schema
	s.collectNestedTypesRecursive(&nestedTypes, make(map[string]bool))
	return nestedTypes
}

// collectNestedTypesRecursive is the recursive helper for collecting nested types.
// It traverses properties in sorted order to ensure deterministic output.
func (s *Schema) collectNestedTypesRecursive(result *[]*Schema, visited map[string]bool) {
	if s == nil || s.Properties == nil {
		return
	}

	// Use GetSortedProperties to ensure deterministic ordering
	for _, prop := range s.GetSortedProperties() {
		if prop == nil {
			continue
		}

		// If this property is an object type with a name and isn't a reference, it should be generated as a nested type
		// Skip references - they have no properties since they're just pointers to other schemas
		if prop.Type == SchemaTypeObject && prop.Name != "" && prop.GoType != "" &&
			!strings.Contains(prop.GoType, ".") && len(prop.Properties) > 0 {
			// Check if we've already visited this type to avoid duplicates
			if !visited[prop.GoType] {
				visited[prop.GoType] = true
				*result = append(*result, prop)
				// Recursively collect nested types from this object
				prop.collectNestedTypesRecursive(result, visited)
			}
		}

		// If it's an array of objects, check the items
		if prop.Type == "array" && prop.Items != nil {
			// If the array item is an object type with a name, it should be generated
			// Skip references - they have no properties since they're just pointers to other schemas
			if prop.Items.Type == SchemaTypeObject && prop.Items.Name != "" &&
				prop.Items.GoType != "" &&
				!strings.Contains(prop.Items.GoType, ".") &&
				len(prop.Items.Properties) > 0 {
				if !visited[prop.Items.GoType] {
					visited[prop.Items.GoType] = true
					*result = append(*result, prop.Items)
					// Recursively collect nested types from this object
					prop.Items.collectNestedTypesRecursive(result, visited)
				}
			} else {
				// Still recursively check for deeper nested types
				prop.Items.collectNestedTypesRecursive(result, visited)
			}
		}
	}
}

// IsOptional returns true if the property is optional (has omitempty tag)
func (s *Schema) IsOptional() bool {
	return strings.Contains(s.JSONTag, ",omitempty")
}

// IsRequired returns true if the property is required (doesn't have omitempty tag)
func (s *Schema) IsRequired() bool {
	return !s.IsOptional()
}

// NeedsPointer returns true if the field type needs to be a pointer
func (s *Schema) NeedsPointer() bool {
	// Optional fields need pointers unless they're already pointers, slices, or maps
	if s.IsOptional() {
		return !strings.HasPrefix(s.GoType, "*") &&
			!strings.HasPrefix(s.GoType, "[]") &&
			!strings.HasPrefix(s.GoType, "map")
	}
	// Non-optional nullable fields need pointers
	return s.Nullable
}

// GetFieldType returns the complete Go type for this field, including pointer if needed
func (s *Schema) GetFieldType(parentSchema *Schema) string {
	baseType := s.GoType
	if len(s.Enum) > 0 && parentSchema != nil {
		baseType = parentSchema.Name + s.Name
	}
	if s.NeedsPointer() {
		return "*" + baseType
	}
	return baseType
}

// GetBaseType returns the base Go type without pointer or qualification
func (s *Schema) GetBaseType(parentSchema *Schema) string {
	if len(s.Enum) > 0 && parentSchema != nil {
		return parentSchema.Name + s.Name
	}
	return s.GoType
}

// NeedsValidation returns true if this field requires validation in constructor
func (s *Schema) NeedsValidation(isEntity bool, parentSchema *Schema) bool {
	// For entities, only validate required non-optional fields that aren't special
	if isEntity {
		if s.IsSpecialField() {
			return false
		}
		if !parentSchema.IsPropertyRequired(s.JSONTag) || s.IsOptional() {
			return false
		}
	} else {
		// For value objects, validate all non-optional fields
		if s.IsOptional() {
			return false
		}
	}

	// Now check if the field type needs validation
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

// IsSpecialField returns true if this is a special entity field
func (s *Schema) IsSpecialField() bool {
	return s.Name == "ID" || s.Name == "CreatedAt" || s.Name == "UpdatedAt"
}

// ShouldIncludeInConstructor returns true if field should be in entity constructor
func (s *Schema) ShouldIncludeInConstructor(parentSchema *Schema) bool {
	// For entities, only include required, non-special fields
	if parentSchema.XCodegenSchemaType == XCodegenSchemaTypeEntity {
		return parentSchema.IsPropertyRequired(s.JSONTag) &&
			!s.IsOptional() &&
			!s.IsSpecialField()
	}
	// For value objects, include all fields
	return true
}

// GetValidationError returns the validation error message for this field
func (s *Schema) GetValidationError(_ *Schema) string {
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

// GetCreateParamValue returns the Go expression for mapping entity field to DB param in Create operations
func (s *Schema) GetCreateParamValue(entityVar string, _ *Schema) string {
	fieldAccess := fmt.Sprintf("%s.%s", entityVar, s.Name)

	// Optional fields - pass through directly
	if s.IsOptional() {
		return fieldAccess
	}

	// Nullable enum - convert to *string
	if s.Nullable && len(s.Enum) > 0 {
		return fmt.Sprintf(
			"func() *string { if %s == nil { return nil }; s := string(*%s); return &s }()",
			fieldAccess, fieldAccess,
		)
	}

	// Regular enum - convert to string
	if len(s.Enum) > 0 {
		return fmt.Sprintf("string(%s)", fieldAccess)
	}

	// Nullable or pointer/slice types - pass through
	if s.Nullable || strings.HasPrefix(s.GoType, "*") ||
		strings.HasPrefix(s.GoType, "[]") {
		return fieldAccess
	}

	// Everything else - pass through
	return fieldAccess
}

// GetUpdateParamValue returns the Go expression for mapping entity field to DB param in Update operations
func (s *Schema) GetUpdateParamValue(entityVar string, _ *Schema) string {
	fieldAccess := fmt.Sprintf("%s.%s", entityVar, s.Name)
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Name[0])), s.Name[1:]) + "Str"

	// Optional fields - pass through directly
	if s.IsOptional() {
		return fieldAccess
	}

	// Nullable enum - convert to *string
	if s.Nullable && len(s.Enum) > 0 {
		return fmt.Sprintf(
			"func() *string { if %s == nil { return nil }; s := string(*%s); return &s }()",
			fieldAccess, fieldAccess,
		)
	}

	// Regular enum - use pre-declared variable
	if len(s.Enum) > 0 {
		return "&" + varName
	}

	// Nullable or already pointer/slice - pass through
	if s.Nullable || strings.HasPrefix(s.GoType, "*") ||
		strings.HasPrefix(s.GoType, "[]") {
		return fieldAccess
	}

	// Primitives need pointer
	return "&" + fieldAccess
}

// NeedsUpdateVarDeclaration returns true if Update operation needs a variable declaration
func (s *Schema) NeedsUpdateVarDeclaration() bool {
	return !s.IsOptional() && len(s.Enum) > 0 && !s.Nullable
}

// GetUpdateVarDeclaration returns the variable declaration for Update operation
func (s *Schema) GetUpdateVarDeclaration(entityVar string) string {
	if !s.NeedsUpdateVarDeclaration() {
		return ""
	}
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Name[0])), s.Name[1:]) + "Str"
	return fmt.Sprintf("%s := string(%s.%s)", varName, entityVar, s.Name)
}

// GetDBMapValue returns the Go expression for mapping DB field to entity field
func (s *Schema) GetDBMapValue(dbVar string, parentSchema *Schema) string {
	fieldAccess := fmt.Sprintf("%s.%s", dbVar, s.Name)
	enumType := ""
	if parentSchema != nil && len(s.Enum) > 0 {
		enumType = fmt.Sprintf("models.%s%s", parentSchema.Name, s.Name)
	}

	// Nullable enum - convert *string to *EnumType
	if s.Nullable && len(s.Enum) > 0 {
		return fmt.Sprintf(
			"func() *%s { if %s == nil { return nil }; v := %s(*%s); return &v }()",
			enumType, fieldAccess, enumType, fieldAccess,
		)
	}

	// Regular enum - cast string to EnumType
	if len(s.Enum) > 0 {
		return fmt.Sprintf("%s(%s)", enumType, fieldAccess)
	}

	// Everything else - pass through
	return fieldAccess
}
