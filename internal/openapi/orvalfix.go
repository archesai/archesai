package openapi

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
	componentsNode := findMapValueOrval(docNode, "components")
	if componentsNode != nil && componentsNode.Kind == yaml.MappingNode {
		pathItemsNode = findMapValueOrval(componentsNode, "pathItems")
	}

	if pathItemsNode == nil || pathItemsNode.Kind != yaml.MappingNode {
		return nil // No pathItems to resolve
	}

	// Find paths section
	pathsNode := findMapValueOrval(docNode, "paths")
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
		refNode := findMapValueOrval(pathValue, "$ref")
		if refNode != nil && refNode.Kind == yaml.ScalarNode {
			refValue := refNode.Value
			if after, ok := strings.CutPrefix(refValue, "#/components/pathItems/"); ok {
				// Extract pathItem name
				pathItemName := after

				// Find the pathItem content
				pathItemContent := findMapValueOrval(pathItemsNode, pathItemName)
				if pathItemContent != nil {
					// Replace the reference with the actual content
					pathsNode.Content[i+1] = pathItemContent
				}
			}
		}
	}

	// Remove pathItems from components
	removeMapKeyOrval(componentsNode, "pathItems")

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

// findMapValueOrval finds a value in a mapping node by key
func findMapValueOrval(node *yaml.Node, key string) *yaml.Node {
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

// removeMapKeyOrval removes a key-value pair from a mapping node
func removeMapKeyOrval(node *yaml.Node, key string) {
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

// cleanupComposedBundle removes duplicate entries created by libopenapi's composed bundler.
// The bundler creates both "Foo: $ref: #/components/.../Foo__suffix" and "Foo__suffix: {...}".
// This function removes the reference entries and renames the suffixed entries to clean names,
// also updating all $ref values throughout the document.
func cleanupComposedBundle(data []byte) ([]byte, error) {
	// Parse YAML
	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if rootNode.Kind != yaml.DocumentNode || len(rootNode.Content) == 0 {
		return nil, fmt.Errorf("invalid YAML structure")
	}

	docNode := rootNode.Content[0]
	if docNode.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected mapping node at root")
	}

	// Find components section
	componentsNode := findMapValueOrval(docNode, "components")
	if componentsNode == nil || componentsNode.Kind != yaml.MappingNode {
		return data, nil // No components, nothing to clean
	}

	// Component types that may have duplicates
	componentTypes := []string{
		"schemas",
		"parameters",
		"responses",
		"headers",
		"requestBodies",
		"securitySchemes",
		"pathItems",
	}

	for _, compType := range componentTypes {
		sectionNode := findMapValueOrval(componentsNode, compType)
		if sectionNode == nil || sectionNode.Kind != yaml.MappingNode {
			continue
		}

		cleanupComponentSection(sectionNode, compType)
	}

	// Update all $ref values throughout the document to remove suffixes
	updateRefs(docNode)

	// Marshal back
	output, err := yaml.Marshal(&rootNode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return output, nil
}

// cleanupComponentSection removes reference entries and renames suffixed entries
func cleanupComponentSection(sectionNode *yaml.Node, compType string) {
	suffix := "__" + compType

	// First pass: identify entries to remove (refs pointing to suffixed names)
	var indicesToRemove []int
	for i := 0; i < len(sectionNode.Content); i += 2 {
		keyNode := sectionNode.Content[i]
		valueNode := sectionNode.Content[i+1]

		// Check if this is a reference to a suffixed version
		if valueNode.Kind == yaml.MappingNode {
			refNode := findMapValueOrval(valueNode, "$ref")
			if refNode != nil && refNode.Kind == yaml.ScalarNode {
				// If ref points to same name with suffix, mark for removal
				expectedRef := fmt.Sprintf("#/components/%s/%s%s", compType, keyNode.Value, suffix)
				if refNode.Value == expectedRef {
					indicesToRemove = append(indicesToRemove, i)
				}
			}
		}
	}

	// Remove reference entries (in reverse order to preserve indices)
	for i := len(indicesToRemove) - 1; i >= 0; i-- {
		idx := indicesToRemove[i]
		sectionNode.Content = append(sectionNode.Content[:idx], sectionNode.Content[idx+2:]...)
	}

	// Second pass: rename suffixed entries to clean names
	for i := 0; i < len(sectionNode.Content); i += 2 {
		keyNode := sectionNode.Content[i]
		keyNode.Value = strings.TrimSuffix(keyNode.Value, suffix)
	}
}

// updateRefs recursively updates all $ref values to remove __suffix patterns
func updateRefs(node *yaml.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if keyNode.Value == "$ref" && valueNode.Kind == yaml.ScalarNode {
				// Remove suffix from ref value
				valueNode.Value = removeRefSuffix(valueNode.Value)
			} else {
				updateRefs(valueNode)
			}
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			updateRefs(child)
		}
	}
}

// removeRefSuffix removes __type suffixes from a $ref value
func removeRefSuffix(ref string) string {
	suffixes := []string{
		"__schemas",
		"__parameters",
		"__responses",
		"__headers",
		"__requestBodies",
		"__securitySchemes",
		"__pathItems",
	}
	for _, suffix := range suffixes {
		ref = strings.ReplaceAll(ref, suffix, "")
	}
	return ref
}
