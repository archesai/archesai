package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// FindProjectRoot walks up the directory tree to find the project root.
// It looks for a go.mod file to identify the root.
// Calls t.Fatal if the root cannot be found.
func FindProjectRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}

// FindProjectRootOrSkip walks up the directory tree to find the project root.
// It skips the test if the root cannot be found instead of failing.
func FindProjectRootOrSkip(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Skipf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Skip("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}
