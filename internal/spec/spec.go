package spec

import (
	"slices"
	"sort"
	"strconv"
	"strings"
)

// Spec represents the entire specification definition.
type Spec struct {
	Operations      []Operation        // All operations in the spec
	Schemas         map[string]*Schema // All schemas defined in the spec (keyed by ref path or name)
	ProjectName     string             // Project name from x-project-name extension
	EnabledIncludes []string           // Names of enabled x-include-* extensions

	// OpenAPI metadata for generation
	Title       string
	Description string
	Version     string
	Tags        []Tag
	Security    map[string]SecScheme

	// Codegen options from x-codegen-* extensions
	CodegenOnly []string // Generators to run (empty = all)
	CodegenLint bool     // Enable strict linting
}

// HasInclude returns true if the named include is enabled (e.g., "auth", "executors").
func (s *Spec) HasInclude(name string) bool {
	return slices.Contains(s.EnabledIncludes, name)
}

// GetSchema returns a schema by name from the registry.
func (s *Spec) GetSchema(name string) *Schema {
	if s == nil || s.Schemas == nil {
		return nil
	}
	return s.Schemas[name]
}

// GetSortedSchemas returns schemas sorted by name for consistent iteration.
func (s *Spec) GetSortedSchemas() []*Schema {
	if s == nil || s.Schemas == nil {
		return nil
	}
	var names []string
	for name := range s.Schemas {
		names = append(names, name)
	}
	sort.Strings(names)
	var sorted []*Schema
	for _, name := range names {
		sorted = append(sorted, s.Schemas[name])
	}
	return sorted
}

// SchemasSlice returns schemas as a slice for backwards compatibility.
func (s *Spec) SchemasSlice() []*Schema {
	return s.GetSortedSchemas()
}

// Operation is the domain type for code generation (converted from OpenAPIOperation by parser).
type Operation struct {
	ID               string       // Operation ID
	Method           string       // HTTP method (GET, POST, etc.)
	Path             string       // URL path
	Summary          string       // Operation summary
	Description      string       // Operation description
	Tag              string       // Primary tag
	Parameters       []Param      // Processed parameters with schema info
	RequestBody      *RequestBody // Processed request body schema
	Responses        []Response   // Processed responses with schema info
	Security         []Security   // Processed security requirements
	EmptySecurity    bool         // Whether this operation has explicitly empty security
	ExplicitSecurity bool         // Whether security was explicitly set
	CustomHandler    bool         // Whether has custom handler (x-codegen-custom-handler)
	PublicEndpoint   bool         // Whether is public endpoint (x-public-endpoint)
	Internal         string       // Internal context (x-internal)
}

// Param represents a parameter in an operation.
type Param struct {
	*Schema        // Embed schema definition
	In      string // Location (path, query, header, cookie)
	Style   string // Parameter style (form, simple, etc.)
	Explode bool   // Whether to explode array/object parameters
}

// RequestBody represents the request body definition for an API operation.
type RequestBody struct {
	*Schema       // Embed schema definition for request body
	Required bool // Whether request body is required
}

// Response represents a response in an operation.
type Response struct {
	*Schema                        // Embed schema definition for response body
	StatusCode  string             // HTTP status code
	ContentType string             // Content-Type for the response
	Headers     map[string]*Schema // Response headers
}

// IsSuccess returns true if the response is a successful one (2xx status code).
func (r *Response) IsSuccess() bool {
	if code, err := strconv.Atoi(r.StatusCode); err == nil {
		return code >= 200 && code < 300
	}
	return false
}

// GetSortedHeaders returns headers sorted by name for consistent iteration.
func (r *Response) GetSortedHeaders() []struct {
	Name   string
	Schema *Schema
} {
	if r.Headers == nil {
		return nil
	}
	var names []string
	for name := range r.Headers {
		names = append(names, name)
	}
	sort.Strings(names)
	var sorted []struct {
		Name   string
		Schema *Schema
	}
	for _, name := range names {
		sorted = append(sorted, struct {
			Name   string
			Schema *Schema
		}{
			Name:   name,
			Schema: r.Headers[name],
		})
	}
	return sorted
}

// IsInternal returns true if this operation should be imported from another package.
func (o *Operation) IsInternal(context string) bool {
	if o.Internal == "" {
		return false
	}
	if context == "" {
		return true
	}
	return o.Internal != context
}

// GetSuccessResponse returns the first successful response (2xx status code).
func (o *Operation) GetSuccessResponse() *Response {
	for i := range o.Responses {
		if o.Responses[i].IsSuccess() {
			return &o.Responses[i]
		}
	}
	return nil
}

// GetErrorResponses returns all error responses (non-2xx status codes).
func (o *Operation) GetErrorResponses() []Response {
	var errors []Response
	for _, resp := range o.Responses {
		if !resp.IsSuccess() {
			errors = append(errors, resp)
		}
	}
	return errors
}

// GetQueryParams returns only the query parameters.
func (o *Operation) GetQueryParams() []Param {
	var queryParams []Param
	for _, p := range o.Parameters {
		if p.In == "query" {
			queryParams = append(queryParams, p)
		}
	}
	return queryParams
}

// GetPathParams returns only the path parameters.
func (o *Operation) GetPathParams() []Param {
	var pathParams []Param
	for _, p := range o.Parameters {
		if p.In == "path" {
			pathParams = append(pathParams, p)
		}
	}
	return pathParams
}

// GetHeaderParams returns only the header parameters.
func (o *Operation) GetHeaderParams() []Param {
	var headerParams []Param
	for _, p := range o.Parameters {
		if p.In == "header" {
			headerParams = append(headerParams, p)
		}
	}
	return headerParams
}

// HasBearerAuth checks if the operation requires bearer token authentication.
func (o *Operation) HasBearerAuth() bool {
	for _, sec := range o.Security {
		if sec.Type == "http" && strings.EqualFold(sec.Scheme, "bearer") {
			return true
		}
	}
	return false
}

// HasCookieAuth checks if the operation requires cookie-based authentication.
func (o *Operation) HasCookieAuth() bool {
	for _, sec := range o.Security {
		if sec.Type == "apiKey" && strings.EqualFold(sec.Scheme, "cookie") {
			return true
		}
	}
	return false
}

// GetOperations returns all operations in the spec.
func (s *Spec) GetOperations() []Operation {
	return s.Operations
}

// GetSortedOperations returns operations sorted by path and method.
func (s *Spec) GetSortedOperations() []Operation {
	if s == nil || len(s.Operations) == 0 {
		return nil
	}
	ops := make([]Operation, len(s.Operations))
	copy(ops, s.Operations)
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Path != ops[j].Path {
			return ops[i].Path < ops[j].Path
		}
		methodOrder := map[string]int{"GET": 0, "POST": 1, "PUT": 2, "PATCH": 3, "DELETE": 4}
		return methodOrder[ops[i].Method] < methodOrder[ops[j].Method]
	})
	return ops
}
