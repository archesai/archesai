package spec

import (
	"github.com/archesai/archesai/internal/strutil"
)

// AutoGenerator handles auto-generation of filters, sorts, pagination, and responses
type AutoGenerator struct {
	schemas map[string]*Schema
}

// NewAutoGenerator creates a new AutoGenerator instance
func NewAutoGenerator(schemas map[string]*Schema) *AutoGenerator {
	return &AutoGenerator{
		schemas: schemas,
	}
}

// GenerateForOperations adds auto-generated parameters and responses to operations
func (g *AutoGenerator) GenerateForOperations(operations []Operation) {
	for i := range operations {
		op := &operations[i]
		g.generateForOperation(op)
	}
}

// generateForOperation handles auto-generation for a single operation
func (g *AutoGenerator) generateForOperation(op *Operation) {
	// Find the response schema reference from processed responses
	var responseSchema *Schema
	var responseType string

	for _, resp := range op.Responses {
		if resp.StatusCode == "200" && resp.Schema != nil {
			responseSchema = resp.Schema
			// Determine type from schema name
			if responseSchema.Name != "" {
				if hasSuffix(responseSchema.Name, "ListResponse") {
					responseType = "list"
				} else if hasSuffix(responseSchema.Name, "Response") {
					responseType = "single"
				}
			}
			break
		}
	}

	if responseSchema == nil {
		return
	}

	// Find the underlying entity schema
	var entitySchema *Schema
	if dataRef, ok := responseSchema.Properties["Data"]; ok && dataRef != nil {
		dataSchema := dataRef.GetOrNil()
		if dataSchema != nil {
			if dataSchema.Type.PrimaryType() == SchemaTypeArray && dataSchema.Items != nil {
				// List response - get the item type
				itemsSchema := dataSchema.Items.GetOrNil()
				if itemsSchema != nil {
					entityName := itemsSchema.Name
					if entityName == "" && itemsSchema.GoType != "" {
						entityName = stripPackagePrefix(itemsSchema.GoType)
					}
					entitySchema = g.schemas[entityName]
				}
			} else {
				// Single response
				entityName := stripPackagePrefix(dataSchema.GoType)
				entitySchema = g.schemas[entityName]
			}
		}
	}

	if entitySchema == nil {
		return
	}

	// Generate parameters based on response type
	switch responseType {
	case "list":
		// Add filter parameter if not present
		if !hasParamNamed(op.Parameters, "filter") {
			op.Parameters = append(
				op.Parameters,
				g.generateFilterParam(entitySchema),
			)
		}
		// Add sort parameter if not present
		if !hasParamNamed(op.Parameters, "sort") {
			op.Parameters = append(
				op.Parameters,
				g.generateSortParam(entitySchema),
			)
		}
		// Add pagination parameter if not present
		if !hasParamNamed(op.Parameters, "page") {
			op.Parameters = append(op.Parameters, g.generatePageParam())
		}
	}
}

// generateFilterParam creates a filter parameter for the given schema
func (g *AutoGenerator) generateFilterParam(schema *Schema) Param {
	filterName := strutil.Pluralize(schema.Name) + "Filter"

	return Param{
		Schema: &Schema{
			Name:    filterName,
			Type:    PropertyType{Types: []string{SchemaTypeObject}},
			GoType:  "*servermodels.FilterNode",
			JSONTag: strutil.CamelCase(filterName),
			YAMLTag: strutil.CamelCase(filterName),
		},
		In:    "query",
		Style: "deepObject",
	}
}

// generateSortParam creates a sort parameter for the given schema
func (g *AutoGenerator) generateSortParam(schema *Schema) Param {
	sortName := strutil.Pluralize(schema.Name) + "Sort"

	return Param{
		Schema: &Schema{
			Name:    sortName,
			Type:    PropertyType{Types: []string{SchemaTypeObject}},
			GoType:  "*servermodels.FilterNode",
			JSONTag: strutil.CamelCase(sortName),
			YAMLTag: strutil.CamelCase(sortName),
		},
		In:    "query",
		Style: "deepObject",
	}
}

// generatePageParam creates a pagination parameter
func (g *AutoGenerator) generatePageParam() Param {
	return Param{
		Schema: &Schema{
			Name:    "Page",
			Type:    PropertyType{Types: []string{SchemaTypeObject}},
			GoType:  "servermodels.Page",
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
		if p.Schema != nil && (p.JSONTag == name || p.Name == name) {
			return true
		}
	}
	return false
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func stripPackagePrefix(goType string) string {
	// Remove "models." or similar package prefix
	if idx := lastIndex(goType, "."); idx >= 0 {
		return goType[idx+1:]
	}
	return goType
}

func lastIndex(s string, sep string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == sep[0] {
			return i
		}
	}
	return -1
}

// GenerateResponseSchemas creates Response and ListResponse wrappers for all entity schemas
func (g *AutoGenerator) GenerateResponseSchemas() {
	for _, schema := range g.schemas {
		if schema.XCodegenSchemaType == SchemaTypeEntity {
			g.generateSingleResponse(schema)
			g.generateListResponse(schema)
		}
	}
}

// generateListResponse creates a list response wrapper schema for an entity
func (g *AutoGenerator) generateListResponse(entitySchema *Schema) *Schema {
	responseName := entitySchema.Name + "ListResponse"

	// Check if already generated
	if existing, ok := g.schemas[responseName]; ok {
		return existing
	}

	schema := &Schema{
		Name: responseName,
		Type: PropertyType{Types: []string{SchemaTypeObject}},
		Properties: map[string]*Ref[Schema]{
			"Data": NewInline(&Schema{
				Name:    "Data",
				Type:    PropertyType{Types: []string{SchemaTypeArray}},
				Items:   NewInline(entitySchema),
				GoType:  "[]models." + entitySchema.Name,
				JSONTag: "data",
				YAMLTag: "data",
			}),
			"Meta": NewInline(&Schema{
				Name:    "Meta",
				Type:    PropertyType{Types: []string{SchemaTypeObject}},
				GoType:  "servermodels.PaginationMeta",
				JSONTag: "meta",
				YAMLTag: "meta",
			}),
		},
		Required: []string{"data", "meta"},
		GoType:   responseName,
	}

	g.schemas[responseName] = schema
	return schema
}

// generateSingleResponse creates a single response wrapper schema for an entity
func (g *AutoGenerator) generateSingleResponse(entitySchema *Schema) *Schema {
	responseName := entitySchema.Name + "Response"

	// Check if already generated
	if existing, ok := g.schemas[responseName]; ok {
		return existing
	}

	schema := &Schema{
		Name: responseName,
		Type: PropertyType{Types: []string{SchemaTypeObject}},
		Properties: map[string]*Ref[Schema]{
			"Data": NewInline(&Schema{
				Name:    "Data",
				Type:    PropertyType{Types: []string{SchemaTypeObject}},
				GoType:  "models." + entitySchema.Name,
				JSONTag: "data",
				YAMLTag: "data",
			}),
		},
		Required: []string{"data"},
		GoType:   responseName,
	}

	g.schemas[responseName] = schema
	return schema
}
