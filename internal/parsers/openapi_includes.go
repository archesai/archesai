package parsers

import (
	"embed"
	"fmt"
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

const (
	yamlKeyRef  = "$ref"
	yamlKeyName = "name"
)

// IncludeSpec defines a spec that can be included via x-include-* extensions
type IncludeSpec struct {
	Name string   // Extension name without "x-include-" prefix (e.g., "auth", "config")
	FS   embed.FS // Embedded filesystem containing the spec/ directory
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

// NewDefaultIncludeMerger creates a new IncludeMerger with all standard includes registered.
func NewDefaultIncludeMerger() *IncludeMerger {
	merger := NewIncludeMerger()

	// Register all standard includes
	merger.RegisterInclude("auth", auth.APISpec)
	merger.RegisterInclude("config", config.APISpec)
	merger.RegisterInclude("server", server.APISpec)
	merger.RegisterInclude("storage", storage.APISpec)
	merger.RegisterInclude("pipelines", pipelines.APISpec)
	merger.RegisterInclude("executor", executor.APISpec)

	return merger
}

// RegisterInclude registers an includable spec with its embedded filesystem
func (m *IncludeMerger) RegisterInclude(name string, fs embed.FS) *IncludeMerger {
	m.includeSpecs[name] = IncludeSpec{
		Name: name,
		FS:   fs,
	}
	return m
}

// ProcessIncludes reads a spec file, finds x-include-* extensions, and prepares
// a working directory with all necessary files for bundling.
// Returns the path to the merged spec file, a cleanup function, and the names of enabled includes.
func (m *IncludeMerger) ProcessIncludes(
	specPath string,
) (mergedSpecPath string, cleanup func(), enabledNames []string, err error) {
	// Read main spec
	mainSpecBytes, err := os.ReadFile(specPath)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to read main spec: %w", err)
	}

	// Parse as YAML node to find includes
	var mainSpec yaml.Node
	if err := yaml.Unmarshal(mainSpecBytes, &mainSpec); err != nil {
		return "", nil, nil, fmt.Errorf("failed to parse main spec: %w", err)
	}

	// Find enabled includes
	enabledIncludes := m.findEnabledIncludes(&mainSpec)

	// Build list of enabled include names
	for _, include := range enabledIncludes {
		enabledNames = append(enabledNames, include.Name)
	}

	// If no includes, return original path with no-op cleanup
	if len(enabledIncludes) == 0 {
		return specPath, func() {}, enabledNames, nil
	}

	// Create temp directory for merged spec
	tempDir, err := os.MkdirTemp("", "archesai-bundle-*")
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	cleanup = func() {
		_ = os.RemoveAll(tempDir)
	}

	// Copy main spec's directory structure to temp
	mainSpecDir := filepath.Dir(specPath)
	if err := m.copyDir(mainSpecDir, tempDir); err != nil {
		cleanup()
		return "", nil, nil, fmt.Errorf("failed to copy main spec directory: %w", err)
	}

	// Extract each included spec's components directly to the main components directories
	// This allows internal refs like #/components/responses/BadRequest to resolve correctly
	for _, include := range enabledIncludes {
		if err := m.extractComponentsToMain(include.FS, tempDir); err != nil {
			cleanup()
			return "", nil, nil, fmt.Errorf(
				"failed to extract components from %s: %w",
				include.Name,
				err,
			)
		}
	}

	// Rewrite main spec to include the merged components
	mainSpecInTemp := filepath.Join(tempDir, filepath.Base(specPath))
	if err := m.mergeAndRewriteSpec(mainSpecInTemp, enabledIncludes, tempDir); err != nil {
		cleanup()
		return "", nil, nil, fmt.Errorf("failed to merge specs: %w", err)
	}

	return mainSpecInTemp, cleanup, enabledNames, nil
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

// copyDir copies a directory recursively
func (m *IncludeMerger) copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

// extractComponentsToMain extracts component files from an included spec to the main spec's
// components directories. This allows internal refs like #/components/responses/BadRequest
// to resolve correctly.
func (m *IncludeMerger) extractComponentsToMain(embedFS embed.FS, mainDir string) error {
	// Component directories to copy
	componentDirs := []string{
		"components/schemas",
		"components/responses",
		"components/parameters",
		"components/headers",
		"components/requestBodies",
		"components/securitySchemes",
		"components/examples",
		"components/links",
		"components/callbacks",
		"components/pathItems",
		"paths",
	}

	for _, compDir := range componentDirs {
		srcPath := filepath.Join("spec", compDir)

		// Check if this directory exists in the embedded FS
		entries, err := embedFS.ReadDir(srcPath)
		if err != nil {
			// Directory doesn't exist, skip
			continue
		}

		// Create target directory if needed
		dstPath := filepath.Join(mainDir, compDir)
		if err := os.MkdirAll(dstPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dstPath, err)
		}

		// Copy each file
		for _, entry := range entries {
			if entry.IsDir() {
				// Recursively copy subdirectories
				subSrc := filepath.Join(srcPath, entry.Name())
				subDst := filepath.Join(dstPath, entry.Name())
				if err := m.extractEmbedFSDir(embedFS, subSrc, subDst); err != nil {
					return err
				}
				continue
			}

			srcFile := filepath.Join(srcPath, entry.Name())
			dstFile := filepath.Join(dstPath, entry.Name())

			// Only copy if file doesn't already exist (don't overwrite main spec's files)
			if _, err := os.Stat(dstFile); err == nil {
				continue
			}

			data, err := embedFS.ReadFile(srcFile)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", srcFile, err)
			}

			if err := os.WriteFile(dstFile, data, 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", dstFile, err)
			}
		}
	}

	return nil
}

// extractEmbedFSDir recursively extracts a directory from embedded FS
func (m *IncludeMerger) extractEmbedFSDir(embedFS embed.FS, srcDir, dstDir string) error {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	entries, err := embedFS.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			if err := m.extractEmbedFSDir(embedFS, srcPath, dstPath); err != nil {
				return err
			}
			continue
		}

		// Only copy if file doesn't already exist
		if _, err := os.Stat(dstPath); err == nil {
			continue
		}

		data, err := embedFS.ReadFile(srcPath)
		if err != nil {
			return err
		}

		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// mergeAndRewriteSpec merges included specs into the main spec.
// Since component files are already copied to the main directories, we just need to
// update the main spec's components section to include refs to those files.
func (m *IncludeMerger) mergeAndRewriteSpec(
	specPath string,
	includes []IncludeSpec,
	tempDir string,
) error {
	// Read main spec
	mainSpecBytes, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to read spec: %w", err)
	}

	var mainSpec yaml.Node
	if err := yaml.Unmarshal(mainSpecBytes, &mainSpec); err != nil {
		return fmt.Errorf("failed to parse spec: %w", err)
	}

	if mainSpec.Kind != yaml.DocumentNode || len(mainSpec.Content) == 0 {
		return fmt.Errorf("invalid spec structure")
	}

	mainRoot := mainSpec.Content[0]

	// For each include, read its spec and merge the component declarations
	for _, include := range includes {
		// Read the included spec from embedded FS
		includeBytes, err := include.FS.ReadFile("spec/openapi.yaml")
		if err != nil {
			return fmt.Errorf("failed to read include spec %s: %w", include.Name, err)
		}

		var includeSpec yaml.Node
		if err := yaml.Unmarshal(includeBytes, &includeSpec); err != nil {
			return fmt.Errorf("failed to parse include spec %s: %w", include.Name, err)
		}

		if includeSpec.Kind != yaml.DocumentNode || len(includeSpec.Content) == 0 {
			continue
		}

		includeRoot := includeSpec.Content[0]

		// DON'T rewrite refs - component files are already in the main directory
		// Just merge the component declarations, paths, and tags
		m.mergePaths(mainRoot, includeRoot)
		m.mergeTags(mainRoot, includeRoot)
		m.mergeComponents(mainRoot, includeRoot)
		m.mergeSecurity(mainRoot, includeRoot)
	}

	// Scan the temp directory and add any component files that aren't already declared
	m.addUndeclaredComponents(mainRoot, tempDir)

	// Remove x-include-* extensions from output
	m.removeIncludeExtensions(mainRoot)

	// Write merged spec back
	mergedBytes, err := yaml.Marshal(&mainSpec)
	if err != nil {
		return fmt.Errorf("failed to marshal merged spec: %w", err)
	}

	return os.WriteFile(specPath, mergedBytes, 0644)
}

// addUndeclaredComponents scans the component directories and adds refs for any
// files that aren't already declared in the spec's components section.
func (m *IncludeMerger) addUndeclaredComponents(mainRoot *yaml.Node, tempDir string) {
	componentTypes := map[string]string{
		"components/schemas":         "schemas",
		"components/responses":       "responses",
		"components/parameters":      "parameters",
		"components/headers":         "headers",
		"components/requestBodies":   "requestBodies",
		"components/securitySchemes": "securitySchemes",
		"components/examples":        "examples",
		"components/links":           "links",
		"components/callbacks":       "callbacks",
		"components/pathItems":       "pathItems",
	}

	// Get or create components section
	components := m.findMapping(mainRoot, "components")
	if components == nil {
		components = &yaml.Node{Kind: yaml.MappingNode}
		mainRoot.Content = append(mainRoot.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "components"},
			components,
		)
	}

	for dirPath, compType := range componentTypes {
		fullPath := filepath.Join(tempDir, dirPath)
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			// Directory doesn't exist, skip
			continue
		}

		// Get or create this component type section
		section := m.findMapping(components, compType)
		if section == nil {
			section = &yaml.Node{Kind: yaml.MappingNode}
			components.Content = append(components.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: compType},
				section,
			)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
				continue
			}

			// Component name is filename without extension
			compName := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")

			// Add if not already declared
			if !m.hasKey(section, compName) {
				refPath := filepath.Join(dirPath, name)
				section.Content = append(section.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: compName},
					&yaml.Node{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "$ref"},
							{Kind: yaml.ScalarNode, Value: refPath},
						},
					},
				)
			}
		}
	}
}

// rewriteRefs rewrites $ref values to prepend the relative path
func (m *IncludeMerger) rewriteRefs(node *yaml.Node, relPath string) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if keyNode.Kind == yaml.ScalarNode && keyNode.Value == yamlKeyRef {
				if valueNode.Kind == yaml.ScalarNode {
					ref := valueNode.Value
					// Only rewrite relative file refs, not internal refs (#/...)
					if !strings.HasPrefix(ref, "#") && !strings.HasPrefix(ref, "http") {
						valueNode.Value = filepath.Join(relPath, ref)
					}
				}
			} else {
				m.rewriteRefs(valueNode, relPath)
			}
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			m.rewriteRefs(child, relPath)
		}
	}
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
