package ref

import (
	"io/fs"
	"path"
	"strings"
)

// Loader is the interface that type-specific loaders must implement.
// It allows the generic Resolver to delegate parsing to type-aware code.
type Loader[T any] interface {
	// Load parses raw bytes into the target type.
	// name is typically the filename without extension, used as a default identifier.
	Load(data []byte, name string) (*T, error)
}

// Resolver orchestrates reference resolution using a FileResolver and type-specific Loader.
// It caches resolved values to avoid re-parsing the same files.
type Resolver[T any] struct {
	file   *FileResolver
	loader Loader[T]
	cache  map[string]*T
}

// NewResolver creates a new Resolver with the given filesystem, base directory, and loader.
func NewResolver[T any](fsys fs.FS, baseDir string, loader Loader[T]) *Resolver[T] {
	return &Resolver[T]{
		file:   NewFileResolver(fsys, baseDir),
		loader: loader,
		cache:  make(map[string]*T),
	}
}

// FileResolver returns the underlying file resolver.
func (r *Resolver[T]) FileResolver() *FileResolver {
	return r.file
}

// Resolve resolves a reference path and returns the loaded value.
// Results are cached by the resolved path.
func (r *Resolver[T]) Resolve(refPath string) (*T, error) {
	// Check cache first
	if cached, ok := r.cache[refPath]; ok {
		return cached, nil
	}

	// Read the file
	data, err := r.file.ReadFile(refPath)
	if err != nil {
		return nil, err
	}

	// Extract name from path for the loader
	name := extractName(refPath)

	// Load using the type-specific loader
	value, err := r.loader.Load(data, name)
	if err != nil {
		return nil, err
	}

	// Cache the result
	r.cache[refPath] = value

	return value, nil
}

// ResolveFrom resolves a reference path relative to another file.
func (r *Resolver[T]) ResolveFrom(fromPath, refPath string) (*T, error) {
	// Build the resolved path for caching
	fromDir := path.Dir(fromPath)
	resolvedPath := path.Join(fromDir, refPath)
	resolvedPath = path.Clean(resolvedPath)

	// Check cache first
	if cached, ok := r.cache[resolvedPath]; ok {
		return cached, nil
	}

	// Read the file relative to fromPath
	data, err := r.file.ReadFileFrom(fromPath, refPath)
	if err != nil {
		return nil, err
	}

	// Extract name from path for the loader
	name := extractName(refPath)

	// Load using the type-specific loader
	value, err := r.loader.Load(data, name)
	if err != nil {
		return nil, err
	}

	// Cache the result
	r.cache[resolvedPath] = value

	return value, nil
}

// Get retrieves a cached value by path, or nil if not cached.
func (r *Resolver[T]) Get(refPath string) *T {
	return r.cache[refPath]
}

// Set adds or updates a value in the cache.
func (r *Resolver[T]) Set(refPath string, value *T) {
	r.cache[refPath] = value
}

// Cache returns all cached values.
func (r *Resolver[T]) Cache() map[string]*T {
	return r.cache
}

// extractName extracts a name from a file path (filename without extension).
func extractName(refPath string) string {
	base := path.Base(refPath)
	ext := path.Ext(base)
	return strings.TrimSuffix(base, ext)
}
