package spec

import (
	"testing"

	"github.com/archesai/archesai/pkg/storage"
)

func TestDiscoverPaths(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected []string
	}{
		{
			name: "finds yaml files in paths directory",
			files: map[string]string{
				"paths/users.yaml":       "x-path: /users",
				"paths/auth_login.yaml":  "x-path: /auth/login",
				"paths/health.yaml":      "x-path: /health",
				"paths/.hidden.yaml":     "x-path: /hidden",
				"paths/README.md":        "# Paths",
				"other/file.yaml":        "some content",
				"components/schema.yaml": "type: object",
			},
			expected: []string{
				"paths/auth_login.yaml",
				"paths/health.yaml",
				"paths/users.yaml",
				"paths/.hidden.yaml",
			},
		},
		{
			name: "finds yml extension files",
			files: map[string]string{
				"paths/users.yml":  "x-path: /users",
				"paths/health.yml": "x-path: /health",
			},
			expected: []string{
				"paths/health.yml",
				"paths/users.yml",
			},
		},
		{
			name:     "returns empty when paths directory does not exist",
			files:    map[string]string{},
			expected: nil,
		},
		{
			name: "returns empty when paths directory is empty",
			files: map[string]string{
				"paths/.gitkeep": "",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := storage.NewMemoryStorageWithFiles(tt.files)

			paths, err := DiscoverPaths(fsys)
			if err != nil {
				t.Fatalf("DiscoverPaths() error = %v", err)
			}

			if len(paths) != len(tt.expected) {
				t.Errorf("DiscoverPaths() got %d paths, want %d", len(paths), len(tt.expected))
				t.Errorf("got: %v", paths)
				t.Errorf("want: %v", tt.expected)
				return
			}

			// Convert to map for easier comparison (order doesn't matter)
			gotMap := make(map[string]bool)
			for _, p := range paths {
				gotMap[p] = true
			}

			for _, expected := range tt.expected {
				if !gotMap[expected] {
					t.Errorf("expected path %q not found in results", expected)
				}
			}
		})
	}
}

func TestDiscoverComponents(t *testing.T) {
	files := map[string]string{
		"components/schemas/User.yaml":               "title: User",
		"components/schemas/Session.yaml":            "title: Session",
		"components/responses/UserResponse.yaml":     "description: User",
		"components/parameters/PageQuery.yaml":       "name: page",
		"components/headers/RateLimit.yaml":          "description: Rate limit",
		"components/securitySchemes/BearerAuth.yaml": "type: http",
	}
	fsys := storage.NewMemoryStorageWithFiles(files)

	tests := []struct {
		kind     ComponentKind
		expected map[string]string
	}{
		{
			kind: ComponentSchemas,
			expected: map[string]string{
				"User":    "components/schemas/User.yaml",
				"Session": "components/schemas/Session.yaml",
			},
		},
		{
			kind: ComponentResponses,
			expected: map[string]string{
				"UserResponse": "components/responses/UserResponse.yaml",
			},
		},
		{
			kind: ComponentParameters,
			expected: map[string]string{
				"PageQuery": "components/parameters/PageQuery.yaml",
			},
		},
		{
			kind: ComponentHeaders,
			expected: map[string]string{
				"RateLimit": "components/headers/RateLimit.yaml",
			},
		},
		{
			kind: ComponentSecuritySchemes,
			expected: map[string]string{
				"BearerAuth": "components/securitySchemes/BearerAuth.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.kind), func(t *testing.T) {
			result, err := DiscoverComponents(fsys, tt.kind)
			if err != nil {
				t.Fatalf("DiscoverComponents(%s) error = %v", tt.kind, err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("got %d items, want %d", len(result), len(tt.expected))
			}

			for name, path := range tt.expected {
				if gotPath, ok := result[name]; !ok {
					t.Errorf("expected %q not found", name)
				} else if gotPath != path {
					t.Errorf("%q path = %q, want %q", name, gotPath, path)
				}
			}
		})
	}
}
