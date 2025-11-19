package codegen

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/archesai/archesai/pkg/executor"
)

// OrvalInput represents the input for orval generation
type OrvalInput struct {
	OpenAPI string `json:"openapi"` // OpenAPI specification content
}

// OrvalOutput represents the output from orval generation
type OrvalOutput struct {
	Files map[string]string `json:"files"` // Map of file paths to contents
}

// GenerateJSClient generates a JavaScript/TypeScript client using orval in a container
func (g *Generator) GenerateJSClient(specPath string, outputDir string) error {
	// Read the OpenAPI specification
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	// Get absolute path for orval.ts mount
	orvalTsPath, err := filepath.Abs("./deployments/containers/runners/node/src/orval.ts")
	if err != nil {
		return fmt.Errorf("failed to resolve orval.ts path: %w", err)
	}

	// Get absolute path for Dockerfile
	dockerfilePath, err := filepath.Abs("./deployments/containers/runners/node/Dockerfile")
	if err != nil {
		return fmt.Errorf("failed to resolve Dockerfile path: %w", err)
	}

	fetcherTsPath, err := filepath.Abs("./web/client/src/fetcher.ts")
	if err != nil {
		return fmt.Errorf("failed to resolve fetcher.ts path: %w", err)
	}

	// Create builder for auto-building the image
	builder, err := executor.NewBuilder()
	if err != nil {
		return fmt.Errorf("failed to create builder: %w", err)
	}
	defer func() {
		if closeErr := builder.Close(); closeErr != nil {
			slog.Warn("failed to close builder", "error", closeErr)
		}
	}()

	// Define input and output schemas
	inputSchema := []byte(`{
		"type": "object",
		"required": ["openapi"],
		"properties": {
			"openapi": {
				"type": "string",
				"description": "OpenAPI specification content"
			}
		}
	}`)

	outputSchema := []byte(`{
		"type": "object",
		"required": ["files"],
		"properties": {
			"files": {
				"type": "object",
				"additionalProperties": { "type": "string" }
			}
		}
	}`)

	// Build configuration for node runner
	additionalPackages := "orval"
	buildConfig := executor.ImageConfig{
		Name:           "node-runner",
		DockerfilePath: dockerfilePath,
		BuildArgs: map[string]*string{
			"ADDITIONAL_PACKAGES": &additionalPackages,
		},
		Tags: []string{"archesai/runner-node:latest"},
	}

	// Create executor configuration with auto-build
	config := executor.ContainerConfig{
		Image:       "archesai/runner-node:latest",
		AutoBuild:   true,
		BuildConfig: &buildConfig,
		ForceBuild:  true,              // Use cached image if available
		DisableNet:  false,             // Orval may need network for downloading types
		ReadOnlyFS:  false,             // Orval needs to write temporary files
		MemoryBytes: 512 * 1024 * 1024, // 512MB for orval processing
		CPUShares:   1024,
		Config: executor.Config{
			Timeout:   60 * time.Second, // 60 seconds
			SchemaIn:  inputSchema,
			SchemaOut: outputSchema,
		},
		Mounts: []executor.Mount{
			{
				Source:   orvalTsPath,
				Target:   "/app/src/execute.ts",
				ReadOnly: true,
			},
			{
				Source:   fetcherTsPath,
				Target:   "/app/src/fetcher.ts",
				ReadOnly: true,
			},
		},
	}

	// Create container executor with auto-build
	exec, err := executor.NewContainerExecutor[OrvalInput, OrvalOutput](config, builder)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	// Prepare orval configuration (mimicking the existing orval.config.ts)
	input := OrvalInput{
		OpenAPI: string(specContent),
	}

	// Execute orval in container
	ctx := context.Background()
	output, err := exec.Execute(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to execute orval: %w", err)
	}

	// Write generated files to output directory
	for filePath, content := range output.Files {
		fullPath := filepath.Join(outputDir, filePath)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}

	return nil
}
