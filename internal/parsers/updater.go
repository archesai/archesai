package parsers

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// UpdateRequiredFields updates the required field in a YAML schema to include properties with defaults
func UpdateRequiredFields(node *yaml.Node) error {
	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return fmt.Errorf("invalid document structure")
	}

	root := node.Content[0]
	if root.Kind != yaml.MappingNode {
		return fmt.Errorf("root is not a mapping node")
	}

	// Find properties with defaults
	propsWithDefaults := findPropertiesWithDefaults(root)
	if len(propsWithDefaults) == 0 {
		// No properties with defaults, nothing to update
		return nil
	}

	// Find or create the required field
	requiredNode := findOrCreateRequiredField(root)
	if requiredNode == nil {
		return fmt.Errorf("failed to find or create required field")
	}

	// Get existing required fields
	existingRequired := getExistingRequired(requiredNode)

	// Combine and deduplicate
	allRequired := combineAndDeduplicate(existingRequired, propsWithDefaults)

	// Update the required node
	updateRequiredNode(requiredNode, allRequired)

	return nil
}

// findPropertiesWithDefaults finds all properties that have default values
func findPropertiesWithDefaults(root *yaml.Node) []string {
	var propsWithDefaults []string

	// Find the properties node
	var propertiesNode *yaml.Node
	for i := 0; i < len(root.Content)-1; i += 2 {
		if root.Content[i].Value == "properties" {
			propertiesNode = root.Content[i+1]
			break
		}
	}

	if propertiesNode == nil || propertiesNode.Kind != yaml.MappingNode {
		return propsWithDefaults
	}

	// Iterate through properties
	for i := 0; i < len(propertiesNode.Content)-1; i += 2 {
		propName := propertiesNode.Content[i].Value
		propDef := propertiesNode.Content[i+1]

		if hasDefaultInNode(propDef) {
			propsWithDefaults = append(propsWithDefaults, propName)
		}
	}

	return propsWithDefaults
}

// hasDefaultInNode checks if a property definition has a default value
func hasDefaultInNode(propNode *yaml.Node) bool {
	if propNode.Kind != yaml.MappingNode {
		return false
	}

	for i := 0; i < len(propNode.Content)-1; i += 2 {
		if propNode.Content[i].Value == "default" {
			return true
		}
	}

	return false
}

// findOrCreateRequiredField finds the required field or creates it
func findOrCreateRequiredField(root *yaml.Node) *yaml.Node {
	// Look for existing required field
	for i := 0; i < len(root.Content)-1; i += 2 {
		if root.Content[i].Value == "required" {
			return root.Content[i+1]
		}
	}

	// Find position to insert (before additionalProperties if it exists)
	insertPos := len(root.Content)
	for i := 0; i < len(root.Content)-1; i += 2 {
		if root.Content[i].Value == "additionalProperties" {
			insertPos = i
			break
		}
	}

	// Create new required field
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "required",
	}
	valueNode := &yaml.Node{
		Kind: yaml.SequenceNode,
	}

	// Insert at the right position
	newContent := make([]*yaml.Node, 0, len(root.Content)+2)
	newContent = append(newContent, root.Content[:insertPos]...)
	newContent = append(newContent, keyNode, valueNode)
	newContent = append(newContent, root.Content[insertPos:]...)
	root.Content = newContent

	return valueNode
}

// getExistingRequired gets the current required fields
func getExistingRequired(requiredNode *yaml.Node) []string {
	var required []string

	if requiredNode.Kind != yaml.SequenceNode {
		return required
	}

	for _, item := range requiredNode.Content {
		if item.Kind == yaml.ScalarNode {
			required = append(required, item.Value)
		}
	}

	return required
}

// combineAndDeduplicate combines two slices and removes duplicates
func combineAndDeduplicate(existing, new []string) []string {
	seen := make(map[string]bool)
	var result []string

	// Add existing
	for _, s := range existing {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	// Add new
	for _, s := range new {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

// updateRequiredNode updates the required node with new values
func updateRequiredNode(requiredNode *yaml.Node, required []string) {
	requiredNode.Kind = yaml.SequenceNode
	requiredNode.Content = make([]*yaml.Node, 0, len(required))

	for _, req := range required {
		node := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: req,
		}
		requiredNode.Content = append(requiredNode.Content, node)
	}
}
