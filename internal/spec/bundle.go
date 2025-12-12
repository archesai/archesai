package spec

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/archesai/archesai/internal/yamlutil"
	"github.com/archesai/archesai/pkg/server"
)

// YAML node helper functions

func addKeyValue(node *yaml.Node, key string, value any) {
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		valueToNode(value),
	)
}

func addKeyValueNode(node *yaml.Node, key string, valueNode *yaml.Node) {
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		valueNode,
	)
}

func valueToNode(value any) *yaml.Node {
	switch v := value.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: v}
	case bool:
		if v {
			return &yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"}
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Value: "false", Tag: "!!bool"}
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%d", v)}
	default:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", v)}
	}
}

// findMappingValue finds a value in a mapping node by key.
func findMappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// findStringValue finds a string value in a mapping node by key.
func findStringValue(node *yaml.Node, key string) string {
	n := findMappingValue(node, key)
	if n == nil || n.Kind != yaml.ScalarNode {
		return ""
	}
	return n.Value
}

// getContentNode returns the content node from a document node.
func getContentNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0]
	}
	return node
}

// Bundler creates a bundled OpenAPI document from source files.
// Operates directly on yaml.Node - no domain model conversion.
// This preserves all $ref names and structures from the source.
type Bundler struct {
	doc  *OpenAPIDocument
	spec *Spec // Optional: for generating filter/sort/page parameters

	// Collected components to add to the bundled output
	schemas         map[string]*yaml.Node
	responses       map[string]*yaml.Node
	parameters      map[string]*yaml.Node
	headers         map[string]*yaml.Node
	securitySchemes map[string]*yaml.Node

	// Track entities that have list operations (for filter/sort params)
	listEntities map[string]*Schema
}

// NewBundler creates a new Bundler.
func NewBundler(doc *OpenAPIDocument) *Bundler {
	return &Bundler{
		doc:             doc,
		schemas:         make(map[string]*yaml.Node),
		responses:       make(map[string]*yaml.Node),
		parameters:      make(map[string]*yaml.Node),
		headers:         make(map[string]*yaml.Node),
		securitySchemes: make(map[string]*yaml.Node),
		listEntities:    make(map[string]*Schema),
	}
}

// WithSpec sets the spec for generating filter/sort/page parameters.
func (b *Bundler) WithSpec(s *Spec) *Bundler {
	b.spec = s

	// Identify entities with list operations
	if s != nil {
		for _, op := range s.GetOperations() {
			// Check responses for 200 with list response schema
			for _, resp := range op.Responses {
				if resp.StatusCode == "200" && resp.Schema != nil {
					schemaName := resp.Name
					if strings.HasSuffix(schemaName, "ListResponse") {
						entityName := strings.TrimSuffix(schemaName, "ListResponse")
						if entitySchema := b.findSpecSchema(entityName); entitySchema != nil {
							b.listEntities[entityName] = entitySchema
						}
					}
				}
			}
		}
	}

	return b
}

func (b *Bundler) findSpecSchema(name string) *Schema {
	if b.spec == nil {
		return nil
	}
	return b.spec.GetSchema(name)
}

// Bundle creates a single bundled YAML document.
// All file $refs are resolved and converted to internal component refs.
func (b *Bundler) Bundle() (*yaml.Node, error) {
	// First, collect all components from discovered files
	if err := b.collectComponents(); err != nil {
		return nil, fmt.Errorf("failed to collect components: %w", err)
	}

	// Build the bundled document
	root := &yaml.Node{Kind: yaml.MappingNode}

	// Add standard fields
	addKeyValue(root, "openapi", "3.1.0")
	addKeyValue(root, "x-project-name", b.doc.ProjectName())

	// Info section
	title, desc, version := b.doc.Info()
	infoNode := &yaml.Node{Kind: yaml.MappingNode}
	addKeyValue(infoNode, "title", title)
	if desc != "" {
		addKeyValue(infoNode, "description", desc)
	}
	addKeyValue(infoNode, "version", version)
	addKeyValueNode(root, "info", infoNode)

	// Tags section
	tagsNode, err := b.buildTagsNode()
	if err != nil {
		return nil, fmt.Errorf("failed to build tags: %w", err)
	}
	addKeyValueNode(root, "tags", tagsNode)

	// Paths section
	pathsNode, err := b.buildPathsNode()
	if err != nil {
		return nil, fmt.Errorf("failed to build paths: %w", err)
	}
	addKeyValueNode(root, "paths", pathsNode)

	// Components section
	componentsNode, err := b.buildComponentsNode()
	if err != nil {
		return nil, fmt.Errorf("failed to build components: %w", err)
	}
	addKeyValueNode(root, "components", componentsNode)

	return root, nil
}

// BundleToYAML returns the bundled document as YAML bytes.
func (b *Bundler) BundleToYAML() ([]byte, error) {
	node, err := b.Bundle()
	if err != nil {
		return nil, err
	}
	return yamlutil.MarshalOpenAPI(node)
}

// collectComponents discovers and collects all components from the filesystem.
func (b *Bundler) collectComponents() error {
	fsys := b.doc.FS()

	// Collect schemas
	schemas, err := DiscoverSchemas(fsys)
	if err != nil {
		return err
	}
	for name, filePath := range schemas {
		if err := b.loadSchemaFile(name, filePath); err != nil {
			return fmt.Errorf("failed to load schema %s: %w", name, err)
		}
	}

	// Collect responses
	responses, err := DiscoverResponses(fsys)
	if err != nil {
		return err
	}
	for name, filePath := range responses {
		if err := b.loadResponseFile(name, filePath); err != nil {
			return fmt.Errorf("failed to load response %s: %w", name, err)
		}
	}

	// Collect parameters
	parameters, err := DiscoverParameters(fsys)
	if err != nil {
		return err
	}
	for name, filePath := range parameters {
		if err := b.loadParameterFile(name, filePath); err != nil {
			return fmt.Errorf("failed to load parameter %s: %w", name, err)
		}
	}

	// Collect headers
	headers, err := DiscoverHeaders(fsys)
	if err != nil {
		return err
	}
	for name, filePath := range headers {
		if err := b.loadHeaderFile(name, filePath); err != nil {
			return fmt.Errorf("failed to load header %s: %w", name, err)
		}
	}

	// Collect security schemes
	securitySchemes, err := DiscoverSecuritySchemes(fsys)
	if err != nil {
		return err
	}
	for name, filePath := range securitySchemes {
		if err := b.loadSecuritySchemeFile(name, filePath); err != nil {
			return fmt.Errorf("failed to load security scheme %s: %w", name, err)
		}
	}

	return nil
}

func (b *Bundler) loadSchemaFile(name, filePath string) error {
	data, err := b.doc.ResolveFileRef(".", filePath)
	if err != nil {
		return err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	contentNode := getContentNode(&node)

	// Use title if available, otherwise the filename
	schemaName := name
	if title := findStringValue(contentNode, "title"); title != "" {
		schemaName = title
	}

	b.schemas[schemaName] = contentNode
	return nil
}

func (b *Bundler) loadResponseFile(name, filePath string) error {
	data, err := b.doc.ResolveFileRef(".", filePath)
	if err != nil {
		return err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	b.responses[name] = getContentNode(&node)
	return nil
}

func (b *Bundler) loadParameterFile(name, filePath string) error {
	data, err := b.doc.ResolveFileRef(".", filePath)
	if err != nil {
		return err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	b.parameters[name] = getContentNode(&node)
	return nil
}

func (b *Bundler) loadHeaderFile(name, filePath string) error {
	data, err := b.doc.ResolveFileRef(".", filePath)
	if err != nil {
		return err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	b.headers[name] = getContentNode(&node)
	return nil
}

func (b *Bundler) loadSecuritySchemeFile(name, filePath string) error {
	data, err := b.doc.ResolveFileRef(".", filePath)
	if err != nil {
		return err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	b.securitySchemes[name] = getContentNode(&node)
	return nil
}

func (b *Bundler) buildTagsNode() (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.SequenceNode}

	// Get tags from document
	tags := b.doc.Tags()

	// Add server include tags if enabled
	if b.doc.HasInclude("server") {
		serverTags := b.loadServerTags()
		tags = append(tags, serverTags...)
	}

	for _, tag := range tags {
		tagNode := &yaml.Node{Kind: yaml.MappingNode}
		addKeyValue(tagNode, "name", tag.Name)
		if tag.Description != "" {
			addKeyValue(tagNode, "description", tag.Description)
		}
		node.Content = append(node.Content, tagNode)
	}

	return node, nil
}

func (b *Bundler) loadServerTags() []Tag {
	data, err := fs.ReadFile(server.API, "api/openapi.yaml")
	if err != nil {
		return nil
	}

	var info struct {
		Tags []Tag `yaml:"tags"`
	}
	if err := yaml.Unmarshal(data, &info); err != nil {
		return nil
	}

	return info.Tags
}

func (b *Bundler) buildPathsNode() (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Collect paths from document and discovered files
	paths, err := b.collectPaths()
	if err != nil {
		return nil, err
	}

	// Sort paths
	sortedPaths := b.sortPaths(paths)

	for _, pathStr := range sortedPaths {
		pathNode := paths[pathStr]

		// If this is a $ref, resolve it
		if ref := findStringValue(pathNode, "$ref"); ref != "" {
			data, err := b.doc.ResolveFileRef(".", ref)
			if err != nil {
				continue
			}

			var resolved yaml.Node
			if err := yaml.Unmarshal(data, &resolved); err != nil {
				continue
			}

			pathNode = getContentNode(&resolved)
		}

		// Build the path item node
		pathItemNode, err := b.buildPathItemNode(pathNode)
		if err != nil {
			return nil, err
		}

		addKeyValueNode(node, pathStr, pathItemNode)
	}

	return node, nil
}

// collectPaths gathers all paths from the document and discovered files.
func (b *Bundler) collectPaths() (map[string]*yaml.Node, error) {
	result := make(map[string]*yaml.Node)

	// Get paths from the document
	rootNode := b.doc.contentNode()
	pathsNode := findMappingValue(rootNode, "paths")
	if pathsNode != nil && pathsNode.Kind == yaml.MappingNode {
		for i := 0; i < len(pathsNode.Content); i += 2 {
			if i+1 >= len(pathsNode.Content) {
				break
			}
			pathStr := pathsNode.Content[i].Value
			result[pathStr] = pathsNode.Content[i+1]
		}
	}

	// Auto-discover paths from paths/ directory
	pathFiles, err := DiscoverPaths(b.doc.FS())
	if err != nil {
		return nil, err
	}

	for _, filePath := range pathFiles {
		data, err := b.doc.ResolveFileRef(".", filePath)
		if err != nil {
			continue
		}

		var node yaml.Node
		if err := yaml.Unmarshal(data, &node); err != nil {
			continue
		}

		contentNode := getContentNode(&node)
		xPath := findStringValue(contentNode, "x-path")
		if xPath == "" {
			continue
		}

		// Skip if already loaded via explicit ref
		if _, exists := result[xPath]; exists {
			continue
		}

		result[xPath] = contentNode
	}

	return result, nil
}

// sortPaths returns path strings in sorted order.
func (b *Bundler) sortPaths(paths map[string]*yaml.Node) []string {
	var pathStrs []string
	for pathStr := range paths {
		pathStrs = append(pathStrs, pathStr)
	}

	// Sort: parameterized paths first, then alphabetical, special paths last
	sort.Slice(pathStrs, func(i, j int) bool {
		baseI := strings.Split(pathStrs[i], "/{")[0]
		baseJ := strings.Split(pathStrs[j], "/{")[0]

		specialPaths := map[string]bool{"/health": true}
		isSpecialI := specialPaths[pathStrs[i]]
		isSpecialJ := specialPaths[pathStrs[j]]
		if isSpecialI != isSpecialJ {
			return !isSpecialI
		}

		if baseI == baseJ {
			hasParamI := strings.Contains(pathStrs[i], "{")
			hasParamJ := strings.Contains(pathStrs[j], "{")
			if hasParamI != hasParamJ {
				return hasParamI
			}
		}
		return baseI < baseJ
	})

	return pathStrs
}

func (b *Bundler) buildPathItemNode(pathNode *yaml.Node) (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Process operations in standard method order
	methodOrder := []string{"get", "post", "put", "patch", "delete"}
	for _, method := range methodOrder {
		opNode := findMappingValue(pathNode, method)
		if opNode == nil {
			continue
		}

		builtOp, err := b.buildOperationNode(opNode)
		if err != nil {
			return nil, err
		}
		addKeyValueNode(node, method, builtOp)
	}

	return node, nil
}

func (b *Bundler) buildOperationNode(opNode *yaml.Node) (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Add operationId first
	addKeyValue(node, "operationId", findStringValue(opNode, "operationId"))

	if summary := findStringValue(opNode, "summary"); summary != "" {
		addKeyValue(node, "summary", summary)
	}

	if desc := findStringValue(opNode, "description"); desc != "" {
		addKeyValue(node, "description", desc)
	}

	// Security - preserve empty security pattern (before tags per OpenAPI convention)
	securityNode := findMappingValue(opNode, "security")
	if securityNode != nil {
		if b.isEmptySecurity(securityNode) {
			secNode := &yaml.Node{Kind: yaml.SequenceNode}
			emptyObjNode := &yaml.Node{Kind: yaml.MappingNode}
			secNode.Content = append(secNode.Content, emptyObjNode)
			addKeyValueNode(node, "security", secNode)
		} else {
			addKeyValueNode(node, "security", b.cloneNode(securityNode))
		}
	}

	// Tags - extract first tag
	if tagsNode := findMappingValue(opNode, "tags"); tagsNode != nil &&
		tagsNode.Kind == yaml.SequenceNode {
		if len(tagsNode.Content) > 0 {
			tag := tagsNode.Content[0].Value
			newTagsNode := &yaml.Node{Kind: yaml.SequenceNode}
			newTagsNode.Content = append(
				newTagsNode.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: tag},
			)
			addKeyValueNode(node, "tags", newTagsNode)
		}
	}

	// Parameters - rewrite file $refs to internal refs
	if paramsNode := findMappingValue(opNode, "parameters"); paramsNode != nil &&
		paramsNode.Kind == yaml.SequenceNode {
		newParamsNode := b.buildParametersNode(paramsNode)
		if len(newParamsNode.Content) > 0 {
			addKeyValueNode(node, "parameters", newParamsNode)
		}
	}

	// RequestBody - clone as-is (refs are already internal)
	if reqBodyNode := findMappingValue(opNode, "requestBody"); reqBodyNode != nil {
		addKeyValueNode(node, "requestBody", b.cloneNode(reqBodyNode))
	}

	// Responses - preserve original $refs
	respNode, err := b.buildResponsesNode(opNode)
	if err != nil {
		return nil, err
	}
	addKeyValueNode(node, "responses", respNode)

	// Add extensions at the end
	if findStringValue(opNode, "x-codegen-custom-handler") == "true" {
		addKeyValue(node, "x-codegen-custom-handler", true)
	}
	if xInternal := findStringValue(opNode, "x-internal"); xInternal != "" {
		addKeyValue(node, "x-internal", xInternal)
	}

	return node, nil
}

// isEmptySecurity checks if security is explicitly empty: security: - {}
func (b *Bundler) isEmptySecurity(secNode *yaml.Node) bool {
	if secNode == nil || secNode.Kind != yaml.SequenceNode {
		return false
	}
	// Check for single empty object: - {}
	if len(secNode.Content) == 1 {
		item := secNode.Content[0]
		if item.Kind == yaml.MappingNode && len(item.Content) == 0 {
			return true
		}
	}
	return false
}

// buildParametersNode builds the parameters node (refs are already internal).
func (b *Bundler) buildParametersNode(sourceNode *yaml.Node) *yaml.Node {
	result := &yaml.Node{Kind: yaml.SequenceNode}

	for _, paramItem := range sourceNode.Content {
		if paramItem.Kind == yaml.MappingNode {
			result.Content = append(result.Content, b.cloneNode(paramItem))
		}
	}

	return result
}

func (b *Bundler) buildResponsesNode(opNode *yaml.Node) (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}

	responsesNode := findMappingValue(opNode, "responses")
	if responsesNode == nil || responsesNode.Kind != yaml.MappingNode {
		return node, nil
	}

	// Collect responses as map
	responses := make(map[string]*yaml.Node)
	for i := 0; i < len(responsesNode.Content); i += 2 {
		if i+1 >= len(responsesNode.Content) {
			break
		}
		statusCode := responsesNode.Content[i].Value
		respNode := responsesNode.Content[i+1]
		responses[statusCode] = respNode
	}

	// Sort status codes
	var statusCodes []string
	for code := range responses {
		statusCodes = append(statusCodes, code)
	}
	// Success codes first, then error codes
	sort.Slice(statusCodes, func(i, j int) bool {
		codeI, codeJ := statusCodes[i], statusCodes[j]
		isSuccessI := codeI[0] == '2'
		isSuccessJ := codeJ[0] == '2'
		if isSuccessI != isSuccessJ {
			return isSuccessI
		}
		return codeI < codeJ
	})

	for _, statusCode := range statusCodes {
		respNode := responses[statusCode]

		statusNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: statusCode,
			Style: yaml.SingleQuotedStyle,
		}

		// Check if this response has a $ref
		if ref := findStringValue(respNode, "$ref"); ref != "" {
			// Convert file ref to internal ref
			internalRef := b.convertRefToInternal(ref)
			refNode := &yaml.Node{Kind: yaml.MappingNode}
			addKeyValue(refNode, "$ref", internalRef)
			node.Content = append(node.Content, statusNode, refNode)
		} else {
			// Inline response - clone as-is
			node.Content = append(node.Content, statusNode, b.cloneNode(respNode))
		}
	}

	return node, nil
}

func (b *Bundler) buildComponentsNode() (*yaml.Node, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Per OpenAPI spec, components order: schemas, responses, parameters, examples, requestBodies, headers, securitySchemes, links, callbacks, pathItems

	// Schemas
	if len(b.schemas) > 0 {
		schemasNode := &yaml.Node{Kind: yaml.MappingNode}
		var schemaNames []string
		for name := range b.schemas {
			schemaNames = append(schemaNames, name)
		}
		sort.Strings(schemaNames)
		for _, name := range schemaNames {
			addKeyValueNode(schemasNode, name, b.cloneNode(b.schemas[name]))
		}
		addKeyValueNode(node, "schemas", schemasNode)
	}

	// Responses
	if len(b.responses) > 0 {
		responsesNode := &yaml.Node{Kind: yaml.MappingNode}
		var respNames []string
		for name := range b.responses {
			respNames = append(respNames, name)
		}
		sort.Strings(respNames)
		for _, name := range respNames {
			addKeyValueNode(responsesNode, name, b.cloneNode(b.responses[name]))
		}
		addKeyValueNode(node, "responses", responsesNode)
	}

	// Parameters
	if len(b.parameters) > 0 {
		paramsNode := &yaml.Node{Kind: yaml.MappingNode}
		var paramNames []string
		for name := range b.parameters {
			paramNames = append(paramNames, name)
		}
		sort.Strings(paramNames)
		for _, name := range paramNames {
			addKeyValueNode(paramsNode, name, b.cloneNode(b.parameters[name]))
		}
		addKeyValueNode(node, "parameters", paramsNode)
	}

	// Headers
	if len(b.headers) > 0 {
		headersNode := &yaml.Node{Kind: yaml.MappingNode}
		var headerNames []string
		for name := range b.headers {
			headerNames = append(headerNames, name)
		}
		sort.Strings(headerNames)
		for _, name := range headerNames {
			addKeyValueNode(headersNode, name, b.cloneNode(b.headers[name]))
		}
		addKeyValueNode(node, "headers", headersNode)
	}

	// Security schemes
	if len(b.securitySchemes) > 0 {
		securitySchemesNode := &yaml.Node{Kind: yaml.MappingNode}
		var secSchemeNames []string
		for name := range b.securitySchemes {
			secSchemeNames = append(secSchemeNames, name)
		}
		sort.Strings(secSchemeNames)
		for _, name := range secSchemeNames {
			addKeyValueNode(securitySchemesNode, name, b.cloneNode(b.securitySchemes[name]))
		}
		addKeyValueNode(node, "securitySchemes", securitySchemesNode)
	}

	return node, nil
}

// cloneNode creates a deep copy of a yaml.Node, converting file refs to internal refs.
func (b *Bundler) cloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}

	clone := &yaml.Node{
		Kind:        node.Kind,
		Style:       node.Style,
		Tag:         node.Tag,
		Value:       node.Value,
		Anchor:      node.Anchor,
		Alias:       nil, // Don't copy alias
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
		Line:        0, // Reset position
		Column:      0,
	}

	if len(node.Content) > 0 {
		clone.Content = make([]*yaml.Node, len(node.Content))
		for i, child := range node.Content {
			clone.Content[i] = b.cloneNode(child)
		}

		// Convert file $refs to internal refs
		if node.Kind == yaml.MappingNode {
			for i := 0; i < len(clone.Content); i += 2 {
				if i+1 < len(clone.Content) && clone.Content[i].Value == "$ref" {
					clone.Content[i+1].Value = b.convertRefToInternal(clone.Content[i+1].Value)
				}
			}
		}
	}

	return clone
}

// convertRefToInternal converts a file $ref to an internal document ref.
// E.g., "../schemas/User.yaml" -> "#/components/schemas/User"
func (b *Bundler) convertRefToInternal(ref string) string {
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
		// For refs without a clear component type path, assume schemas
		// This handles cases like "./User.yaml" within the schemas directory
		return "#/components/schemas/" + name
	}
}
