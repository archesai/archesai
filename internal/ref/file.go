package ref

import (
	"fmt"
	"io/fs"
	"path"
)

// FileResolver handles reading files from a filesystem.
// It wraps fs.FS for testability and handles relative path resolution.
type FileResolver struct {
	fsys    fs.FS
	baseDir string // base directory within the filesystem (can be "." or subdirectory)
}

// NewFileResolver creates a new FileResolver with the given filesystem and base directory.
func NewFileResolver(fsys fs.FS, baseDir string) *FileResolver {
	if baseDir == "" {
		baseDir = "."
	}
	return &FileResolver{fsys: fsys, baseDir: baseDir}
}

// FS returns the underlying filesystem.
func (r *FileResolver) FS() fs.FS {
	return r.fsys
}

// BaseDir returns the base directory.
func (r *FileResolver) BaseDir() string {
	return r.baseDir
}

// ReadFile reads the file at the given path relative to the base directory.
func (r *FileResolver) ReadFile(filePath string) ([]byte, error) {
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

// ReadFileFrom reads a file relative to another file or directory.
func (r *FileResolver) ReadFileFrom(fromPath, filePath string) ([]byte, error) {
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
