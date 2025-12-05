package openapi

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/archesai/archesai/pkg/auth"
	"github.com/archesai/archesai/pkg/config"
	"github.com/archesai/archesai/pkg/executor"
	"github.com/archesai/archesai/pkg/pipelines"
	"github.com/archesai/archesai/pkg/server"
	"github.com/archesai/archesai/pkg/storage"
)

// IncludeFS is an interface for embedded filesystems used by includes.
// Both embed.FS and fs.Sub-created filesystems implement this interface.
type IncludeFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

const (
	yamlKeyName = "name"
)

// IncludeSpec defines a spec that can be included via x-include-* extensions
type IncludeSpec struct {
	Name string    // Extension name without "x-include-" prefix (e.g., "auth", "config")
	FS   IncludeFS // Embedded filesystem containing the spec/ directory
}

// IncludeMerger handles merging of OpenAPI specs based on x-include-* extensions
type IncludeMerger struct {
	includeSpecs map[string]IncludeSpec
}

// NewIncludeMerger creates a new IncludeMerger
func NewIncludeMerger() *IncludeMerger {
	return &IncludeMerger{
		includeSpecs: make(map[string]IncludeSpec),
	}
}

// NewDefaultIncludeMerger creates an IncludeMerger with all standard includes registered.
func NewDefaultIncludeMerger() *IncludeMerger {
	merger := NewIncludeMerger()
	merger.RegisterInclude("auth", auth.APISpec)
	merger.RegisterInclude("config", config.APISpec)
	merger.RegisterInclude("server", server.APISpec)
	merger.RegisterInclude("storage", storage.APISpec)
	merger.RegisterInclude("pipelines", pipelines.APISpec)
	merger.RegisterInclude("executor", executor.APISpec)
	return merger
}

// RegisterInclude registers an includable spec with its embedded filesystem.
// The filesystem should contain a spec/openapi.yaml file at its root.
func (m *IncludeMerger) RegisterInclude(name string, fsys IncludeFS) *IncludeMerger {
	m.includeSpecs[name] = IncludeSpec{
		Name: name,
		FS:   fsys,
	}
	return m
}

// BuildCompositeFS creates a composite filesystem from the spec path and enabled includes.
// This allows libopenapi to resolve $ref references from both the local spec directory
// and embedded include filesystems.
func (m *IncludeMerger) BuildCompositeFS(specPath string) (*CompositeFS, []string, error) {
	// Read main spec to find enabled includes
	mainSpecBytes, err := os.ReadFile(specPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read main spec: %w", err)
	}

	// Parse as YAML node to find includes
	var mainSpec yaml.Node
	if err := yaml.Unmarshal(mainSpecBytes, &mainSpec); err != nil {
		return nil, nil, fmt.Errorf("failed to parse main spec: %w", err)
	}

	// Find enabled includes
	enabledIncludes := m.findEnabledIncludes(&mainSpec)

	// Build list of enabled include names
	var enabledNames []string
	for _, include := range enabledIncludes {
		enabledNames = append(enabledNames, include.Name)
	}

	// Create base filesystem from the spec's directory
	specDir := filepath.Dir(specPath)
	baseFS := os.DirFS(specDir)

	// Create composite filesystem
	composite := NewCompositeFS(baseFS)

	// Add each enabled include
	for _, include := range enabledIncludes {
		composite.AddInclude(include.Name, include.FS)
	}

	return composite, enabledNames, nil
}

// MergeSpec reads a spec file, finds x-include-* extensions, and returns merged spec bytes.
// The composite filesystem should be used with libopenapi to resolve $ref references.
func (m *IncludeMerger) MergeSpec(specPath string) ([]byte, []string, error) {
	// Read main spec
	mainSpecBytes, err := os.ReadFile(specPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read main spec: %w", err)
	}

	// Parse as YAML node
	var mainSpec yaml.Node
	if err := yaml.Unmarshal(mainSpecBytes, &mainSpec); err != nil {
		return nil, nil, fmt.Errorf("failed to parse main spec: %w", err)
	}

	// Find enabled includes
	enabledIncludes := m.findEnabledIncludes(&mainSpec)

	// Build list of enabled include names
	var enabledNames []string
	for _, include := range enabledIncludes {
		enabledNames = append(enabledNames, include.Name)
	}

	// If no includes, return original spec bytes
	if len(enabledIncludes) == 0 {
		return mainSpecBytes, enabledNames, nil
	}

	if mainSpec.Kind != yaml.DocumentNode || len(mainSpec.Content) == 0 {
		return nil, nil, fmt.Errorf("invalid spec structure")
	}

	mainRoot := mainSpec.Content[0]

	// For each include, read its spec and merge paths/tags/components/security
	for _, include := range enabledIncludes {
		includeBytes, err := include.FS.ReadFile("spec/openapi.yaml")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read include spec %s: %w", include.Name, err)
		}

		var includeSpec yaml.Node
		if err := yaml.Unmarshal(includeBytes, &includeSpec); err != nil {
			return nil, nil, fmt.Errorf("failed to parse include spec %s: %w", include.Name, err)
		}

		if includeSpec.Kind != yaml.DocumentNode || len(includeSpec.Content) == 0 {
			continue
		}

		includeRoot := includeSpec.Content[0]

		m.mergePaths(mainRoot, includeRoot)
		m.mergeTags(mainRoot, includeRoot)
		m.mergeComponents(mainRoot, includeRoot)
		m.mergeSecurity(mainRoot, includeRoot)
	}

	// Remove x-include-* extensions from output
	m.removeIncludeExtensions(mainRoot)

	// Marshal merged spec
	mergedBytes, err := yaml.Marshal(&mainSpec)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal merged spec: %w", err)
	}

	return mergedBytes, enabledNames, nil
}

// findEnabledIncludes finds all x-include-* extensions set to true
func (m *IncludeMerger) findEnabledIncludes(doc *yaml.Node) []IncludeSpec {
	var enabled []IncludeSpec

	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return enabled
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return enabled
	}

	for i := 0; i < len(root.Content)-1; i += 2 {
		keyNode := root.Content[i]
		valueNode := root.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			continue
		}

		key := keyNode.Value
		if !strings.HasPrefix(key, "x-include-") {
			continue
		}

		// Check if value is true
		if valueNode.Kind == yaml.ScalarNode && valueNode.Value == "true" {
			includeName := strings.TrimPrefix(key, "x-include-")

			// Find matching include spec
			if spec, ok := m.includeSpecs[includeName]; ok {
				enabled = append(enabled, spec)
			}
		}
	}

	return enabled
}

// mergePaths merges paths from included spec into main spec
func (m *IncludeMerger) mergePaths(mainRoot, includeRoot *yaml.Node) {
	includePaths := m.findMapping(includeRoot, "paths")
	if includePaths == nil {
		return
	}

	mainPaths := m.findMapping(mainRoot, "paths")
	if mainPaths == nil {
		mainPaths = &yaml.Node{Kind: yaml.MappingNode}
		mainRoot.Content = append(mainRoot.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "paths"},
			mainPaths,
		)
	}

	for i := 0; i < len(includePaths.Content)-1; i += 2 {
		pathKey := includePaths.Content[i]
		pathValue := includePaths.Content[i+1]

		if !m.hasKey(mainPaths, pathKey.Value) {
			mainPaths.Content = append(mainPaths.Content,
				m.cloneNode(pathKey),
				m.cloneNode(pathValue),
			)
		}
	}
}

// mergeTags merges tags from included spec into main spec
func (m *IncludeMerger) mergeTags(mainRoot, includeRoot *yaml.Node) {
	includeTags := m.findSequence(includeRoot, "tags")
	if includeTags == nil {
		return
	}

	mainTags := m.findSequence(mainRoot, "tags")
	if mainTags == nil {
		mainTags = &yaml.Node{Kind: yaml.SequenceNode}
		mainRoot.Content = append(mainRoot.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "tags"},
			mainTags,
		)
	}

	existingTags := make(map[string]bool)
	for _, tag := range mainTags.Content {
		if tag.Kind == yaml.MappingNode {
			for i := 0; i < len(tag.Content)-1; i += 2 {
				if tag.Content[i].Value == yamlKeyName {
					existingTags[tag.Content[i+1].Value] = true
					break
				}
			}
		}
	}

	for _, tag := range includeTags.Content {
		if tag.Kind == yaml.MappingNode {
			for i := 0; i < len(tag.Content)-1; i += 2 {
				if tag.Content[i].Value == yamlKeyName {
					tagName := tag.Content[i+1].Value
					if !existingTags[tagName] {
						mainTags.Content = append(mainTags.Content, m.cloneNode(tag))
						existingTags[tagName] = true
					}
					break
				}
			}
		}
	}
}

// mergeComponents merges components from included spec into main spec
func (m *IncludeMerger) mergeComponents(mainRoot, includeRoot *yaml.Node) {
	includeComponents := m.findMapping(includeRoot, "components")
	if includeComponents == nil {
		return
	}

	mainComponents := m.findMapping(mainRoot, "components")
	if mainComponents == nil {
		mainComponents = &yaml.Node{Kind: yaml.MappingNode}
		mainRoot.Content = append(mainRoot.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "components"},
			mainComponents,
		)
	}

	componentTypes := []string{
		"schemas", "parameters", "responses", "headers",
		"securitySchemes", "requestBodies", "examples",
		"links", "callbacks", "pathItems",
	}

	for _, compType := range componentTypes {
		includeSection := m.findMapping(includeComponents, compType)
		if includeSection == nil {
			continue
		}

		mainSection := m.findMapping(mainComponents, compType)
		if mainSection == nil {
			mainSection = &yaml.Node{Kind: yaml.MappingNode}
			mainComponents.Content = append(mainComponents.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: compType},
				mainSection,
			)
		}

		for i := 0; i < len(includeSection.Content)-1; i += 2 {
			itemKey := includeSection.Content[i]
			itemValue := includeSection.Content[i+1]

			if !m.hasKey(mainSection, itemKey.Value) {
				mainSection.Content = append(mainSection.Content,
					m.cloneNode(itemKey),
					m.cloneNode(itemValue),
				)
			}
		}
	}
}

// mergeSecurity merges security definitions from included spec into main spec
func (m *IncludeMerger) mergeSecurity(mainRoot, includeRoot *yaml.Node) {
	includeSecurity := m.findSequence(includeRoot, "security")
	if includeSecurity == nil {
		return
	}

	mainSecurity := m.findSequence(mainRoot, "security")
	if mainSecurity == nil {
		mainSecurity = &yaml.Node{Kind: yaml.SequenceNode}
		mainRoot.Content = append(mainRoot.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "security"},
			mainSecurity,
		)
	}

	// Build set of existing security requirement keys
	existingKeys := make(map[string]bool)
	for _, sec := range mainSecurity.Content {
		if sec.Kind == yaml.MappingNode && len(sec.Content) >= 2 {
			existingKeys[sec.Content[0].Value] = true
		}
	}

	for _, sec := range includeSecurity.Content {
		if sec.Kind == yaml.MappingNode && len(sec.Content) >= 2 {
			key := sec.Content[0].Value
			if !existingKeys[key] {
				mainSecurity.Content = append(mainSecurity.Content, m.cloneNode(sec))
				existingKeys[key] = true
			}
		}
	}
}

// removeIncludeExtensions removes x-include-* extensions from the root
func (m *IncludeMerger) removeIncludeExtensions(root *yaml.Node) {
	if root.Kind != yaml.MappingNode {
		return
	}

	var newContent []*yaml.Node
	for i := 0; i < len(root.Content)-1; i += 2 {
		keyNode := root.Content[i]
		valueNode := root.Content[i+1]

		if keyNode.Kind == yaml.ScalarNode && strings.HasPrefix(keyNode.Value, "x-include-") {
			continue
		}
		newContent = append(newContent, keyNode, valueNode)
	}
	root.Content = newContent
}

// Helper methods

func (m *IncludeMerger) findMapping(parent *yaml.Node, key string) *yaml.Node {
	if parent.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(parent.Content)-1; i += 2 {
		if parent.Content[i].Kind == yaml.ScalarNode && parent.Content[i].Value == key {
			if parent.Content[i+1].Kind == yaml.MappingNode {
				return parent.Content[i+1]
			}
		}
	}
	return nil
}

func (m *IncludeMerger) findSequence(parent *yaml.Node, key string) *yaml.Node {
	if parent.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(parent.Content)-1; i += 2 {
		if parent.Content[i].Kind == yaml.ScalarNode && parent.Content[i].Value == key {
			if parent.Content[i+1].Kind == yaml.SequenceNode {
				return parent.Content[i+1]
			}
		}
	}
	return nil
}

func (m *IncludeMerger) hasKey(mapping *yaml.Node, key string) bool {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Kind == yaml.ScalarNode && mapping.Content[i].Value == key {
			return true
		}
	}
	return false
}

func (m *IncludeMerger) cloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}

	clone := &yaml.Node{
		Kind:        node.Kind,
		Style:       node.Style,
		Tag:         node.Tag,
		Value:       node.Value,
		Anchor:      node.Anchor,
		Alias:       node.Alias,
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
		Line:        node.Line,
		Column:      node.Column,
	}

	if len(node.Content) > 0 {
		clone.Content = make([]*yaml.Node, len(node.Content))
		for i, child := range node.Content {
			clone.Content[i] = m.cloneNode(child)
		}
	}

	return clone
}

// CompositeFS is a filesystem that overlays multiple filesystems.
// It first checks the base filesystem (local disk), then falls back to embedded includes.
// This allows libopenapi to resolve references from both local files and embedded specs.
type CompositeFS struct {
	base     fs.FS           // Primary filesystem (local disk)
	includes []includeFSInfo // Embedded filesystems with their mount points
}

// includeFSInfo holds an embedded filesystem and its virtual mount point.
type includeFSInfo struct {
	name string    // Include name (e.g., "auth", "config")
	fsys IncludeFS // The embedded filesystem
}

// NewCompositeFS creates a new composite filesystem.
// The base filesystem is checked first, then includes are checked in order.
func NewCompositeFS(base fs.FS) *CompositeFS {
	return &CompositeFS{
		base:     base,
		includes: make([]includeFSInfo, 0),
	}
}

// AddInclude adds an embedded filesystem as an overlay.
// Files from the include's spec/ directory are accessible at the root.
func (c *CompositeFS) AddInclude(name string, fsys IncludeFS) {
	c.includes = append(c.includes, includeFSInfo{
		name: name,
		fsys: fsys,
	})
}

// Open implements fs.FS by first checking the base filesystem, then includes.
func (c *CompositeFS) Open(name string) (fs.File, error) {
	// Clean the path
	name = filepath.Clean(name)
	name = strings.TrimPrefix(name, "/")

	// Try base filesystem first
	if f, err := c.base.Open(name); err == nil {
		return f, nil
	}

	// Try each include's embedded filesystem
	// Include files are under spec/ in the embed, but we want them at root
	for _, inc := range c.includes {
		embedPath := filepath.Join("spec", name)
		if f, err := inc.fsys.Open(embedPath); err == nil {
			return f, nil
		}
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// ReadDir implements fs.ReadDirFS by merging directory entries from all filesystems.
func (c *CompositeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = filepath.Clean(name)
	name = strings.TrimPrefix(name, "/")

	entries := make(map[string]fs.DirEntry)

	// Read from base filesystem
	if baseDir, ok := c.base.(fs.ReadDirFS); ok {
		if dirEntries, err := baseDir.ReadDir(name); err == nil {
			for _, e := range dirEntries {
				entries[e.Name()] = e
			}
		}
	}

	// Read from each include (only add if not already present)
	for _, inc := range c.includes {
		embedPath := filepath.Join("spec", name)
		if dirEntries, err := inc.fsys.ReadDir(embedPath); err == nil {
			for _, e := range dirEntries {
				if _, exists := entries[e.Name()]; !exists {
					entries[e.Name()] = e
				}
			}
		}
	}

	if len(entries) == 0 {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}

	// Convert map to slice
	result := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		result = append(result, e)
	}
	return result, nil
}

// ReadFile implements fs.ReadFileFS.
func (c *CompositeFS) ReadFile(name string) (data []byte, err error) {
	f, err := c.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	return io.ReadAll(f)
}

// Stat implements fs.StatFS.
func (c *CompositeFS) Stat(name string) (info fs.FileInfo, err error) {
	f, err := c.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	return f.Stat()
}
