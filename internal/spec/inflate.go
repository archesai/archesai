package spec

import (
	"maps"
	"strings"

	"github.com/archesai/archesai/internal/ref"
	"github.com/archesai/archesai/internal/schema"
	"github.com/archesai/archesai/internal/strutil"
)

// InflateConfig controls what gets generated during inflation.
type InflateConfig struct {
	ResponseWrappers bool // Generate {Entity}Response and {Entity}ListResponse
	ListParams       bool // Generate filter/sort/page params for list operations
}

// DefaultInflateConfig returns the default inflation configuration.
func DefaultInflateConfig() InflateConfig {
	return InflateConfig{
		ResponseWrappers: true,
		ListParams:       true,
	}
}

// Inflated contains the generated content from inflation.
type Inflated struct {
	ResponseSchemas map[string]*schema.Schema // Generated response wrappers
	ListParams      map[string][]Param        // Generated params per operation ID
}

// Inflater generates additional schemas and params without mutating the input.
type Inflater struct {
	config InflateConfig
}

// NewInflater creates a new Inflater with the given configuration.
func NewInflater(config InflateConfig) *Inflater {
	return &Inflater{config: config}
}

// Inflate generates additional content based on the schemas and operations.
// Returns new content that can be merged with the spec - does not mutate inputs.
func (i *Inflater) Inflate(schemas map[string]*schema.Schema, operations []Operation) *Inflated {
	result := &Inflated{
		ResponseSchemas: make(map[string]*schema.Schema),
		ListParams:      make(map[string][]Param),
	}

	if i.config.ResponseWrappers {
		i.generateResponseSchemas(schemas, result)
	}

	if i.config.ListParams {
		// Combine original schemas with generated response schemas for lookup
		allSchemas := make(map[string]*schema.Schema, len(schemas)+len(result.ResponseSchemas))
		maps.Copy(allSchemas, schemas)
		maps.Copy(allSchemas, result.ResponseSchemas)
		i.generateListParams(operations, allSchemas, result)
	}

	return result
}

// generateResponseSchemas creates Response and ListResponse wrappers for entity schemas.
func (i *Inflater) generateResponseSchemas(schemas map[string]*schema.Schema, result *Inflated) {
	for _, s := range schemas {
		if s.XCodegenSchemaType == schema.TypeEntity {
			singleResp := i.buildSingleResponse(s)
			listResp := i.buildListResponse(s)
			result.ResponseSchemas[singleResp.Title] = singleResp
			result.ResponseSchemas[listResp.Title] = listResp
		}
	}
}

// generateListParams generates filter/sort/page params for list operations.
func (i *Inflater) generateListParams(
	operations []Operation,
	schemas map[string]*schema.Schema,
	result *Inflated,
) {
	for _, op := range operations {
		params := i.generateParamsForOperation(&op, schemas)
		if len(params) > 0 {
			result.ListParams[op.ID] = params
		}
	}
}

// generateParamsForOperation generates params for a single operation if it's a list operation.
func (i *Inflater) generateParamsForOperation(
	op *Operation,
	schemas map[string]*schema.Schema,
) []Param {
	// Find the response schema reference from processed responses
	var responseSchema *schema.Schema
	var responseType string

	for _, resp := range op.Responses {
		if resp.StatusCode == "200" && resp.Schema != nil {
			responseSchema = resp.Schema
			if responseSchema.Title != "" {
				if strings.HasSuffix(responseSchema.Title, "ListResponse") {
					responseType = "list"
				} else if strings.HasSuffix(responseSchema.Title, "Response") {
					responseType = "single"
				}
			}
			break
		}
	}

	if responseSchema == nil || responseType != "list" {
		return nil
	}

	// Find the underlying entity schema
	var entitySchema *schema.Schema
	if dataRef, ok := responseSchema.Properties["Data"]; ok && dataRef != nil {
		dataSchema := dataRef.GetOrNil()
		if dataSchema != nil {
			if dataSchema.Type.PrimaryType() == schema.TypeArray && dataSchema.Items != nil {
				itemsSchema := dataSchema.Items.GetOrNil()
				if itemsSchema != nil {
					entityName := itemsSchema.Title
					if entityName == "" && itemsSchema.GoType != "" {
						entityName = stripPackagePrefix(itemsSchema.GoType)
					}
					entitySchema = schemas[entityName]
				}
			}
		}
	}

	if entitySchema == nil {
		return nil
	}

	// Generate params, skipping any that already exist
	var params []Param
	if !hasParamNamed(op.Parameters, "filter") {
		params = append(params, buildFilterParam(entitySchema))
	}
	if !hasParamNamed(op.Parameters, "sort") {
		params = append(params, buildSortParam(entitySchema))
	}
	if !hasParamNamed(op.Parameters, "page") {
		params = append(params, buildPageParam())
	}

	return params
}

// buildSingleResponse creates a single response wrapper schema for an entity.
func (i *Inflater) buildSingleResponse(entitySchema *schema.Schema) *schema.Schema {
	responseName := entitySchema.Title + "Response"
	return &schema.Schema{
		Title: responseName,
		Type:  schema.PropertyType{Types: []string{schema.TypeObject}},
		Properties: map[string]*ref.Ref[schema.Schema]{
			"Data": ref.NewInline(&schema.Schema{
				Title:   "Data",
				Type:    schema.PropertyType{Types: []string{schema.TypeObject}},
				GoType:  "schemas." + entitySchema.Title,
				JSONTag: "data",
				YAMLTag: "data",
			}),
		},
		Required: []string{"data"},
		GoType:   responseName,
	}
}

// buildListResponse creates a list response wrapper schema for an entity.
func (i *Inflater) buildListResponse(entitySchema *schema.Schema) *schema.Schema {
	responseName := entitySchema.Title + "ListResponse"
	return &schema.Schema{
		Title: responseName,
		Type:  schema.PropertyType{Types: []string{schema.TypeObject}},
		Properties: map[string]*ref.Ref[schema.Schema]{
			"Data": ref.NewInline(&schema.Schema{
				Title:   "Data",
				Type:    schema.PropertyType{Types: []string{schema.TypeArray}},
				Items:   ref.NewInline(entitySchema),
				GoType:  "[]schemas." + entitySchema.Title,
				JSONTag: "data",
				YAMLTag: "data",
			}),
			"Meta": ref.NewInline(&schema.Schema{
				Title:   "Meta",
				Type:    schema.PropertyType{Types: []string{schema.TypeObject}},
				GoType:  "serverschemas.PaginationMeta",
				JSONTag: "meta",
				YAMLTag: "meta",
			}),
		},
		Required: []string{"data", "meta"},
		GoType:   responseName,
	}
}

// buildFilterParam creates a filter parameter for the given schema.
func buildFilterParam(s *schema.Schema) Param {
	filterName := strutil.Pluralize(s.Title) + "Filter"
	return Param{
		Schema: &schema.Schema{
			Title:   filterName,
			Type:    schema.PropertyType{Types: []string{schema.TypeObject}},
			GoType:  "*serverschemas.FilterNode",
			JSONTag: strutil.CamelCase(filterName),
			YAMLTag: strutil.CamelCase(filterName),
		},
		In:    "query",
		Style: "deepObject",
	}
}

// buildSortParam creates a sort parameter for the given schema.
func buildSortParam(s *schema.Schema) Param {
	sortName := strutil.Pluralize(s.Title) + "Sort"
	return Param{
		Schema: &schema.Schema{
			Title:   sortName,
			Type:    schema.PropertyType{Types: []string{schema.TypeObject}},
			GoType:  "*serverschemas.FilterNode",
			JSONTag: strutil.CamelCase(sortName),
			YAMLTag: strutil.CamelCase(sortName),
		},
		In:    "query",
		Style: "deepObject",
	}
}

// buildPageParam creates a pagination parameter.
func buildPageParam() Param {
	return Param{
		Schema: &schema.Schema{
			Title:   "Page",
			Type:    schema.PropertyType{Types: []string{schema.TypeObject}},
			GoType:  "serverschemas.Page",
			JSONTag: "page",
			YAMLTag: "page",
		},
		In:    "query",
		Style: "deepObject",
	}
}

// Helper functions

func hasParamNamed(params []Param, name string) bool {
	for _, p := range params {
		if p.Schema != nil && (p.JSONTag == name || p.Title == name) {
			return true
		}
	}
	return false
}

func stripPackagePrefix(goType string) string {
	if idx := strings.LastIndex(goType, "."); idx >= 0 {
		return goType[idx+1:]
	}
	return goType
}

// ProblemSchema returns the RFC 7807 Problem schema for error responses.
func ProblemSchema() *schema.Schema {
	return &schema.Schema{
		Title: "Problem",
		Type:  schema.PropertyType{Types: []string{schema.TypeObject}},
		Properties: map[string]*ref.Ref[schema.Schema]{
			"Type": ref.NewInline(&schema.Schema{
				Title:   "Type",
				Type:    schema.PropertyType{Types: []string{schema.TypeString}},
				GoType:  schema.GoTypeString,
				JSONTag: "type",
			}),
			"Title": ref.NewInline(&schema.Schema{
				Title:   "Title",
				Type:    schema.PropertyType{Types: []string{schema.TypeString}},
				GoType:  schema.GoTypeString,
				JSONTag: "title",
			}),
			"Status": ref.NewInline(&schema.Schema{
				Title:   "Status",
				Type:    schema.PropertyType{Types: []string{schema.TypeInteger}},
				GoType:  schema.GoTypeInt,
				JSONTag: "status",
			}),
			"Detail": ref.NewInline(&schema.Schema{
				Title:   "Detail",
				Type:    schema.PropertyType{Types: []string{schema.TypeString}},
				GoType:  schema.GoTypeString,
				JSONTag: "detail",
			}),
		},
		GoType: "Problem",
	}
}

// StandardErrorResponses returns the standard error responses for an operation.
// hasResourceID determines whether to include a 404 response.
func StandardErrorResponses(hasResourceID bool) []Response {
	responses := []Response{
		{StatusCode: "400", ContentType: "application/problem+json", Schema: ProblemSchema()},
		{StatusCode: "401", ContentType: "application/problem+json", Schema: ProblemSchema()},
	}
	// Only add 404 for operations with a resource ID (GET/PATCH/DELETE by ID)
	if hasResourceID {
		responses = append(
			responses,
			Response{
				StatusCode:  "404",
				ContentType: "application/problem+json",
				Schema:      ProblemSchema(),
			},
		)
	}
	responses = append(
		responses,
		Response{
			StatusCode:  "422",
			ContentType: "application/problem+json",
			Schema:      ProblemSchema(),
		},
		Response{
			StatusCode:  "429",
			ContentType: "application/problem+json",
			Schema:      ProblemSchema(),
		},
		Response{
			StatusCode:  "500",
			ContentType: "application/problem+json",
			Schema:      ProblemSchema(),
		},
	)
	return responses
}
