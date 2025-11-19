package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"sync"
	"time"
)

// ContainerExecutor runs code in Docker containers with JSON Schema validation
type ContainerExecutor[A any, B any] struct {
	config    ContainerConfig
	validator *SchemaValidator
	builder   *Builder
}

// Mount represents a volume mount for the container
type Mount struct {
	Source   string // Host path
	Target   string // Container path
	ReadOnly bool   // Mount as read-only
}

// ContainerConfig holds configuration for a container executor
type ContainerConfig = struct {
	Config

	// Container configuration
	Image   string   // Container image to use
	Cmd     []string // Optional command override
	WorkDir string   // Working directory in container
	Env     []string // Environment variables
	Mounts  []Mount  // Volume mounts

	// Resource limits
	CPUShares   int64 // CPU shares (relative weight)
	MemoryBytes int64 // Memory limit in bytes

	// Security settings
	DisableNet bool   // Disable network access
	ReadOnlyFS bool   // Mount filesystem as read-only
	User       string // User to run as (e.g., "1000:1000")

	// Auto-build settings
	AutoBuild   bool         // Enable automatic image building
	BuildConfig *ImageConfig // Build configuration (required if AutoBuild is true)
	ForceBuild  bool         // Force rebuild even if image exists
}

// ContainerRequest is the JSON structure sent to the container via stdin
type ContainerRequest struct {
	SchemaIn  json.RawMessage `json:"schema_in"`  // JSON Schema for input validation
	SchemaOut json.RawMessage `json:"schema_out"` // JSON Schema for output validation
	Input     json.RawMessage `json:"input"`      // The actual input data
}

// ContainerResponse is the JSON structure returned by the container via stdout
type ContainerResponse struct {
	OK     bool            `json:"ok"`               // Whether execution was successful
	Output json.RawMessage `json:"output,omitempty"` // The output data (if ok=true)
	Error  *ContainerError `json:"error,omitempty"`  // Error details (if ok=false)
}

// ContainerError represents an error from the container
type ContainerError struct {
	Message string `json:"message"`           // Error message
	Details string `json:"details,omitempty"` // Optional additional details
}

// NewContainerExecutor creates a new container-based executor with optional auto-build support
func NewContainerExecutor[A any, B any](
	config ContainerConfig,
	builder *Builder,
) (*ContainerExecutor[A, B], error) {
	// Validate configuration
	if config.Image == "" {
		return nil, fmt.Errorf("container image is required")
	}

	// Validate auto-build configuration
	if config.AutoBuild {
		if builder == nil {
			return nil, fmt.Errorf("builder is required when AutoBuild is enabled")
		}
		if config.BuildConfig == nil {
			return nil, fmt.Errorf("BuildConfig is required when AutoBuild is enabled")
		}
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MemoryBytes == 0 {
		config.MemoryBytes = 256 * 1024 * 1024 // 256MB default
	}
	if config.CPUShares == 0 {
		config.CPUShares = 512 // Medium priority
	}

	// Create schema validator if schemas are provided
	var validator *SchemaValidator
	if len(config.SchemaIn) > 0 && len(config.SchemaOut) > 0 {
		var err error
		validator, err = NewSchemaValidator(config.SchemaIn, config.SchemaOut)
		if err != nil {
			return nil, fmt.Errorf("create schema validator: %w", err)
		}
	}

	return &ContainerExecutor[A, B]{
		config:    config,
		validator: validator,
		builder:   builder,
	}, nil
}

// Execute runs the container with the given input and returns the output
func (e *ContainerExecutor[A, B]) Execute(ctx context.Context, input A) (B, error) {
	var zero B

	// Auto-build image if configured
	if e.config.AutoBuild && e.builder != nil {
		slog.Debug("Auto-build enabled for container executor",
			"image", e.config.Image,
			"force_build", e.config.ForceBuild)

		// Check if image exists (unless force build)
		shouldBuild := e.config.ForceBuild
		if !shouldBuild {
			exists, err := e.builder.ImageExists(ctx, e.config.Image)
			if err != nil {
				slog.Error("Failed to check image existence",
					"image", e.config.Image,
					"error", err)
				return zero, fmt.Errorf("check image existence: %w", err)
			}
			shouldBuild = !exists
		}

		// Build image if needed
		if shouldBuild {
			result, err := e.builder.BuildImage(ctx, *e.config.BuildConfig)
			if err != nil {
				slog.Error("Failed to build container image",
					"image", e.config.Image,
					"error", err)
				return zero, fmt.Errorf("build image: %w", err)
			}
			if result.Error != nil {
				slog.Error("Container image build failed",
					"image", e.config.Image,
					"error", result.Error,
					"output", result.Output)
				return zero, fmt.Errorf("image build failed: %w", result.Error)
			}
			slog.Debug("Container image built successfully", "image", e.config.Image)
		} else {
			slog.Debug("Using existing container image", "image", e.config.Image)
		}
	}

	slog.Debug("Executing container",
		"image", e.config.Image,
		"timeout", e.config.Timeout)

	// Marshal input to JSON
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return zero, fmt.Errorf("marshal input: %w", err)
	}

	// Validate input against schema if validator is configured
	if e.validator != nil {
		var inputAny any
		if err := json.Unmarshal(inputBytes, &inputAny); err != nil {
			return zero, fmt.Errorf("unmarshal input for validation: %w", err)
		}
		if err := e.validator.ValidateInput(inputAny); err != nil {
			return zero, fmt.Errorf("input validation failed: %w", err)
		}
	}

	// Prepare container request
	request := ContainerRequest{
		Input: json.RawMessage(inputBytes),
	}

	// Include schemas if validator is configured
	if e.validator != nil {
		request.SchemaIn = json.RawMessage(e.validator.GetInputSchema())
		request.SchemaOut = json.RawMessage(e.validator.GetOutputSchema())
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return zero, fmt.Errorf("marshal container request: %w", err)
	}

	// Run container
	output, err := e.runContainer(ctx, requestBytes)
	if err != nil {
		return zero, err
	}

	// Parse container response
	var response ContainerResponse
	if err := json.Unmarshal(output, &response); err != nil {
		// Log partial output for debugging (first and last parts)
		outputPreview := string(output)
		if len(outputPreview) > 500 {
			outputPreview = fmt.Sprintf("%s... [truncated %d bytes] ...%s",
				outputPreview[:200],
				len(output)-400,
				outputPreview[len(outputPreview)-200:])
		}
		slog.Error("Failed to parse container response",
			"image", e.config.Image,
			"error", err,
			"output_length", len(output),
			"output_preview", outputPreview,
		)
		return zero, fmt.Errorf("parse container response: %w", err)
	}

	// Check if container reported an error
	if !response.OK {
		if response.Error != nil {
			return zero, fmt.Errorf("container error: %s (details=%s)",
				response.Error.Message, response.Error.Details)
		}
		return zero, fmt.Errorf("container execution failed with ok=false")
	}

	// Validate output against schema if validator is configured
	if e.validator != nil {
		var outputAny any
		if err := json.Unmarshal(response.Output, &outputAny); err != nil {
			return zero, fmt.Errorf("unmarshal output for validation: %w", err)
		}
		if err := e.validator.ValidateOutput(outputAny); err != nil {
			return zero, fmt.Errorf("output validation failed: %w", err)
		}
	}

	// Unmarshal output to type B
	var result B
	if err := json.Unmarshal(response.Output, &result); err != nil {
		return zero, fmt.Errorf("unmarshal output: %w", err)
	}

	return result, nil
}

// runContainer executes the Docker container and returns the output
func (e *ContainerExecutor[A, B]) runContainer(ctx context.Context, input []byte) ([]byte, error) {
	// Build docker run command
	args := []string{"run", "--rm", "-i"}

	// Security settings
	if e.config.DisableNet {
		args = append(args, "--network", "none")
	}
	if e.config.ReadOnlyFS {
		args = append(args, "--read-only")
	}
	if e.config.User != "" {
		args = append(args, "--user", e.config.User)
	}

	// Resource limits
	if e.config.MemoryBytes > 0 {
		args = append(args, "--memory", fmt.Sprintf("%d", e.config.MemoryBytes))
	}
	if e.config.CPUShares > 0 {
		args = append(args, "--cpu-shares", fmt.Sprintf("%d", e.config.CPUShares))
	}

	// Working directory
	if e.config.WorkDir != "" {
		args = append(args, "-w", e.config.WorkDir)
	}

	// Environment variables
	for _, env := range e.config.Env {
		args = append(args, "-e", env)
	}

	// Volume mounts
	for _, mount := range e.config.Mounts {
		mountStr := fmt.Sprintf("%s:%s", mount.Source, mount.Target)
		if mount.ReadOnly {
			mountStr += ":ro"
		}
		args = append(args, "-v", mountStr)
	}

	// Container image
	args = append(args, e.config.Image)

	// Optional command override
	if len(e.config.Cmd) > 0 {
		args = append(args, e.config.Cmd...)
	}

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	// Execute Docker command
	cmd := exec.CommandContext(execCtx, "docker", args...)
	cmd.Stdin = bytes.NewReader(input)

	// Get pipes for stdout/stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("create stderr pipe: %w", err)
	}

	// Start the command
	slog.Debug("Running Docker container command",
		"image", e.config.Image,
		"args", args)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	// Read stdout and stderr concurrently
	// Critical: We must read the pipes BEFORE calling Wait() to avoid deadlock
	var stdoutBytes, stderrBytes []byte
	var stdoutErr, stderrErr error
	var wg sync.WaitGroup
	wg.Add(2)

	// Read all stdout data
	go func() {
		defer wg.Done()
		// Use io.ReadAll which reads until EOF
		stdoutBytes, stdoutErr = io.ReadAll(stdoutPipe)
	}()

	// Read all stderr data
	go func() {
		defer wg.Done()
		// Use io.ReadAll which reads until EOF
		stderrBytes, stderrErr = io.ReadAll(stderrPipe)
	}()

	// Wait for goroutines to finish reading
	wg.Wait()

	// Now wait for the command to exit
	cmdErr := cmd.Wait()

	// Check for read errors
	if stdoutErr != nil {
		slog.Warn("Error reading stdout", "error", stdoutErr)
	}
	if stderrErr != nil {
		slog.Warn("Error reading stderr", "error", stderrErr)
	}

	// Convert stderr to string for logging
	stderrStr := string(stderrBytes)

	// Log output sizes
	if len(stdoutBytes) > 32768 { // Log if output is > 32KB
		slog.Debug("Large container output received",
			"image", e.config.Image,
			"stdout_size", len(stdoutBytes),
			"stderr_size", len(stderrBytes))
	} else {
		slog.Debug("Container output received",
			"image", e.config.Image,
			"stdout_size", len(stdoutBytes),
			"stderr_size", len(stderrBytes))
	}

	// Log stderr output (container logs) if present
	if stderrStr != "" {
		slog.Debug("Container stderr output",
			"image", e.config.Image,
			"stderr", stderrStr)
	}

	// Handle command errors
	if cmdErr != nil {
		// Check if it was a timeout
		if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
			slog.Error("Container execution timed out",
				"image", e.config.Image,
				"timeout", e.config.Timeout,
				"stderr", stderrStr)
			return nil, fmt.Errorf("container execution timed out after %s", e.config.Timeout)
		}
		// Include stderr in error message
		slog.Error("Container execution failed",
			"image", e.config.Image,
			"error", cmdErr,
			"stdout_size", len(stdoutBytes),
			"stderr", stderrStr)
		return nil, fmt.Errorf("container execution failed: %v; stderr: %s", cmdErr, stderrStr)
	}

	slog.Debug("Container execution completed successfully",
		"image", e.config.Image,
		"stdout_size", len(stdoutBytes))
	return stdoutBytes, nil
}
