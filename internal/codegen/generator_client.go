package codegen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/archesai/archesai/pkg/executor"
)

// OrvalInput is the input structure for the Orval TypeScript client generator.
type OrvalInput struct {
	OpenAPI string `json:"openapi"`
}

// OrvalOutput is the output structure from the Orval TypeScript client generator.
type OrvalOutput struct {
	Files map[string]string `json:"files"`
}

// ClientGenerator generates TypeScript API client code using Orval.
type ClientGenerator struct{}

// Name returns the generator name.
func (g *ClientGenerator) Name() string { return "client" }

// Priority returns the generator priority.
func (g *ClientGenerator) Priority() int { return PriorityNormal }

// Generate creates TypeScript API client code from the OpenAPI spec.
func (g *ClientGenerator) Generate(ctx *GeneratorContext) error {
	specContent, err := os.ReadFile(ctx.SpecPath)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	orvalTsPath, err := filepath.Abs("./deployments/containers/runners/node/src/orval.ts")
	if err != nil {
		return fmt.Errorf("failed to resolve orval.ts path: %w", err)
	}

	dockerfilePath, err := filepath.Abs("./deployments/containers/runners/node/Dockerfile")
	if err != nil {
		return fmt.Errorf("failed to resolve Dockerfile path: %w", err)
	}

	fetcherTsPath, err := filepath.Abs("./pkg/client/src/fetcher.ts")
	if err != nil {
		return fmt.Errorf("failed to resolve fetcher.ts path: %w", err)
	}

	builder, err := executor.NewBuilder()
	if err != nil {
		return fmt.Errorf("failed to create builder: %w", err)
	}
	defer func() { _ = builder.Close() }()

	inputSchema := []byte(
		`{"type":"object","required":["openapi"],"properties":{"openapi":{"type":"string"}}}`,
	)
	outputSchema := []byte(
		`{"type":"object","required":["files"],"properties":{"files":{"type":"object","additionalProperties":{"type":"string"}}}}`,
	)

	additionalPackages := "orval"
	config := executor.ContainerConfig{
		Image:     "archesai/runner-node:latest",
		AutoBuild: true,
		BuildConfig: &executor.ImageConfig{
			Name:           "node-runner",
			DockerfilePath: dockerfilePath,
			BuildArgs:      map[string]*string{"ADDITIONAL_PACKAGES": &additionalPackages},
			Tags:           []string{"archesai/runner-node:latest"},
		},
		ForceBuild:  true,
		DisableNet:  false,
		ReadOnlyFS:  false,
		MemoryBytes: 512 * 1024 * 1024,
		CPUShares:   1024,
		Config: executor.Config{
			Timeout:   60 * time.Second,
			SchemaIn:  inputSchema,
			SchemaOut: outputSchema,
		},
		Mounts: []executor.Mount{
			{Source: orvalTsPath, Target: "/app/src/execute.ts", ReadOnly: true},
			{Source: fetcherTsPath, Target: "/app/src/fetcher.ts", ReadOnly: true},
		},
	}

	exec, err := executor.NewContainerExecutor[OrvalInput, OrvalOutput](config, builder)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	output, err := exec.Execute(context.Background(), OrvalInput{OpenAPI: string(specContent)})
	if err != nil {
		return fmt.Errorf("failed to execute orval: %w", err)
	}

	for filePath, content := range output.Files {
		fullPath := filepath.Join(ctx.Storage.BaseDir(), "src", "lib", "client", filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}

	return nil
}
