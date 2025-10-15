package parsers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// Constants for various type mappings and validations
const (
	// SQL Dialects
	SQLDialectPostgres = "postgresql"
	SQLDialectSQLite   = "sqlite"

	// Default SQL Types
	SQLTypeText      = "TEXT"
	SQLTypeInteger   = "INTEGER"
	SQLTypeBigInt    = "BIGINT"
	SQLTypeBoolean   = "BOOLEAN"
	SQLTypeTimestamp = "TIMESTAMPTZ"
	SQLTypeDateTime  = "DATETIME"
	SQLTypeDate      = "DATE"
	SQLTypeUUID      = "UUID"
	SQLTypeJSONB     = "JSONB"
	SQLTypeNumeric   = "NUMERIC"
	SQLTypeReal      = "REAL"
	SQLTypeDouble    = "DOUBLE PRECISION"
)

// OperationDef represents an API operation
type OperationDef struct {
	ID                    string          // Original operation ID from OpenAPI
	Method                string          // HTTP method (GET, POST, etc.)
	Path                  string          // URL path
	Description           string          // Operation description
	Tag                   string          // Operation tags
	Parameters            []ParamDef      // All parameters (backward compat)
	Responses             []ResponseDef   // All responses
	Security              []SecurityDef   // Security requirements
	RequestBody           *RequestBodyDef // Processed request body schema
	XCodegenCustomHandler bool            // Whether this operation has a custom handler implementation
	XCodegenRepository    string          // Custom repository name from x-codegen-repository extension
}

// GetSuccessResponse returns the first successful response (2xx status code)
func (o *OperationDef) GetSuccessResponse() *ResponseDef {
	for _, resp := range o.Responses {
		if resp.IsSuccess() {
			return &resp
		}
	}
	return nil
}

// GetErrorResponses returns all error responses (non-2xx status codes)
func (o *OperationDef) GetErrorResponses() []ResponseDef {
	var errors []ResponseDef
	for _, resp := range o.Responses {
		if resp.IsSuccess() {
			continue
		}
		errors = append(errors, resp)
	}
	return errors
}

// GetQueryParams returns only the query parameters
func (o *OperationDef) GetQueryParams() []ParamDef {
	var queryParams []ParamDef
	for _, p := range o.Parameters {
		if p.In == "query" {
			queryParams = append(queryParams, p)
		}
	}
	return queryParams
}

// GetPathParams returns only the path parameters
func (o *OperationDef) GetPathParams() []ParamDef {
	var pathParams []ParamDef
	for _, p := range o.Parameters {
		if p.In == "path" {
			pathParams = append(pathParams, p)
		}
	}
	return pathParams
}

// GetHeaderParams returns only the header parameters
func (o *OperationDef) GetHeaderParams() []ParamDef {
	var headerParams []ParamDef
	for _, p := range o.Parameters {
		if p.In == "header" {
			headerParams = append(headerParams, p)
		}
	}
	return headerParams
}

// HasBearerAuth checks if the operation requires bearer token authentication
func (o *OperationDef) HasBearerAuth() bool {
	for _, sec := range o.Security {
		if sec.Type == "http" && strings.EqualFold(sec.Scheme, "bearer") {
			return true
		}
	}
	return false
}

// HasCookieAuth checks if the operation requires cookie-based authentication
func (o *OperationDef) HasCookieAuth() bool {
	for _, sec := range o.Security {
		if sec.Type == "apiKey" && strings.EqualFold(sec.Scheme, "cookie") {
			return true
		}
	}
	return false
}

// ResponseDef represents a response in an operation
type ResponseDef struct {
	*SchemaDef                        // Embed schema definition for response body
	StatusCode  string                // HTTP status code
	ContentType string                // Content-Type for the response (e.g., "application/json")
	Headers     map[string]*SchemaDef // Response headers
}

// IsSuccess returns true if the response is a successful one (2xx status code)
func (r *ResponseDef) IsSuccess() bool {
	if code, err := strconv.Atoi(r.StatusCode); err == nil {
		return code >= 200 && code < 300
	}
	return false
}

// GetSortedHeaders returns headers sorted by name for consistent iteration
func (r *ResponseDef) GetSortedHeaders() []struct {
	Name   string
	Schema *SchemaDef
} {
	if r.Headers == nil {
		return nil
	}

	// Create a slice of header names and sort them
	var names []string
	for name := range r.Headers {
		names = append(names, name)
	}
	sort.Strings(names)

	// Build the sorted slice
	var sorted []struct {
		Name   string
		Schema *SchemaDef
	}
	for _, name := range names {
		sorted = append(sorted, struct {
			Name   string
			Schema *SchemaDef
		}{
			Name:   name,
			Schema: r.Headers[name],
		})
	}
	return sorted
}

// RequestBodyDef represents the request body definition for an API operation
type RequestBodyDef struct {
	*SchemaDef      // Embed schema definition for request body
	Required   bool // Whether request body is required
}

// ParamDef represents a parameter in an operation
type ParamDef struct {
	*SchemaDef        // Embed schema definition
	In         string // Location (path, query, header, cookie)
	Style      string // Parameter style (form, simple, etc.)
	Explode    bool   // Whether to explode array/object parameters
}

// SecurityDef represents a security requirement
type SecurityDef struct {
	Name   string   // Security scheme name
	Type   string   // Security type (http, apiKey, oauth2)
	Scheme string   // Security scheme (bearer for http, cookie for apiKey)
	Scopes []string // Required scopes
}

// SchemaDef represents any schema - from a full object to a single field
type SchemaDef struct {
	// Identity
	Name        string
	Description string

	// Type information
	Type   string // "object", "string", "integer", "array", "boolean", "number"
	Format string // "uuid", "email", "date-time", etc.

	// For objects: nested properties
	Properties map[string]*SchemaDef
	Required   []string // List of required property names

	// For arrays: item schema
	Items *SchemaDef

	// For all types
	Enum         []string
	DefaultValue any
	Nullable     bool

	// Code generation
	GoType  string // Computed Go type name
	JSONTag string // JSON tag value
	YAMLTag string // YAML tag value

	// Extensions
	XCodegen *XCodegenExtension

	// Original OpenAPI schema reference
	Schema *base.Schema
}

// IsEnum returns true if the schema is an enum
func (s *SchemaDef) IsEnum() bool {
	return len(s.Enum) > 0
}

// HasDomainEvents returns true if the schema has domain events
func (s *SchemaDef) HasDomainEvents() bool {
	return s.GetSchemaType() == string(XCodegenExtensionSchemaTypeEntity)
}

// GetSchemaType returns the schema type as a string
func (s *SchemaDef) GetSchemaType() string {
	if s.XCodegen != nil && s.XCodegen.SchemaType != "" {
		return s.XCodegen.SchemaType.String()
	}
	return string(XCodegenExtensionSchemaTypeValueobject)
}

// GetRequiredProperties returns only the required properties
func (s *SchemaDef) GetRequiredProperties() map[string]*SchemaDef {
	required := make(map[string]*SchemaDef)
	for _, name := range s.Required {
		if prop, exists := s.Properties[name]; exists {
			required[name] = prop
		}
	}
	return required
}

// GetOptionalProperties returns only the optional properties
func (s *SchemaDef) GetOptionalProperties() map[string]*SchemaDef {
	optional := make(map[string]*SchemaDef)
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
func (s *SchemaDef) GetSortedProperties() []*SchemaDef {
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
	var sorted []*SchemaDef
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
func (s *SchemaDef) IsPropertyRequired(propName string) bool {
	for _, req := range s.Required {
		if req == propName {
			return true
		}
	}
	return false
}

// CollectNestedTypes walks the schema tree and collects all nested object types that should be generated as separate structs.
// Returns nested types in the order they appear in the root struct (depth-first traversal).
func (s *SchemaDef) CollectNestedTypes() []*SchemaDef {
	var nestedTypes []*SchemaDef
	s.collectNestedTypesRecursive(&nestedTypes, make(map[string]bool))
	return nestedTypes
}

// collectNestedTypesRecursive is the recursive helper for collecting nested types.
// It traverses properties in sorted order to ensure deterministic output.
func (s *SchemaDef) collectNestedTypesRecursive(result *[]*SchemaDef, visited map[string]bool) {
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
		if prop.Type == schemaTypeObject && prop.Name != "" && prop.GoType != "" &&
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
			if prop.Items.Type == "object" && prop.Items.Name != "" && prop.Items.GoType != "" &&
				!strings.Contains(prop.Items.GoType, ".") && len(prop.Items.Properties) > 0 {
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

// Template helper methods for cleaner code generation

// IsOptional returns true if the property is optional (has omitempty tag)
func (s *SchemaDef) IsOptional() bool {
	return strings.Contains(s.JSONTag, ",omitempty")
}

// IsRequired returns true if the property is required (doesn't have omitempty tag)
func (s *SchemaDef) IsRequired() bool {
	return !s.IsOptional()
}

// NeedsPointer returns true if the field type needs to be a pointer
func (s *SchemaDef) NeedsPointer() bool {
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
func (s *SchemaDef) GetFieldType(parentSchema *SchemaDef) string {
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
func (s *SchemaDef) GetBaseType(parentSchema *SchemaDef) string {
	if len(s.Enum) > 0 && parentSchema != nil {
		return parentSchema.Name + s.Name
	}
	return s.GoType
}

// NeedsValidation returns true if this field requires validation in constructor
func (s *SchemaDef) NeedsValidation(isEntity bool, parentSchema *SchemaDef) bool {
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
func (s *SchemaDef) IsSpecialField() bool {
	return s.Name == "ID" || s.Name == "CreatedAt" || s.Name == "UpdatedAt"
}

// ShouldIncludeInConstructor returns true if field should be in entity constructor
func (s *SchemaDef) ShouldIncludeInConstructor(parentSchema *SchemaDef) bool {
	// For entities, only include required, non-special fields
	if parentSchema.GetSchemaType() == string(XCodegenExtensionSchemaTypeEntity) {
		return parentSchema.IsPropertyRequired(s.JSONTag) &&
			!s.IsOptional() &&
			!s.IsSpecialField()
	}
	// For value objects, include all fields
	return true
}

// GetValidationError returns the validation error message for this field
func (s *SchemaDef) GetValidationError(_ *SchemaDef) string {
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

// Repository template helpers for database operations

// GetCreateParamValue returns the Go expression for mapping entity field to DB param in Create operations
func (s *SchemaDef) GetCreateParamValue(entityVar string, _ *SchemaDef) string {
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
func (s *SchemaDef) GetUpdateParamValue(entityVar string, _ *SchemaDef) string {
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
func (s *SchemaDef) NeedsUpdateVarDeclaration() bool {
	return !s.IsOptional() && len(s.Enum) > 0 && !s.Nullable
}

// GetUpdateVarDeclaration returns the variable declaration for Update operation
func (s *SchemaDef) GetUpdateVarDeclaration(entityVar string) string {
	if !s.NeedsUpdateVarDeclaration() {
		return ""
	}
	varName := fmt.Sprintf("%s%s", strings.ToLower(string(s.Name[0])), s.Name[1:]) + "Str"
	return fmt.Sprintf("%s := string(%s.%s)", varName, entityVar, s.Name)
}

// GetDBMapValue returns the Go expression for mapping DB field to entity field
func (s *SchemaDef) GetDBMapValue(dbVar string, parentSchema *SchemaDef) string {
	fieldAccess := fmt.Sprintf("%s.%s", dbVar, s.Name)
	enumType := ""
	if parentSchema != nil && len(s.Enum) > 0 {
		enumType = fmt.Sprintf("entities.%s%s", parentSchema.Name, s.Name)
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
