package spec

import (
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
)

// CompositeFS overlays multiple fs.FS instances.
// Later layers take precedence for individual files.
// For directories, entries are merged from all layers.
type CompositeFS struct {
	layers []fs.FS
}

// NewCompositeFS creates a new composite filesystem from multiple layers.
// Later layers take precedence for individual files.
func NewCompositeFS(layers ...fs.FS) *CompositeFS {
	return &CompositeFS{layers: layers}
}

// Open opens the named file from the first layer that contains it.
func (c *CompositeFS) Open(name string) (fs.File, error) {
	// Try layers in reverse order (later layers have priority)
	for i := len(c.layers) - 1; i >= 0; i-- {
		f, err := c.layers[i].Open(name)
		if err == nil {
			return f, nil
		}
		if !isNotExist(err) {
			return nil, err
		}
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// ReadFile reads the named file from the first layer that contains it.
func (c *CompositeFS) ReadFile(name string) ([]byte, error) {
	// Try layers in reverse order (later layers have priority)
	for i := len(c.layers) - 1; i >= 0; i-- {
		if rfs, ok := c.layers[i].(fs.ReadFileFS); ok {
			data, err := rfs.ReadFile(name)
			if err == nil {
				return data, nil
			}
			if !isNotExist(err) {
				return nil, err
			}
		} else {
			// Fallback to Open + Read
			f, err := c.layers[i].Open(name)
			if err == nil {
				defer f.Close()
				return readAll(f)
			}
			if !isNotExist(err) {
				return nil, err
			}
		}
	}
	return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
}

// ReadDir reads the named directory and merges entries from all layers.
func (c *CompositeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	seen := make(map[string]fs.DirEntry)
	var foundAny bool

	// Go through layers in order (later layers override)
	for _, layer := range c.layers {
		if rfs, ok := layer.(fs.ReadDirFS); ok {
			entries, err := rfs.ReadDir(name)
			if err == nil {
				foundAny = true
				for _, e := range entries {
					seen[e.Name()] = e
				}
			} else if !isNotExist(err) {
				return nil, err
			}
		}
	}

	if !foundAny {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}

	// Convert map to sorted slice
	result := make([]fs.DirEntry, 0, len(seen))
	for _, e := range seen {
		result = append(result, e)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() < result[j].Name()
	})

	return result, nil
}

// Stat returns file info for the named file from the first layer that contains it.
func (c *CompositeFS) Stat(name string) (fs.FileInfo, error) {
	for i := len(c.layers) - 1; i >= 0; i-- {
		if sfs, ok := c.layers[i].(fs.StatFS); ok {
			info, err := sfs.Stat(name)
			if err == nil {
				return info, nil
			}
			if !isNotExist(err) {
				return nil, err
			}
		} else {
			// Fallback to Open + Stat
			f, err := c.layers[i].Open(name)
			if err == nil {
				defer f.Close()
				return f.Stat()
			}
			if !isNotExist(err) {
				return nil, err
			}
		}
	}
	return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrNotExist}
}

// isNotExist checks if an error indicates a file doesn't exist.
func isNotExist(err error) bool {
	return os.IsNotExist(err) || err == fs.ErrNotExist
}

// readAll reads all data from a file.
func readAll(f fs.File) ([]byte, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, info.Size())
	_, err = f.(interface{ Read([]byte) (int, error) }).Read(data)
	return data, err
}

// DiscoverPaths finds all YAML files in the paths/ directory.
func DiscoverPaths(fsys fs.FS) ([]string, error) {
	return discoverYAMLFiles(fsys, "paths")
}

// DiscoverSchemas finds all YAML files in components/schemas/.
func DiscoverSchemas(fsys fs.FS) (map[string]string, error) {
	return discoverNamedFiles(fsys, "components/schemas")
}

// DiscoverResponses finds all YAML files in components/responses/.
func DiscoverResponses(fsys fs.FS) (map[string]string, error) {
	return discoverNamedFiles(fsys, "components/responses")
}

// DiscoverParameters finds all YAML files in components/parameters/.
func DiscoverParameters(fsys fs.FS) (map[string]string, error) {
	return discoverNamedFiles(fsys, "components/parameters")
}

// DiscoverHeaders finds all YAML files in components/headers/.
func DiscoverHeaders(fsys fs.FS) (map[string]string, error) {
	return discoverNamedFiles(fsys, "components/headers")
}

// DiscoverSecuritySchemes finds all YAML files in components/securitySchemes/.
func DiscoverSecuritySchemes(fsys fs.FS) (map[string]string, error) {
	return discoverNamedFiles(fsys, "components/securitySchemes")
}

// discoverYAMLFiles returns all YAML file paths in a directory.
func discoverYAMLFiles(fsys fs.FS, dir string) ([]string, error) {
	if _, err := fs.Stat(fsys, dir); err != nil {
		if isNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			files = append(files, path.Join(dir, name))
		}
	}

	return files, nil
}

// discoverNamedFiles returns a map of name (filename without extension) to file path.
func discoverNamedFiles(fsys fs.FS, dir string) (map[string]string, error) {
	files, err := discoverYAMLFiles(fsys, dir)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, file := range files {
		base := path.Base(file)
		ext := path.Ext(base)
		name := strings.TrimSuffix(base, ext)
		result[name] = file
	}

	return result, nil
}
