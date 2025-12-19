package spec

import (
	"io/fs"
	"os"
	"path"
	"strings"
)

// ComponentKind represents a type of OpenAPI component.
type ComponentKind string

// OpenAPI component kind constants.
const (
	ComponentSchemas         ComponentKind = "components/schemas"
	ComponentResponses       ComponentKind = "components/responses"
	ComponentParameters      ComponentKind = "components/parameters"
	ComponentHeaders         ComponentKind = "components/headers"
	ComponentSecuritySchemes ComponentKind = "components/securitySchemes"
)

// DiscoverComponents finds all YAML files for a given component kind.
// Returns a map of name (filename without extension) to file path.
func DiscoverComponents(fsys fs.FS, kind ComponentKind) (map[string]string, error) {
	return discoverNamedFiles(fsys, string(kind))
}

// DiscoverPaths finds all YAML files in the paths/ directory.
// Returns a slice of file paths (not a map, since paths use x-path for naming).
func DiscoverPaths(fsys fs.FS) ([]string, error) {
	return discoverYAMLFiles(fsys, "paths")
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

// isNotExist checks if an error indicates a file doesn't exist.
func isNotExist(err error) bool {
	return os.IsNotExist(err) || err == fs.ErrNotExist
}
