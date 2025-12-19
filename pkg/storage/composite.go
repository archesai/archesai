package storage

import (
	"io/fs"
	"os"
	"sort"
)

// Composite overlays multiple fs.FS instances.
// Later layers take precedence for individual files.
// For directories, entries are merged from all layers.
type Composite struct {
	layers []fs.FS
}

// NewComposite creates a new composite filesystem from multiple layers.
// Later layers take precedence for individual files.
func NewComposite(layers ...fs.FS) *Composite {
	return &Composite{layers: layers}
}

// Open opens the named file from the first layer that contains it.
func (c *Composite) Open(name string) (fs.File, error) {
	// Try layers in reverse order (later layers have priority)
	for i := len(c.layers) - 1; i >= 0; i-- {
		f, err := c.layers[i].Open(name)
		if err == nil {
			return f, nil
		}
		if !compositeIsNotExist(err) {
			return nil, err
		}
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// ReadFile reads the named file from the first layer that contains it.
func (c *Composite) ReadFile(name string) ([]byte, error) {
	// Try layers in reverse order (later layers have priority)
	for i := len(c.layers) - 1; i >= 0; i-- {
		if rfs, ok := c.layers[i].(fs.ReadFileFS); ok {
			data, err := rfs.ReadFile(name)
			if err == nil {
				return data, nil
			}
			if !compositeIsNotExist(err) {
				return nil, err
			}
		} else {
			// Fallback to Open + Read
			f, err := c.layers[i].Open(name)
			if err == nil {
				defer func() { _ = f.Close() }()
				return compositeReadAll(f)
			}
			if !compositeIsNotExist(err) {
				return nil, err
			}
		}
	}
	return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
}

// ReadDir reads the named directory and merges entries from all layers.
func (c *Composite) ReadDir(name string) ([]fs.DirEntry, error) {
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
			} else if !compositeIsNotExist(err) {
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
func (c *Composite) Stat(name string) (fs.FileInfo, error) {
	for i := len(c.layers) - 1; i >= 0; i-- {
		if sfs, ok := c.layers[i].(fs.StatFS); ok {
			info, err := sfs.Stat(name)
			if err == nil {
				return info, nil
			}
			if !compositeIsNotExist(err) {
				return nil, err
			}
		} else {
			// Fallback to Open + Stat
			f, err := c.layers[i].Open(name)
			if err == nil {
				defer func() { _ = f.Close() }()
				return f.Stat()
			}
			if !compositeIsNotExist(err) {
				return nil, err
			}
		}
	}
	return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrNotExist}
}

// compositeIsNotExist checks if an error indicates a file doesn't exist.
func compositeIsNotExist(err error) bool {
	return os.IsNotExist(err) || err == fs.ErrNotExist
}

// compositeReadAll reads all data from a file.
func compositeReadAll(f fs.File) ([]byte, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, info.Size())
	_, err = f.(interface{ Read([]byte) (int, error) }).Read(data)
	return data, err
}
