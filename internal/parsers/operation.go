package parsers

import "strings"

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
	XInternal             string          // When set (e.g., "server", "config"), this operation should be imported not generated
}

// IsInternal returns true if this operation should be imported from another package instead of generated.
// If context is empty, it returns true whenever XInternal is set.
// If context is provided, it returns true only when XInternal is set AND doesn't match the context.
func (o *OperationDef) IsInternal(context string) bool {
	if o.XInternal == "" {
		return false
	}
	// If no context provided, treat all internal operations as internal
	if context == "" {
		return true
	}
	// Only internal if the internal tag doesn't match the current context
	return o.XInternal != context
}

// NeedsServerModels returns true if this operation references any types from the server package.
func (o *OperationDef) NeedsServerModels() bool {
	// Check responses
	for _, resp := range o.Responses {
		if resp.SchemaDef != nil && resp.NeedsServerModels() {
			return true
		}
	}
	// Check request body
	if o.RequestBody != nil && o.RequestBody.SchemaDef != nil && o.RequestBody.NeedsServerModels() {
		return true
	}
	// Check parameters
	for _, param := range o.Parameters {
		if param.SchemaDef != nil && param.NeedsServerModels() {
			return true
		}
	}
	return false
}

// ParamDef represents a parameter in an operation
type ParamDef struct {
	*SchemaDef        // Embed schema definition
	In         string // Location (path, query, header, cookie)
	Style      string // Parameter style (form, simple, etc.)
	Explode    bool   // Whether to explode array/object parameters
}

// RequestBodyDef represents the request body definition for an API operation
type RequestBodyDef struct {
	*SchemaDef      // Embed schema definition for request body
	Required   bool // Whether request body is required
}

// SecurityDef represents a security requirement
type SecurityDef struct {
	Name   string   // Security scheme name
	Type   string   // Security type (http, apiKey, oauth2)
	Scheme string   // Security scheme (bearer for http, cookie for apiKey)
	Scopes []string // Required scopes
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

// NeedsUUID returns true if the operation requires the UUID package.
// This is true if the operation uses bearer/cookie auth or has UUID path parameters.
func (o *OperationDef) NeedsUUID() bool {
	if o.HasBearerAuth() || o.HasCookieAuth() {
		return true
	}
	for _, p := range o.GetPathParams() {
		if p.GoType == goTypeUUID {
			return true
		}
	}
	return false
}
