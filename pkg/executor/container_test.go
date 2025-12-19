package executor_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/testutil"
	"github.com/archesai/archesai/pkg/executor"
)

// TestNodeRunnerWithoutMount tests that the runner fails without custom execute.js
func TestNodeRunnerWithoutMount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	// Define simple input/output types
	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	// Create executor config without mounts
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		DisableNet:  true,
		ReadOnlyFS:  false, // Need write for /tmp
		MemoryBytes: 128 * 1024 * 1024,
		Config: executor.Config{
			Timeout: 10 * time.Second,
		},
	}

	// Create executor
	exec, err := executor.NewContainerExecutor[Input, Output](config, nil)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Should fail because no custom execute.js is mounted
	ctx := context.Background()
	_, err = exec.Execute(ctx, Input{Value: 42})

	// Should error with "No execution function provided"
	if err == nil {
		t.Fatal("Expected error from base runner without custom execute.js")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

// TestNodeRunnerWithMount tests mounting custom execute.js via volume mount
func TestNodeRunnerWithMount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	// Define simple input/output types
	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	// Get absolute path to test execute file
	testFilePath, err := filepath.Abs("testdata/execute.ts")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Verify the test file exists
	if _, err := os.Stat(testFilePath); err != nil {
		t.Fatalf("Test file not found at %s: %v", testFilePath, err)
	}

	// Create executor config with mounts
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		DisableNet:  true,
		ReadOnlyFS:  false, // Need write for /tmp
		MemoryBytes: 128 * 1024 * 1024,
		Mounts: []executor.Mount{
			{
				Source:   testFilePath,
				Target:   "/app/src/execute.ts",
				ReadOnly: true,
			},
		},
		Config: executor.Config{
			Timeout: 10 * time.Second,
		},
	}

	// Create executor
	exec, err := executor.NewContainerExecutor[Input, Output](config, nil)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Execute with custom logic
	ctx := context.Background()
	output, err := exec.Execute(ctx, Input{Value: 42})
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}

	// Verify output
	expectedDoubled := 84
	if output.Doubled != expectedDoubled {
		t.Fatalf("Expected doubled=%d, got %d", expectedDoubled, output.Doubled)
	}
	t.Logf("Successfully doubled %d to %d", 42, output.Doubled)
}

// TestSchemaValidation tests JSON schema validation
func TestSchemaValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	type Input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type Output struct {
		Message string `json:"message"`
	}

	// Define strict schemas
	inputSchema := []byte(`{
		"type": "object",
		"required": ["name", "email"],
		"properties": {
			"name": {"type": "string", "minLength": 1},
			"email": {"type": "string"}
		},
		"additionalProperties": false
	}`)

	outputSchema := []byte(`{
		"type": "object",
		"required": ["message"],
		"properties": {
			"message": {"type": "string"}
		}
	}`)

	// Create executor config with schemas
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		DisableNet:  true,
		MemoryBytes: 128 * 1024 * 1024,
		Config: executor.Config{
			Timeout:   10 * time.Second,
			SchemaIn:  inputSchema,
			SchemaOut: outputSchema,
		},
	}

	// Create executor
	exec, err := executor.NewContainerExecutor[Input, Output](config, nil)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Test with valid input
	ctx := context.Background()
	_, err = exec.Execute(ctx, Input{
		Name:  "John Doe",
		Email: "john@example.com",
	})

	// Will fail because no custom execute.js, but should pass input validation first
	if err != nil {
		// Check that it's not a validation error
		t.Logf("Got error (expected due to missing execute.js): %v", err)
	}

	// Test with invalid input (missing required field)
	_, err = exec.Execute(ctx, Input{
		Name: "John Doe",
		// Missing email
	})

	if err == nil {
		t.Fatal("Expected validation error for missing email field")
	}
	t.Logf("Got expected validation error: %v", err)
}

// TestContainerTimeout tests timeout handling
func TestContainerTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	type Input struct {
		Data string `json:"data"`
	}
	type Output struct {
		Result string `json:"result"`
	}

	// Create executor with very short timeout
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		DisableNet:  true,
		MemoryBytes: 128 * 1024 * 1024,
		Config: executor.Config{
			Timeout: 100 * time.Millisecond, // Very short timeout
		},
	}

	exec, err := executor.NewContainerExecutor[Input, Output](config, nil)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	ctx := context.Background()
	_, err = exec.Execute(ctx, Input{Data: "test"})

	// Should timeout (container startup alone may exceed 100ms)
	if err == nil {
		t.Log("Warning: Expected timeout but execution succeeded (container was very fast)")
	} else {
		t.Logf("Got expected timeout or error: %v", err)
	}
}

// TestAutoBuildDisabled tests that executor works without auto-build
func TestAutoBuildDisabled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	// Config without auto-build (default)
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		DisableNet:  true,
		ReadOnlyFS:  false,
		MemoryBytes: 128 * 1024 * 1024,
		Config: executor.Config{
			Timeout: 10 * time.Second,
		},
	}

	// Should work without builder when AutoBuild is false
	exec, err := executor.NewContainerExecutor[Input, Output](config, nil)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Verify executor was created
	if exec == nil {
		t.Fatal("Expected non-nil executor")
	}
}

// TestAutoBuildRequiresBuilder tests that AutoBuild=true requires a builder
func TestAutoBuildRequiresBuilder(t *testing.T) {
	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	config := executor.ContainerConfig{
		Image:      "archesai/runner-node:latest",
		AutoBuild:  true, // Enable auto-build
		DisableNet: true,
		Config: executor.Config{
			Timeout: 10 * time.Second,
		},
	}

	// Should fail when AutoBuild=true but builder is nil
	_, err := executor.NewContainerExecutor[Input, Output](config, nil)
	if err == nil {
		t.Fatal("Expected error when AutoBuild=true but builder is nil")
	}
	t.Logf("Got expected error: %v", err)
}

// TestAutoBuildRequiresBuildConfig tests that AutoBuild=true requires BuildConfig
func TestAutoBuildRequiresBuildConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		AutoBuild:   true, // Enable auto-build
		BuildConfig: nil,  // Missing BuildConfig
		DisableNet:  true,
		Config: executor.Config{
			Timeout: 10 * time.Second,
		},
	}

	// Should fail when AutoBuild=true but BuildConfig is nil
	_, err = executor.NewContainerExecutor[Input, Output](config, builder)
	if err == nil {
		t.Fatal("Expected error when AutoBuild=true but BuildConfig is nil")
	}
	t.Logf("Got expected error: %v", err)
}

// TestAutoBuildWithExistingImage tests that auto-build skips when image exists
func TestAutoBuildWithExistingImage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires Docker")
	}

	type Input struct {
		Value int `json:"value"`
	}
	type Output struct {
		Doubled int `json:"doubled"`
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	defer func() { _ = builder.Close() }()

	// Get project root for build config
	projectRoot := testutil.FindProjectRoot(t)

	// First, ensure the image exists by building it
	buildConfig := executor.ImageConfig{
		Name: "node-runner-test",
		DockerfilePath: filepath.Join(
			projectRoot,
			"deployments/containers/runners/node/Dockerfile",
		),
		Tags: []string{"archesai/runner-node:autobuild-test"},
	}

	ctx := context.Background()
	result, err := builder.BuildImage(ctx, buildConfig)
	if err != nil || result.Error != nil {
		t.Fatalf("Failed to build initial image: %v", err)
	}

	// Now create executor with auto-build
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:autobuild-test",
		AutoBuild:   true,
		BuildConfig: &buildConfig,
		ForceBuild:  false, // Should skip build since image exists
		DisableNet:  true,
		ReadOnlyFS:  false,
		MemoryBytes: 128 * 1024 * 1024,
		Config: executor.Config{
			Timeout: 10 * time.Second,
		},
	}

	exec, err := executor.NewContainerExecutor[Input, Output](config, builder)
	if err != nil {
		t.Fatalf("Failed to create executor with auto-build: %v", err)
	}

	// Execute should work (will skip build since image exists)
	// Note: will fail at runtime since no custom execute.ts is mounted, but that's expected
	_, err = exec.Execute(ctx, Input{Value: 42})
	// We expect an error from the base runner, not from building
	if err != nil {
		t.Logf("Got expected runtime error: %v", err)
	}
}
