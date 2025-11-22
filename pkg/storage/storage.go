// Package storage provides file and object storage abstractions.
package storage

import (
	"os"
	"path/filepath"
)

// Storage abstracts file system operations
type Storage interface {
	// WriteFile writes data to a file at the given path
	WriteFile(path string, data []byte, perm os.FileMode) error

	// ReadFile reads the entire file at the given path
	ReadFile(path string) ([]byte, error)

	// Exists checks if a file or directory exists at the given path
	Exists(path string) (bool, error)

	// MkdirAll creates a directory along with any necessary parents
	MkdirAll(path string, perm os.FileMode) error

	// Remove removes a file or empty directory
	Remove(path string) error

	// RemoveAll removes a path and any children it contains
	RemoveAll(path string) error

	// Stat returns file info for the given path
	Stat(path string) (os.FileInfo, error)

	// Walk walks the file tree rooted at root
	Walk(root string, fn filepath.WalkFunc) error
}
