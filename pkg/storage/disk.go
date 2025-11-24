package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

var _ Storage = (*DiskStorage)(nil)

// DiskStorage implements Storage using the actual filesystem
type DiskStorage struct {
	baseDir string // Optional base directory for all operations
}

// NewDiskStorage creates a new DiskStorage
func NewDiskStorage(baseDir string) *DiskStorage {
	return &DiskStorage{
		baseDir: baseDir,
	}
}

// resolvePath resolves a path relative to the base directory
func (d *DiskStorage) resolvePath(path string) string {
	if d.baseDir == "" {
		return path
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(d.baseDir, path)
}

// WriteFile writes data to a file at the given path
func (d *DiskStorage) WriteFile(path string, data []byte, perm os.FileMode) error {
	fullPath := d.resolvePath(path)
	dir := filepath.Dir(fullPath)

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return os.WriteFile(fullPath, data, perm)
}

// ReadFile reads the entire file at the given path
func (d *DiskStorage) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(d.resolvePath(path))
}

// Exists checks if a file or directory exists at the given path
func (d *DiskStorage) Exists(path string) (bool, error) {
	_, err := os.Stat(d.resolvePath(path))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MkdirAll creates a directory along with any necessary parents
func (d *DiskStorage) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(d.resolvePath(path), perm)
}

// Remove removes a file or empty directory
func (d *DiskStorage) Remove(path string) error {
	return os.Remove(d.resolvePath(path))
}

// RemoveAll removes a path and any children it contains
func (d *DiskStorage) RemoveAll(path string) error {
	return os.RemoveAll(d.resolvePath(path))
}

// Stat returns file info for the given path
func (d *DiskStorage) Stat(path string) (os.FileInfo, error) {
	return os.Stat(d.resolvePath(path))
}

// Walk walks the file tree rooted at root
func (d *DiskStorage) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(d.resolvePath(root), fn)
}

// BaseDir returns the base directory of the storage
func (d *DiskStorage) BaseDir() string {
	return d.baseDir
}
