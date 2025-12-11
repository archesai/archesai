package storage

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

// NewMemoryStorageWithFiles creates a MemoryStorage pre-populated with files.
// Useful for testing.
func NewMemoryStorageWithFiles(files map[string]string) *MemoryStorage {
	m := NewMemoryStorage()
	for path, content := range files {
		m.files[path] = []byte(content)
		// Add parent directories
		dir := filepath.Dir(path)
		for dir != "." && dir != "/" && dir != "" {
			m.dirs[dir] = true
			dir = filepath.Dir(dir)
		}
	}
	return m
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
func (fi *memFileInfo) Sys() any           { return nil }

// Open opens the named file for reading (implements fs.FS).
func (m *MemoryStorage) Open(name string) (fs.File, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	name = filepath.Clean(name)

	// Check if it's a file
	if data, ok := m.files[name]; ok {
		return &memFile{
			name:   name,
			data:   bytes.NewReader(data),
			info:   &memFileInfo{name: filepath.Base(name), size: int64(len(data)), mode: 0644},
			isDir:  false,
			parent: m,
		}, nil
	}

	// Check if it's a directory
	if _, ok := m.dirs[name]; ok {
		return &memFile{
			name:   name,
			info:   &memFileInfo{name: filepath.Base(name), mode: os.ModeDir | 0755, dir: true},
			isDir:  true,
			parent: m,
		}, nil
	}

	// Check if it's the root
	if name == "." {
		return &memFile{
			name:   name,
			info:   &memFileInfo{name: ".", mode: os.ModeDir | 0755, dir: true},
			isDir:  true,
			parent: m,
		}, nil
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// ReadDir reads the named directory (implements fs.ReadDirFS).
func (m *MemoryStorage) ReadDir(name string) ([]fs.DirEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	name = filepath.Clean(name)
	if name == "." {
		name = ""
	}

	// Check if it's a directory (or root)
	if name != "" {
		if _, ok := m.dirs[name]; !ok {
			// Also check if any files are under this path
			hasChildren := false
			prefix := name + "/"
			for p := range m.files {
				if strings.HasPrefix(p, prefix) {
					hasChildren = true
					break
				}
			}
			if !hasChildren {
				return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
			}
		}
	}

	// Collect immediate children
	seen := make(map[string]fs.DirEntry)
	prefix := ""
	if name != "" {
		prefix = name + "/"
	}

	for p, data := range m.files {
		if name == "" || strings.HasPrefix(p, prefix) {
			rel := strings.TrimPrefix(p, prefix)
			if idx := strings.Index(rel, "/"); idx >= 0 {
				// It's in a subdirectory
				dirName := rel[:idx]
				if _, ok := seen[dirName]; !ok {
					seen[dirName] = &memDirEntry{name: dirName, isDir: true}
				}
			} else {
				// Direct child file
				seen[rel] = &memDirEntry{
					name:  rel,
					isDir: false,
					size:  int64(len(data)),
				}
			}
		}
	}

	for p := range m.dirs {
		if name == "" || strings.HasPrefix(p, prefix) {
			rel := strings.TrimPrefix(p, prefix)
			if rel != "" && !strings.Contains(rel, "/") {
				if _, ok := seen[rel]; !ok {
					seen[rel] = &memDirEntry{name: rel, isDir: true}
				}
			}
		}
	}

	result := make([]fs.DirEntry, 0, len(seen))
	for _, e := range seen {
		result = append(result, e)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})

	return result, nil
}

// memFile implements fs.File for in-memory files.
type memFile struct {
	name   string
	data   *bytes.Reader
	info   *memFileInfo
	isDir  bool
	parent *MemoryStorage
}

func (f *memFile) Stat() (fs.FileInfo, error) { return f.info, nil }
func (f *memFile) Close() error               { return nil }

func (f *memFile) Read(b []byte) (int, error) {
	if f.isDir {
		return 0, &fs.PathError{Op: "read", Path: f.name, Err: fs.ErrInvalid}
	}
	return f.data.Read(b)
}

func (f *memFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if !f.isDir {
		return nil, &fs.PathError{Op: "readdir", Path: f.name, Err: fs.ErrInvalid}
	}
	entries, err := f.parent.ReadDir(f.name)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return entries, nil
	}
	if n > len(entries) {
		n = len(entries)
	}
	if n == 0 {
		return nil, io.EOF
	}
	return entries[:n], nil
}

// memDirEntry implements fs.DirEntry for in-memory directories.
type memDirEntry struct {
	name  string
	isDir bool
	size  int64
}

func (e *memDirEntry) Name() string { return e.name }
func (e *memDirEntry) IsDir() bool  { return e.isDir }
func (e *memDirEntry) Type() fs.FileMode {
	if e.isDir {
		return fs.ModeDir
	}
	return 0
}
func (e *memDirEntry) Info() (fs.FileInfo, error) {
	mode := os.FileMode(0644)
	if e.isDir {
		mode = os.ModeDir | 0755
	}
	return &memFileInfo{name: e.name, size: e.size, mode: mode, dir: e.isDir}, nil
}
