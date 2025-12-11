package storage

import "testing"

func TestNewTrackedStorage(t *testing.T) {
	mem := NewMemoryStorage()
	tracked := NewTrackedStorage(mem)

	if tracked == nil {
		t.Fatal("NewTrackedStorage() returned nil")
	}
	if tracked.Storage == nil {
		t.Error("Storage should not be nil")
	}
	if tracked.writtenFiles == nil {
		t.Error("writtenFiles should not be nil")
	}
}

func TestTrackedStorage_WriteFile(t *testing.T) {
	mem := NewMemoryStorage()
	tracked := NewTrackedStorage(mem)

	err := tracked.WriteFile("test.txt", []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Verify file was written to underlying storage
	data, err := mem.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("file content = %q, want %q", string(data), "hello")
	}
}

func TestTrackedStorage_WasFileWritten(t *testing.T) {
	mem := NewMemoryStorage()
	tracked := NewTrackedStorage(mem)

	// Before writing
	if tracked.WasFileWritten("test.txt") {
		t.Error("WasFileWritten() should return false before writing")
	}

	// After writing
	if err := tracked.WriteFile("test.txt", []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if !tracked.WasFileWritten("test.txt") {
		t.Error("WasFileWritten() should return true after writing")
	}

	// Different file
	if tracked.WasFileWritten("other.txt") {
		t.Error("WasFileWritten() should return false for unwritten file")
	}
}

func TestTrackedStorage_WrittenFiles(t *testing.T) {
	mem := NewMemoryStorage()
	tracked := NewTrackedStorage(mem)

	// No files written
	files := tracked.WrittenFiles()
	if len(files) != 0 {
		t.Errorf("WrittenFiles() = %v, want empty slice", files)
	}

	// Write some files
	if err := tracked.WriteFile("a.txt", []byte("a"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := tracked.WriteFile("b.txt", []byte("b"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := tracked.WriteFile("c.txt", []byte("c"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	files = tracked.WrittenFiles()
	if len(files) != 3 {
		t.Errorf("len(WrittenFiles()) = %d, want 3", len(files))
	}

	// Verify all files are present
	fileSet := make(map[string]bool)
	for _, f := range files {
		fileSet[f] = true
	}

	for _, expected := range []string{"a.txt", "b.txt", "c.txt"} {
		if !fileSet[expected] {
			t.Errorf("WrittenFiles() missing %q", expected)
		}
	}
}

func TestTrackedStorage_DelegatesOtherOperations(t *testing.T) {
	mem := NewMemoryStorage()
	tracked := NewTrackedStorage(mem)

	// Write a file
	if err := tracked.WriteFile("test.txt", []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// ReadFile should delegate
	data, err := tracked.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("ReadFile() = %q, want %q", string(data), "hello")
	}

	// Exists should delegate
	exists, err := tracked.Exists("test.txt")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() should return true for existing file")
	}

	// BaseDir should delegate
	baseDir := tracked.BaseDir()
	if baseDir != "" {
		t.Errorf("BaseDir() = %q, want empty string for MemoryStorage", baseDir)
	}
}
