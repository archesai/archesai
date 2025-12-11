package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/archesai/archesai/internal/testutil"
)

func TestBundler_Integration(t *testing.T) {
	// Find examples/basic/api/openapi.yaml from repo root
	rootDir := testutil.FindProjectRootOrSkip(t)
	specPath := filepath.Join(rootDir, "examples", "basic", "api", "openapi.yaml")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Skipf("spec file not found: %s", specPath)
	}

	baseFS := os.DirFS(filepath.Dir(specPath))
	doc, err := NewOpenAPIDocumentFromFS(baseFS, filepath.Base(specPath))
	if err != nil {
		t.Fatalf("failed to load document: %v", err)
	}

	bundler := NewBundler(doc)
	yamlBytes, err := bundler.BundleToYAML()
	if err != nil {
		t.Fatalf("failed to bundle: %v", err)
	}

	yaml := string(yamlBytes)

	// Verify structure
	if !strings.Contains(yaml, "openapi: 3.1.0") {
		t.Error("missing openapi version")
	}
	if !strings.Contains(yaml, "paths:") {
		t.Error("missing paths section")
	}
	if !strings.Contains(yaml, "components:") {
		t.Error("missing components section")
	}
	if !strings.Contains(yaml, "schemas:") {
		t.Error("missing schemas section")
	}

	// Verify refs are converted to internal format
	if strings.Contains(yaml, "../") || strings.Contains(yaml, "./") {
		t.Error("file refs not converted to internal refs")
	}

	t.Logf("Bundled YAML (first 2000 chars):\n%s", yaml[:min(2000, len(yaml))])
}
