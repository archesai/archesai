package located

import "path/filepath"

// Located wraps a value with its file path, providing path-related helpers.
type Located[T any] struct {
	Value *T
	Path  string
}

// Dir returns the directory containing the file.
func (l *Located[T]) Dir() string {
	if l.Path == "" {
		return "."
	}
	return filepath.Dir(l.Path)
}
