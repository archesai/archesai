package spec

import (
	"fmt"
	"io/fs"
	"maps"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/archesai/archesai/internal/yamlutil"
)

// Bundler creates a bundled OpenAPI document from scattered source files.
// It works directly with yaml.Node - no intermediate storage needed.
type Bundler struct {
	doc *OpenAPIDocument
}

// NewBundler creates a new Bundler.
func NewBundler(doc *OpenAPIDocument) *Bundler {
	return &Bundler{doc: doc}
}

// Bundle creates a single bundled YAML document.
// All file $refs are resolved and converted to internal component refs.
func (b *Bundler) Bundle() (*yaml.Node, error) {
	root := &yaml.Node{Kind: yaml.MappingNode}

	// Standard fields
	yamlutil.AddKeyValue(root, "openapi", "3.1.0")
	yamlutil.AddKeyValue(root, "x-project-name", b.doc.doc.XProjectName)

	// Info
	info := b.doc.doc.Info
	infoNode := &yaml.Node{Kind: yaml.MappingNode}
	yamlutil.AddKeyValue(infoNode, "title", info.Title)
	if info.Description != "" {
		yamlutil.AddKeyValue(infoNode, "description", info.Description)
	}
	yamlutil.AddKeyValue(infoNode, "version", info.Version)
	yamlutil.AddKeyValueNode(root, "info", infoNode)

	// Tags
	yamlutil.AddKeyValueNode(root, "tags", b.buildTags())

	// Paths - discover and inline
	yamlutil.AddKeyValueNode(root, "paths", b.buildPaths())

	// Components - discover and inline
	yamlutil.AddKeyValueNode(root, "components", b.buildComponents())

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

// BundleToJSON returns the bundled document as JSON bytes.
func (b *Bundler) BundleToJSON() ([]byte, error) {
	yamlBytes, err := b.BundleToYAML()
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := yaml.Unmarshal(yamlBytes, &doc); err != nil {
		return nil, err
	}
	return yamlToJSON(doc)
}

// Render returns the bundled document in the specified format.
func (b *Bundler) Render(format string) ([]byte, error) {
	switch format {
	case RenderFormatYAML:
		return b.BundleToYAML()
	case RenderFormatJSON:
		return b.BundleToJSON()
	default:
		return nil, fmt.Errorf("unsupported render format: %s", format)
	}
}

func (b *Bundler) buildTags() *yaml.Node {
	node := &yaml.Node{Kind: yaml.SequenceNode}

	// Merge tags from includes and root document
	// Root document tags take precedence over include tags
	tagMap := make(map[string]Tag)

	// First, discover tags from tags/*.yaml files (includes)
	for _, tag := range b.discoverTags() {
		tagMap[tag.Name] = tag
	}

	// Then, add root document tags (override includes with same name)
	for _, tag := range b.doc.doc.Tags {
		tagMap[tag.Name] = tag
	}

	// Sort tag names for deterministic output
	names := make([]string, 0, len(tagMap))
	for name := range tagMap {
		names = append(names, name)
	}
	slices.Sort(names)

	// Build YAML nodes
	for _, name := range names {
		tag := tagMap[name]
		tagNode := &yaml.Node{Kind: yaml.MappingNode}
		yamlutil.AddKeyValue(tagNode, "name", tag.Name)
		if tag.Description != "" {
			yamlutil.AddKeyValue(tagNode, "description", tag.Description)
		}
		node.Content = append(node.Content, tagNode)
	}

	return node
}

// discoverTags loads tags from tags/*.yaml files in the composite filesystem.
// This discovers tags from all enabled includes automatically.
func (b *Bundler) discoverTags() []Tag {
	var tags []Tag

	// Read tags directory from composite FS
	entries, err := fs.ReadDir(b.doc.fsys, "tags")
	if err != nil {
		return tags
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		data, err := fs.ReadFile(b.doc.fsys, "tags/"+entry.Name())
		if err != nil {
			continue
		}

		// Parse tags from file
		var fileTags []Tag
		if err := yaml.Unmarshal(data, &fileTags); err != nil {
			continue
		}
		tags = append(tags, fileTags...)
	}

	return tags
}

func (b *Bundler) buildPaths() *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode}
	paths := make(map[string]*yaml.Node)

	// Get paths from root document
	rootNode := yamlutil.GetContentNode(b.doc.root)
	if pathsNode := yamlutil.FindMappingValue(rootNode, "paths"); pathsNode != nil {
		for i := 0; i < len(pathsNode.Content); i += 2 {
			if i+1 >= len(pathsNode.Content) {
				break
			}
			pathStr := pathsNode.Content[i].Value
			pathNode := b.resolvePathNode(pathsNode.Content[i+1])
			paths[pathStr] = pathNode
		}
	}

	// Discover paths from paths/ directory
	pathFiles, _ := DiscoverPaths(b.doc.fsys)
	for _, filePath := range pathFiles {
		data, err := b.doc.resolver.ReadFile(filePath)
		if err != nil {
			continue
		}
		var pathNode yaml.Node
		if err := yaml.Unmarshal(data, &pathNode); err != nil {
			continue
		}
		contentNode := yamlutil.GetContentNode(&pathNode)

		// Get x-path
		xPath := yamlutil.FindStringValue(contentNode, "x-path")
		if xPath == "" {
			continue
		}
		if _, exists := paths[xPath]; exists {
			continue
		}

		b.resolveRefsInNode(contentNode)
		paths[xPath] = contentNode
	}

	// Sort and add to output
	for _, pathStr := range sortPathKeys(paths) {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: pathStr},
			paths[pathStr],
		)
	}

	return node
}

func (b *Bundler) resolvePathNode(node *yaml.Node) *yaml.Node {
	// If it's a $ref, load and resolve the file
	if ref := yamlutil.FindStringValue(node, "$ref"); ref != "" {
		data, err := b.doc.resolver.ReadFile(ref)
		if err != nil {
			return node
		}
		var resolved yaml.Node
		if err := yaml.Unmarshal(data, &resolved); err != nil {
			return node
		}
		contentNode := yamlutil.GetContentNode(&resolved)
		b.resolveRefsInNode(contentNode)
		return contentNode
	}

	// Clone and resolve refs in place
	cloned := yamlutil.CloneNode(node)
	b.resolveRefsInNode(cloned)
	return cloned
}

func (b *Bundler) buildComponents() *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode}

	// Discover and add each component type
	componentTypes := []struct {
		name string
		kind ComponentKind
	}{
		{"schemas", ComponentSchemas},
		{"responses", ComponentResponses},
		{"parameters", ComponentParameters},
		{"headers", ComponentHeaders},
		{"securitySchemes", ComponentSecuritySchemes},
	}

	for _, ct := range componentTypes {
		components := b.discoverComponentNodes(ct.kind)
		if len(components) > 0 {
			sectionNode := &yaml.Node{Kind: yaml.MappingNode}
			names := slices.Sorted(maps.Keys(components))
			for _, name := range names {
				sectionNode.Content = append(sectionNode.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: name},
					components[name],
				)
			}
			yamlutil.AddKeyValueNode(node, ct.name, sectionNode)
		}
	}

	return node
}

func (b *Bundler) discoverComponentNodes(kind ComponentKind) map[string]*yaml.Node {
	result := make(map[string]*yaml.Node)

	files, err := DiscoverComponents(b.doc.fsys, kind)
	if err != nil {
		return result
	}

	for name, filePath := range files {
		data, err := b.doc.resolver.ReadFile(filePath)
		if err != nil {
			continue
		}

		var node yaml.Node
		if err := yaml.Unmarshal(data, &node); err != nil {
			continue
		}

		contentNode := yamlutil.GetContentNode(&node)
		b.resolveRefsInNode(contentNode)

		// For schemas, use title if available
		componentName := name
		if kind == ComponentSchemas {
			if title := yamlutil.FindStringValue(contentNode, "title"); title != "" {
				componentName = title
			}
		}

		result[componentName] = contentNode
	}

	return result
}

// resolveRefsInNode recursively converts file $refs to internal refs.
func (b *Bundler) resolveRefsInNode(node *yaml.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 >= len(node.Content) {
				break
			}
			key := node.Content[i]
			value := node.Content[i+1]

			if key.Value == "$ref" && value.Kind == yaml.ScalarNode {
				value.Value = FileRefToInternalRef(value.Value, "")
			} else {
				b.resolveRefsInNode(value)
			}
		}
	case yaml.SequenceNode:
		for _, item := range node.Content {
			b.resolveRefsInNode(item)
		}
	}
}

func sortPathKeys(paths map[string]*yaml.Node) []string {
	keys := slices.Collect(maps.Keys(paths))

	slices.SortFunc(keys, func(a, b string) int {
		baseA := strings.Split(a, "/{")[0]
		baseB := strings.Split(b, "/{")[0]

		// Special paths (/health) go last
		specialPaths := map[string]bool{"/health": true}
		if specialPaths[a] != specialPaths[b] {
			if specialPaths[a] {
				return 1
			}
			return -1
		}

		// Same base: parameterized paths first
		if baseA == baseB {
			hasParamA := strings.Contains(a, "{")
			hasParamB := strings.Contains(b, "{")
			if hasParamA != hasParamB {
				if hasParamA {
					return -1
				}
				return 1
			}
		}

		return strings.Compare(baseA, baseB)
	})

	return keys
}
