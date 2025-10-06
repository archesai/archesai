package parsers

import (
	"sort"
	"strconv"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
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
	*SchemaDef        // Embed schema definition for response body
	StatusCode string // HTTP status code
}

// IsSuccess returns true if the response is a successful one (2xx status code)
func (r *ResponseDef) IsSuccess() bool {
	if code, err := strconv.Atoi(r.StatusCode); err == nil {
		return code >= 200 && code < 300
	}
	return false
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
	Schema *oas3.Schema
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
	// Sort alphabetically but put ID first if it exists
	sort.Slice(names, func(i, j int) bool {
		if strings.EqualFold(names[i], "ID") {
			return true
		}
		if strings.EqualFold(names[j], "ID") {
			return false
		}
		return names[i] < names[j]
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
