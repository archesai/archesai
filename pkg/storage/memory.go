package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var _ Storage = (*MemoryStorage)(nil)

// MemoryStorage implements Storage in memory (useful for testing)
type MemoryStorage struct {
	mu      sync.RWMutex
	files   map[string][]byte
	dirs    map[string]bool
	baseDir string
}

// NewMemoryStorage creates a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

// WriteFile writes data to a file at the given path
func (m *MemoryStorage) WriteFile(path string, data []byte, _ os.FileMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clean the path
	path = filepath.Clean(path)

	// Mark all parent directories as existing
	dir := filepath.Dir(path)
	for dir != "." && dir != "/" && dir != "" {
		m.dirs[dir] = true
		dir = filepath.Dir(dir)
	}

	// Store the file content (make a copy to avoid external modifications)
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	m.files[path] = dataCopy

	return nil
}

// ReadFile reads the entire file at the given path
func (m *MemoryStorage) ReadFile(path string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path = filepath.Clean(path)
	data, ok := m.files[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	// Return a copy to avoid external modifications
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	return dataCopy, nil
}

// Exists checks if a file or directory exists at the given path
func (m *MemoryStorage) Exists(path string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path = filepath.Clean(path)

	// Check if it's a file
	if _, ok := m.files[path]; ok {
		return true, nil
	}

	// Check if it's a directory
	if _, ok := m.dirs[path]; ok {
		return true, nil
	}

	return false, nil
}

// MkdirAll creates a directory along with any necessary parents
func (m *MemoryStorage) MkdirAll(path string, _ os.FileMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	path = filepath.Clean(path)

	// Mark the directory and all parents as existing
	for path != "." && path != "/" && path != "" {
		m.dirs[path] = true
		path = filepath.Dir(path)
	}

	return nil
}

// Remove removes a file or empty directory
func (m *MemoryStorage) Remove(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	path = filepath.Clean(path)

	// Try to remove as a file
	if _, ok := m.files[path]; ok {
		delete(m.files, path)
		return nil
	}

	// Try to remove as a directory (check if empty)
	if _, ok := m.dirs[path]; ok {
		// Check if directory is empty
		prefix := path + string(filepath.Separator)
		for p := range m.files {
			if strings.HasPrefix(p, prefix) {
				return fmt.Errorf("directory not empty: %s", path)
			}
		}
		for p := range m.dirs {
			if p != path && strings.HasPrefix(p, prefix) {
				return fmt.Errorf("directory not empty: %s", path)
			}
		}
		delete(m.dirs, path)
		return nil
	}

	return os.ErrNotExist
}

// RemoveAll removes a path and any children it contains
func (m *MemoryStorage) RemoveAll(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	path = filepath.Clean(path)
	prefix := path + string(filepath.Separator)

	// Remove the path itself if it's a file
	delete(m.files, path)

	// Remove all files under this path
	for p := range m.files {
		if p == path || strings.HasPrefix(p, prefix) {
			delete(m.files, p)
		}
	}

	// Remove all directories under this path
	delete(m.dirs, path)
	for p := range m.dirs {
		if strings.HasPrefix(p, prefix) {
			delete(m.dirs, p)
		}
	}

	return nil
}

// Stat returns file info for the given path
func (m *MemoryStorage) Stat(path string) (os.FileInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path = filepath.Clean(path)

	// Check if it's a file
	if data, ok := m.files[path]; ok {
		return &memFileInfo{
			name: filepath.Base(path),
			size: int64(len(data)),
			mode: 0644,
			dir:  false,
		}, nil
	}

	// Check if it's a directory
	if _, ok := m.dirs[path]; ok {
		return &memFileInfo{
			name: filepath.Base(path),
			size: 0,
			mode: os.ModeDir | 0755,
			dir:  true,
		}, nil
	}

	return nil, os.ErrNotExist
}

// Walk walks the file tree rooted at root
func (m *MemoryStorage) Walk(root string, fn filepath.WalkFunc) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	root = filepath.Clean(root)

	// Check if root exists
	rootExists := false
	if _, ok := m.files[root]; ok {
		rootExists = true
	} else if _, ok := m.dirs[root]; ok {
		rootExists = true
	}

	if !rootExists {
		return fn(root, nil, os.ErrNotExist)
	}

	// Build a sorted list of all paths
	paths := make(map[string]bool)
	paths[root] = true

	for path := range m.files {
		if path == root || strings.HasPrefix(path, root+string(filepath.Separator)) {
			paths[path] = true
			// Add all parent directories
			dir := filepath.Dir(path)
			for dir != "." && dir != "/" && dir != "" && (dir == root || strings.HasPrefix(dir, root+string(filepath.Separator))) {
				paths[dir] = true
				dir = filepath.Dir(dir)
			}
		}
	}

	for path := range m.dirs {
		if path == root || strings.HasPrefix(path, root+string(filepath.Separator)) {
			paths[path] = true
		}
	}

	// Walk through paths in sorted order
	for path := range paths {
		info, err := m.Stat(path)
		if err := fn(path, info, err); err != nil {
			return err
		}
	}

	return nil
}

// GetFiles returns all files in the memory storage (useful for testing)
func (m *MemoryStorage) GetFiles() map[string][]byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string][]byte)
	for path, data := range m.files {
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		result[path] = dataCopy
	}
	return result
}

// BaseDir returns the base directory of the storage
func (m *MemoryStorage) BaseDir() string {
	return m.baseDir
}

// memFileInfo implements os.FileInfo for in-memory files
type memFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	dir     bool
}

func (fi *memFileInfo) Name() string       { return fi.name }
func (fi *memFileInfo) Size() int64        { return fi.size }
func (fi *memFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *memFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *memFileInfo) IsDir() bool        { return fi.dir }
func (fi *memFileInfo) Sys() interface{}   { return nil }
