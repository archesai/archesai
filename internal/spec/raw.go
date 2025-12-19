package spec

import (
	"github.com/archesai/archesai/internal/ref"
	"github.com/archesai/archesai/internal/schema"
)

// RawDocument represents the root openapi.yaml structure for YAML parsing.
type RawDocument struct {
	OpenAPI      string                `yaml:"openapi"`
	Info         RawInfo               `yaml:"info"`
	Paths        map[string]RawPathRef `yaml:"paths"`
	Components   RawComponents         `yaml:"components"`
	Security     []map[string][]string `yaml:"security"`
	Tags         []Tag                 `yaml:"tags"`
	XProjectName string                `yaml:"x-project-name"`
}

// RawInfo represents the info section of an OpenAPI document.
type RawInfo struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// RawComponents represents the components section.
type RawComponents struct {
	Schemas         map[string]ref.Ref[schema.Schema] `yaml:"schemas"`
	Responses       map[string]RawSimpleRef           `yaml:"responses"`
	Parameters      map[string]RawSimpleRef           `yaml:"parameters"`
	Headers         map[string]RawSimpleRef           `yaml:"headers"`
	SecuritySchemes map[string]SecScheme              `yaml:"securitySchemes"`
}

// RawSimpleRef represents a simple $ref structure.
type RawSimpleRef struct {
	Ref string `yaml:"$ref"`
}

// RawPathRef represents a path that can be either a $ref or an inline definition.
type RawPathRef struct {
	Ref        string         `yaml:"$ref,omitempty"`
	Get        *RawOperation  `yaml:"get,omitempty"`
	Post       *RawOperation  `yaml:"post,omitempty"`
	Put        *RawOperation  `yaml:"put,omitempty"`
	Delete     *RawOperation  `yaml:"delete,omitempty"`
	Patch      *RawOperation  `yaml:"patch,omitempty"`
	Parameters []RawParameter `yaml:"parameters,omitempty"`
}

// IsRef returns true if this is a $ref to another path file.
func (p *RawPathRef) IsRef() bool {
	return p.Ref != ""
}

// IsInline returns true if this is an inline path definition.
func (p *RawPathRef) IsInline() bool {
	return p.Ref == "" &&
		(p.Get != nil || p.Post != nil || p.Put != nil || p.Delete != nil || p.Patch != nil)
}

// ToPathItem converts a RawPathRef to a RawPathItem.
func (p *RawPathRef) ToPathItem() *RawPathItem {
	return &RawPathItem{
		Ref:        p.Ref,
		Get:        p.Get,
		Post:       p.Post,
		Put:        p.Put,
		Delete:     p.Delete,
		Patch:      p.Patch,
		Parameters: p.Parameters,
	}
}

// RawPathItem represents a path item (get/post/put/delete/patch).
type RawPathItem struct {
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
	Ref         string                  `yaml:"$ref"`
	Description string                  `yaml:"description"`
	Required    bool                    `yaml:"required"`
	Content     map[string]RawMediaType `yaml:"content"`
}

// RawMediaType represents a media type in request/response.
type RawMediaType struct {
	Schema  *ref.Ref[schema.Schema] `yaml:"schema"`
	Example any                     `yaml:"example"`
}

// RawResponse represents a response (inline or $ref) for YAML parsing.
type RawResponse struct {
	Ref         string                  `yaml:"$ref"`
	Description string                  `yaml:"description"`
	Content     map[string]RawMediaType `yaml:"content"`
	Headers     map[string]RawHeader    `yaml:"headers"`
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

// RawParameterFile represents a parameter file in components/parameters/.
type RawParameterFile struct {
	Name        string              `yaml:"name"`
	In          string              `yaml:"in"`
	Description string              `yaml:"description"`
	Required    bool                `yaml:"required"`
	Schema      *RawParameterSchema `yaml:"schema"`
}

// RawParameterSchema represents the schema within a parameter file.
type RawParameterSchema struct {
	Ref    string `yaml:"$ref"`
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

// RawResponseFile represents a response file in components/responses/.
type RawResponseFile struct {
	Description string                     `yaml:"description"`
	Content     map[string]RawMediaTypeRef `yaml:"content"`
	Headers     map[string]RawSimpleRef    `yaml:"headers"`
}

// RawMediaTypeRef represents a media type with a schema reference.
type RawMediaTypeRef struct {
	Schema RawSimpleRef `yaml:"schema"`
}

// RawHeaderFile represents a header file in components/headers/.
type RawHeaderFile struct {
	Description string                `yaml:"description"`
	Schema      RawHeaderSchemaInline `yaml:"schema"`
	Example     string                `yaml:"example"`
}

// RawHeaderSchemaInline represents the schema within a header file.
type RawHeaderSchemaInline struct {
	Type    string `yaml:"type"`
	Format  string `yaml:"format"`
	Minimum *int64 `yaml:"minimum"`
	Maximum *int64 `yaml:"maximum"`
	Example any    `yaml:"example"`
}

// RawOpenAPIInfo represents the tags section of an openapi.yaml file.
type RawOpenAPIInfo struct {
	Tags []Tag `yaml:"tags"`
}

// ConvertFileRefs converts file $refs to internal refs in the response.
func (r *RawResponseFile) ConvertFileRefs() {
	if r == nil {
		return
	}
	for contentType, mediaType := range r.Content {
		if mediaType.Schema.Ref != "" {
			mediaType.Schema.Ref = FileRefToInternalRef(mediaType.Schema.Ref, "")
			r.Content[contentType] = mediaType
		}
	}
	for headerName, header := range r.Headers {
		if header.Ref != "" {
			header.Ref = FileRefToInternalRef(header.Ref, "")
			r.Headers[headerName] = header
		}
	}
}

// ConvertFileRefs converts file $refs to internal refs in the parameter.
func (p *RawParameterFile) ConvertFileRefs() {
	if p == nil {
		return
	}
	if p.Schema != nil && p.Schema.Ref != "" {
		p.Schema.Ref = FileRefToInternalRef(p.Schema.Ref, "")
	}
}
