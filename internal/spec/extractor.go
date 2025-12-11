package spec

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/archesai/archesai/internal/strutil"
	"github.com/archesai/archesai/internal/yamlutil"
	"github.com/archesai/archesai/pkg/server"
)

// ExtractorConfig holds configuration for the extractor.
type ExtractorConfig struct {
	SpecPath string
	Output   string
	DryRun   bool
	Force    bool
	Verbose  bool
	Includes []string // Includes from arches.yaml config
}

// ExtractionResult holds the results of an extraction operation.
type ExtractionResult struct {
	SchemasExtracted       int
	RequestBodiesExtracted int
	ResponsesExtracted     int
	ParametersExtracted    int
	HeadersExtracted       int
	FilesModified          int
	Skipped                int
	ExtractedFiles         []string
	ModifiedFiles          []string
}

// InlineDefinition represents a detected inline definition.
type InlineDefinition struct {
	Type       string     // "schema", "response", "parameter", "header", "requestBody"
	Name       string     // Generated name
	SourceFile string     // File where the inline definition was found
	SourcePath []string   // JSON path within the file
	Node       *yaml.Node // The YAML node to extract
	TargetPath string     // Target file path for the extracted component
	RefPath    string     // The $ref path to use in the replacement
}

// Extractor handles extracting inline OpenAPI definitions.
type Extractor struct {
	config    ExtractorConfig
	specDir   string
	outputDir string
	fsys      fs.FS
	existing  map[string]bool // Track existing component files
}

// NewExtractor creates a new Extractor instance.
func NewExtractor(config ExtractorConfig) *Extractor {
	specDir := filepath.Dir(config.SpecPath)
	outputDir := config.Output
	if outputDir == "" {
		outputDir = specDir
	}

	return &Extractor{
		config:    config,
		specDir:   specDir,
		outputDir: outputDir,
		existing:  make(map[string]bool),
	}
}

// getServerIncludeFS returns the embedded server API filesystem.
// The returned FS is already stripped of the "api/" prefix.
func (e *Extractor) getServerIncludeFS() fs.FS {
	subFS, err := fs.Sub(server.API, "api")
	if err != nil {
		return nil
	}
	return subFS
}

// hasInclude checks if a specific include is in the config includes list.
func (e *Extractor) hasInclude(name string) bool {
	for _, inc := range e.config.Includes {
		if inc == name {
			return true
		}
	}
	return false
}

// Extract performs the extraction process.
func (e *Extractor) Extract() (*ExtractionResult, error) {
	result := &ExtractionResult{}

	// Initialize filesystem - use CompositeFS if server include is enabled
	if e.hasInclude("server") {
		serverFS := e.getServerIncludeFS()
		if serverFS != nil {
			e.fsys = NewCompositeFS(
				serverFS,
				os.DirFS(e.specDir),
			)
			if e.config.Verbose {
				fmt.Println("Using composite filesystem with embedded server spec")
			}
		} else {
			e.fsys = os.DirFS(e.specDir)
		}
	} else {
		e.fsys = os.DirFS(e.specDir)
	}

	// 1. Discover existing components to avoid conflicts
	if err := e.discoverExisting(); err != nil {
		return nil, fmt.Errorf("failed to discover existing components: %w", err)
	}

	// 2. Find all inline definitions in paths
	definitions, err := e.findInlineDefinitions()
	if err != nil {
		return nil, fmt.Errorf("failed to find inline definitions: %w", err)
	}

	if e.config.Verbose {
		fmt.Printf("Found %d inline definitions to extract\n", len(definitions))
	}

	// 3. Process each definition
	modifiedFiles := make(map[string]bool)
	for _, def := range definitions {
		extracted, err := e.processDefinition(def, result)
		if err != nil {
			return nil, fmt.Errorf("failed to process %s %s: %w", def.Type, def.Name, err)
		}
		if extracted {
			modifiedFiles[def.SourceFile] = true
		}
	}

	result.FilesModified = len(modifiedFiles)
	for file := range modifiedFiles {
		result.ModifiedFiles = append(result.ModifiedFiles, file)
	}

	return result, nil
}

// discoverExisting populates the existing map with known component files.
func (e *Extractor) discoverExisting() error {
	// Discover schemas
	schemas, err := DiscoverSchemas(e.fsys)
	if err != nil {
		return err
	}
	for name := range schemas {
		e.existing["schemas/"+name] = true
	}

	// Discover responses
	responses, err := DiscoverResponses(e.fsys)
	if err != nil {
		return err
	}
	for name := range responses {
		e.existing["responses/"+name] = true
	}

	// Discover parameters
	parameters, err := DiscoverParameters(e.fsys)
	if err != nil {
		return err
	}
	for name := range parameters {
		e.existing["parameters/"+name] = true
	}

	// Discover headers
	headers, err := DiscoverHeaders(e.fsys)
	if err != nil {
		return err
	}
	for name := range headers {
		e.existing["headers/"+name] = true
	}

	// Discover request bodies
	requestBodies, err := e.discoverRequestBodies()
	if err != nil {
		return err
	}
	for name := range requestBodies {
		e.existing["requestBodies/"+name] = true
	}

	return nil
}

// discoverRequestBodies finds all YAML files in components/requestBodies/.
func (e *Extractor) discoverRequestBodies() (map[string]string, error) {
	dir := "components/requestBodies"
	// Check if directory exists
	if _, err := fs.Stat(e.fsys, dir); err != nil {
		if isNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	entries, err := fs.ReadDir(e.fsys, dir)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			ext := path.Ext(name)
			baseName := strings.TrimSuffix(name, ext)
			result[baseName] = path.Join(dir, name)
		}
	}

	return result, nil
}

// findInlineDefinitions scans all files for inline definitions.
func (e *Extractor) findInlineDefinitions() ([]InlineDefinition, error) {
	var definitions []InlineDefinition

	// Get all path files
	pathFiles, err := DiscoverPaths(e.fsys)
	if err != nil {
		return nil, err
	}

	for _, pathFile := range pathFiles {
		fileDefs, err := e.scanPathFile(pathFile)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s: %w", pathFile, err)
		}
		definitions = append(definitions, fileDefs...)
	}

	// Scan existing response files for inline schemas
	responseFiles, err := DiscoverResponses(e.fsys)
	if err != nil {
		return nil, err
	}

	for name, filePath := range responseFiles {
		fileDefs, err := e.scanResponseFile(name, filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s: %w", filePath, err)
		}
		definitions = append(definitions, fileDefs...)
	}

	// Note: We don't scan requestBody files for inline schemas because
	// the code generator handles request body schemas inline in the handlers.

	return definitions, nil
}

// scanPathFile scans a single path file for inline definitions.
func (e *Extractor) scanPathFile(filePath string) ([]InlineDefinition, error) {
	var definitions []InlineDefinition

	data, err := fs.ReadFile(e.fsys, filePath)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	// The root is a document node, get the actual content
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, nil
	}

	pathNode := root.Content[0]
	if pathNode.Kind != yaml.MappingNode {
		return nil, nil
	}

	// Scan each HTTP method (get, post, put, patch, delete)
	methods := []string{"get", "post", "put", "patch", "delete", "options", "head"}
	for _, method := range methods {
		methodNode := e.getMapValue(pathNode, method)
		if methodNode == nil {
			continue
		}

		operationID := e.getStringValue(methodNode, "operationId")

		// Note: We don't scan requestBody for inline definitions because
		// the code generator handles request body schemas inline in the handlers.

		// Scan responses for inline definitions
		respDefs := e.scanResponses(methodNode, filePath, method, operationID)
		definitions = append(definitions, respDefs...)

		// Scan parameters for inline definitions
		paramDefs := e.scanParameters(methodNode, filePath, method, operationID)
		definitions = append(definitions, paramDefs...)
	}

	return definitions, nil
}

// scanResponses scans responses for inline definitions.
func (e *Extractor) scanResponses(
	methodNode *yaml.Node,
	filePath, method, operationID string,
) []InlineDefinition {
	var definitions []InlineDefinition

	responses := e.getMapValue(methodNode, "responses")
	if responses == nil {
		return nil
	}

	// Iterate through all response status codes
	for i := 0; i < len(responses.Content); i += 2 {
		statusCode := responses.Content[i].Value
		respNode := responses.Content[i+1]

		// Skip if it's a $ref
		if e.hasRef(respNode) {
			continue
		}

		// Check if this response has content (indicating it's an inline response)
		content := e.getMapValue(respNode, "content")
		if content == nil {
			continue
		}

		name := e.deriveResponseName(operationID, statusCode, method)

		definitions = append(definitions, InlineDefinition{
			Type:       "response",
			Name:       name,
			SourceFile: filePath,
			SourcePath: []string{method, "responses", statusCode},
			Node:       respNode,
			TargetPath: path.Join("components", "responses", name+".yaml"),
			RefPath:    "../components/responses/" + name + ".yaml",
		})
	}

	return definitions
}

// scanParameters scans parameters for inline definitions.
func (e *Extractor) scanParameters(
	methodNode *yaml.Node,
	filePath, method, operationID string,
) []InlineDefinition {
	var definitions []InlineDefinition

	params := e.getSequenceValue(methodNode, "parameters")
	if params == nil {
		return nil
	}

	for idx, paramNode := range params.Content {
		// Skip if it's a $ref
		if e.hasRef(paramNode) {
			continue
		}

		// Check if this parameter has name and in (required for inline parameter)
		paramName := e.getStringValue(paramNode, "name")
		if paramName == "" {
			continue
		}

		name := strutil.PascalCase(paramName)

		definitions = append(definitions, InlineDefinition{
			Type:       "parameter",
			Name:       name,
			SourceFile: filePath,
			SourcePath: []string{method, "parameters", fmt.Sprintf("%d", idx)},
			Node:       paramNode,
			TargetPath: path.Join("components", "parameters", name+".yaml"),
			RefPath:    "../components/parameters/" + name + ".yaml",
		})
	}

	return definitions
}

// scanResponseFile scans an existing response file for inline schemas.
func (e *Extractor) scanResponseFile(responseName, filePath string) ([]InlineDefinition, error) {
	var definitions []InlineDefinition

	data, err := fs.ReadFile(e.fsys, filePath)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, nil
	}

	respNode := root.Content[0]
	if respNode.Kind != yaml.MappingNode {
		return nil, nil
	}

	// Check if this response has inline schema in content.application/json.schema
	content := e.getMapValue(respNode, "content")
	if content == nil {
		return nil, nil
	}

	jsonContent := e.getMapValue(content, "application/json")
	if jsonContent == nil {
		return nil, nil
	}

	schema := e.getMapValue(jsonContent, "schema")
	if schema == nil {
		return nil, nil
	}

	// Check if the schema is inline (not a $ref)
	if e.hasRef(schema) {
		return nil, nil
	}

	// Check if it has properties.data with an inline object
	// This is the common pattern: { data: { ...actual schema... } }
	properties := e.getMapValue(schema, "properties")
	if properties != nil {
		dataNode := e.getMapValue(properties, "data")
		if dataNode != nil && !e.hasRef(dataNode) {
			// Check if data is an inline object (not an array)
			if e.isInlineObjectSchema(dataNode) {
				// Extract the data property's inline schema
				schemaName := responseName
				if strings.HasSuffix(schemaName, "Response") {
					schemaName = strings.TrimSuffix(schemaName, "Response")
				}

				definitions = append(definitions, InlineDefinition{
					Type:       "schema",
					Name:       schemaName,
					SourceFile: filePath,
					SourcePath: []string{
						"content",
						"application/json",
						"schema",
						"properties",
						"data",
					},
					Node:       dataNode,
					TargetPath: path.Join("components", "schemas", schemaName+".yaml"),
					RefPath:    "../schemas/" + schemaName + ".yaml",
				})
				return definitions, nil
			}

			// Check if data is an array with inline items
			if e.getStringValue(dataNode, "type") == "array" {
				items := e.getMapValue(dataNode, "items")
				if items != nil && e.isInlineObjectSchema(items) {
					// Extract the array items schema
					schemaName := responseName
					if strings.HasSuffix(schemaName, "ListResponse") {
						schemaName = strings.TrimSuffix(schemaName, "ListResponse") + "Item"
					} else if strings.HasSuffix(schemaName, "Response") {
						schemaName = strings.TrimSuffix(schemaName, "Response") + "Item"
					}

					definitions = append(definitions, InlineDefinition{
						Type:       "schema",
						Name:       schemaName,
						SourceFile: filePath,
						SourcePath: []string{
							"content",
							"application/json",
							"schema",
							"properties",
							"data",
							"items",
						},
						Node:       items,
						TargetPath: path.Join("components", "schemas", schemaName+".yaml"),
						RefPath:    "../schemas/" + schemaName + ".yaml",
					})
					return definitions, nil
				}
			}
		}
	}

	// Do NOT extract the entire response wrapper schema - the code generator handles
	// list responses (data + meta) and single responses (data only) specially.
	// Only inline object schemas that aren't wrapped in data should be extracted.

	return definitions, nil
}

// isInlineObjectSchema checks if a node is an inline object schema.
func (e *Extractor) isInlineObjectSchema(node *yaml.Node) bool {
	if node.Kind != yaml.MappingNode {
		return false
	}

	// Must have properties to be considered an object schema
	hasProperties := e.getMapValue(node, "properties") != nil
	if !hasProperties {
		return false
	}

	// Must NOT already be a $ref
	if e.hasRef(node) {
		return false
	}

	return true
}

// processDefinition extracts a single inline definition.
func (e *Extractor) processDefinition(
	def InlineDefinition,
	result *ExtractionResult,
) (bool, error) {
	// Check for conflicts
	key := def.Type + "s/" + def.Name
	if def.Type == "requestBody" {
		key = "requestBodies/" + def.Name
	}
	if e.existing[key] && !e.config.Force {
		if e.config.Verbose {
			fmt.Printf("Skipping %s %s: already exists\n", def.Type, def.Name)
		}
		result.Skipped++
		return false, nil
	}

	if e.config.DryRun {
		fmt.Printf("Would extract %s: %s -> %s\n", def.Type, def.Name, def.TargetPath)
		// For schemas, also check for nested inline objects
		if def.Type == "schema" {
			nestedDefs := e.findNestedInlineSchemas(def.Node, def.Name, def.TargetPath)
			for _, nestedDef := range nestedDefs {
				fmt.Printf(
					"Would extract %s: %s -> %s\n",
					nestedDef.Type,
					nestedDef.Name,
					nestedDef.TargetPath,
				)
			}
		}
		return false, nil
	}

	if e.config.Verbose {
		fmt.Printf("Extracting %s: %s -> %s\n", def.Type, def.Name, def.TargetPath)
	}

	// For schemas, first extract nested inline objects
	if def.Type == "schema" {
		if err := e.extractNestedSchemas(def.Node, def.Name, def.TargetPath, result); err != nil {
			return false, fmt.Errorf("failed to extract nested schemas: %w", err)
		}
	}

	// 1. Write the extracted component file
	if err := e.writeComponentFile(def); err != nil {
		return false, fmt.Errorf("failed to write component file: %w", err)
	}

	// 2. Update the source file to use $ref
	if err := e.updateSourceFile(def); err != nil {
		return false, fmt.Errorf("failed to update source file: %w", err)
	}

	// Update result counts
	switch def.Type {
	case "schema":
		result.SchemasExtracted++
	case "requestBody":
		result.RequestBodiesExtracted++
	case "response":
		result.ResponsesExtracted++
	case "parameter":
		result.ParametersExtracted++
	case "header":
		result.HeadersExtracted++
	}

	result.ExtractedFiles = append(result.ExtractedFiles, def.TargetPath)

	// Mark as existing so we don't try to extract again
	if def.Type == "requestBody" {
		e.existing["requestBodies/"+def.Name] = true
	} else {
		e.existing[def.Type+"s/"+def.Name] = true
	}

	return true, nil
}

// findNestedInlineSchemas finds all nested inline object schemas within a schema node.
func (e *Extractor) findNestedInlineSchemas(
	node *yaml.Node,
	parentName, parentPath string,
) []InlineDefinition {
	var definitions []InlineDefinition

	if node.Kind != yaml.MappingNode {
		return nil
	}

	// Check properties for nested inline objects
	properties := e.getMapValue(node, "properties")
	if properties != nil && properties.Kind == yaml.MappingNode {
		for i := 0; i < len(properties.Content); i += 2 {
			propName := properties.Content[i].Value
			propValue := properties.Content[i+1]

			// Check if this property is an inline object schema
			if e.isInlineObjectSchema(propValue) {
				nestedName := parentName + strutil.PascalCase(propName)
				definitions = append(definitions, InlineDefinition{
					Type:       "schema",
					Name:       nestedName,
					SourceFile: parentPath,
					SourcePath: []string{"properties", propName},
					Node:       propValue,
					TargetPath: path.Join("components", "schemas", nestedName+".yaml"),
					RefPath:    nestedName + ".yaml",
				})
			}

			// Check if this is an array with inline object items
			if e.getStringValue(propValue, "type") == "array" {
				items := e.getMapValue(propValue, "items")
				if items != nil && e.isInlineObjectSchema(items) {
					nestedName := parentName + strutil.PascalCase(propName) + "Item"
					definitions = append(definitions, InlineDefinition{
						Type:       "schema",
						Name:       nestedName,
						SourceFile: parentPath,
						SourcePath: []string{"properties", propName, "items"},
						Node:       items,
						TargetPath: path.Join("components", "schemas", nestedName+".yaml"),
						RefPath:    nestedName + ".yaml",
					})
				}
			}
		}
	}

	// Check items for array type (top-level)
	if e.getStringValue(node, "type") == "array" {
		items := e.getMapValue(node, "items")
		if items != nil && e.isInlineObjectSchema(items) {
			nestedName := parentName + "Item"
			definitions = append(definitions, InlineDefinition{
				Type:       "schema",
				Name:       nestedName,
				SourceFile: parentPath,
				SourcePath: []string{"items"},
				Node:       items,
				TargetPath: path.Join("components", "schemas", nestedName+".yaml"),
				RefPath:    nestedName + ".yaml",
			})
		}
	}

	return definitions
}

// extractNestedSchemas extracts all nested inline schemas from a schema node.
func (e *Extractor) extractNestedSchemas(
	node *yaml.Node,
	parentName, parentPath string,
	result *ExtractionResult,
) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	// Check properties for nested inline objects
	properties := e.getMapValue(node, "properties")
	if properties != nil && properties.Kind == yaml.MappingNode {
		for i := 0; i < len(properties.Content); i += 2 {
			propName := properties.Content[i].Value
			propValue := properties.Content[i+1]

			// Check if this property is an inline object schema
			if e.isInlineObjectSchema(propValue) {
				nestedName := parentName + strutil.PascalCase(propName)
				nestedTargetPath := path.Join("components", "schemas", nestedName+".yaml")

				// Recursively extract nested schemas first
				if err := e.extractNestedSchemas(propValue, nestedName, nestedTargetPath, result); err != nil {
					return err
				}

				// Write the nested schema file
				nestedDef := InlineDefinition{
					Type:       "schema",
					Name:       nestedName,
					Node:       propValue,
					TargetPath: nestedTargetPath,
				}
				if err := e.writeComponentFile(nestedDef); err != nil {
					return err
				}

				if e.config.Verbose {
					fmt.Printf("Extracting nested schema: %s -> %s\n", nestedName, nestedTargetPath)
				}
				result.SchemasExtracted++
				result.ExtractedFiles = append(result.ExtractedFiles, nestedTargetPath)

				// Replace the inline object with a $ref
				refNode := &yaml.Node{Kind: yaml.MappingNode}
				refKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "$ref"}
				refVal := &yaml.Node{Kind: yaml.ScalarNode, Value: nestedName + ".yaml"}
				refNode.Content = []*yaml.Node{refKey, refVal}
				properties.Content[i+1] = refNode
			}

			// Check if this is an array with inline object items
			if e.getStringValue(propValue, "type") == "array" {
				items := e.getMapValue(propValue, "items")
				if items != nil && e.isInlineObjectSchema(items) {
					nestedName := parentName + strutil.PascalCase(propName) + "Item"
					nestedTargetPath := path.Join("components", "schemas", nestedName+".yaml")

					// Recursively extract nested schemas first
					if err := e.extractNestedSchemas(items, nestedName, nestedTargetPath, result); err != nil {
						return err
					}

					// Write the nested schema file
					nestedDef := InlineDefinition{
						Type:       "schema",
						Name:       nestedName,
						Node:       items,
						TargetPath: nestedTargetPath,
					}
					if err := e.writeComponentFile(nestedDef); err != nil {
						return err
					}

					if e.config.Verbose {
						fmt.Printf(
							"Extracting nested schema: %s -> %s\n",
							nestedName,
							nestedTargetPath,
						)
					}
					result.SchemasExtracted++
					result.ExtractedFiles = append(result.ExtractedFiles, nestedTargetPath)

					// Replace the inline items with a $ref
					e.replaceItemsWithRef(propValue, nestedName+".yaml")
				}
			}
		}
	}

	return nil
}

// replaceItemsWithRef replaces the items node in an array schema with a $ref.
func (e *Extractor) replaceItemsWithRef(arrayNode *yaml.Node, refPath string) {
	for i := 0; i < len(arrayNode.Content); i += 2 {
		if arrayNode.Content[i].Value == "items" {
			refNode := &yaml.Node{Kind: yaml.MappingNode}
			refKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "$ref"}
			refVal := &yaml.Node{Kind: yaml.ScalarNode, Value: refPath}
			refNode.Content = []*yaml.Node{refKey, refVal}
			arrayNode.Content[i+1] = refNode
			return
		}
	}
}

// writeComponentFile writes an extracted component to a file.
func (e *Extractor) writeComponentFile(def InlineDefinition) error {
	targetPath := filepath.Join(e.outputDir, def.TargetPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Marshal with consistent formatting
	output, err := yamlutil.MarshalOpenAPI(def.Node)
	if err != nil {
		return err
	}

	return os.WriteFile(targetPath, output, 0644)
}

// updateSourceFile updates the source file to replace inline definition with $ref.
func (e *Extractor) updateSourceFile(def InlineDefinition) error {
	sourceAbsPath := filepath.Join(e.specDir, def.SourceFile)

	data, err := os.ReadFile(sourceAbsPath)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	// Navigate to the parent node and replace the inline definition with $ref
	if err := e.replaceWithRef(&root, def); err != nil {
		return err
	}

	// Marshal back
	output, err := yamlutil.MarshalOpenAPI(root.Content[0])
	if err != nil {
		return err
	}

	return os.WriteFile(sourceAbsPath, output, 0644)
}

// replaceWithRef navigates to the target location and replaces it with a $ref node.
func (e *Extractor) replaceWithRef(root *yaml.Node, def InlineDefinition) error {
	current := root.Content[0] // Document content

	// Navigate to parent of target node
	for i := 0; i < len(def.SourcePath)-1; i++ {
		pathPart := def.SourcePath[i]

		if current.Kind == yaml.MappingNode {
			found := false
			for j := 0; j < len(current.Content); j += 2 {
				if current.Content[j].Value == pathPart {
					current = current.Content[j+1]
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("path not found: %s", pathPart)
			}
		} else if current.Kind == yaml.SequenceNode {
			// pathPart is an index
			var idx int
			if _, err := fmt.Sscanf(pathPart, "%d", &idx); err != nil {
				return fmt.Errorf("invalid index: %s", pathPart)
			}
			if idx >= len(current.Content) {
				return fmt.Errorf("index out of range: %s", pathPart)
			}
			current = current.Content[idx]
		}
	}

	// Replace the target key's value with $ref node
	lastKey := def.SourcePath[len(def.SourcePath)-1]

	switch current.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(current.Content); i += 2 {
			if current.Content[i].Value == lastKey {
				refNode := &yaml.Node{Kind: yaml.MappingNode}
				refKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "$ref"}
				refVal := &yaml.Node{Kind: yaml.ScalarNode, Value: def.RefPath}
				refNode.Content = []*yaml.Node{refKey, refVal}
				current.Content[i+1] = refNode
				return nil
			}
		}
	case yaml.SequenceNode:
		// For parameters in a sequence, replace the element at the index
		var idx int
		if _, err := fmt.Sscanf(lastKey, "%d", &idx); err != nil {
			return fmt.Errorf("invalid index: %s", lastKey)
		}
		if idx >= len(current.Content) {
			return fmt.Errorf("index out of range: %s", lastKey)
		}
		refNode := &yaml.Node{Kind: yaml.MappingNode}
		refKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "$ref"}
		refVal := &yaml.Node{Kind: yaml.ScalarNode, Value: def.RefPath}
		refNode.Content = []*yaml.Node{refKey, refVal}
		current.Content[idx] = refNode
		return nil
	}

	return fmt.Errorf("key not found: %s", lastKey)
}

// Helper methods

func (e *Extractor) hasRef(node *yaml.Node) bool {
	return e.getStringValue(node, "$ref") != ""
}

func (e *Extractor) getStringValue(node *yaml.Node, key string) string {
	if node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1].Value
		}
	}
	return ""
}

func (e *Extractor) getMapValue(node *yaml.Node, key string) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func (e *Extractor) getSequenceValue(node *yaml.Node, key string) *yaml.Node {
	val := e.getMapValue(node, key)
	if val != nil && val.Kind == yaml.SequenceNode {
		return val
	}
	return nil
}

func (e *Extractor) deriveResponseName(operationID, statusCode, method string) string {
	if operationID != "" {
		if statusCode == "200" || statusCode == "201" {
			return operationID + "Response"
		}
		return operationID + statusCode + "Response"
	}
	// Fallback to method-based name
	return strutil.PascalCase(method) + statusCode + "Response"
}
