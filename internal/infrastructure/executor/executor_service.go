package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/models"
	"github.com/archesai/archesai/internal/core/repositories"
)

// ExecutorService defines operations for custom executors with generic input/output types
type ExecutorService[A any, B any] interface { //nolint:revive // stutters but needed for backward compatibility
	// Execute runs an executor with typed input and returns typed output
	Execute(ctx context.Context, executorID string, input A) (ExecuteResult[B], error)

	// ValidateInput validates input against executor's input schema
	ValidateInput(ctx context.Context, executorID string, input A) error

	// ValidateOutput validates output against executor's output schema
	ValidateOutput(ctx context.Context, executorID string, output B) error

	// BuildExecutor pre-builds the Docker image for an executor
	BuildExecutor(ctx context.Context, executorID string) error

	// GetExecutor returns the executor configuration (for debugging/testing)
	GetExecutor(ctx context.Context, executorID string) (*models.Executor, error)
}

// ExecuteResult contains the result of an executor execution
type ExecuteResult[B any] struct {
	Output          B
	ExecutionTimeMs int64
	Logs            string
}

// executorService implements the generic ExecutorService interface
type executorService[A any, B any] struct {
	repo    repositories.ExecutorRepository
	builder *Builder

	// Base images for each language
	baseImages map[ExecutorLanguage]string

	// Base dockerfile paths
	dockerfiles map[ExecutorLanguage]string
}

// NewExecutorService creates a new executor service with generic types
func NewExecutorService[A any, B any](
	repo repositories.ExecutorRepository,
	builder *Builder,
) ExecutorService[A, B] {
	return &executorService[A, B]{
		repo:    repo,
		builder: builder,
		baseImages: map[ExecutorLanguage]string{
			ExecutorLanguageNodejs: "archesai/runner-node:latest",
			ExecutorLanguagePython: "archesai/runner-python:latest",
			ExecutorLanguageGo:     "archesai/runner-go:latest",
		},
		dockerfiles: map[ExecutorLanguage]string{
			ExecutorLanguageNodejs: "./deployments/containers/runners/node/Dockerfile",
			ExecutorLanguagePython: "./deployments/containers/runners/python/Dockerfile",
			ExecutorLanguageGo:     "./deployments/containers/runners/go/Dockerfile",
		},
	}
}

func (s *executorService[A, B]) Execute(
	ctx context.Context,
	executorID string,
	input A,
) (ExecuteResult[B], error) {
	var zero ExecuteResult[B]

	// 1. Get executor from database
	exec, err := s.GetExecutor(ctx, executorID)
	if err != nil {
		return zero, fmt.Errorf("get executor: %w", err)
	}

	if !exec.IsActive {
		return zero, fmt.Errorf("executor %s is not active", executorID)
	}

	// 2. Prepare temporary directory for code mounting
	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("executor-%s-*", executorID))
	if err != nil {
		return zero, fmt.Errorf("create temp dir: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// 3. Write execute code to temp file
	executeFile := s.getExecuteFileName(ExecutorLanguage(exec.Language))
	executePath := filepath.Join(tmpDir, executeFile)
	if err := os.WriteFile(executePath, []byte(exec.ExecuteCode), 0644); err != nil {
		return zero, fmt.Errorf("write execute code: %w", err)
	}

	// 4. Prepare mounts
	mounts := []Mount{
		{
			Source:   executePath,
			Target:   fmt.Sprintf("/app/src/%s", executeFile),
			ReadOnly: true,
		},
	}

	// Add extra files if provided
	if exec.ExtraFiles != nil {
		var extraFiles []struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(*exec.ExtraFiles), &extraFiles); err == nil {
			for _, extraFile := range extraFiles {
				filePath := filepath.Join(tmpDir, extraFile.Path)

				// Create subdirectories if needed
				dir := filepath.Dir(filePath)
				if dir != "." && dir != tmpDir {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return zero, fmt.Errorf("create extra file dir: %w", err)
					}
				}

				if err := os.WriteFile(filePath, []byte(extraFile.Content), 0644); err != nil {
					return zero, fmt.Errorf("write extra file %s: %w", extraFile.Path, err)
				}

				mounts = append(mounts, Mount{
					Source:   filePath,
					Target:   fmt.Sprintf("/app/%s", extraFile.Path),
					ReadOnly: true,
				})
			}
		}
	}

	// 5. Prepare container configuration
	imageName := s.getImageName(executorID, exec.Version)

	// Parse dependencies for ADDITIONAL_PACKAGES
	additionalPackages := s.parseDependencies(ExecutorLanguage(exec.Language), exec.Dependencies)

	buildConfig := ImageConfig{
		Name:           fmt.Sprintf("executor-%s", executorID),
		DockerfilePath: s.dockerfiles[ExecutorLanguage(exec.Language)],
		BuildArgs: map[string]*string{
			"ADDITIONAL_PACKAGES": &additionalPackages,
		},
		Tags: []string{imageName},
	}

	// Parse schemas
	var schemaIn, schemaOut []byte
	if exec.SchemaIn != nil {
		schemaIn = []byte(*exec.SchemaIn)
	}
	if exec.SchemaOut != nil {
		schemaOut = []byte(*exec.SchemaOut)
	}

	// 6. Configure container
	containerConfig := ContainerConfig{
		Image:       imageName,
		AutoBuild:   true,
		BuildConfig: &buildConfig,
		ForceBuild:  false, // Use cached image if available
		Mounts:      mounts,
		Config: Config{
			Timeout:   time.Duration(exec.Timeout) * time.Second,
			SchemaIn:  schemaIn,
			SchemaOut: schemaOut,
		},
		MemoryBytes: int64(exec.MemoryMB) * 1024 * 1024,
		CPUShares:   int64(exec.CPUShares),
	}

	// Apply environment variables
	if exec.Env != nil {
		var envVars []string
		if err := json.Unmarshal([]byte(*exec.Env), &envVars); err == nil {
			containerConfig.Env = envVars
		}
	}

	// 7. Create and execute
	startTime := time.Now()
	containerExec, err := NewContainerExecutor[A, B](containerConfig, s.builder)
	if err != nil {
		return zero, fmt.Errorf("create container executor: %w", err)
	}

	output, err := containerExec.Execute(ctx, input)
	if err != nil {
		return zero, fmt.Errorf("execute: %w", err)
	}

	executionTime := time.Since(startTime)

	return ExecuteResult[B]{
		Output:          output,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Logs:            "", // TODO: Capture stderr logs
	}, nil
}

func (s *executorService[A, B]) ValidateInput(
	ctx context.Context,
	executorID string,
	input A,
) error {
	exec, err := s.GetExecutor(ctx, executorID)
	if err != nil {
		return fmt.Errorf("get executor: %w", err)
	}

	if exec.SchemaIn == nil {
		return nil // No schema to validate against
	}

	schemaBytes := []byte(*exec.SchemaIn)
	validator, err := NewSchemaValidator(schemaBytes, []byte("{}"))
	if err != nil {
		return fmt.Errorf("create validator: %w", err)
	}

	return validator.ValidateInput(input)
}

func (s *executorService[A, B]) ValidateOutput(
	ctx context.Context,
	executorID string,
	output B,
) error {
	exec, err := s.GetExecutor(ctx, executorID)
	if err != nil {
		return fmt.Errorf("get executor: %w", err)
	}

	if exec.SchemaOut == nil {
		return nil // No schema to validate against
	}

	schemaBytes := []byte(*exec.SchemaOut)
	validator, err := NewSchemaValidator([]byte("{}"), schemaBytes)
	if err != nil {
		return fmt.Errorf("create validator: %w", err)
	}

	return validator.ValidateOutput(output)
}

func (s *executorService[A, B]) BuildExecutor(ctx context.Context, executorID string) error {
	exec, err := s.GetExecutor(ctx, executorID)
	if err != nil {
		return fmt.Errorf("get executor: %w", err)
	}

	imageName := s.getImageName(executorID, exec.Version)
	additionalPackages := s.parseDependencies(ExecutorLanguage(exec.Language), exec.Dependencies)

	buildConfig := ImageConfig{
		Name:           fmt.Sprintf("executor-%s", executorID),
		DockerfilePath: s.dockerfiles[ExecutorLanguage(exec.Language)],
		BuildArgs: map[string]*string{
			"ADDITIONAL_PACKAGES": &additionalPackages,
		},
		Tags: []string{imageName},
	}

	result, err := s.builder.BuildImage(ctx, buildConfig)
	if err != nil {
		return fmt.Errorf("build image: %w", err)
	}

	if result.Error != nil {
		return fmt.Errorf("image build failed: %w", result.Error)
	}

	return nil
}

func (s *executorService[A, B]) GetExecutor(
	ctx context.Context,
	executorID string,
) (*models.Executor, error) {
	// Parse UUID
	id, err := uuid.Parse(executorID)
	if err != nil {
		return nil, fmt.Errorf("parse executor ID: %w", err)
	}

	exec, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get executor: %w", err)
	}

	return exec, nil
}

// Helper methods

func (s *executorService[A, B]) getExecuteFileName(language ExecutorLanguage) string {
	switch language {
	case ExecutorLanguageNodejs:
		return "execute.ts"
	case ExecutorLanguagePython:
		return "execute.py"
	case ExecutorLanguageGo:
		return "execute.go"
	default:
		return "execute"
	}
}

func (s *executorService[A, B]) getImageName(executorID string, version int32) string {
	return fmt.Sprintf("archesai/executor-%s:v%d", executorID, version)
}

func (s *executorService[A, B]) parseDependencies(
	language ExecutorLanguage,
	deps *string,
) string {
	if deps == nil {
		return ""
	}

	switch language {
	case ExecutorLanguageNodejs:
		// Parse package.json and extract package names
		var pkg map[string]any
		if err := json.Unmarshal([]byte(*deps), &pkg); err == nil {
			if dependencies, ok := pkg["dependencies"].(map[string]any); ok {
				packages := make([]string, 0, len(dependencies))
				for name := range dependencies {
					packages = append(packages, name)
				}
				return strings.Join(packages, " ")
			}
		}
		return ""
	case ExecutorLanguagePython:
		// requirements.txt is already in the right format (one per line)
		// Convert to space-separated for ADDITIONAL_PACKAGES
		lines := strings.Split(*deps, "\n")
		var packages []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				packages = append(packages, line)
			}
		}
		return strings.Join(packages, " ")
	case ExecutorLanguageGo:
		// Go modules are more complex, return empty for now
		return ""
	default:
		return ""
	}
}
