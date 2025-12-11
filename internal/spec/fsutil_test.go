package spec

import (
	"testing"
	"testing/fstest"
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
			fsys := makeTestFS(tt.files)

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

func TestDiscoverSchemas(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected map[string]string
	}{
		{
			name: "finds schemas in components/schemas directory",
			files: map[string]string{
				"components/schemas/User.yaml":         "title: User\ntype: object",
				"components/schemas/Organization.yaml": "title: Organization\ntype: object",
				"components/schemas/Session.yaml":      "title: Session\ntype: object",
			},
			expected: map[string]string{
				"User":         "components/schemas/User.yaml",
				"Organization": "components/schemas/Organization.yaml",
				"Session":      "components/schemas/Session.yaml",
			},
		},
		{
			name: "handles yml extension",
			files: map[string]string{
				"components/schemas/User.yml":    "title: User",
				"components/schemas/Session.yml": "title: Session",
			},
			expected: map[string]string{
				"User":    "components/schemas/User.yml",
				"Session": "components/schemas/Session.yml",
			},
		},
		{
			name:     "returns empty map when directory does not exist",
			files:    map[string]string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := makeTestFS(tt.files)

			schemas, err := DiscoverSchemas(fsys)
			if err != nil {
				t.Fatalf("DiscoverSchemas() error = %v", err)
			}

			if len(schemas) != len(tt.expected) {
				t.Errorf(
					"DiscoverSchemas() got %d schemas, want %d",
					len(schemas),
					len(tt.expected),
				)
				return
			}

			for name, path := range tt.expected {
				if gotPath, ok := schemas[name]; !ok {
					t.Errorf("expected schema %q not found", name)
				} else if gotPath != path {
					t.Errorf("schema %q path = %q, want %q", name, gotPath, path)
				}
			}
		})
	}
}

func TestDiscoverResponses(t *testing.T) {
	files := map[string]string{
		"components/responses/UserResponse.yaml":     "description: User response",
		"components/responses/UserListResponse.yaml": "description: User list response",
		"components/responses/BadRequest.yaml":       "description: Bad request",
	}

	fsys := makeTestFS(files)

	responses, err := DiscoverResponses(fsys)
	if err != nil {
		t.Fatalf("DiscoverResponses() error = %v", err)
	}

	expected := map[string]string{
		"UserResponse":     "components/responses/UserResponse.yaml",
		"UserListResponse": "components/responses/UserListResponse.yaml",
		"BadRequest":       "components/responses/BadRequest.yaml",
	}

	if len(responses) != len(expected) {
		t.Errorf("got %d responses, want %d", len(responses), len(expected))
	}

	for name, path := range expected {
		if gotPath, ok := responses[name]; !ok {
			t.Errorf("expected response %q not found", name)
		} else if gotPath != path {
			t.Errorf("response %q path = %q, want %q", name, gotPath, path)
		}
	}
}

func TestDiscoverParameters(t *testing.T) {
	files := map[string]string{
		"components/parameters/ResourceID.yaml": "name: id\nin: path",
		"components/parameters/PageQuery.yaml":  "name: page\nin: query",
	}

	fsys := makeTestFS(files)

	params, err := DiscoverParameters(fsys)
	if err != nil {
		t.Fatalf("DiscoverParameters() error = %v", err)
	}

	if len(params) != 2 {
		t.Errorf("got %d parameters, want 2", len(params))
	}

	if _, ok := params["ResourceID"]; !ok {
		t.Error("expected ResourceID parameter not found")
	}
	if _, ok := params["PageQuery"]; !ok {
		t.Error("expected PageQuery parameter not found")
	}
}

func TestDiscoverHeaders(t *testing.T) {
	files := map[string]string{
		"components/headers/RateLimitLimit.yaml":     "description: Rate limit",
		"components/headers/RateLimitRemaining.yaml": "description: Remaining",
	}

	fsys := makeTestFS(files)

	headers, err := DiscoverHeaders(fsys)
	if err != nil {
		t.Fatalf("DiscoverHeaders() error = %v", err)
	}

	if len(headers) != 2 {
		t.Errorf("got %d headers, want 2", len(headers))
	}

	if _, ok := headers["RateLimitLimit"]; !ok {
		t.Error("expected RateLimitLimit header not found")
	}
}

// makeTestFS creates a testing filesystem from a map of path -> content
func makeTestFS(files map[string]string) fstest.MapFS {
	fs := make(fstest.MapFS)
	for path, content := range files {
		fs[path] = &fstest.MapFile{Data: []byte(content)}
	}
	return fs
}
