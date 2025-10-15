package executor_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/archesai/archesai/internal/infrastructure/executor"
)

func TestNewBuilder(t *testing.T) {
	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	if builder == nil {
		t.Fatal("Expected non-nil builder")
	}
}

func TestBuildImage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	// Get project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Build node runner
	config := executor.ImageConfig{
		Name: "node-runner-test",
		DockerfilePath: filepath.Join(
			projectRoot,
			"deployments/containers/runners/node/Dockerfile",
		),
		Tags: []string{"archesai/runner-node:test"},
	}

	ctx := context.Background()
	result, err := builder.BuildImage(ctx, config)
	if err != nil {
		t.Fatalf("Failed to build image: %v\nOutput: %s", err, result.Output)
	}

	if result.Error != nil {
		t.Fatalf("Build failed: %v\nOutput: %s", result.Error, result.Output)
	}

	t.Logf("Successfully built %s with tags %v", result.Name, result.Tags)
}

func TestBuildImages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	// Get project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Build all runners in parallel
	configs := executor.RunnerConfigs(projectRoot)

	ctx := context.Background()
	results := builder.BuildImages(ctx, configs)

	// Check results
	var failed []string
	for _, result := range results {
		if result.Error != nil {
			failed = append(failed, result.Name)
			t.Logf("Failed to build %s: %v", result.Name, result.Error)
		} else {
			t.Logf("Successfully built %s with tags %v", result.Name, result.Tags)
		}
	}

	if len(failed) > 0 {
		t.Fatalf("Failed to build %d images: %v", len(failed), failed)
	}
}

func TestValidateImageConfig(t *testing.T) {
	tests := []struct {
		name   string
		config executor.ImageConfig
	}{
		{
			name: "missing name",
			config: executor.ImageConfig{
				DockerfilePath: "/path/to/Dockerfile",
				Tags:           []string{"test:latest"},
			},
		},
		{
			name: "missing dockerfile path",
			config: executor.ImageConfig{
				Name: "test",
				Tags: []string{"test:latest"},
			},
		},
		{
			name: "missing tags",
			config: executor.ImageConfig{
				Name:           "test",
				DockerfilePath: "/path/to/Dockerfile",
			},
		},
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := builder.BuildImage(ctx, tt.config)

			// All test cases should produce validation errors
			if err == nil {
				t.Error("Expected error but got none")
			}
			if result.Error == nil {
				t.Error("Expected result.Error but got none")
			}
		})
	}
}

func TestRunnerConfigs(t *testing.T) {
	projectRoot := "/test/project"
	configs := executor.RunnerConfigs(projectRoot)

	if len(configs) != 3 {
		t.Fatalf("Expected 3 runner configs, got %d", len(configs))
	}

	// Check that each runner has required fields
	runners := map[string]bool{
		"go-runner":     false,
		"node-runner":   false,
		"python-runner": false,
	}

	for _, cfg := range configs {
		if cfg.Name == "" {
			t.Error("Found config with empty name")
		}
		if cfg.DockerfilePath == "" {
			t.Error("Found config with empty DockerfilePath")
		}
		if len(cfg.Tags) == 0 {
			t.Error("Found config with no tags")
		}

		runners[cfg.Name] = true
	}

	for name, found := range runners {
		if !found {
			t.Errorf("Missing runner config for %s", name)
		}
	}
}

func TestImageExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	ctx := context.Background()

	// Test with a common image that should exist
	exists, err := builder.ImageExists(ctx, "alpine:latest")
	if err != nil {
		t.Fatalf("Failed to check image existence: %v", err)
	}

	// Pull alpine if it doesn't exist (for test consistency)
	if !exists {
		t.Skip("Skipping test - alpine:latest not available and auto-pull not implemented")
	}

	// Test with an image that definitely doesn't exist
	nonExistent := "archesai/nonexistent-image-12345:latest"
	exists, err = builder.ImageExists(ctx, nonExistent)
	if err != nil {
		t.Fatalf("Failed to check non-existent image: %v", err)
	}
	if exists {
		t.Fatalf("Expected image %s to not exist", nonExistent)
	}
}

// findProjectRoot walks up the directory tree to find the project root
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if go.mod exists in current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding go.mod
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
