package spec

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/archesai/archesai/internal/ref"
	schemalib "github.com/archesai/archesai/internal/schema"
	"github.com/archesai/archesai/internal/strutil"
)

// Content type constants.
const (
	contentTypeJSON        = "application/json"
	contentTypeProblemJSON = "application/problem+json"
)

// Parser reads and parses OpenAPI specifications.
type Parser struct {
	doc                  *OpenAPIDocument  // source document (provides FS, resolver, parsed doc)
	schemaLoader         *schemalib.Loader // handles schema property processing
	schemas              map[string]*schemalib.Schema
	refToName            map[string]string // maps ref name (filename) to actual schema name (title)
	specResult           *Spec
	responseContentTypes map[string]string // maps response name to its content type
	configCodegenOnly    []string          // codegen only from arches.yaml
	configCodegenLint    bool              // codegen lint from arches.yaml
	configCodegenOutput  string            // codegen output from arches.yaml
	includes             []string          // enabled includes passed by caller
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

// WithCodegenOutput sets the output directory from external config (arches.yaml).
// When set, tag-based generators output directly to this directory instead of tag subdirectories.
func (p *Parser) WithCodegenOutput(output string) *Parser {
	p.configCodegenOutput = output
	return p
}

// WithIncludes sets the enabled includes for the spec.
// These are passed through to the Spec.EnabledIncludes field.
func (p *Parser) WithIncludes(includes []string) *Parser {
	p.includes = includes
	return p
}

// NewParser creates a new Parser that will parse the given document.
func NewParser(doc *OpenAPIDocument) *Parser {
	return &Parser{
		doc:                  doc,
		schemaLoader:         schemalib.NewLoader(""),
		schemas:              make(map[string]*schemalib.Schema),
		refToName:            make(map[string]string),
		responseContentTypes: make(map[string]string),
	}
}

// FS returns the composite filesystem used by the parser.
// This includes the project filesystem with include filesystems layered underneath.
func (p *Parser) FS() fs.FS {
	return p.doc.fsys
}

// Parse parses the OpenAPI document and returns a Spec.
func (p *Parser) Parse() (*Spec, error) {
	doc := p.doc.doc

	// Load response content types for later lookup
	p.loadResponseContentTypes()

	// Load schemas and operations from OpenAPI spec (including from merged includes)
	operations, err := p.loadSpec(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	// Auto-generate filter, sort, and page parameters for list operations
	listInflater := NewInflater(InflateConfig{ResponseWrappers: false, ListParams: true})
	listInflated := listInflater.Inflate(p.schemas, operations)
	for i := range operations {
		if params, ok := listInflated.ListParams[operations[i].ID]; ok {
			operations[i].Parameters = append(operations[i].Parameters, params...)
		}
	}

	// Convert Tags and Security to spec types
	var specTags []Tag
	for _, t := range doc.Tags {
		specTags = append(specTags, Tag{Name: t.Name, Description: t.Description})
	}

	specSecurity := make(map[string]SecScheme)
	for name, s := range doc.Components.SecuritySchemes {
		specSecurity[name] = SecScheme{
			Type:        s.Type,
			Scheme:      s.Scheme,
			Name:        s.Name,
			In:          s.In,
			Description: s.Description,
		}
	}

	// Build the spec with all metadata
	projectName := doc.XProjectName
	if projectName == "" {
		projectName = doc.Info.Title
	}

	// Determine codegen options (config takes precedence)
	codegenOnly := p.configCodegenOnly
	codegenLint := p.configCodegenLint

	p.specResult = &Spec{
		ProjectName:     projectName,
		Operations:      operations,
		Schemas:         p.getSchemas(),
		EnabledIncludes: p.includes,
		Title:           doc.Info.Title,
		Description:     doc.Info.Description,
		Version:         doc.Info.Version,
		Tags:            specTags,
		Security:        specSecurity,
		CodegenOnly:     codegenOnly,
		CodegenLint:     codegenLint,
		CodegenOutput:   p.configCodegenOutput,
	}

	return p.specResult, nil
}

// loadResponseContentTypes discovers and caches the content types for all response definitions.
func (p *Parser) loadResponseContentTypes() {
	responses, err := DiscoverComponents(p.doc.fsys, ComponentResponses)
	if err != nil {
		return
	}

	for name, filePath := range responses {
		data, err := fs.ReadFile(p.doc.fsys, filePath)
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
		if _, ok := resp.Content[contentTypeProblemJSON]; ok {
			p.responseContentTypes[name] = contentTypeProblemJSON
		} else if _, ok := resp.Content[contentTypeJSON]; ok {
			p.responseContentTypes[name] = contentTypeJSON
		}
	}
}

// loadSpec extracts schemas and operations from the parsed RawDocument.
func (p *Parser) loadSpec(doc *RawDocument) ([]Operation, error) {
	resolver := p.doc.resolver

	// Load schemas from components/schemas (both explicit refs and auto-discovered)
	if err := p.loadSchemas(doc, resolver); err != nil {
		return nil, fmt.Errorf("failed to load schemas: %w", err)
	}

	// Second pass: resolve refs to actual schema names
	for _, schema := range p.schemas {
		p.resolveSchemaRefs(schema)
	}

	// Auto-generate response schemas for all entities
	inflater := NewInflater(InflateConfig{ResponseWrappers: true, ListParams: false})
	inflated := inflater.Inflate(p.schemas, nil)
	for name, schema := range inflated.ResponseSchemas {
		p.schemas[name] = schema
	}

	// Load operations from paths (both explicit refs and auto-discovered)
	operations, err := p.loadPaths(doc, resolver)
	if err != nil {
		return nil, fmt.Errorf("failed to load paths: %w", err)
	}

	return operations, nil
}

// loadSchemas loads schemas from components/schemas section and auto-discovered files.
func (p *Parser) loadSchemas(doc *RawDocument, resolver *ref.FileResolver) error {
	// First, load schemas from components/schemas (both $ref and inline)
	for schemaName, schemaRef := range doc.Components.Schemas {
		if schemaRef.IsRef() {
			// $ref to external file
			data, err := resolver.ReadFile(schemaRef.RefPath)
			if err != nil {
				return fmt.Errorf("failed to resolve schema ref %s: %w", schemaRef.RefPath, err)
			}

			schema, err := p.loadSchemaFile(data, schemaName)
			if err != nil {
				return fmt.Errorf("failed to load schema %s: %w", schemaName, err)
			}

			// Map the filename to the actual schema name
			refName := ExtractSchemaNameFromRef(schemaRef.RefPath)
			p.refToName[refName] = schema.Title

			p.schemas[schema.Title] = schema
		} else if schemaRef.IsInline() {
			// Inline schema definition - the value is already a *Schema
			inlineSchema := schemaRef.GetOrNil()
			if inlineSchema == nil {
				continue
			}
			// Set the name if not already set
			if inlineSchema.Title == "" {
				inlineSchema.Title = schemaName
			}
			// Process the inline schema
			p.processInlineSchema(inlineSchema, schemaName)
			p.schemas[inlineSchema.Title] = inlineSchema
		}
	}

	// Auto-discover schemas from components/schemas/ directory
	discovered, err := DiscoverComponents(p.doc.fsys, ComponentSchemas)
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

		data, err := fs.ReadFile(p.doc.fsys, filePath)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", filePath, err)
		}

		schema, err := p.loadSchemaFile(data, schemaName)
		if err != nil {
			return fmt.Errorf("failed to load schema %s: %w", schemaName, err)
		}

		p.refToName[schemaName] = schema.Title
		p.schemas[schema.Title] = schema
	}

	// Second pass: resolve ref properties now that all schemas are loaded
	for _, schema := range p.schemas {
		p.resolveRefProperties(schema)
	}

	return nil
}

// processInlineSchema processes an inline Schema parsed from YAML.
// It normalizes property names, sets computed fields like JSONTag/GoType.
func (p *Parser) processInlineSchema(schema *schemalib.Schema, defaultName string) {
	// Derive name from title or default
	if schema.Title == "" {
		schema.Title = defaultName
	}

	// Set default type if not set
	if schema.Type.PrimaryType() == "" {
		schema.Type = schemalib.PropertyType{Types: []string{schemalib.TypeObject}}
	}

	// Initialize properties map if nil
	if schema.Properties == nil {
		schema.Properties = make(map[string]*ref.Ref[schemalib.Schema])
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
					p.processPropertySchema(propSchema, fieldName, propName, schema.Title)
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
		processedProps := make(map[string]*ref.Ref[schemalib.Schema])
		for propName, propRef := range schema.Properties {
			fieldName := strutil.PascalCase(propName)
			propSchema := propRef.GetOrNil()
			if propSchema != nil {
				p.processPropertySchema(propSchema, fieldName, propName, schema.Title)
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
	if schema.XCodegenSchemaType == schemalib.TypeEntity {
		p.addBaseFieldsToSchema(schema)
	}

	// Set GoType with schemas alias prefix if this schema comes from an include
	if schema.XInternal != "" {
		alias := schemalib.GetSchemasAlias(schema.XInternal)
		if alias != "" {
			schema.GoType = alias + "." + schema.Title
		} else {
			schema.GoType = schema.Title
		}
	} else {
		schema.GoType = schema.Title
	}
}

// processPropertySchema sets computed fields on a property schema.
// Delegates to SchemaLoader.ProcessProperty.
func (p *Parser) processPropertySchema(
	propSchema *schemalib.Schema,
	fieldName, jsonName, parentName string,
) {
	p.schemaLoader.ProcessProperty(propSchema, fieldName, jsonName, parentName)
}

// loadSchemaFile parses an OpenAPI schema file.
func (p *Parser) loadSchemaFile(
	data []byte,
	defaultName string,
) (*schemalib.Schema, error) {
	var sf schemalib.Schema
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Derive name from title or default
	name := sf.Title
	if name == "" {
		name = defaultName
	}

	schema := &schemalib.Schema{
		Title:              name,
		Description:        sf.Description,
		Type:               schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
		Properties:         make(map[string]*ref.Ref[schemalib.Schema]),
		Required:           sf.Required,
		XCodegen:           sf.XCodegen,
		XCodegenSchemaType: sf.XCodegenSchemaType,
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
					propSchema := &schemalib.Schema{
						Title:   fieldName,
						GoType:  refSchema.Title,
						Type:    refSchema.Type,
						JSONTag: propName,
						YAMLTag: propName,
					}
					schema.Properties[fieldName] = ref.NewInline(propSchema)
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
	if sf.XCodegenSchemaType == schemalib.TypeEntity {
		p.addBaseFieldsToSchema(schema)
	}

	// Set GoType with schemas alias prefix if this schema comes from an include
	if schema.XInternal != "" {
		alias := schemalib.GetSchemasAlias(schema.XInternal)
		if alias != "" {
			schema.GoType = alias + "." + schema.Title
		} else {
			schema.GoType = schema.Title
		}
	} else {
		schema.GoType = schema.Title
	}

	return schema, nil
}

// loadPaths loads operations from the paths section and auto-discovered files.
func (p *Parser) loadPaths(
	doc *RawDocument,
	resolver *ref.FileResolver,
) ([]Operation, error) {
	var operations []Operation
	loadedPaths := make(map[string]bool) // Track which paths have been loaded

	// First, load explicitly referenced paths (both $ref and inline)
	for pathStr, pathRef := range doc.Paths {
		var pathItem *RawPathItem
		var pathFileDir string

		if pathRef.IsRef() {
			// Resolve the $ref to load the path file
			data, err := resolver.ReadFile(pathRef.Ref)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path ref %s: %w", pathRef.Ref, err)
			}
			var pi RawPathItem
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
		ops := p.processPathItem(pathStr, pathItem, pathFileDir, resolver)
		operations = append(operations, ops...)
		loadedPaths[pathStr] = true
	}

	// Auto-discover paths from paths/ directory
	pathFiles, err := DiscoverPaths(p.doc.fsys)
	if err != nil {
		return nil, fmt.Errorf("failed to discover paths: %w", err)
	}

	for _, filePath := range pathFiles {
		data, err := fs.ReadFile(p.doc.fsys, filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read path file %s: %w", filePath, err)
		}

		var pathItem RawPathItem
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
		ops := p.processPathItem(pathStr, &pathItem, pathFileDir, resolver)
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
	pathItem *RawPathItem,
	pathFileDir string,
	resolver *ref.FileResolver,
) []Operation {
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

		op := p.operationToSpecOperation(pathStr, method, opDef, pathFileDir, resolver)
		operations = append(operations, *op)
	}

	return operations
}

// operationToSpecOperation converts an OpenAPI operation to internal Operation.
func (p *Parser) operationToSpecOperation(
	pathStr, method string,
	opDef *RawOperation,
	pathFileDir string,
	resolver *ref.FileResolver,
) *Operation {
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
	} else if securityReqs == nil && len(p.doc.doc.Components.SecuritySchemes) > 0 {
		// Inherit default security from project
		for name := range p.doc.doc.Components.SecuritySchemes {
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
			if scheme, ok := p.doc.doc.Components.SecuritySchemes[name]; ok {
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
	for _, errResp := range StandardErrorResponses(hasPathParam) {
		if !existingCodes[errResp.StatusCode] {
			op.Responses = append(op.Responses, errResp)
		}
	}

	// Sort responses by status code for consistent output
	sort.Slice(op.Responses, func(i, j int) bool {
		return op.Responses[i].StatusCode < op.Responses[j].StatusCode
	})

	return op
}

// paramToSpecParam converts an OpenAPI parameter to internal Param.
func (p *Parser) paramToSpecParam(param *RawParameter) Param {
	paramType := schemalib.PropertyType{Types: []string{schemalib.TypeString}}
	format := ""
	if param.Schema != nil {
		if param.Schema.Type != "" {
			paramType = schemalib.PropertyType{Types: []string{param.Schema.Type}}
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

	schema := &schemalib.Schema{
		Title:       strutil.PascalCase(param.Name),
		Type:        paramType,
		Format:      format,
		Description: param.Description,
		GoType:      schemalib.ToGoType(paramType.PrimaryType(), format),
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
		if param.In == paramLocationPath {
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
						Schema: &schemalib.Schema{
							Title: strutil.PascalCase(paramName),
							Type: schemalib.PropertyType{
								Types: []string{schemalib.TypeString},
							},
							Format:      schemalib.FormatUUID,
							GoType:      schemalib.GoTypeUUID,
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
	resolver *ref.FileResolver,
) *RequestBody {
	// Handle $ref to request body file
	if reqBody.Ref != "" {
		data, err := resolver.ReadFileFrom(pathFileDir, reqBody.Ref)
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
	resolver *ref.FileResolver,
) *RequestBody {
	var schema *schemalib.Schema

	// Get the JSON content type schema
	if content, ok := reqBody.Content[contentTypeJSON]; ok {
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
				data, err := resolver.ReadFileFrom(pathFileDir, schemaRef.RefPath)
				if err == nil {
					parsedSchema, parseErr := p.loadSchemaFile(data, refName)
					if parseErr == nil {
						schema = parsedSchema
						p.schemas[refName] = schema
					}
				}
				if schema == nil {
					schema = &schemalib.Schema{
						Title:  refName,
						GoType: refName,
						Type:   schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
					}
				}
			}
		} else {
			// Inline schema definition
			inlineSchema := schemaRef.GetOrNil()
			if inlineSchema != nil && len(inlineSchema.Properties) > 0 {
				schemaName := opID + "Request"
				schema = &schemalib.Schema{
					Title:      schemaName,
					Type:       schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
					Properties: make(map[string]*ref.Ref[schemalib.Schema]),
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
	resolver *ref.FileResolver,
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
	resolver *ref.FileResolver,
) *Response {
	// Handle $ref to response file
	if respRef.Ref != "" {
		// Extract response name from ref path
		responseName := ExtractSchemaNameFromRef(respRef.Ref)

		// First check if we have a pre-generated response schema
		if responseSchema := p.resolveSchemaByRef(responseName); responseSchema != nil {
			respContentType := p.responseContentTypes[responseName]
			if respContentType == "" {
				respContentType = contentTypeJSON
			}
			return &Response{
				StatusCode:  statusCode,
				ContentType: respContentType,
				Schema:      responseSchema,
			}
		}

		// Resolve as file reference
		data, err := resolver.ReadFileFrom(pathFileDir, respRef.Ref)
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
func (p *Parser) resolveSchemaByRef(refName string) *schemalib.Schema {
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
	_ string,
	_ *ref.FileResolver,
) *Response {
	// Handle content - check for both application/json and application/problem+json
	var content *RawMediaType
	var contentType string

	if c, ok := resp.Content[contentTypeProblemJSON]; ok {
		content = &c
		contentType = contentTypeProblemJSON
	} else if c, ok := resp.Content[contentTypeJSON]; ok {
		content = &c
		contentType = contentTypeJSON
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
				if dataProp != nil && dataProp.Type.PrimaryType() == schemalib.TypeArray &&
					dataProp.Items != nil {
					// This is a list response - look up the pre-generated ListResponse schema
					if dataProp.Items.IsRef() {
						itemRefName := ExtractSchemaNameFromRef(dataProp.Items.RefPath)
						if itemSchema := p.resolveSchemaByRef(itemRefName); itemSchema != nil {
							responseName := itemSchema.Title + "ListResponse"
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
						responseName := itemSchema.Title + "Response"
						if responseSchema := p.resolveSchemaByRef(responseName); responseSchema != nil {
							return &Response{
								StatusCode:  statusCode,
								ContentType: contentType,
								Schema:      responseSchema,
							}
						}
						// Response schema doesn't exist - create it on the fly
						// Use schema name for GoType - templates will add appropriate prefix
						dataGoType := itemSchema.Title
						responseSchema := &schemalib.Schema{
							Title: responseName,
							Type:  schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
							Properties: map[string]*ref.Ref[schemalib.Schema]{
								"Data": ref.NewInline(&schemalib.Schema{
									Title:   "Data",
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
						responseName := itemSchema.Title + "Response"
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
		emptySchema := &schemalib.Schema{
			Title:      "NoContent",
			Type:       schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
			Properties: make(map[string]*ref.Ref[schemalib.Schema]),
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
		emptySchema := &schemalib.Schema{
			Title:      "NoContent",
			Type:       schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
			Properties: make(map[string]*ref.Ref[schemalib.Schema]),
		}
		return &Response{
			StatusCode:  statusCode,
			ContentType: "",
			Schema:      emptySchema,
		}
	}

	customSchema := &schemalib.Schema{
		Title:      opID + "Output",
		Type:       schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
		Properties: make(map[string]*ref.Ref[schemalib.Schema]),
		GoType:     opID + "Output",
	}
	p.schemas[customSchema.Title] = customSchema
	return &Response{
		StatusCode:  statusCode,
		ContentType: contentTypeJSON,
		Schema:      customSchema,
	}
}

// createInlineResponseDef creates a response def from inline schema.
func (p *Parser) createInlineResponseDef(
	statusCode, opID string,
	schema *schemalib.Schema,
	contentType string,
) *Response {
	responseName := opID + "Output"

	// Build required map for property lookup
	requiredMap := make(map[string]bool)
	for _, req := range schema.Required {
		requiredMap[req] = true
	}

	responseSchema := &schemalib.Schema{
		Title:      responseName,
		Type:       schemalib.PropertyType{Types: []string{schemalib.TypeObject}},
		Properties: make(map[string]*ref.Ref[schemalib.Schema]),
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
				propSchema := &schemalib.Schema{
					Title:   fieldName,
					GoType:  refSchema.Title, // Use schema name, template adds schemas. prefix
					Type:    refSchema.Type,
					JSONTag: propName,
					YAMLTag: propName,
				}
				if !requiredMap[propName] {
					propSchema.JSONTag = propName + ",omitempty"
					propSchema.YAMLTag = propName + ",omitempty"
				}
				responseSchema.Properties[fieldName] = ref.NewInline(propSchema)
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
func (p *Parser) resolveRefProperties(schema *schemalib.Schema) {
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
				propSchema := &schemalib.Schema{
					Title:   fieldName,
					GoType:  refSchema.Title,
					Type:    refSchema.Type,
					JSONTag: jsonTag,
					YAMLTag: jsonTag,
				}
				schema.Properties[propName] = ref.NewInline(propSchema)
			}
		}
	}
}

// resolveSchemaRefs resolves ref names to actual schema names for all properties.
func (p *Parser) resolveSchemaRefs(schema *schemalib.Schema) {
	for _, propRef := range schema.Properties {
		p.resolvePropertyRef(propRef.GetOrNil())
	}
}

// resolvePropertyRef resolves a single property's ref to the actual schema name.
func (p *Parser) resolvePropertyRef(prop *schemalib.Schema) {
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
func (p *Parser) addBaseFieldsToSchema(schema *schemalib.Schema) {
	schemalib.AddBaseFields(schema)
}

// getSchemas returns all schemas as a map.
func (p *Parser) getSchemas() map[string]*schemalib.Schema {
	return p.schemas
}

// yamlToJSON converts a YAML document to JSON format.
func yamlToJSON(doc map[string]any) ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}

// FileRefToInternalRef converts a file $ref to an internal document ref.
// E.g., "../schemas/User.yaml" -> "#/components/schemas/User"
// If defaultKind is provided, it's used when the component type can't be determined from the path.
func FileRefToInternalRef(ref string, defaultKind string) string {
	if ref == "" {
		return ref
	}

	// Already an internal ref
	if strings.HasPrefix(ref, "#/") {
		return ref
	}

	// Extract the component name from the file path
	name := ExtractSchemaNameFromRef(ref)
	if name == "" {
		return ref
	}

	// Determine the component type based on the path
	refLower := strings.ToLower(ref)
	switch {
	case strings.Contains(refLower, "/schemas/") || strings.Contains(refLower, "schemas/"):
		return "#/components/schemas/" + name
	case strings.Contains(refLower, "/responses/") || strings.Contains(refLower, "responses/"):
		return "#/components/responses/" + name
	case strings.Contains(refLower, "/parameters/") || strings.Contains(refLower, "parameters/"):
		return "#/components/parameters/" + name
	case strings.Contains(refLower, "/headers/") || strings.Contains(refLower, "headers/"):
		return "#/components/headers/" + name
	case strings.Contains(refLower, "/securityschemes/") || strings.Contains(refLower, "securityschemes/"):
		return "#/components/securitySchemes/" + name
	default:
		// Use default kind if provided, otherwise assume schemas
		if defaultKind != "" {
			return "#/components/" + defaultKind + "/" + name
		}
		return "#/components/schemas/" + name
	}
}

// InternalRefToName extracts the component name from an internal ref.
// E.g., "#/components/schemas/User" -> "User"
func InternalRefToName(ref string) string {
	if !strings.HasPrefix(ref, "#/") {
		return ""
	}
	parts := strings.Split(ref, "/")
	if len(parts) < 4 {
		return ""
	}
	return parts[len(parts)-1]
}

// ExtractSchemaNameFromRef extracts the component name from a file $ref path.
// E.g., "./User.yaml" -> "User", "../schemas/User.yaml" -> "User"
func ExtractSchemaNameFromRef(ref string) string {
	if ref == "" {
		return ""
	}
	// Get the base filename
	base := path.Base(ref)
	// Remove .yaml or .yml extension
	if strings.HasSuffix(base, ".yaml") {
		return strings.TrimSuffix(base, ".yaml")
	}
	if strings.HasSuffix(base, ".yml") {
		return strings.TrimSuffix(base, ".yml")
	}
	return base
}

// Render format constants.
const (
	RenderFormatYAML = "yaml"
	RenderFormatJSON = "json"
)
