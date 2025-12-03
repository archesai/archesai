package parsers

import (
	"gopkg.in/yaml.v3"
)

// injectDefaults adds default values to the bundled OpenAPI spec.
// This includes:
// - Adding 500 InternalServerError response to all operations
// - Ensuring Problem schema exists for error responses
func injectDefaults(data []byte) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return data, nil
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return data, nil
	}

	// Ensure components.schemas.Problem exists
	ensureProblemSchema(root)

	// Ensure components.responses.InternalServerError exists
	ensureInternalServerErrorResponse(root)

	// Add 500 response to all operations
	injectDefaultResponses(root)

	return yaml.Marshal(&doc)
}

// ensureProblemSchema ensures the Problem schema exists in components.schemas
func ensureProblemSchema(root *yaml.Node) {
	components := findOrCreateMapKey(root, "components")
	if components == nil {
		return
	}

	schemas := findOrCreateMapKey(components, "schemas")
	if schemas == nil {
		return
	}

	// Check if Problem already exists
	if findMapKey(schemas, "Problem") != nil {
		return
	}

	// Add Problem schema
	problemSchema := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "type"},
			{Kind: yaml.ScalarNode, Value: "object"},
			{Kind: yaml.ScalarNode, Value: "description"},
			{Kind: yaml.ScalarNode, Value: "RFC 7807 Problem Details"},
			{Kind: yaml.ScalarNode, Value: "properties"},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "type"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "type"},
						{Kind: yaml.ScalarNode, Value: "string"},
					}},
					{Kind: yaml.ScalarNode, Value: "title"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "type"},
						{Kind: yaml.ScalarNode, Value: "string"},
					}},
					{Kind: yaml.ScalarNode, Value: "status"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "type"},
						{Kind: yaml.ScalarNode, Value: "integer"},
					}},
					{Kind: yaml.ScalarNode, Value: "detail"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "type"},
						{Kind: yaml.ScalarNode, Value: "string"},
					}},
					{Kind: yaml.ScalarNode, Value: "instance"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "type"},
						{Kind: yaml.ScalarNode, Value: "string"},
					}},
				},
			},
		},
	}

	schemas.Content = append(schemas.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "Problem"},
		problemSchema,
	)
}

// ensureInternalServerErrorResponse ensures the InternalServerError response exists
func ensureInternalServerErrorResponse(root *yaml.Node) {
	components := findOrCreateMapKey(root, "components")
	if components == nil {
		return
	}

	responses := findOrCreateMapKey(components, "responses")
	if responses == nil {
		return
	}

	// Check if InternalServerError already exists
	if findMapKey(responses, "InternalServerError") != nil {
		return
	}

	// Add InternalServerError response
	errorResponse := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "description"},
			{Kind: yaml.ScalarNode, Value: "Internal server error"},
			{Kind: yaml.ScalarNode, Value: "content"},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "application/problem+json"},
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "schema"},
							{
								Kind: yaml.MappingNode,
								Content: []*yaml.Node{
									{Kind: yaml.ScalarNode, Value: "$ref"},
									{Kind: yaml.ScalarNode, Value: "#/components/schemas/Problem"},
								},
							},
						},
					},
				},
			},
		},
	}

	responses.Content = append(responses.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "InternalServerError"},
		errorResponse,
	)
}

// injectDefaultResponses adds 500 response to all operations that don't have one
func injectDefaultResponses(root *yaml.Node) {
	paths := findMapKey(root, "paths")
	if paths == nil {
		return
	}

	// Iterate through all paths
	for i := 0; i < len(paths.Content); i += 2 {
		if i+1 >= len(paths.Content) {
			break
		}
		pathItem := paths.Content[i+1]
		if pathItem.Kind != yaml.MappingNode {
			continue
		}

		// Check each HTTP method
		methods := []string{"get", "post", "put", "patch", "delete", "head", "options"}
		for _, method := range methods {
			operation := findMapKey(pathItem, method)
			if operation == nil {
				continue
			}

			responses := findOrCreateMapKey(operation, "responses")
			if responses == nil {
				continue
			}

			// Check if 500 response exists
			if findMapKey(responses, "500") != nil {
				continue
			}

			// Add 500 response reference
			responses.Content = append(responses.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "500"},
				&yaml.Node{
					Kind: yaml.MappingNode,
					Content: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "$ref"},
						{
							Kind:  yaml.ScalarNode,
							Value: "#/components/responses/InternalServerError",
						},
					},
				},
			)
		}
	}
}

// findMapKey finds a key in a mapping node and returns its value
func findMapKey(node *yaml.Node, key string) *yaml.Node {
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

// findOrCreateMapKey finds or creates a key in a mapping node
func findOrCreateMapKey(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	// Try to find existing key
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	// Create new key
	newValue := &yaml.Node{Kind: yaml.MappingNode}
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		newValue,
	)
	return newValue
}
