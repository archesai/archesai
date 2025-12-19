package yamlutil

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// AddKeyValue adds a key-value pair to a mapping node.
func AddKeyValue(node *yaml.Node, key string, value any) {
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		ValueToNode(value),
	)
}

// AddKeyValueNode adds a key with a pre-built value node to a mapping node.
func AddKeyValueNode(node *yaml.Node, key string, valueNode *yaml.Node) {
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		valueNode,
	)
}

// ValueToNode converts a Go value to a yaml.Node.
func ValueToNode(value any) *yaml.Node {
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

// FindMappingValue finds a value in a mapping node by key.
func FindMappingValue(node *yaml.Node, key string) *yaml.Node {
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

// FindStringValue finds a string value in a mapping node by key.
func FindStringValue(node *yaml.Node, key string) string {
	n := FindMappingValue(node, key)
	if n == nil || n.Kind != yaml.ScalarNode {
		return ""
	}
	return n.Value
}

// GetContentNode returns the content node from a document node.
func GetContentNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0]
	}
	return node
}

// MapToNode converts a map[string]any to a yaml.Node with sorted keys.
func MapToNode(m map[string]any) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)

	for _, key := range keys {
		value := m[key]
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
		valueNode := AnyToNode(value)
		node.Content = append(node.Content, keyNode, valueNode)
	}

	return node
}

// AnyToNode converts any value to a yaml.Node.
func AnyToNode(value any) *yaml.Node {
	switch v := value.(type) {
	case map[string]any:
		return MapToNode(v)
	case []any:
		node := &yaml.Node{Kind: yaml.SequenceNode}
		for _, item := range v {
			node.Content = append(node.Content, AnyToNode(item))
		}
		return node
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: v}
	case bool:
		node := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool"}
		if v {
			node.Value = "true"
		} else {
			node.Value = "false"
		}
		return node
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", v)}
	case int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", v)}
	case float64:
		// Check if it's actually an integer
		if v == float64(int64(v)) {
			return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", int64(v))}
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%g", v)}
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}
	default:
		// Fallback: marshal and unmarshal
		data, err := yaml.Marshal(v)
		if err != nil {
			return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", v)}
		}
		var node yaml.Node
		if err := yaml.Unmarshal(data, &node); err != nil {
			return &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", v)}
		}
		if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
			return node.Content[0]
		}
		return &node
	}
}

// sortStrings sorts a slice of strings in place.
func sortStrings(s []string) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// CloneNode creates a deep copy of a yaml.Node.
func CloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	clone := &yaml.Node{
		Kind:        node.Kind,
		Style:       node.Style,
		Tag:         node.Tag,
		Value:       node.Value,
		Anchor:      node.Anchor,
		Alias:       CloneNode(node.Alias),
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
		Line:        node.Line,
		Column:      node.Column,
	}
	if len(node.Content) > 0 {
		clone.Content = make([]*yaml.Node, len(node.Content))
		for i, c := range node.Content {
			clone.Content[i] = CloneNode(c)
		}
	}
	return clone
}
