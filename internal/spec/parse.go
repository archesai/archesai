package spec

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/archesai/archesai/internal/strutil"
	"github.com/archesai/archesai/pkg/auth"
	"github.com/archesai/archesai/pkg/config"
	"github.com/archesai/archesai/pkg/executor"
	"github.com/archesai/archesai/pkg/pipelines"
	"github.com/archesai/archesai/pkg/server"
	"github.com/archesai/archesai/pkg/storage"
)

// Parser reads and parses OpenAPI specifications.
type Parser struct {
	projectDir           string
	fsys                 fs.FS          // filesystem for reading files (defaults to os.DirFS)
	project              *projectConfig // internal config built from OpenAPI
	schemas              map[string]*Schema
	refToName            map[string]string // maps ref name (filename) to actual schema name (title)
	specResult           *Spec
	responseContentTypes map[string]string // maps response name to its content type
	configIncludes       []string          // includes from arches.yaml (takes precedence over x-include-*)
	configCodegenOnly    []string          // codegen only from arches.yaml
	configCodegenLint    bool              // codegen lint from arches.yaml
}

// WithFS sets an alternative filesystem for reading files.
func (p *Parser) WithFS(fsys fs.FS) *Parser {
	p.fsys = fsys
	return p
}

// WithIncludes sets the includes list from external config (arches.yaml).
// These take precedence over x-include-* extensions in the OpenAPI spec.
func (p *Parser) WithIncludes(includes []string) *Parser {
	p.configIncludes = includes
	return p
}

// WithCodegenOnly sets the codegen only list from external config (arches.yaml).
// This takes precedence over x-generate.only extension in the OpenAPI spec.
func (p *Parser) WithCodegenOnly(only []string) *Parser {
	p.configCodegenOnly = only
	return p
}

// WithCodegenLint sets the codegen lint flag from external config (arches.yaml).
// This takes precedence over x-generate.lint extension in the OpenAPI spec.
func (p *Parser) WithCodegenLint(lint bool) *Parser {
	p.configCodegenLint = lint
	return p
}

// projectConfig is an internal structure for holding parsed config from OpenAPI.
type projectConfig struct {
	Name        string
	Title       string
	Description string
	Version     string
	Includes    []string
	Security    map[string]SecScheme
	Tags        []Tag
	CodegenOnly []string
	CodegenLint bool
}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{
		schemas:              make(map[string]*Schema),
		refToName:            make(map[string]string),
		responseContentTypes: make(map[string]string),
	}
}

// FS returns the composite filesystem used by the parser.
// This includes the project filesystem with include filesystems layered underneath.
// Must be called after Parse().
func (p *Parser) FS() fs.FS {
	return p.fsys
}

// Parse reads an openapi.yaml file and returns a Spec.
// The openapiPath should point directly to the openapi.yaml file.
func (p *Parser) Parse(openapiPath string) (*Spec, error) {
	p.projectDir = filepath.Dir(openapiPath)
	specDir := filepath.Dir(openapiPath)

	// Initialize base filesystem if not set via WithFS
	baseFS := p.fsys
	if baseFS == nil {
		baseFS = os.DirFS(specDir)
	}

	// Temporarily use base filesystem to load the document and detect includes
	p.fsys = baseFS

	// Load OpenAPI document and build project config from it
	if err := p.loadDocument(filepath.Base(openapiPath)); err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI document: %w", err)
	}

	// Build a composite filesystem with include filesystems as base layers
	// and project spec as the top layer (so project files override includes)
	p.fsys = p.buildIncludeFS(baseFS)

	// Load response content types for later lookup
	p.loadResponseContentTypes()

	// Load schemas and operations from OpenAPI spec (including from merged includes)
	operations, err := p.loadSpec(filepath.Base(openapiPath))
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// Auto-generate filter, sort, and page parameters for list operations
	autoGen := NewAutoGenerator(p.schemas)
	autoGen.GenerateForOperations(operations)

	// Convert Tags and Security to spec types
	var specTags []Tag
	for _, t := range p.project.Tags {
		specTags = append(specTags, Tag{Name: t.Name, Description: t.Description})
	}

	specSecurity := make(map[string]SecScheme)
	for name, s := range p.project.Security {
		specSecurity[name] = SecScheme{
			Type:        s.Type,
			Scheme:      s.Scheme,
			Name:        s.Name,
			In:          s.In,
			Description: s.Description,
		}
	}

	// Build the spec with all metadata
	p.specResult = &Spec{
		ProjectName:     p.project.Name,
		Operations:      operations,
		Schemas:         p.getSchemas(),
		EnabledIncludes: p.project.Includes,
		Title:           p.project.Title,
		Description:     p.project.Description,
		Version:         p.project.Version,
		Tags:            specTags,
		Security:        specSecurity,
		CodegenOnly:     p.project.CodegenOnly,
		CodegenLint:     p.project.CodegenLint,
	}

	return p.specResult, nil
}

// loadResponseContentTypes discovers and caches the content types for all response definitions.
func (p *Parser) loadResponseContentTypes() {
	responses, err := DiscoverResponses(p.fsys)
	if err != nil {
		return
	}

	for name, filePath := range responses {
		data, err := fs.ReadFile(p.fsys, filePath)
		if err != nil {
			continue
		}

		// Parse just enough to get the content type
		var resp struct {
			Content map[string]struct{} `yaml:"content"`
		}
		if err := yaml.Unmarshal(data, &resp); err != nil {
			continue
		}

		// Check for content types in priority order
		if _, ok := resp.Content["application/problem+json"]; ok {
			p.responseContentTypes[name] = "application/problem+json"
		} else if _, ok := resp.Content["application/json"]; ok {
			p.responseContentTypes[name] = "application/json"
		}
	}
}

// loadDocument reads the openapi.yaml file and builds project config from it.
func (p *Parser) loadDocument(filePath string) error {
	data, err := fs.ReadFile(p.fsys, filePath)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI document: %w", err)
	}

	var doc Document
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	// Build project config from OpenAPI document
	p.project = &projectConfig{
		Name:        doc.XProjectName,
		Title:       doc.Info.Title,
		Description: doc.Info.Description,
		Version:     doc.Info.Version,
		Tags:        doc.Tags,
		Security:    doc.Components.SecuritySchemes,
	}

	// Use config values from arches.yaml
	p.project.Includes = p.configIncludes
	p.project.CodegenOnly = p.configCodegenOnly
	p.project.CodegenLint = p.configCodegenLint

	return nil
}

// buildIncludeFS creates a composite filesystem with include filesystems as base layers.
// The project filesystem is on top so its files override include files with the same name.
func (p *Parser) buildIncludeFS(projectFS fs.FS) fs.FS {
	if len(p.project.Includes) == 0 {
		return projectFS
	}

	// Build layers: include filesystems first (as base), project FS last (as overlay)
	var layers []fs.FS

	for _, include := range p.project.Includes {
		includeFS := p.getIncludeFS(include)
		if includeFS != nil {
			layers = append(layers, includeFS)
		}
	}

	// Project FS is the top layer (overrides includes)
	layers = append(layers, projectFS)

	return NewCompositeFS(layers...)
}

// getIncludeFS returns the embedded filesystem for an include package.
// The returned FS is already stripped of the "api/" prefix.
func (p *Parser) getIncludeFS(include string) fs.FS {
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

// getModelsAlias returns the Go models alias for an include package.
func getModelsAlias(include string) string {
	switch include {
	case "server":
		return "servermodels"
	case "auth":
		return "authmodels"
	case "config":
		return "configmodels"
	case "pipelines":
		return "pipelinesmodels"
	case "executor":
		return "executormodels"
	case "storage":
		return "storagemodels"
	default:
		return ""
	}
}

// loadSpec reads the OpenAPI spec file and extracts schemas and operations.
func (p *Parser) loadSpec(specPath string) ([]Operation, error) {
	specDir := path.Dir(specPath)
	if specDir == "." {
		specDir = ""
	}

	data, err := fs.ReadFile(p.fsys, specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	var doc Document
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Load schemas from components/schemas (both explicit refs and auto-discovered)
	resolver := NewResolver(p.fsys, specDir)
	if err := p.loadSchemas(&doc, resolver); err != nil {
		return nil, fmt.Errorf("failed to load schemas: %w", err)
	}

	// Second pass: resolve refs to actual schema names
	for _, schema := range p.schemas {
		p.resolveSchemaRefs(schema)
	}

	// Auto-generate response schemas for all entities
	autoGen := NewAutoGenerator(p.schemas)
	autoGen.GenerateResponseSchemas()

	// Load operations from paths (both explicit refs and auto-discovered)
	operations, err := p.loadPaths(&doc, resolver)
	if err != nil {
		return nil, fmt.Errorf("failed to load paths: %w", err)
	}

	return operations, nil
}

// loadSchemas loads schemas from components/schemas section and auto-discovered files.
func (p *Parser) loadSchemas(doc *Document, resolver *Resolver) error {
	// First, load schemas from components/schemas (both $ref and inline)
	for schemaName, schemaRef := range doc.Components.Schemas {
		if schemaRef.IsRef() {
			// $ref to external file
			data, err := resolver.ResolveFile(schemaRef.RefPath)
			if err != nil {
				return fmt.Errorf("failed to resolve schema ref %s: %w", schemaRef.RefPath, err)
			}

			schema, err := p.loadSchemaFile(data, schemaName, path.Dir(schemaRef.RefPath), resolver)
			if err != nil {
				return fmt.Errorf("failed to load schema %s: %w", schemaName, err)
			}

			// Map the filename to the actual schema name
			refName := ExtractSchemaNameFromRef(schemaRef.RefPath)
			p.refToName[refName] = schema.Name

			p.schemas[schema.Name] = schema
		} else if schemaRef.IsInline() {
			// Inline schema definition - the value is already a *Schema
			inlineSchema := schemaRef.GetOrNil()
			if inlineSchema == nil {
				continue
			}
			// Set the name if not already set
			if inlineSchema.Name == "" {
				inlineSchema.Name = schemaName
			}
			// Process the inline schema
			p.processInlineSchema(inlineSchema, schemaName)
			p.schemas[inlineSchema.Name] = inlineSchema
		}
	}

	// Auto-discover schemas from components/schemas/ directory
	discovered, err := DiscoverSchemas(p.fsys)
	if err != nil {
		return fmt.Errorf("failed to discover schemas: %w", err)
	}

	for schemaName, filePath := range discovered {
		// Skip if already loaded via explicit ref or inline
		if _, exists := p.schemas[schemaName]; exists {
			continue
		}
		// Also check if this schema was loaded under a different name (via title)
		if _, exists := p.refToName[schemaName]; exists {
			continue
		}

		data, err := fs.ReadFile(p.fsys, filePath)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", filePath, err)
		}

		schema, err := p.loadSchemaFile(data, schemaName, path.Dir(filePath), resolver)
		if err != nil {
			return fmt.Errorf("failed to load schema %s: %w", schemaName, err)
		}

		p.refToName[schemaName] = schema.Name
		p.schemas[schema.Name] = schema
	}

	// Second pass: resolve ref properties now that all schemas are loaded
	for _, schema := range p.schemas {
		p.resolveRefProperties(schema)
	}

	return nil
}

// processInlineSchema processes an inline Schema parsed from YAML.
// It normalizes property names, sets computed fields like JSONTag/GoType.
func (p *Parser) processInlineSchema(schema *Schema, defaultName string) {
	// Derive name from title or default
	if schema.Name == "" {
		if schema.Title != "" {
			schema.Name = schema.Title
		} else {
			schema.Name = defaultName
		}
	}

	// Set default type if not set
	if schema.Type.PrimaryType() == "" {
		schema.Type = PropertyType{Types: []string{SchemaTypeObject}}
	}

	// Initialize properties map if nil
	if schema.Properties == nil {
		schema.Properties = make(map[string]*Ref[Schema])
	}

	// Handle allOf composition (common pattern: allOf: [$ref: ./Base.yaml, {type: object, properties: ...}])
	if len(schema.AllOf) > 0 {
		for _, item := range schema.AllOf {
			if item.IsRef() {
				// This is a reference to a base schema (e.g., ./Base.yaml)
				// The base fields (id, createdAt, updatedAt) will be added later
				continue
			}
			// This is the inline object definition with properties
			inlineSchema := item.GetOrNil()
			if inlineSchema == nil {
				continue
			}
			for propName, propRef := range inlineSchema.Properties {
				// Normalize property name
				fieldName := strutil.PascalCase(propName)
				propSchema := propRef.GetOrNil()
				if propSchema != nil {
					p.processPropertySchema(propSchema, fieldName, propName, schema.Name)
				}
				// Move property to parent schema with normalized name
				if fieldName != propName {
					schema.Properties[fieldName] = propRef
				} else {
					schema.Properties[propName] = propRef
				}
			}
			// Merge required fields
			schema.Required = append(schema.Required, inlineSchema.Required...)
		}
	} else {
		// Regular schema - process properties in place
		processedProps := make(map[string]*Ref[Schema])
		for propName, propRef := range schema.Properties {
			fieldName := strutil.PascalCase(propName)
			propSchema := propRef.GetOrNil()
			if propSchema != nil {
				p.processPropertySchema(propSchema, fieldName, propName, schema.Name)
			}
			processedProps[fieldName] = propRef
		}
		schema.Properties = processedProps
	}

	// Build required map and set omitempty on optional fields
	requiredMap := make(map[string]bool)
	for _, r := range schema.Required {
		requiredMap[r] = true
	}
	for propName, propRef := range schema.Properties {
		propSchema := propRef.GetOrNil()
		if propSchema == nil {
			continue
		}
		jsonName := propSchema.JSONTag
		if jsonName == "" {
			jsonName = strutil.CamelCase(propName)
		}
		if !requiredMap[jsonName] && !requiredMap[propName] {
			propSchema.JSONTag = jsonName + ",omitempty"
			propSchema.YAMLTag = jsonName + ",omitempty"
		}
	}

	// Add base fields for entity schemas
	if schema.XCodegenSchemaType == SchemaTypeEntity {
		p.addBaseFieldsToSchema(schema)
	}

	// Set GoType with models alias prefix if this schema comes from an include
	if schema.XInternal != "" {
		alias := getModelsAlias(schema.XInternal)
		if alias != "" {
			schema.GoType = alias + "." + schema.Name
		} else {
			schema.GoType = schema.Name
		}
	} else {
		schema.GoType = schema.Name
	}
}

// processPropertySchema sets computed fields on a property schema.
func (p *Parser) processPropertySchema(propSchema *Schema, fieldName, jsonName, parentName string) {
	propSchema.Name = fieldName
	propSchema.JSONTag = jsonName
	propSchema.YAMLTag = jsonName
	propSchema.Nullable = propSchema.Type.Nullable

	// Handle oneOf for nullable types with constraints
	// e.g., oneOf: [{type: string, minLength: 1}, {type: 'null'}]
	if len(propSchema.OneOf) == 2 {
		var nonNullSchema *Schema
		hasNull := false
		for _, ref := range propSchema.OneOf {
			s := ref.GetOrNil()
			if s == nil {
				continue
			}
			pt := s.Type.PrimaryType()
			if pt == SchemaTypeNull || pt == "" && s.Type.Nullable {
				hasNull = true
			} else if nonNullSchema == nil {
				nonNullSchema = s
			}
		}
		// If it's a nullable union type, use the non-null schema's type
		if hasNull && nonNullSchema != nil {
			propSchema.Nullable = true
			propSchema.Type = nonNullSchema.Type
			propSchema.Format = nonNullSchema.Format
			propSchema.MinLength = nonNullSchema.MinLength
			propSchema.MaxLength = nonNullSchema.MaxLength
			propSchema.Pattern = nonNullSchema.Pattern
			propSchema.Minimum = nonNullSchema.Minimum
			propSchema.Maximum = nonNullSchema.Maximum
			propSchema.Enum = nonNullSchema.Enum
		}
	}

	// Determine Go type based on the schema type
	primaryType := propSchema.Type.PrimaryType()
	switch primaryType {
	case SchemaTypeArray:
		if propSchema.Items == nil {
			propSchema.GoType = "[]any"
		} else if propSchema.Items.IsRef() {
			// Array items are a $ref - use the referenced schema name
			itemName := ExtractSchemaNameFromRef(propSchema.Items.RefPath)
			propSchema.GoType = "[]" + itemName
		} else {
			// Array items are inline
			itemSchema := propSchema.Items.GetOrNil()
			if itemSchema == nil {
				propSchema.GoType = "[]any"
			} else if itemSchema.Type.PrimaryType() == SchemaTypeObject && len(itemSchema.Properties) > 0 {
				// Inline object items - generate a type name
				itemTypeName := parentName + fieldName + "Item"
				itemSchema.GoType = itemTypeName
				// Process nested properties
				processedProps := make(map[string]*Ref[Schema])
				for nestedPropName, nestedPropRef := range itemSchema.Properties {
					nestedFieldName := strutil.PascalCase(nestedPropName)
					nestedPropSchema := nestedPropRef.GetOrNil()
					if nestedPropSchema != nil {
						p.processPropertySchema(nestedPropSchema, nestedFieldName, nestedPropName, itemTypeName)
					}
					processedProps[nestedFieldName] = nestedPropRef
				}
				itemSchema.Properties = processedProps
				propSchema.GoType = "[]" + itemTypeName
			} else {
				// Primitive array items
				itemGoType := SchemaToGoType(itemSchema.Type.PrimaryType(), itemSchema.Format)
				propSchema.GoType = "[]" + itemGoType
			}
		}
	case SchemaTypeObject:
		// Objects with no properties are generic maps
		if len(propSchema.Properties) == 0 {
			propSchema.GoType = GoTypeMapString
		} else {
			// Nested object type name includes parent name to avoid conflicts
			nestedTypeName := parentName + fieldName
			propSchema.GoType = nestedTypeName
			// Recursively process nested object properties
			processedProps := make(map[string]*Ref[Schema])
			for nestedPropName, nestedPropRef := range propSchema.Properties {
				nestedFieldName := strutil.PascalCase(nestedPropName)
				nestedPropSchema := nestedPropRef.GetOrNil()
				if nestedPropSchema != nil {
					p.processPropertySchema(nestedPropSchema, nestedFieldName, nestedPropName, nestedTypeName)
				}
				processedProps[nestedFieldName] = nestedPropRef
			}
			propSchema.Properties = processedProps
		}
	default:
		propSchema.GoType = SchemaToGoType(primaryType, propSchema.Format)
	}
}

// loadSchemaFile parses an OpenAPI schema file.
func (p *Parser) loadSchemaFile(
	data []byte,
	defaultName, schemaDir string,
	resolver *Resolver,
) (*Schema, error) {
	var sf SchemaFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Derive name from title or default
	name := sf.Title
	if name == "" {
		name = defaultName
	}

	schema := &Schema{
		Name:               name,
		Description:        sf.Description,
		Type:               PropertyType{Types: []string{SchemaTypeObject}},
		Properties:         make(map[string]*Ref[Schema]),
		Required:           sf.Required,
		XCodegen:           sf.XCodegen,
		XCodegenSchemaType: SchemaType(sf.XCodegenSchemaType),
		XInternal:          sf.XInternal,
	}

	// Handle allOf composition (common pattern: allOf: [$ref: Base.yaml, {type: object, properties: ...}])
	if len(sf.AllOf) > 0 {
		for _, itemRef := range sf.AllOf {
			if itemRef.IsRef() {
				// This is a reference to a base schema (e.g., Base.yaml)
				// The base fields (id, createdAt, updatedAt) will be added later
				continue
			}
			// This is the inline object definition with properties
			item := itemRef.GetOrNil()
			if item == nil {
				continue
			}
			for propName, propRef := range item.Properties {
				fieldName := strutil.PascalCase(propName)
				propSchema := propRef.GetOrNil()
				if propSchema != nil {
					p.processPropertySchema(propSchema, fieldName, propName, name)
				}
				schema.Properties[fieldName] = propRef
			}
			// Merge required fields
			schema.Required = append(schema.Required, item.Required...)
		}
	} else {
		// Regular schema without allOf
		for propName, propRef := range sf.Properties {
			fieldName := strutil.PascalCase(propName)

			// Handle ref properties by resolving them
			if propRef.IsRef() {
				refName := ExtractSchemaNameFromRef(propRef.RefPath)
				if refSchema := p.resolveSchemaByRef(refName); refSchema != nil {
					// Create inline schema that references the resolved schema
					propSchema := &Schema{
						Name:    fieldName,
						GoType:  refSchema.Name,
						Type:    refSchema.Type,
						JSONTag: propName,
						YAMLTag: propName,
					}
					schema.Properties[fieldName] = NewInline(propSchema)
					continue
				}
			}

			// Handle inline properties
			propSchema := propRef.GetOrNil()
			if propSchema != nil {
				p.processPropertySchema(propSchema, fieldName, propName, name)
			}
			schema.Properties[fieldName] = propRef
		}
	}

	// Build required map and set omitempty on optional fields
	requiredMap := make(map[string]bool)
	for _, r := range schema.Required {
		requiredMap[r] = true
	}
	for propName, propRef := range schema.Properties {
		propSchema := propRef.GetOrNil()
		if propSchema == nil {
			continue
		}
		jsonName := propSchema.JSONTag
		if jsonName == "" {
			jsonName = strutil.CamelCase(propName)
		}
		if !requiredMap[jsonName] && !requiredMap[propName] {
			propSchema.JSONTag = jsonName + ",omitempty"
			propSchema.YAMLTag = jsonName + ",omitempty"
		}
	}

	// Add base fields for entity schemas
	if sf.XCodegenSchemaType == "entity" {
		p.addBaseFieldsToSchema(schema)
	}

	// Set GoType with models alias prefix if this schema comes from an include
	if schema.XInternal != "" {
		alias := getModelsAlias(schema.XInternal)
		if alias != "" {
			schema.GoType = alias + "." + schema.Name
		} else {
			schema.GoType = schema.Name
		}
	} else {
		schema.GoType = schema.Name
	}

	return schema, nil
}

// loadPaths loads operations from the paths section and auto-discovered files.
func (p *Parser) loadPaths(
	doc *Document,
	resolver *Resolver,
) ([]Operation, error) {
	var operations []Operation
	loadedPaths := make(map[string]bool) // Track which paths have been loaded

	// First, load explicitly referenced paths (both $ref and inline)
	for pathStr, pathRef := range doc.Paths {
		var pathItem *PathItem
		var pathFileDir string

		if pathRef.IsRef() {
			// Resolve the $ref to load the path file
			data, err := resolver.ResolveFile(pathRef.Ref)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path ref %s: %w", pathRef.Ref, err)
			}
			var pi PathItem
			if err := yaml.Unmarshal(data, &pi); err != nil {
				return nil, fmt.Errorf("failed to parse path file %s: %w", pathRef.Ref, err)
			}
			pathItem = &pi
			pathFileDir = path.Dir(pathRef.Ref)
		} else if pathRef.IsInline() {
			// Inline path definition
			pathItem = pathRef.ToPathItem()
			pathFileDir = ""
		} else {
			// Empty path entry, skip
			continue
		}

		// Process each HTTP method
		ops, err := p.processPathItem(pathStr, pathItem, pathFileDir, resolver)
		if err != nil {
			return nil, err
		}
		operations = append(operations, ops...)
		loadedPaths[pathStr] = true
	}

	// Auto-discover paths from paths/ directory
	pathFiles, err := DiscoverPaths(p.fsys)
	if err != nil {
		return nil, fmt.Errorf("failed to discover paths: %w", err)
	}

	for _, filePath := range pathFiles {
		data, err := fs.ReadFile(p.fsys, filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read path file %s: %w", filePath, err)
		}

		var pathItem PathItem
		if err := yaml.Unmarshal(data, &pathItem); err != nil {
			return nil, fmt.Errorf("failed to parse path file %s: %w", filePath, err)
		}

		// x-path is required for auto-discovered paths
		if pathItem.XPath == "" {
			return nil, fmt.Errorf("path file %s is missing required x-path extension", filePath)
		}

		pathStr := pathItem.XPath

		// Skip if already loaded via explicit ref
		if loadedPaths[pathStr] {
			continue
		}

		// Get path directory (paths from discovery are already relative to fs root)
		pathFileDir := path.Dir(filePath)
		ops, err := p.processPathItem(pathStr, &pathItem, pathFileDir, resolver)
		if err != nil {
			return nil, err
		}
		operations = append(operations, ops...)
		loadedPaths[pathStr] = true
	}

	// Sort operations by path then method for consistent ordering
	sort.Slice(operations, func(i, j int) bool {
		if operations[i].Path != operations[j].Path {
			return operations[i].Path < operations[j].Path
		}
		return operations[i].Method < operations[j].Method
	})

	return operations, nil
}

// processPathItem extracts operations from a path item.
func (p *Parser) processPathItem(
	pathStr string,
	pathItem *PathItem,
	pathFileDir string,
	resolver *Resolver,
) ([]Operation, error) {
	var operations []Operation

	methods := map[string]*RawOperation{
		"GET":    pathItem.Get,
		"POST":   pathItem.Post,
		"PUT":    pathItem.Put,
		"DELETE": pathItem.Delete,
		"PATCH":  pathItem.Patch,
	}

	for method, opDef := range methods {
		if opDef == nil {
			continue
		}

		op, err := p.operationToSpecOperation(pathStr, method, opDef, pathFileDir, resolver)
		if err != nil {
			return nil, fmt.Errorf("failed to process %s %s: %w", method, pathStr, err)
		}
		operations = append(operations, *op)
	}

	return operations, nil
}

// operationToSpecOperation converts an OpenAPI operation to internal Operation.
func (p *Parser) operationToSpecOperation(
	pathStr, method string,
	opDef *RawOperation,
	pathFileDir string,
	resolver *Resolver,
) (*Operation, error) {
	tag := ""
	if len(opDef.Tags) > 0 {
		tag = opDef.Tags[0]
	}

	op := &Operation{
		ID:             opDef.OperationID,
		Method:         method,
		Path:           pathStr,
		Summary:        opDef.Summary,
		Description:    opDef.Description,
		Tag:            tag,
		CustomHandler:  opDef.XCustomHandler,
		PublicEndpoint: opDef.XPublicEndpoint,
		Internal:       opDef.XInternal,
	}

	// Handle security
	securityReqs := opDef.Security
	explicitSecurity := len(opDef.Security) > 0

	// Check for public endpoint marker
	if opDef.XPublicEndpoint {
		// Operation is explicitly public
		securityReqs = []map[string][]string{}
	} else if securityReqs == nil && len(p.project.Security) > 0 {
		// Inherit default security from project
		for name := range p.project.Security {
			securityReqs = append(securityReqs, map[string][]string{name: {}})
		}
	}

	op.ExplicitSecurity = explicitSecurity

	// Check if this is explicitly empty security (security: - {})
	// This pattern means the endpoint is public/unauthenticated
	if explicitSecurity && len(securityReqs) == 1 && len(securityReqs[0]) == 0 {
		op.EmptySecurity = true
	}

	for _, secReq := range securityReqs {
		for name, scopes := range secReq {
			sec := Security{
				Name:   name,
				Scopes: scopes,
			}
			if scheme, ok := p.project.Security[name]; ok {
				sec.Type = scheme.Type
				sec.Scheme = scheme.Scheme
				if scheme.Type == "apiKey" && scheme.In == "cookie" {
					sec.Scheme = "cookie"
				}
			}
			op.Security = append(op.Security, sec)
		}
	}

	// Convert parameters
	for _, param := range opDef.Parameters {
		if param.Ref != "" {
			// Skip $ref parameters for now - they're usually filter/sort/page which are auto-generated
			continue
		}
		op.Parameters = append(op.Parameters, p.paramToSpecParam(&param))
	}

	// Auto-extract path parameters from path template
	op.Parameters = append(op.Parameters, p.extractPathParamsFromPath(pathStr, op.Parameters)...)

	// Convert request body
	if opDef.RequestBody != nil {
		isUpdate := method == "PATCH"
		op.RequestBody = p.requestBodyToSpecRequestBody(
			opDef.RequestBody,
			opDef.OperationID,
			isUpdate,
			pathFileDir,
			resolver,
		)
	}

	// Convert responses
	op.Responses = p.responsesToSpecResponses(
		opDef.Responses,
		opDef.OperationID,
		pathFileDir,
		resolver,
	)

	// Determine if operation has a path parameter (for 404 response)
	hasPathParam := strings.Contains(pathStr, "{")

	// Build set of existing response status codes to avoid duplicates
	existingCodes := make(map[string]bool)
	for _, resp := range op.Responses {
		existingCodes[resp.StatusCode] = true
	}

	// Add standard error responses only if not already defined
	for _, errResp := range p.standardErrorResponses(hasPathParam) {
		if !existingCodes[errResp.StatusCode] {
			op.Responses = append(op.Responses, errResp)
		}
	}

	// Sort responses by status code for consistent output
	sort.Slice(op.Responses, func(i, j int) bool {
		return op.Responses[i].StatusCode < op.Responses[j].StatusCode
	})

	return op, nil
}

// paramToSpecParam converts an OpenAPI parameter to internal Param.
func (p *Parser) paramToSpecParam(param *RawParameter) Param {
	paramType := PropertyType{Types: []string{SchemaTypeString}}
	format := ""
	if param.Schema != nil {
		if param.Schema.Type != "" {
			paramType = PropertyType{Types: []string{param.Schema.Type}}
		}
		if param.Schema.Format != "" {
			format = param.Schema.Format
		}
		if param.Schema.Ref != "" {
			// Handle $ref to a schema - extract the type from the referenced schema
			refName := ExtractSchemaNameFromRef(param.Schema.Ref)
			if refSchema := p.resolveSchemaByRef(refName); refSchema != nil {
				paramType = refSchema.Type
				format = refSchema.Format
			}
		}
	}

	// Build JSON tag - add omitempty for optional params
	jsonTag := param.Name
	if !param.Required {
		jsonTag += ",omitempty"
	}

	schema := &Schema{
		Name:        strutil.PascalCase(param.Name),
		Type:        paramType,
		Format:      format,
		Description: param.Description,
		GoType:      SchemaToGoType(paramType.PrimaryType(), format),
		JSONTag:     jsonTag,
		YAMLTag:     param.Name,
	}

	explode := false
	if param.Explode != nil {
		explode = *param.Explode
	}

	return Param{
		Schema:  schema,
		In:      param.In,
		Style:   param.Style,
		Explode: explode,
	}
}

// extractPathParamsFromPath extracts path parameters from a path template.
func (p *Parser) extractPathParamsFromPath(
	pathStr string,
	existingParams []Param,
) []Param {
	var params []Param

	// Build a set of already defined path parameters
	definedParams := make(map[string]bool)
	for _, param := range existingParams {
		if param.In == "path" {
			definedParams[param.JSONTag] = true
		}
	}

	// Find all {param} patterns in the path
	for i := 0; i < len(pathStr); i++ {
		if pathStr[i] == '{' {
			end := strings.Index(pathStr[i:], "}")
			if end > 0 {
				paramName := pathStr[i+1 : i+end]
				if !definedParams[paramName] {
					params = append(params, Param{
						Schema: &Schema{
							Name:        strutil.PascalCase(paramName),
							Type:        PropertyType{Types: []string{SchemaTypeString}},
							Format:      FormatUUID,
							GoType:      GoTypeUUID,
							JSONTag:     paramName,
							YAMLTag:     paramName,
							Description: "Resource identifier",
						},
						In: "path",
					})
					definedParams[paramName] = true
				}
				i += end
			}
		}
	}

	return params
}

// requestBodyToSpecRequestBody converts an OpenAPI request body to internal RequestBody.
func (p *Parser) requestBodyToSpecRequestBody(
	reqBody *RawRequestBody,
	opID string,
	isUpdate bool,
	pathFileDir string,
	resolver *Resolver,
) *RequestBody {
	// Handle $ref to request body file
	if reqBody.Ref != "" {
		data, err := resolver.ResolveFileFrom(pathFileDir, reqBody.Ref)
		if err != nil {
			return nil
		}

		var reqBodyFile RawRequestBody
		if err := yaml.Unmarshal(data, &reqBodyFile); err != nil {
			return nil
		}

		// Parse the resolved request body file
		return p.parseRequestBodyContent(
			&reqBodyFile,
			opID,
			isUpdate,
			path.Dir(reqBody.Ref),
			resolver,
		)
	}

	return p.parseRequestBodyContent(reqBody, opID, isUpdate, pathFileDir, resolver)
}

// parseRequestBodyContent parses the actual request body content.
func (p *Parser) parseRequestBodyContent(
	reqBody *RawRequestBody,
	opID string,
	isUpdate bool,
	pathFileDir string,
	resolver *Resolver,
) *RequestBody {
	var schema *Schema

	// Get the JSON content type schema
	if content, ok := reqBody.Content["application/json"]; ok {
		schemaRef := content.Schema
		if schemaRef == nil {
			return nil
		}

		if schemaRef.IsRef() {
			// Reference to existing schema - resolve it
			refName := ExtractSchemaNameFromRef(schemaRef.RefPath)
			if refSchema, ok := p.schemas[refName]; ok {
				schema = refSchema
			} else {
				// Try to resolve the schema file
				data, err := resolver.ResolveFileFrom(pathFileDir, schemaRef.RefPath)
				if err == nil {
					parsedSchema, parseErr := p.loadSchemaFile(data, refName, path.Dir(schemaRef.RefPath), resolver)
					if parseErr == nil {
						schema = parsedSchema
						p.schemas[refName] = schema
					}
				}
				if schema == nil {
					schema = &Schema{
						Name:   refName,
						GoType: refName,
						Type:   PropertyType{Types: []string{SchemaTypeObject}},
					}
				}
			}
		} else {
			// Inline schema definition
			inlineSchema := schemaRef.GetOrNil()
			if inlineSchema != nil && len(inlineSchema.Properties) > 0 {
				schemaName := opID + "Request"
				schema = &Schema{
					Name:       schemaName,
					Type:       PropertyType{Types: []string{SchemaTypeObject}},
					Properties: make(map[string]*Ref[Schema]),
					GoType:     schemaName,
				}

				requiredMap := make(map[string]bool)
				for _, r := range inlineSchema.Required {
					requiredMap[r] = true
				}

				if !isUpdate {
					schema.Required = inlineSchema.Required
				}

				for propName, propRef := range inlineSchema.Properties {
					fieldName := strutil.PascalCase(propName)
					propSchema := propRef.GetOrNil()
					if propSchema != nil {
						p.processPropertySchema(propSchema, fieldName, propName, schemaName)

						if isUpdate {
							propSchema.JSONTag = propName + ",omitempty"
							propSchema.YAMLTag = propName + ",omitempty"
						} else if !requiredMap[propName] {
							propSchema.JSONTag = propName + ",omitempty"
							propSchema.YAMLTag = propName + ",omitempty"
						}
					}

					schema.Properties[fieldName] = propRef
				}

				p.schemas[schemaName] = schema
			}
		}
	}

	return &RequestBody{
		Schema:   schema,
		Required: !isUpdate,
	}
}

// responsesToSpecResponses converts OpenAPI responses to internal Response list.
func (p *Parser) responsesToSpecResponses(
	responses map[string]RawResponse,
	opID string,
	pathFileDir string,
	resolver *Resolver,
) []Response {
	var result []Response

	for statusCode, respRef := range responses {
		respDef := p.responseRefToSpecResponse(statusCode, &respRef, opID, pathFileDir, resolver)
		if respDef != nil {
			result = append(result, *respDef)
		}
	}

	return result
}

// responseRefToSpecResponse converts a single OpenAPI response to Response.
func (p *Parser) responseRefToSpecResponse(
	statusCode string,
	respRef *RawResponse,
	opID string,
	pathFileDir string,
	resolver *Resolver,
) *Response {
	// Handle $ref to response file
	if respRef.Ref != "" {
		// Extract response name from ref path
		responseName := ExtractSchemaNameFromRef(respRef.Ref)

		// First check if we have a pre-generated response schema
		if responseSchema := p.resolveSchemaByRef(responseName); responseSchema != nil {
			respContentType := p.responseContentTypes[responseName]
			if respContentType == "" {
				respContentType = "application/json"
			}
			return &Response{
				StatusCode:  statusCode,
				ContentType: respContentType,
				Schema:      responseSchema,
			}
		}

		// Resolve as file reference
		data, err := resolver.ResolveFileFrom(pathFileDir, respRef.Ref)
		if err != nil {
			return p.createCustomResponseDef(statusCode, opID)
		}

		var respFile RawResponse
		if err := yaml.Unmarshal(data, &respFile); err != nil {
			return p.createCustomResponseDef(statusCode, opID)
		}

		return p.parseResponseFile(statusCode, &respFile, opID, path.Dir(respRef.Ref), resolver)
	}

	return p.parseResponseFile(statusCode, respRef, opID, pathFileDir, resolver)
}

// resolveSchemaByRef looks up a schema by its ref name, using refToName mapping if needed.
func (p *Parser) resolveSchemaByRef(refName string) *Schema {
	// First try direct lookup
	if schema, ok := p.schemas[refName]; ok {
		return schema
	}
	// Then try mapped name
	if actualName, ok := p.refToName[refName]; ok {
		if schema, ok := p.schemas[actualName]; ok {
			return schema
		}
	}
	return nil
}

// parseResponseFile parses an OpenAPI response definition.
func (p *Parser) parseResponseFile(
	statusCode string,
	resp *RawResponse,
	opID string,
	respFileDir string,
	resolver *Resolver,
) *Response {
	// Handle content - check for both application/json and application/problem+json
	var content *MediaType
	var contentType string

	if c, ok := resp.Content["application/problem+json"]; ok {
		content = &c
		contentType = "application/problem+json"
	} else if c, ok := resp.Content["application/json"]; ok {
		content = &c
		contentType = "application/json"
	}

	if content != nil {
		schemaRef := content.Schema
		if schemaRef == nil {
			return nil
		}

		if schemaRef.IsRef() {
			// Direct reference to a schema - return the schema directly (properties become output fields)
			refName := ExtractSchemaNameFromRef(schemaRef.RefPath)
			if refSchema := p.resolveSchemaByRef(refName); refSchema != nil {
				return &Response{
					StatusCode:  statusCode,
					ContentType: contentType,
					Schema:      refSchema,
				}
			}
		}

		// Get the inline schema
		inlineSchema := schemaRef.GetOrNil()
		if inlineSchema == nil {
			return nil
		}

		// Handle inline response schema with data/meta pattern
		if len(inlineSchema.Properties) > 0 {
			// Check if this is a list response (has data array and meta)
			if dataPropRef, hasData := inlineSchema.Properties["data"]; hasData {
				dataProp := dataPropRef.GetOrNil()
				if dataProp != nil && dataProp.Type.PrimaryType() == SchemaTypeArray &&
					dataProp.Items != nil {
					// This is a list response - look up the pre-generated ListResponse schema
					if dataProp.Items.IsRef() {
						itemRefName := ExtractSchemaNameFromRef(dataProp.Items.RefPath)
						if itemSchema := p.resolveSchemaByRef(itemRefName); itemSchema != nil {
							responseName := itemSchema.Name + "ListResponse"
							if responseSchema := p.resolveSchemaByRef(responseName); responseSchema != nil {
								return &Response{
									StatusCode:  statusCode,
									ContentType: contentType,
									Schema:      responseSchema,
								}
							}
						}
					}
				} else if dataPropRef.IsRef() {
					// Single item response with data wrapper - look up the pre-generated Response schema
					itemRefName := ExtractSchemaNameFromRef(dataPropRef.RefPath)
					if itemSchema := p.resolveSchemaByRef(itemRefName); itemSchema != nil {
						responseName := itemSchema.Name + "Response"
						if responseSchema := p.resolveSchemaByRef(responseName); responseSchema != nil {
							return &Response{
								StatusCode:  statusCode,
								ContentType: contentType,
								Schema:      responseSchema,
							}
						}
						// Response schema doesn't exist - create it on the fly
						// Use schema name for GoType - templates will add appropriate prefix
						dataGoType := itemSchema.Name
						responseSchema := &Schema{
							Name: responseName,
							Type: PropertyType{Types: []string{SchemaTypeObject}},
							Properties: map[string]*Ref[Schema]{
								"Data": NewInline(&Schema{
									Name:    "Data",
									GoType:  dataGoType,
									Type:    itemSchema.Type,
									JSONTag: "data",
									YAMLTag: "data",
								}),
							},
							Required: []string{"data"},
							GoType:   responseName,
						}
						p.schemas[responseName] = responseSchema
						return &Response{
							StatusCode:  statusCode,
							ContentType: contentType,
							Schema:      responseSchema,
						}
					}
				}
			}

			// Check if response has meta (list pattern) but data is not array
			if _, hasMeta := inlineSchema.Properties["meta"]; hasMeta {
				if dataPropRef, hasData := inlineSchema.Properties["data"]; hasData &&
					dataPropRef.IsRef() {
					// Single item with meta - look up the pre-generated Response schema
					itemRefName := ExtractSchemaNameFromRef(dataPropRef.RefPath)
					if itemSchema := p.resolveSchemaByRef(itemRefName); itemSchema != nil {
						responseName := itemSchema.Name + "Response"
						if responseSchema := p.resolveSchemaByRef(responseName); responseSchema != nil {
							return &Response{
								StatusCode:  statusCode,
								ContentType: contentType,
								Schema:      responseSchema,
							}
						}
					}
				}
			}

			// Custom inline response - create output types for any response with properties
			if len(inlineSchema.Properties) > 0 {
				return p.createInlineResponseDef(statusCode, opID, inlineSchema, contentType)
			}
		}
	}

	// No content or empty response
	if statusCode == "204" {
		emptySchema := &Schema{
			Name:       "NoContent",
			Type:       PropertyType{Types: []string{SchemaTypeObject}},
			Properties: make(map[string]*Ref[Schema]),
		}
		return &Response{
			StatusCode:  statusCode,
			ContentType: "",
			Schema:      emptySchema,
		}
	}

	// For operations without defined responses, return nil instead of creating empty output
	return nil
}

// createCustomResponseDef creates a response def for custom handlers.
func (p *Parser) createCustomResponseDef(statusCode, opID string) *Response {
	// 204 responses should not create output types - they have no content
	if statusCode == "204" {
		emptySchema := &Schema{
			Name:       "NoContent",
			Type:       PropertyType{Types: []string{SchemaTypeObject}},
			Properties: make(map[string]*Ref[Schema]),
		}
		return &Response{
			StatusCode:  statusCode,
			ContentType: "",
			Schema:      emptySchema,
		}
	}

	customSchema := &Schema{
		Name:       opID + "Output",
		Type:       PropertyType{Types: []string{SchemaTypeObject}},
		Properties: make(map[string]*Ref[Schema]),
		GoType:     opID + "Output",
	}
	p.schemas[customSchema.Name] = customSchema
	return &Response{
		StatusCode:  statusCode,
		ContentType: "application/json",
		Schema:      customSchema,
	}
}

// createInlineResponseDef creates a response def from inline schema.
func (p *Parser) createInlineResponseDef(
	statusCode, opID string,
	schema *Schema,
	contentType string,
) *Response {
	responseName := opID + "Output"

	// Build required map for property lookup
	requiredMap := make(map[string]bool)
	for _, req := range schema.Required {
		requiredMap[req] = true
	}

	responseSchema := &Schema{
		Name:       responseName,
		Type:       PropertyType{Types: []string{SchemaTypeObject}},
		Properties: make(map[string]*Ref[Schema]),
		Required:   schema.Required,
		GoType:     responseName,
	}

	for propName, propRef := range schema.Properties {
		fieldName := strutil.PascalCase(propName)

		// Handle ref properties by looking up the referenced schema
		if propRef.IsRef() {
			refName := ExtractSchemaNameFromRef(propRef.RefPath)
			if refSchema := p.resolveSchemaByRef(refName); refSchema != nil {
				// Create a new inline schema that references the resolved schema
				propSchema := &Schema{
					Name:    fieldName,
					GoType:  refSchema.Name, // Use schema name, template adds models. prefix
					Type:    refSchema.Type,
					JSONTag: propName,
					YAMLTag: propName,
				}
				if !requiredMap[propName] {
					propSchema.JSONTag = propName + ",omitempty"
					propSchema.YAMLTag = propName + ",omitempty"
				}
				responseSchema.Properties[fieldName] = NewInline(propSchema)
				continue
			}
		}

		// Handle inline properties
		propSchema := propRef.GetOrNil()
		if propSchema != nil {
			p.processPropertySchema(propSchema, fieldName, propName, responseName)

			// Mark optional properties with omitempty
			if !requiredMap[propName] {
				propSchema.JSONTag = propName + ",omitempty"
				propSchema.YAMLTag = propName + ",omitempty"
			}
		}

		responseSchema.Properties[fieldName] = propRef
	}

	p.schemas[responseName] = responseSchema
	return &Response{
		StatusCode:  statusCode,
		ContentType: contentType,
		Schema:      responseSchema,
	}
}

// resolveRefProperties resolves ref properties that couldn't be resolved during initial loading.
// This is called in a second pass after all schemas have been loaded.
func (p *Parser) resolveRefProperties(schema *Schema) {
	for propName, propRef := range schema.Properties {
		// Skip if already resolved (has inline schema)
		if propRef.GetOrNil() != nil {
			continue
		}

		// If it's a ref, try to resolve it now
		if propRef.IsRef() {
			refName := ExtractSchemaNameFromRef(propRef.RefPath)
			if refSchema := p.resolveSchemaByRef(refName); refSchema != nil {
				// propName is PascalCase (from Properties map key), convert back to camelCase for JSON tag
				fieldName := propName
				jsonName := strutil.CamelCase(propName)

				// Check if this property is required (required list uses camelCase names)
				isRequired := false
				for _, req := range schema.Required {
					if req == jsonName {
						isRequired = true
						break
					}
				}

				jsonTag := jsonName
				if !isRequired {
					jsonTag = jsonName + ",omitempty"
				}

				// Create inline schema that references the resolved schema
				propSchema := &Schema{
					Name:    fieldName,
					GoType:  refSchema.Name,
					Type:    refSchema.Type,
					JSONTag: jsonTag,
					YAMLTag: jsonTag,
				}
				schema.Properties[propName] = NewInline(propSchema)
			}
		}
	}
}

// resolveSchemaRefs resolves ref names to actual schema names for all properties.
func (p *Parser) resolveSchemaRefs(schema *Schema) {
	for _, propRef := range schema.Properties {
		p.resolvePropertyRef(propRef.GetOrNil())
	}
}

// resolvePropertyRef resolves a single property's ref to the actual schema name.
func (p *Parser) resolvePropertyRef(prop *Schema) {
	if prop == nil {
		return
	}

	// If this property has a GoType that matches a ref name, resolve it
	if actualName, ok := p.refToName[prop.GoType]; ok {
		prop.GoType = actualName
	}

	// Recursively resolve nested properties
	for _, nestedRef := range prop.Properties {
		p.resolvePropertyRef(nestedRef.GetOrNil())
	}

	// Resolve array item refs
	if prop.Items != nil {
		p.resolvePropertyRef(prop.Items.GetOrNil())
	}
}

// addBaseFieldsToSchema adds or updates id, createdAt, updatedAt fields on an entity schema.
func (p *Parser) addBaseFieldsToSchema(schema *Schema) {
	// Define base field defaults
	baseFields := map[string]*Schema{
		"ID": {
			Name:        "ID",
			Description: "Unique identifier for the resource",
			Type:        PropertyType{Types: []string{SchemaTypeString}},
			Format:      FormatUUID,
			GoType:      GoTypeUUID,
			JSONTag:     "id",
			YAMLTag:     "id",
		},
		"CreatedAt": {
			Name:        "CreatedAt",
			Description: "The date and time when the resource was created",
			Type:        PropertyType{Types: []string{SchemaTypeString}},
			Format:      FormatDateTime,
			GoType:      GoTypeTime,
			JSONTag:     "createdAt",
			YAMLTag:     "createdAt",
		},
		"UpdatedAt": {
			Name:        "UpdatedAt",
			Description: "The date and time when the resource was last updated",
			Type:        PropertyType{Types: []string{SchemaTypeString}},
			Format:      FormatDateTime,
			GoType:      GoTypeTime,
			JSONTag:     "updatedAt",
			YAMLTag:     "updatedAt",
		},
	}

	// Add or update base fields
	for fieldName, defaults := range baseFields {
		if existingRef, ok := schema.Properties[fieldName]; ok {
			// Field exists - add description if missing
			existing := existingRef.GetOrNil()
			if existing.Description == "" {
				existing.Description = defaults.Description
			}
		} else {
			// Field doesn't exist - add it
			schema.Properties[fieldName] = NewInline(defaults)
		}
	}

	// Add base fields to required list if not already present
	requiredFields := []string{"id", "createdAt", "updatedAt"}
	requiredMap := make(map[string]bool)
	for _, r := range schema.Required {
		requiredMap[r] = true
	}
	for _, bf := range requiredFields {
		if !requiredMap[bf] {
			schema.Required = append(schema.Required, bf)
		}
	}
}

// standardErrorResponses returns standard error response definitions.
// hasResourceID indicates whether the operation has a path parameter (for 404 response).
func (p *Parser) standardErrorResponses(hasResourceID bool) []Response {
	responses := []Response{
		{StatusCode: "400", ContentType: "application/problem+json", Schema: p.problemSchema()},
		{StatusCode: "401", ContentType: "application/problem+json", Schema: p.problemSchema()},
	}
	// Only add 404 for operations with a resource ID (GET/PATCH/DELETE by ID)
	if hasResourceID {
		responses = append(
			responses,
			Response{
				StatusCode:  "404",
				ContentType: "application/problem+json",
				Schema:      p.problemSchema(),
			},
		)
	}
	responses = append(
		responses,
		Response{
			StatusCode:  "422",
			ContentType: "application/problem+json",
			Schema:      p.problemSchema(),
		},
		Response{
			StatusCode:  "429",
			ContentType: "application/problem+json",
			Schema:      p.problemSchema(),
		},
		Response{
			StatusCode:  "500",
			ContentType: "application/problem+json",
			Schema:      p.problemSchema(),
		},
	)
	return responses
}

// problemSchema returns the RFC 7807 Problem schema.
func (p *Parser) problemSchema() *Schema {
	return &Schema{
		Name: "Problem",
		Type: PropertyType{Types: []string{SchemaTypeObject}},
		Properties: map[string]*Ref[Schema]{
			"Type": NewInline(&Schema{
				Name:    "Type",
				Type:    PropertyType{Types: []string{SchemaTypeString}},
				GoType:  GoTypeString,
				JSONTag: "type",
			}),
			"Title": NewInline(&Schema{
				Name:    "Title",
				Type:    PropertyType{Types: []string{SchemaTypeString}},
				GoType:  GoTypeString,
				JSONTag: "title",
			}),
			"Status": NewInline(&Schema{
				Name:    "Status",
				Type:    PropertyType{Types: []string{SchemaTypeInteger}},
				GoType:  GoTypeInt,
				JSONTag: "status",
			}),
			"Detail": NewInline(&Schema{
				Name:    "Detail",
				Type:    PropertyType{Types: []string{SchemaTypeString}},
				GoType:  GoTypeString,
				JSONTag: "detail",
			}),
		},
		GoType: "Problem",
	}
}

// getSchemas returns all schemas as a map.
func (p *Parser) getSchemas() map[string]*Schema {
	return p.schemas
}

// Stats holds statistics about an OpenAPI document.
type Stats struct {
	Title                string
	Version              string
	TotalPaths           int
	TotalOperations      int
	TotalSchemas         int
	TotalParameters      int
	TotalResponses       int
	TotalSecuritySchemes int
}

// GetStats computes and returns statistics about the parsed OpenAPI document.
func (p *Parser) GetStats() (*Stats, error) {
	if p.specResult == nil {
		return nil, fmt.Errorf("spec not parsed, call Parse() first")
	}

	// Count paths (unique paths from operations)
	pathSet := make(map[string]struct{})
	for _, op := range p.specResult.Operations {
		pathSet[op.Path] = struct{}{}
	}

	// Count components using discovery
	schemas, _ := DiscoverSchemas(p.fsys)
	responses, _ := DiscoverResponses(p.fsys)
	parameters, _ := DiscoverParameters(p.fsys)

	stats := &Stats{
		Title:                p.specResult.Title,
		Version:              p.specResult.Version,
		TotalPaths:           len(pathSet),
		TotalOperations:      len(p.specResult.Operations),
		TotalSchemas:         len(schemas),
		TotalParameters:      len(parameters),
		TotalResponses:       len(responses),
		TotalSecuritySchemes: len(p.specResult.Security),
	}

	return stats, nil
}

// RenderDocument renders the bundled OpenAPI document in the specified format.
// Parse must be called before RenderDocument.
func (p *Parser) RenderDocument(format string) ([]byte, error) {
	if p.specResult == nil {
		return nil, fmt.Errorf("spec not parsed, call Parse() first")
	}

	// Look for bundled file
	data, err := fs.ReadFile(p.fsys, "openapi.bundled.yaml")
	if err != nil {
		return nil, fmt.Errorf("bundled spec not found: %w", err)
	}

	switch format {
	case RenderFormatYAML:
		return data, nil
	case RenderFormatJSON:
		// Convert YAML to JSON
		var doc map[string]any
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("failed to parse bundled YAML: %w", err)
		}
		return yamlToJSON(doc)
	default:
		return nil, fmt.Errorf("unsupported render format: %s", format)
	}
}

// yamlToJSON converts a YAML document to JSON format.
func yamlToJSON(doc map[string]any) ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}

// Render format constants.
const (
	RenderFormatYAML = "yaml"
	RenderFormatJSON = "json"
)
