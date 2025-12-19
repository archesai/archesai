package spec

import (
	"sort"
	"strconv"
	"strings"

	"github.com/archesai/archesai/internal/schema"
)

// HTTP method constants.
const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodPATCH  = "PATCH"
	methodDELETE = "DELETE"
)

// Parameter location constants.
const paramLocationPath = "path"

// InternalContext returns the last segment of the project name.
// This is used to determine which operations/schemas belong to this package.
func (s *Spec) InternalContext() string {
	if s.ProjectName == "" {
		return ""
	}
	parts := strings.Split(s.ProjectName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// OwnOperations returns operations that belong to this package, sorted by ID.
// An operation belongs to this package if x-internal is empty or matches InternalContext.
func (s *Spec) OwnOperations() []Operation {
	internalContext := s.InternalContext()
	var operations []Operation
	for _, op := range s.Operations {
		if op.Internal == "" || op.Internal == internalContext {
			operations = append(operations, op)
		}
	}
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].ID < operations[j].ID
	})
	return operations
}

// ComposedPackages returns the unique x-internal package names from operations
// that belong to OTHER packages (not this one).
func (s *Spec) ComposedPackages() []string {
	internalContext := s.InternalContext()
	pkgMap := make(map[string]bool)
	for _, op := range s.Operations {
		if op.Internal != "" && op.Internal != internalContext {
			pkgMap[op.Internal] = true
		}
	}

	var packages []string
	for pkg := range pkgMap {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	return packages
}

// IsStandalone returns true if this project has its own entities or operations.
func (s *Spec) IsStandalone() bool {
	return len(s.OwnOperations()) > 0 || len(s.OwnEntitySchemas()) > 0
}

// IsComposition returns true if this project composes other packages.
func (s *Spec) IsComposition() bool {
	return len(s.ComposedPackages()) > 0
}

// NeedsPublisher returns true if any non-GET, non-custom operation exists.
func (s *Spec) NeedsPublisher() bool {
	for _, op := range s.OwnOperations() {
		if !op.CustomHandler && op.Method != methodGET {
			return true
		}
	}
	return false
}

// Repositories returns the unique repository names (tags) from own operations.
func (s *Spec) Repositories() []string {
	repoMap := make(map[string]bool)
	for _, op := range s.OwnOperations() {
		if !op.CustomHandler && op.Tag != "" {
			repoMap[op.Tag] = true
		}
	}

	var repos []string
	for repo := range repoMap {
		repos = append(repos, repo)
	}
	sort.Strings(repos)
	return repos
}

// ShouldSkipWire returns true if there's nothing to generate for wire.
// This includes packages that only have custom handlers (no entities, no repositories).
func (s *Spec) ShouldSkipWire() bool {
	if s.IsComposition() {
		return false // Composition apps need wire
	}
	// Skip if no entities and no repositories (only custom handlers)
	if len(s.OwnEntitySchemas()) == 0 && len(s.Repositories()) == 0 {
		return true
	}
	return !s.IsStandalone()
}

// InternalPackage represents an internal package for composition.
type InternalPackage struct {
	Name           string
	Alias          string
	ImportPath     string
	Repositories   []string // Repository names (tags) for this package
	NeedsPublisher bool     // Whether this package has non-GET operations
}

// InternalPackageBase is the base import path for internal packages.
const InternalPackageBase = "github.com/archesai/archesai/pkg"

// InternalPackages returns internal package info for composed packages.
func (s *Spec) InternalPackages() []InternalPackage {
	pkgNames := s.ComposedPackages()

	// Group repositories and publisher needs by package
	reposByPkg := make(map[string]map[string]bool)
	pkgNeedsPublisher := make(map[string]bool)

	for _, op := range s.Operations {
		if op.Internal == "" || op.CustomHandler {
			continue
		}
		if reposByPkg[op.Internal] == nil {
			reposByPkg[op.Internal] = make(map[string]bool)
		}
		if op.Tag != "" {
			reposByPkg[op.Internal][op.Tag] = true
		}
		if op.Method != methodGET {
			pkgNeedsPublisher[op.Internal] = true
		}
	}

	var packages []InternalPackage
	for _, name := range pkgNames {
		var repos []string
		for repo := range reposByPkg[name] {
			repos = append(repos, repo)
		}
		sort.Strings(repos)

		packages = append(packages, InternalPackage{
			Name:           name,
			Alias:          name,
			ImportPath:     InternalPackageBase + "/" + name,
			Repositories:   repos,
			NeedsPublisher: pkgNeedsPublisher[name],
		})
	}
	return packages
}

// OwnEntitySchemas returns entity schemas that belong to this package, sorted by name.
// Filters by SchemaType == "entity" and excludes internal schemas.
func (s *Spec) OwnEntitySchemas() []*schema.Schema {
	internalContext := s.InternalContext()
	var entities []*schema.Schema
	for _, sch := range s.Schemas {
		if sch.XCodegenSchemaType != schema.TypeEntity {
			continue
		}
		if sch.IsInternal(internalContext) {
			continue
		}
		entities = append(entities, sch)
	}
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Title < entities[j].Title
	})
	return entities
}

// AllEntitySchemas returns all entity schemas regardless of x-internal, sorted by name.
// Used by database generators that need to generate repositories
// for all entities including those from included packages.
func (s *Spec) AllEntitySchemas() []*schema.Schema {
	var entities []*schema.Schema
	for _, sch := range s.Schemas {
		if sch.XCodegenSchemaType != schema.TypeEntity {
			continue
		}
		entities = append(entities, sch)
	}
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Title < entities[j].Title
	})
	return entities
}

// Spec represents the entire specification definition.
type Spec struct {
	Operations      []Operation               // All operations in the spec
	Schemas         map[string]*schema.Schema // All schemas defined in the spec (keyed by ref path or name)
	ProjectName     string                    // Project name from x-project-name extension
	EnabledIncludes []string                  // Names of enabled x-include-* extensions

	// OpenAPI metadata for generation
	Title       string
	Description string
	Version     string
	Tags        []Tag
	Security    map[string]SecScheme

	// Codegen options from x-codegen-* extensions
	CodegenOnly   []string // Generators to run (empty = all)
	CodegenLint   bool     // Enable strict linting
	CodegenOutput string   // Output directory for tag-based generators (empty = use tag subdirectories)
}

// GetSchema returns a schema by name from the registry.
func (s *Spec) GetSchema(name string) *schema.Schema {
	if s == nil || s.Schemas == nil {
		return nil
	}
	return s.Schemas[name]
}

// GetSortedSchemas returns schemas sorted by name for consistent iteration.
func (s *Spec) GetSortedSchemas() []*schema.Schema {
	if s == nil || s.Schemas == nil {
		return nil
	}
	var names []string
	for name := range s.Schemas {
		names = append(names, name)
	}
	sort.Strings(names)
	var sorted []*schema.Schema
	for _, name := range names {
		sorted = append(sorted, s.Schemas[name])
	}
	return sorted
}

// SchemasSlice returns schemas as a slice for backwards compatibility.
func (s *Spec) SchemasSlice() []*schema.Schema {
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
	*schema.Schema        // Embed schema definition
	In             string // Location (path, query, header, cookie)
	Style          string // Parameter style (form, simple, etc.)
	Explode        bool   // Whether to explode array/object parameters
}

// RequestBody represents the request body definition for an API operation.
type RequestBody struct {
	*schema.Schema      // Embed schema definition for request body
	Required       bool // Whether request body is required
}

// Response represents a response in an operation.
type Response struct {
	*schema.Schema                           // Embed schema definition for response body
	StatusCode     string                    // HTTP status code
	ContentType    string                    // Content-Type for the response
	Headers        map[string]*schema.Schema // Response headers
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
	Schema *schema.Schema
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
		Schema *schema.Schema
	}
	for _, name := range names {
		sorted = append(sorted, struct {
			Name   string
			Schema *schema.Schema
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
		if p.In == paramLocationPath {
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

// IsVoidExecutor returns true if the operation's success response is 204 No Content.
func (o *Operation) IsVoidExecutor() bool {
	successResp := o.GetSuccessResponse()
	return successResp != nil && successResp.StatusCode == "204"
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
		methodOrder := map[string]int{
			methodGET:    0,
			methodPOST:   1,
			methodPUT:    2,
			methodPATCH:  3,
			methodDELETE: 4,
		}
		return methodOrder[ops[i].Method] < methodOrder[ops[j].Method]
	})
	return ops
}

// TagGroup represents all operations and schemas belonging to a single tag.
type TagGroup struct {
	Name       string           // Tag name (e.g., "Pipeline", "User")
	Package    string           // Go package name (lowercase, e.g., "pipeline", "user")
	Operations []Operation      // Operations in this tag
	Entities   []*schema.Schema // Entity schemas (x-codegen-schema-type: entity)
}

// GetOperationsByTag returns operations grouped by their primary tag.
func (s *Spec) GetOperationsByTag() []TagGroup {
	if s == nil || len(s.Operations) == 0 {
		return nil
	}

	// Group operations by tag
	tagOps := make(map[string][]Operation)
	for _, op := range s.Operations {
		tag := op.Tag
		if tag == "" {
			tag = "Default"
		}
		tagOps[tag] = append(tagOps[tag], op)
	}

	// Build tag groups
	var groups []TagGroup
	for tag, ops := range tagOps {
		// Sort operations within each tag
		sort.Slice(ops, func(i, j int) bool {
			if ops[i].Path != ops[j].Path {
				return ops[i].Path < ops[j].Path
			}
			methodOrder := map[string]int{
				methodGET:    0,
				methodPOST:   1,
				methodPUT:    2,
				methodPATCH:  3,
				methodDELETE: 4,
			}
			return methodOrder[ops[i].Method] < methodOrder[ops[j].Method]
		})

		// Find entities for this tag by matching schema title to tag name
		// e.g., "Pipeline" schema belongs to "Pipeline" tag
		var entities []*schema.Schema
		for _, sch := range s.Schemas {
			if sch.IsEntity() && strings.EqualFold(sch.Title, tag) {
				entities = append(entities, sch)
			}
		}
		sort.Slice(entities, func(i, j int) bool {
			return entities[i].Title < entities[j].Title
		})

		groups = append(groups, TagGroup{
			Name:       tag,
			Package:    strings.ToLower(tag),
			Operations: ops,
			Entities:   entities,
		})
	}

	// Sort groups by name for consistent output
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups
}
