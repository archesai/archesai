package spec

import (
	"testing"
)

func TestResolver_ResolveFile(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		files       map[string]string
		filePath    string
		wantContent string
		wantErr     bool
	}{
		{
			name:    "resolves file relative to base directory",
			baseDir: ".",
			files: map[string]string{
				"components/schemas/User.yaml": "title: User\ntype: object",
			},
			filePath:    "components/schemas/User.yaml",
			wantContent: "title: User\ntype: object",
			wantErr:     false,
		},
		{
			name:    "resolves file with subdirectory base",
			baseDir: "api",
			files: map[string]string{
				"api/components/schemas/User.yaml": "title: User",
			},
			filePath:    "components/schemas/User.yaml",
			wantContent: "title: User",
			wantErr:     false,
		},
		{
			name:     "returns error for empty file path",
			baseDir:  ".",
			files:    map[string]string{},
			filePath: "",
			wantErr:  true,
		},
		{
			name:     "returns error for non-existent file",
			baseDir:  ".",
			files:    map[string]string{},
			filePath: "components/schemas/NotFound.yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := makeTestFS(tt.files)
			r := NewResolver(fsys, tt.baseDir)

			data, err := r.ResolveFile(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(data) != tt.wantContent {
				t.Errorf("ResolveFile() = %q, want %q", string(data), tt.wantContent)
			}
		})
	}
}

func TestResolver_ResolveFileFrom(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		files       map[string]string
		fromPath    string
		filePath    string
		wantContent string
		wantErr     bool
	}{
		{
			name:    "resolves file relative to source file",
			baseDir: ".",
			files: map[string]string{
				"paths/users.yaml":             "x-path: /users",
				"components/schemas/User.yaml": "title: User",
			},
			fromPath:    "paths/users.yaml",
			filePath:    "../components/schemas/User.yaml",
			wantContent: "title: User",
			wantErr:     false,
		},
		{
			name:    "resolves file relative to source directory",
			baseDir: ".",
			files: map[string]string{
				"components/responses/UserResponse.yaml": "description: User",
				"components/schemas/User.yaml":           "title: User",
			},
			fromPath:    "components/responses",
			filePath:    "../schemas/User.yaml",
			wantContent: "title: User",
			wantErr:     false,
		},
		{
			name:    "resolves sibling file",
			baseDir: ".",
			files: map[string]string{
				"components/schemas/User.yaml":    "title: User",
				"components/schemas/Session.yaml": "title: Session",
			},
			fromPath:    "components/schemas/User.yaml",
			filePath:    "Session.yaml",
			wantContent: "title: Session",
			wantErr:     false,
		},
		{
			name:     "returns error for empty file path",
			baseDir:  ".",
			files:    map[string]string{},
			fromPath: "paths/users.yaml",
			filePath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := makeTestFS(tt.files)
			r := NewResolver(fsys, tt.baseDir)

			data, err := r.ResolveFileFrom(tt.fromPath, tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveFileFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(data) != tt.wantContent {
				t.Errorf("ResolveFileFrom() = %q, want %q", string(data), tt.wantContent)
			}
		})
	}
}

func TestExtractSchemaNameFromRef(t *testing.T) {
	tests := []struct {
		ref  string
		want string
	}{
		{"./User.yaml", "User"},
		{"../schemas/Organization.yaml", "Organization"},
		{"../../components/responses/BadRequest.yaml", "BadRequest"},
		{"../parameters/ResourceID.yml", "ResourceID"},
		{"User.yaml", "User"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			if got := ExtractSchemaNameFromRef(tt.ref); got != tt.want {
				t.Errorf("ExtractSchemaNameFromRef(%q) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}
