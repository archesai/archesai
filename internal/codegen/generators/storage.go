package generators

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/archesai/archesai/internal/spec"
)

// Storage abstracts file system operations for code generation.
type Storage interface {
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) ([]byte, error)
	Exists(path string) (bool, error)
	MkdirAll(path string, perm os.FileMode) error
	Remove(path string) error
	RemoveAll(path string) error
	Stat(path string) (os.FileInfo, error)
	Walk(root string, fn filepath.WalkFunc) error
	BaseDir() string
}

// LocalStorage implements Storage using the local filesystem.
type LocalStorage struct {
	baseDir      string
	writtenFiles map[string]bool
	mu           sync.Mutex
}

// NewLocalStorage creates a new LocalStorage with the given base directory.
func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{
		baseDir:      baseDir,
		writtenFiles: make(map[string]bool),
	}
}

func (s *LocalStorage) WriteFile(path string, data []byte, perm os.FileMode) error {
	fullPath := filepath.Join(s.baseDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(fullPath, data, perm); err != nil {
		return err
	}
	s.mu.Lock()
	s.writtenFiles[path] = true
	s.mu.Unlock()
	return nil
}

// WasFileWritten returns true if the file was written during this generation run.
func (s *LocalStorage) WasFileWritten(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writtenFiles[path]
}

// WrittenFiles returns all files that were written during this generation run.
func (s *LocalStorage) WrittenFiles() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	files := make([]string, 0, len(s.writtenFiles))
	for path := range s.writtenFiles {
		files = append(files, path)
	}
	return files
}

func (s *LocalStorage) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.baseDir, path))
}

func (s *LocalStorage) Exists(path string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.baseDir, path))
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (s *LocalStorage) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(filepath.Join(s.baseDir, path), perm)
}

func (s *LocalStorage) Remove(path string) error {
	return os.Remove(filepath.Join(s.baseDir, path))
}

func (s *LocalStorage) RemoveAll(path string) error {
	return os.RemoveAll(filepath.Join(s.baseDir, path))
}

func (s *LocalStorage) Stat(path string) (os.FileInfo, error) {
	return os.Stat(filepath.Join(s.baseDir, path))
}

func (s *LocalStorage) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(filepath.Join(s.baseDir, root), fn)
}

func (s *LocalStorage) BaseDir() string {
	return s.baseDir
}

// DatabaseTemplateData holds the data for rendering database repository templates.
// Used by database generators to generate concrete implementations.
type DatabaseTemplateData struct {
	Entity          *spec.Schema
	ProjectName     string
	ModelImportPath string
}

// getDatabaseImportPaths returns the model import path for database generators.
// Since repository interfaces are now in models, this only returns the models path.
func getDatabaseImportPaths(ctx *GeneratorContext, schema *spec.Schema) string {
	internalContext := ctx.InternalContext()
	if schema.IsInternal(internalContext) && schema.XInternal != "" {
		return InternalPackageModelsPath(schema.XInternal)
	}
	return ctx.ProjectName + "/models"
}
