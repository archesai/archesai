package parsers

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

// resolvePathItems resolves pathItems references by inlining content into paths
func resolvePathItems(filePath string) error {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML into a Node to preserve order
	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// The root node is a Document node, get the actual content
	if rootNode.Kind != yaml.DocumentNode || len(rootNode.Content) == 0 {
		return fmt.Errorf("invalid YAML structure")
	}

	docNode := rootNode.Content[0]
	if docNode.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node at root")
	}

	// Find components.pathItems
	var pathItemsNode *yaml.Node
	componentsNode := findMapValue(docNode, "components")
	if componentsNode != nil && componentsNode.Kind == yaml.MappingNode {
		pathItemsNode = findMapValue(componentsNode, "pathItems")
	}

	if pathItemsNode == nil || pathItemsNode.Kind != yaml.MappingNode {
		return nil // No pathItems to resolve
	}

	// Find paths section
	pathsNode := findMapValue(docNode, "paths")
	if pathsNode == nil || pathsNode.Kind != yaml.MappingNode {
		return nil // No paths section
	}

	// Iterate through paths and resolve references
	for i := 0; i < len(pathsNode.Content); i += 2 {
		pathValue := pathsNode.Content[i+1]
		if pathValue.Kind != yaml.MappingNode {
			continue
		}

		// Check if this is a $ref node
		refNode := findMapValue(pathValue, "$ref")
		if refNode != nil && refNode.Kind == yaml.ScalarNode {
			refValue := refNode.Value
			if after, ok := strings.CutPrefix(refValue, "#/components/pathItems/"); ok {
				// Extract pathItem name
				pathItemName := after

				// Find the pathItem content
				pathItemContent := findMapValue(pathItemsNode, pathItemName)
				if pathItemContent != nil {
					// Replace the reference with the actual content
					pathsNode.Content[i+1] = pathItemContent
				}
			}
		}
	}

	// Remove pathItems from components
	removeMapKey(componentsNode, "pathItems")

	// Marshal back to YAML
	output, err := yaml.Marshal(&rootNode)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(filePath, output, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// findMapValue finds a value in a mapping node by key
func findMapValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// removeMapKey removes a key-value pair from a mapping node
func removeMapKey(node *yaml.Node, key string) {
	if node == nil || node.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			// Remove both key and value
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			return
		}
	}
}
