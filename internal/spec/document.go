package spec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/archesai/archesai/pkg/auth"
	"github.com/archesai/archesai/pkg/config"
	"github.com/archesai/archesai/pkg/executor"
	"github.com/archesai/archesai/pkg/pipelines"
	"github.com/archesai/archesai/pkg/server"
	"github.com/archesai/archesai/pkg/storage"
)

// OpenAPIDocument wraps the raw YAML representation of an OpenAPI spec.
// This is the single source of truth - all operations query this directly.
type OpenAPIDocument struct {
	root           *yaml.Node        // Root document node (for bundling)
	doc            *Document         // Structured document (for config access)
	fsys           fs.FS             // Filesystem for resolving $refs
	baseDir        string            // Base directory for relative refs (absolute path)
	refCache       map[string][]byte // Cache for resolved file refs
	configIncludes []string          // Includes from arches.yaml config (takes precedence over x-include-*)
}

// Document represents the root openapi.yaml structure for YAML parsing.
type Document struct {
	OpenAPI      string                `yaml:"openapi"`
	Info         Info                  `yaml:"info"`
	Paths        map[string]PathRef    `yaml:"paths"`
	Components   Components            `yaml:"components"`
	Security     []map[string][]string `yaml:"security"`
	Tags         []Tag                 `yaml:"tags"`
	XProjectName string                `yaml:"x-project-name"`
}

// Info represents the info section of an OpenAPI document.
type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// Components represents the components section.
type Components struct {
	Schemas         map[string]Ref[Schema] `yaml:"schemas"`
	Responses       map[string]SimpleRef   `yaml:"responses"`
	Parameters      map[string]SimpleRef   `yaml:"parameters"`
	Headers         map[string]SimpleRef   `yaml:"headers"`
	SecuritySchemes map[string]SecScheme   `yaml:"securitySchemes"`
}

// SimpleRef represents a simple $ref structure.
type SimpleRef struct {
	Ref string `yaml:"$ref"`
}

// PathRef represents a path that can be either a $ref or an inline definition.
type PathRef struct {
	Ref        string         `yaml:"$ref,omitempty"`
	Get        *RawOperation  `yaml:"get,omitempty"`
	Post       *RawOperation  `yaml:"post,omitempty"`
	Put        *RawOperation  `yaml:"put,omitempty"`
	Delete     *RawOperation  `yaml:"delete,omitempty"`
	Patch      *RawOperation  `yaml:"patch,omitempty"`
	Parameters []RawParameter `yaml:"parameters,omitempty"`
}

// IsRef returns true if this is a $ref to another path file.
func (p *PathRef) IsRef() bool {
	return p.Ref != ""
}

// IsInline returns true if this is an inline path definition.
func (p *PathRef) IsInline() bool {
	return p.Ref == "" &&
		(p.Get != nil || p.Post != nil || p.Put != nil || p.Delete != nil || p.Patch != nil)
}

// ToPathItem converts a PathRef to a PathItem.
func (p *PathRef) ToPathItem() *PathItem {
	return &PathItem{
		Ref:        p.Ref,
		Get:        p.Get,
		Post:       p.Post,
		Put:        p.Put,
		Delete:     p.Delete,
		Patch:      p.Patch,
		Parameters: p.Parameters,
	}
}

// PathItem represents a path item (get/post/put/delete/patch).
type PathItem struct {
	Ref        string         `yaml:"$ref"`
	XPath      string         `yaml:"x-path"` // Required when auto-discovered from paths/ directory
	Get        *RawOperation  `yaml:"get"`
	Post       *RawOperation  `yaml:"post"`
	Put        *RawOperation  `yaml:"put"`
	Delete     *RawOperation  `yaml:"delete"`
	Patch      *RawOperation  `yaml:"patch"`
	Parameters []RawParameter `yaml:"parameters"`
}

// RawOperation represents an operation within a path item (raw YAML parsing type).
// This is the YAML structure; use Operation for the domain model.
type RawOperation struct {
	OperationID     string                 `yaml:"operationId"`
	Summary         string                 `yaml:"summary"`
	Description     string                 `yaml:"description"`
	Tags            []string               `yaml:"tags"`
	Security        []map[string][]string  `yaml:"security"`
	Parameters      []RawParameter         `yaml:"parameters"`
	RequestBody     *RawRequestBody        `yaml:"requestBody"`
	Responses       map[string]RawResponse `yaml:"responses"`
	XInternal       string                 `yaml:"x-internal"`
	XCustomHandler  bool                   `yaml:"x-codegen-custom-handler"`
	XPublicEndpoint bool                   `yaml:"x-public-endpoint"`
}

// RawParameter represents a parameter (inline or $ref) for YAML parsing.
type RawParameter struct {
	Ref         string           `yaml:"$ref"`
	Name        string           `yaml:"name"`
	In          string           `yaml:"in"`
	Description string           `yaml:"description"`
	Required    bool             `yaml:"required"`
	Style       string           `yaml:"style"`
	Explode     *bool            `yaml:"explode"`
	Schema      *RawSchemaInline `yaml:"schema"`
}

// RawRequestBody represents a request body (inline or $ref) for YAML parsing.
type RawRequestBody struct {
	Ref         string               `yaml:"$ref"`
	Description string               `yaml:"description"`
	Required    bool                 `yaml:"required"`
	Content     map[string]MediaType `yaml:"content"`
}

// MediaType represents a media type in request/response.
type MediaType struct {
	Schema  *Ref[Schema] `yaml:"schema"`
	Example any          `yaml:"example"`
}

// RawResponse represents a response (inline or $ref) for YAML parsing.
type RawResponse struct {
	Ref         string               `yaml:"$ref"`
	Description string               `yaml:"description"`
	Content     map[string]MediaType `yaml:"content"`
	Headers     map[string]RawHeader `yaml:"headers"`
}

// RawHeader represents a header (inline or $ref) for YAML parsing.
type RawHeader struct {
	Ref         string           `yaml:"$ref"`
	Description string           `yaml:"description"`
	Schema      *RawSchemaInline `yaml:"schema"`
}

// RawSchemaInline represents an inline schema definition (used in parameters).
type RawSchemaInline struct {
	Ref    string `yaml:"$ref"`
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

// SchemaFile represents a schema file in components/schemas/.
type SchemaFile struct {
	Title                 string                  `yaml:"title"`
	Description           string                  `yaml:"description"`
	Type                  PropertyType            `yaml:"type"`
	Properties            map[string]*Ref[Schema] `yaml:"properties"`
	Required              []string                `yaml:"required"`
	AllOf                 []*Ref[Schema]          `yaml:"allOf"`
	XCodegenSchemaType    string                  `yaml:"x-codegen-schema-type"`
	XCodegen              *XCodegenExtension      `yaml:"x-codegen"`
	XInternal             string                  `yaml:"x-internal"`
	UnevaluatedProperties *bool                   `yaml:"unevaluatedProperties"`
	AdditionalProperties  *bool                   `yaml:"additionalProperties"`
}

// ParameterFile represents a parameter file in components/parameters/.
type ParameterFile struct {
	Name        string           `yaml:"name"`
	In          string           `yaml:"in"`
	Description string           `yaml:"description"`
	Required    bool             `yaml:"required"`
	Schema      *ParameterSchema `yaml:"schema"`
}

// ParameterSchema represents the schema within a parameter file.
type ParameterSchema struct {
	Ref    string `yaml:"$ref"`
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

// ResponseFile represents a response file in components/responses/.
type ResponseFile struct {
	Description string                  `yaml:"description"`
	Content     map[string]MediaTypeRef `yaml:"content"`
	Headers     map[string]SimpleRef    `yaml:"headers"`
}

// MediaTypeRef represents a media type with a schema reference.
type MediaTypeRef struct {
	Schema SimpleRef `yaml:"schema"`
}

// HeaderFile represents a header file in components/headers/.
type HeaderFile struct {
	Description string             `yaml:"description"`
	Schema      HeaderSchemaInline `yaml:"schema"`
	Example     string             `yaml:"example"`
}

// HeaderSchemaInline represents the schema within a header file.
type HeaderSchemaInline struct {
	Type    string `yaml:"type"`
	Format  string `yaml:"format"`
	Minimum *int64 `yaml:"minimum"`
	Maximum *int64 `yaml:"maximum"`
	Example any    `yaml:"example"`
}

// OpenAPIInfo represents the tags section of an openapi.yaml file.
type OpenAPIInfo struct {
	Tags []Tag `yaml:"tags"`
}

// NewOpenAPIDocument loads an OpenAPI specification from a file path.
// The path should point directly to the openapi.yaml file.
func NewOpenAPIDocument(openapiPath string) (*OpenAPIDocument, error) {
	return NewOpenAPIDocumentWithIncludes(openapiPath, nil)
}

// NewOpenAPIDocumentWithIncludes loads an OpenAPI specification with explicit includes.
// If includes is non-empty, it takes precedence over x-include-* extensions in the spec.
func NewOpenAPIDocumentWithIncludes(
	openapiPath string,
	includes []string,
) (*OpenAPIDocument, error) {
	absPath, err := filepath.Abs(openapiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	baseDir := filepath.Dir(absPath)
	baseFS := os.DirFS(baseDir)

	// Load the root document
	data, err := fs.ReadFile(baseFS, filepath.Base(absPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI document: %w", err)
	}

	// Parse into yaml.Node (for bundling)
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	// Parse into Document struct (for config access)
	var parsedDoc Document
	if err := yaml.Unmarshal(data, &parsedDoc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	d := &OpenAPIDocument{
		root:           &root,
		doc:            &parsedDoc,
		baseDir:        baseDir,
		fsys:           baseFS,
		refCache:       make(map[string][]byte),
		configIncludes: includes,
	}
	// Build composite filesystem with includes
	d.fsys = d.buildIncludeFS(baseFS)

	return d, nil
}

// WithFS sets an alternative filesystem for reading files.
func (d *OpenAPIDocument) WithFS(fsys fs.FS) *OpenAPIDocument {
	d.fsys = fsys
	return d
}

// Node returns the underlying yaml.Node for direct access.
func (d *OpenAPIDocument) Node() *yaml.Node {
	return d.root
}

// FS returns the filesystem used by the document.
func (d *OpenAPIDocument) FS() fs.FS {
	return d.fsys
}

// Includes returns enabled include names from arches.yaml config.
func (d *OpenAPIDocument) Includes() []string {
	return d.configIncludes
}

// Info returns info section values.
func (d *OpenAPIDocument) Info() (title, description, version string) {
	return d.doc.Info.Title, d.doc.Info.Description, d.doc.Info.Version
}

// ProjectName returns x-project-name value.
func (d *OpenAPIDocument) ProjectName() string {
	return d.doc.XProjectName
}

// Tags returns all tags from the document.
func (d *OpenAPIDocument) Tags() []Tag {
	return d.doc.Tags
}

// Security returns security schemes from components.
func (d *OpenAPIDocument) Security() map[string]SecScheme {
	return d.doc.Components.SecuritySchemes
}

// ResolveFileRef resolves a file $ref and returns its content.
// Results are cached for efficiency.
func (d *OpenAPIDocument) ResolveFileRef(fromDir, refPath string) ([]byte, error) {
	cacheKey := fromDir + ":" + refPath
	if data, ok := d.refCache[cacheKey]; ok {
		return data, nil
	}

	resolver := NewResolver(d.fsys, fromDir)
	data, err := resolver.ResolveFile(refPath)
	if err != nil {
		return nil, err
	}

	d.refCache[cacheKey] = data
	return data, nil
}

// buildIncludeFS creates a composite filesystem with include filesystems as base layers.
func (d *OpenAPIDocument) buildIncludeFS(projectFS fs.FS) fs.FS {
	includes := d.Includes()
	if len(includes) == 0 {
		return projectFS
	}

	var layers []fs.FS
	for _, include := range includes {
		includeFS := d.getIncludeFS(include)
		if includeFS != nil {
			layers = append(layers, includeFS)
		}
	}

	// Project FS is the top layer (overrides includes)
	layers = append(layers, projectFS)
	return NewCompositeFS(layers...)
}

// getIncludeFS returns the embedded filesystem for an include package.
func (d *OpenAPIDocument) getIncludeFS(include string) fs.FS {
	var embedFS fs.FS
	switch include {
	case "server":
		embedFS = server.API
	case "auth":
		embedFS = auth.API
	case "config":
		embedFS = config.API
	case "pipelines":
		embedFS = pipelines.API
	case "executor":
		embedFS = executor.API
	case "storage":
		embedFS = storage.API
	default:
		return nil
	}

	// Strip the "api/" prefix from the embedded filesystem
	subFS, err := fs.Sub(embedFS, "api")
	if err != nil {
		return nil
	}
	return subFS
}

// HasInclude checks if a specific include is enabled in arches.yaml config.
func (d *OpenAPIDocument) HasInclude(name string) bool {
	for _, inc := range d.configIncludes {
		if inc == name {
			return true
		}
	}
	return false
}

// contentNode returns the content node of the root document.
// Used by bundler for low-level yaml.Node access.
func (d *OpenAPIDocument) contentNode() *yaml.Node {
	if d.root == nil {
		return nil
	}
	// The root node is typically a DocumentNode containing the actual content
	if d.root.Kind == yaml.DocumentNode && len(d.root.Content) > 0 {
		return d.root.Content[0]
	}
	return d.root
}
