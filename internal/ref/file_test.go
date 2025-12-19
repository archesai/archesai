package ref

import (
	"testing"

	"github.com/archesai/archesai/pkg/storage"
)

func TestNewFileResolver(t *testing.T) {
	fsys := storage.NewMemoryStorageWithFiles(map[string]string{
		"file.txt": "content",
	})

	t.Run("with empty baseDir defaults to dot", func(t *testing.T) {
		r := NewFileResolver(fsys, "")
		if r.BaseDir() != "." {
			t.Errorf("BaseDir() = %q, want %q", r.BaseDir(), ".")
		}
	})

	t.Run("with specified baseDir", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		if r.BaseDir() != "api" {
			t.Errorf("BaseDir() = %q, want %q", r.BaseDir(), "api")
		}
	})
}

func TestFileResolver_FS(t *testing.T) {
	fsys := storage.NewMemoryStorageWithFiles(map[string]string{
		"file.txt": "content",
	})
	r := NewFileResolver(fsys, ".")

	if r.FS() != fsys {
		t.Error("FS() should return the underlying filesystem")
	}
}

func TestFileResolver_BaseDir(t *testing.T) {
	fsys := storage.NewMemoryStorageWithFiles(map[string]string{})

	tests := []struct {
		name    string
		baseDir string
		want    string
	}{
		{"empty", "", "."},
		{"dot", ".", "."},
		{"subdirectory", "api", "api"},
		{"nested", "api/v1", "api/v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFileResolver(fsys, tt.baseDir)
			if got := r.BaseDir(); got != tt.want {
				t.Errorf("BaseDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFileResolver_ReadFile(t *testing.T) {
	fsys := storage.NewMemoryStorageWithFiles(map[string]string{
		"openapi.yaml":         "openapi: 3.1.0",
		"api/openapi.yaml":     "openapi: 3.0.0",
		"api/paths/users.yaml": "get: listUsers",
		"components/User.yaml": "type: object",
	})

	t.Run("reads file relative to baseDir", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		data, err := r.ReadFile("openapi.yaml")
		if err != nil {
			t.Fatalf("ReadFile error: %v", err)
		}
		if string(data) != "openapi: 3.0.0" {
			t.Errorf("ReadFile = %q, want %q", string(data), "openapi: 3.0.0")
		}
	})

	t.Run("reads nested file", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		data, err := r.ReadFile("paths/users.yaml")
		if err != nil {
			t.Fatalf("ReadFile error: %v", err)
		}
		if string(data) != "get: listUsers" {
			t.Errorf("ReadFile = %q, want %q", string(data), "get: listUsers")
		}
	})

	t.Run("reads from root with dot baseDir", func(t *testing.T) {
		r := NewFileResolver(fsys, ".")
		data, err := r.ReadFile("openapi.yaml")
		if err != nil {
			t.Fatalf("ReadFile error: %v", err)
		}
		if string(data) != "openapi: 3.1.0" {
			t.Errorf("ReadFile = %q, want %q", string(data), "openapi: 3.1.0")
		}
	})

	t.Run("returns error for empty path", func(t *testing.T) {
		r := NewFileResolver(fsys, ".")
		_, err := r.ReadFile("")
		if err == nil {
			t.Error("expected error for empty path")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		r := NewFileResolver(fsys, ".")
		_, err := r.ReadFile("nonexistent.yaml")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}

func TestFileResolver_ReadFileFrom(t *testing.T) {
	fsys := storage.NewMemoryStorageWithFiles(map[string]string{
		"api/openapi.yaml":            "openapi: 3.1.0",
		"api/paths/users.yaml":        "$ref: ../components/User.yaml",
		"api/components/User.yaml":    "type: object",
		"api/components/Session.yaml": "type: object",
	})

	t.Run("reads file relative to another file", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		// From paths/users.yaml, read ../components/User.yaml
		data, err := r.ReadFileFrom("paths/users.yaml", "../components/User.yaml")
		if err != nil {
			t.Fatalf("ReadFileFrom error: %v", err)
		}
		if string(data) != "type: object" {
			t.Errorf("ReadFileFrom = %q, want %q", string(data), "type: object")
		}
	})

	t.Run("reads file relative to directory", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		// From components directory, read User.yaml
		data, err := r.ReadFileFrom("components", "User.yaml")
		if err != nil {
			t.Fatalf("ReadFileFrom error: %v", err)
		}
		if string(data) != "type: object" {
			t.Errorf("ReadFileFrom = %q, want %q", string(data), "type: object")
		}
	})

	t.Run("reads sibling file", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		// From components/User.yaml, read Session.yaml (sibling)
		data, err := r.ReadFileFrom("components/User.yaml", "Session.yaml")
		if err != nil {
			t.Fatalf("ReadFileFrom error: %v", err)
		}
		if string(data) != "type: object" {
			t.Errorf("ReadFileFrom = %q, want %q", string(data), "type: object")
		}
	})

	t.Run("returns error for empty target path", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		_, err := r.ReadFileFrom("openapi.yaml", "")
		if err == nil {
			t.Error("expected error for empty target path")
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		r := NewFileResolver(fsys, "api")
		_, err := r.ReadFileFrom("openapi.yaml", "nonexistent.yaml")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}
