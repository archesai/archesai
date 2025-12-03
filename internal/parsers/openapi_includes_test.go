package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestIncludeMerger_FindEnabledIncludes(t *testing.T) {
	tests := []struct {
		name     string
		specYAML string
		want     []string
	}{
		{
			name: "no includes",
			specYAML: `
openapi: 3.1.0
info:
  title: Test API
  version: 1.0.0
paths: {}
`,
			want: []string{},
		},
		{
			name: "single include enabled",
			specYAML: `
openapi: 3.1.0
x-include-auth: true
info:
  title: Test API
  version: 1.0.0
paths: {}
`,
			want: []string{"auth"},
		},
		{
			name: "multiple includes enabled",
			specYAML: `
openapi: 3.1.0
x-include-auth: true
x-include-config: true
x-include-server: true
info:
  title: Test API
  version: 1.0.0
paths: {}
`,
			want: []string{"auth", "config", "server"},
		},
		{
			name: "include disabled with false",
			specYAML: `
openapi: 3.1.0
x-include-auth: false
x-include-config: true
info:
  title: Test API
  version: 1.0.0
paths: {}
`,
			want: []string{"config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merger := NewDefaultIncludeMerger()

			var doc yaml.Node
			err := yaml.Unmarshal([]byte(tt.specYAML), &doc)
			require.NoError(t, err)

			enabled := merger.findEnabledIncludes(&doc)
			enabledNames := make([]string, len(enabled))
			for i, inc := range enabled {
				enabledNames[i] = inc.Name
			}

			assert.ElementsMatch(t, tt.want, enabledNames)
		})
	}
}

func TestIncludeMerger_RewriteRefs(t *testing.T) {
	merger := NewIncludeMerger()

	tests := []struct {
		name    string
		input   string
		relPath string
		want    string
	}{
		{
			name: "rewrite relative ref",
			input: `
paths:
  /test:
    $ref: paths/test.yaml
`,
			relPath: "_includes/auth",
			want:    "_includes/auth/paths/test.yaml",
		},
		{
			name: "preserve internal ref",
			input: `
schema:
  $ref: '#/components/schemas/User'
`,
			relPath: "_includes/auth",
			want:    "#/components/schemas/User",
		},
		{
			name: "preserve http ref",
			input: `
schema:
  $ref: https://example.com/schema.yaml
`,
			relPath: "_includes/auth",
			want:    "https://example.com/schema.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node yaml.Node
			err := yaml.Unmarshal([]byte(tt.input), &node)
			require.NoError(t, err)

			merger.rewriteRefs(node.Content[0], tt.relPath)

			// Find the $ref value
			ref := findRefValue(node.Content[0])
			assert.Equal(t, tt.want, ref)
		})
	}
}

func findRefValue(node *yaml.Node) string {
	if node == nil {
		return ""
	}

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content)-1; i += 2 {
			if node.Content[i].Value == "$ref" {
				return node.Content[i+1].Value
			}
			if ref := findRefValue(node.Content[i+1]); ref != "" {
				return ref
			}
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			if ref := findRefValue(child); ref != "" {
				return ref
			}
		}
	}
	return ""
}

func TestIncludeMerger_MergePaths(t *testing.T) {
	merger := NewIncludeMerger()

	mainYAML := `
openapi: 3.1.0
paths:
  /existing:
    get:
      summary: Existing endpoint
`
	includeYAML := `
openapi: 3.1.0
paths:
  /new:
    get:
      summary: New endpoint
  /existing:
    get:
      summary: Should not override
`

	var mainDoc, includeDoc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(mainYAML), &mainDoc))
	require.NoError(t, yaml.Unmarshal([]byte(includeYAML), &includeDoc))

	mainRoot := mainDoc.Content[0]
	includeRoot := includeDoc.Content[0]

	merger.mergePaths(mainRoot, includeRoot)

	// Check that /new was added
	paths := merger.findMapping(mainRoot, "paths")
	require.NotNil(t, paths)

	assert.True(t, merger.hasKey(paths, "/existing"))
	assert.True(t, merger.hasKey(paths, "/new"))
}

func TestIncludeMerger_MergeComponents(t *testing.T) {
	merger := NewIncludeMerger()

	mainYAML := `
openapi: 3.1.0
components:
  schemas:
    ExistingSchema:
      type: object
`
	includeYAML := `
openapi: 3.1.0
components:
  schemas:
    NewSchema:
      type: object
  parameters:
    NewParam:
      name: test
      in: query
`

	var mainDoc, includeDoc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(mainYAML), &mainDoc))
	require.NoError(t, yaml.Unmarshal([]byte(includeYAML), &includeDoc))

	mainRoot := mainDoc.Content[0]
	includeRoot := includeDoc.Content[0]

	merger.mergeComponents(mainRoot, includeRoot)

	// Check schemas
	components := merger.findMapping(mainRoot, "components")
	require.NotNil(t, components)

	schemas := merger.findMapping(components, "schemas")
	require.NotNil(t, schemas)
	assert.True(t, merger.hasKey(schemas, "ExistingSchema"))
	assert.True(t, merger.hasKey(schemas, "NewSchema"))

	// Check parameters was added
	params := merger.findMapping(components, "parameters")
	require.NotNil(t, params)
	assert.True(t, merger.hasKey(params, "NewParam"))
}

func TestIncludeMerger_MergeTags(t *testing.T) {
	merger := NewIncludeMerger()

	mainYAML := `
openapi: 3.1.0
tags:
  - name: Existing
    description: Existing tag
`
	includeYAML := `
openapi: 3.1.0
tags:
  - name: New
    description: New tag
  - name: Existing
    description: Should not duplicate
`

	var mainDoc, includeDoc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(mainYAML), &mainDoc))
	require.NoError(t, yaml.Unmarshal([]byte(includeYAML), &includeDoc))

	mainRoot := mainDoc.Content[0]
	includeRoot := includeDoc.Content[0]

	merger.mergeTags(mainRoot, includeRoot)

	// Check tags
	tags := merger.findSequence(mainRoot, "tags")
	require.NotNil(t, tags)

	// Should have 2 tags (Existing and New, no duplicate)
	assert.Len(t, tags.Content, 2)

	// Verify tag names
	tagNames := make([]string, 0)
	for _, tag := range tags.Content {
		if tag.Kind == yaml.MappingNode {
			for i := 0; i < len(tag.Content)-1; i += 2 {
				if tag.Content[i].Value == "name" {
					tagNames = append(tagNames, tag.Content[i+1].Value)
				}
			}
		}
	}
	assert.ElementsMatch(t, []string{"Existing", "New"}, tagNames)
}

func TestIncludeMerger_RemoveIncludeExtensions(t *testing.T) {
	merger := NewIncludeMerger()

	specYAML := `
openapi: 3.1.0
x-include-auth: true
x-include-config: true
x-project-name: test
info:
  title: Test
`

	var doc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(specYAML), &doc))

	root := doc.Content[0]
	merger.removeIncludeExtensions(root)

	// x-include-* should be removed
	assert.False(t, merger.hasKey(root, "x-include-auth"))
	assert.False(t, merger.hasKey(root, "x-include-config"))

	// Other extensions should remain
	assert.True(t, merger.hasKey(root, "x-project-name"))
	assert.True(t, merger.hasKey(root, "openapi"))
	assert.True(t, merger.hasKey(root, "info"))
}

func TestIncludeMerger_ProcessIncludes_NoIncludes(t *testing.T) {
	// Create a temp spec file without includes
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "openapi.yaml")

	specContent := `
openapi: 3.1.0
info:
  title: Test API
  version: 1.0.0
paths: {}
`
	require.NoError(t, os.WriteFile(specPath, []byte(specContent), 0644))

	merger := NewIncludeMerger()
	mergedPath, cleanup, enabledNames, err := merger.ProcessIncludes(specPath)
	require.NoError(t, err)
	defer cleanup()

	// Should return original path when no includes
	assert.Equal(t, specPath, mergedPath)
	assert.Empty(t, enabledNames)
}

func TestIncludeMerger_ProcessIncludes_WithIncludes(t *testing.T) {
	// This test uses the real embedded specs
	merger := NewDefaultIncludeMerger()

	// Create a temp spec that includes auth
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "openapi.yaml")

	specContent := `
openapi: 3.1.0
x-include-server: true
info:
  title: Test API
  version: 1.0.0
paths: {}
`
	require.NoError(t, os.WriteFile(specPath, []byte(specContent), 0644))

	mergedPath, cleanup, enabledNames, err := merger.ProcessIncludes(specPath)
	require.NoError(t, err)
	defer cleanup()

	// Should return a different (temp) path
	assert.NotEqual(t, specPath, mergedPath)

	// Should return enabled include names
	assert.Contains(t, enabledNames, "server")

	// Read the merged spec
	mergedContent, err := os.ReadFile(mergedPath)
	require.NoError(t, err)

	// Parse and verify it contains merged content
	var doc yaml.Node
	require.NoError(t, yaml.Unmarshal(mergedContent, &doc))

	root := doc.Content[0]

	// Should not have x-include-server anymore
	assert.False(t, merger.hasKey(root, "x-include-server"))

	// Should have paths from server (like /health)
	paths := merger.findMapping(root, "paths")
	require.NotNil(t, paths, "paths should exist after merge")
	assert.True(t, merger.hasKey(paths, "/health"), "should have /health path from server include")

	// Should have tags from server
	tags := merger.findSequence(root, "tags")
	require.NotNil(t, tags, "tags should exist after merge")
}

func TestIncludeMerger_CloneNode(t *testing.T) {
	merger := NewIncludeMerger()

	original := &yaml.Node{
		Kind:  yaml.MappingNode,
		Value: "test",
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "key"},
			{Kind: yaml.ScalarNode, Value: "value"},
		},
	}

	clone := merger.cloneNode(original)

	// Should be a deep copy
	assert.Equal(t, original.Kind, clone.Kind)
	assert.Equal(t, original.Value, clone.Value)
	assert.Len(t, clone.Content, 2)

	// Modifying clone should not affect original
	clone.Content[1].Value = "modified"
	assert.Equal(t, "value", original.Content[1].Value)
}
