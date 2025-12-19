package storage

import (
	"os"
	"sync"
)

// TrackedStorage wraps a Storage implementation to track which files were written.
type TrackedStorage struct {
	Storage
	writtenFiles map[string]bool
	mu           sync.Mutex
}

// NewTrackedStorage creates a new TrackedStorage wrapping the given Storage.
func NewTrackedStorage(s Storage) *TrackedStorage {
	return &TrackedStorage{
		Storage:      s,
		writtenFiles: make(map[string]bool),
	}
}

// WriteFile writes a file and tracks that it was written.
func (s *TrackedStorage) WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := s.Storage.WriteFile(path, data, perm); err != nil {
		return err
	}
	s.mu.Lock()
	s.writtenFiles[path] = true
	s.mu.Unlock()
	return nil
}

// WasFileWritten returns true if the file was written during this generation run.
func (s *TrackedStorage) WasFileWritten(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writtenFiles[path]
}

// WrittenFiles returns all files that were written during this generation run.
func (s *TrackedStorage) WrittenFiles() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	files := make([]string, 0, len(s.writtenFiles))
	for path := range s.writtenFiles {
		files = append(files, path)
	}
	return files
}
