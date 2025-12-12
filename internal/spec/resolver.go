package spec

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// Resolver handles resolution of file paths in OpenAPI documents.
type Resolver struct {
	fsys    fs.FS
	baseDir string // base directory within the filesystem (can be "." or subdirectory)
}

// NewResolver creates a new Resolver with the given filesystem and base directory.
func NewResolver(fsys fs.FS, baseDir string) *Resolver {
	if baseDir == "" {
		baseDir = "."
	}
	return &Resolver{fsys: fsys, baseDir: baseDir}
}

// ResolveFile reads the file at the given path relative to the base directory.
func (r *Resolver) ResolveFile(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, fmt.Errorf("empty file path")
	}

	// Resolve the path relative to base directory
	fullPath := path.Join(r.baseDir, filePath)
	fullPath = path.Clean(fullPath)

	data, err := fs.ReadFile(r.fsys, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return data, nil
}

// ResolveFileFrom reads a file relative to another file or directory.
func (r *Resolver) ResolveFileFrom(fromPath, filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, fmt.Errorf("empty file path")
	}

	// fromPath can be a file path or a directory path
	// If it's a file path, get its directory; if it's a directory, use it directly
	fromFullPath := path.Join(r.baseDir, fromPath)
	info, err := fs.Stat(r.fsys, fromFullPath)
	var fromDir string
	if err == nil && info.IsDir() {
		fromDir = fromFullPath
	} else {
		fromDir = path.Dir(fromFullPath)
	}

	// Resolve the path relative to the source directory
	fullPath := path.Join(fromDir, filePath)
	fullPath = path.Clean(fullPath)

	data, err := fs.ReadFile(r.fsys, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s (from %s): %w", filePath, fromPath, err)
	}

	return data, nil
}

// ExtractSchemaNameFromRef extracts the component name from a file $ref path.
// E.g., "./User.yaml" -> "User", "../schemas/User.yaml" -> "User"
func ExtractSchemaNameFromRef(ref string) string {
	if ref == "" {
		return ""
	}
	// Get the base filename
	base := path.Base(ref)
	// Remove .yaml or .yml extension
	if strings.HasSuffix(base, ".yaml") {
		return strings.TrimSuffix(base, ".yaml")
	}
	if strings.HasSuffix(base, ".yml") {
		return strings.TrimSuffix(base, ".yml")
	}
	return base
}
