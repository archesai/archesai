package executor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"sync"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// ImageConfig defines configuration for building a Docker image
type ImageConfig struct {
	Name           string             // Human-readable name for logging
	DockerfilePath string             // Path to Dockerfile
	Tags           []string           // Image tags
	BuildArgs      map[string]*string // Build arguments
	Target         string             // Build target (for multi-stage builds)
	NoCache        bool               // Disable build cache
}

// Builder handles building Docker images
type Builder struct {
	cli *client.Client
}

// BuildResult contains the result of a build operation
type BuildResult struct {
	Name   string
	Tags   []string
	Error  error
	Output string
}

// BuildMessage represents a message from Docker build output
type BuildMessage struct {
	Stream      string `json:"stream,omitempty"`
	Error       string `json:"error,omitempty"`
	ErrorDetail struct {
		Message string `json:"message,omitempty"`
	} `json:"errorDetail,omitempty"`
}

// NewBuilder creates a new Docker image builder
func NewBuilder() (*Builder, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}

	return &Builder{cli: cli}, nil
}

// Close closes the Docker client connection
func (b *Builder) Close() error {
	if b.cli != nil {
		return b.cli.Close()
	}
	return nil
}

// BuildImage builds a single Docker image
func (b *Builder) BuildImage(ctx context.Context, cfg ImageConfig) (*BuildResult, error) {
	slog.Debug("Starting Docker image build",
		"name", cfg.Name,
		"tags", cfg.Tags,
		"dockerfile", cfg.DockerfilePath)

	result := &BuildResult{
		Name: cfg.Name,
		Tags: cfg.Tags,
	}

	// Validate configuration
	if err := validateImageConfig(cfg); err != nil {
		result.Error = err
		slog.Error("Image configuration validation failed",
			"name", cfg.Name,
			"error", err)
		return result, err
	}

	// Create build context
	slog.Debug("Creating build context",
		"name", cfg.Name,
		"context_dir", filepath.Dir(cfg.DockerfilePath))
	//nolint:staticcheck // archive.TarWithOptions is the correct function to use
	buildContext, err := archive.TarWithOptions(
		filepath.Dir(cfg.DockerfilePath),
		&archive.TarOptions{}, //nolint:staticcheck // TarOptions is the correct type
	)
	if err != nil {
		result.Error = fmt.Errorf("create build context: %w", err)
		slog.Error("Failed to create build context",
			"name", cfg.Name,
			"error", err)
		return result, result.Error
	}
	defer func() {
		if closeErr := buildContext.Close(); closeErr != nil {
			slog.Warn("Failed to close build context",
				"name", cfg.Name,
				"error", closeErr)
		}
	}()

	// Prepare build options
	//nolint:staticcheck // types.ImageBuildOptions is the correct type for Docker client API
	buildOpts := types.ImageBuildOptions{
		Dockerfile:  filepath.Base(cfg.DockerfilePath),
		Tags:        cfg.Tags,
		BuildArgs:   cfg.BuildArgs,
		Target:      cfg.Target,
		Remove:      true,
		ForceRemove: true,
		NoCache:     cfg.NoCache,
	}

	slog.Debug("Build options configured",
		"name", cfg.Name,
		"no_cache", cfg.NoCache,
		"target", cfg.Target)

	// Build the image
	buildResp, err := b.cli.ImageBuild(ctx, buildContext, buildOpts)
	if err != nil {
		result.Error = fmt.Errorf("build image: %w", err)
		slog.Error("Docker image build failed",
			"name", cfg.Name,
			"error", err)
		return result, result.Error
	}
	defer func() {
		if closeErr := buildResp.Body.Close(); closeErr != nil {
			slog.Warn("Failed to close build response",
				"name", cfg.Name,
				"error", closeErr)
		}
	}()

	// Parse build output
	slog.Debug("Parsing build output", "name", cfg.Name)
	output, err := parseBuildOutput(buildResp.Body)
	result.Output = output
	if err != nil {
		result.Error = fmt.Errorf("build failed: %w", err)
		slog.Error("Docker image build failed",
			"name", cfg.Name,
			"error", err,
			"output", output)
		return result, result.Error
	}

	slog.Debug("Docker image built successfully",
		"name", cfg.Name,
		"tags", cfg.Tags)
	return result, nil
}

// BuildImages builds multiple Docker images in parallel
func (b *Builder) BuildImages(ctx context.Context, configs []ImageConfig) []*BuildResult {
	results := make([]*BuildResult, len(configs))
	var wg sync.WaitGroup

	for i, cfg := range configs {
		wg.Add(1)
		go func(idx int, config ImageConfig) {
			defer wg.Done()
			results[idx], _ = b.BuildImage(ctx, config)
		}(i, cfg)
	}

	wg.Wait()
	return results
}

// validateImageConfig validates the image configuration
func validateImageConfig(cfg ImageConfig) error {
	if cfg.DockerfilePath == "" {
		return fmt.Errorf("dockerfile path is required")
	}
	if len(cfg.Tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// parseBuildOutput parses Docker build output and returns any errors
func parseBuildOutput(reader io.Reader) (string, error) {
	var output string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Bytes()
		var msg BuildMessage

		if err := json.Unmarshal(line, &msg); err != nil {
			// If we can't parse JSON, just append the raw line
			output += string(line) + "\n"
			continue
		}

		// Check for error
		if msg.Error != "" {
			return output, fmt.Errorf("%s", msg.Error)
		}
		if msg.ErrorDetail.Message != "" {
			return output, fmt.Errorf("%s", msg.ErrorDetail.Message)
		}

		// Append stream output
		if msg.Stream != "" {
			output += msg.Stream
		}
	}

	if err := scanner.Err(); err != nil {
		return output, fmt.Errorf("read build output: %w", err)
	}

	return output, nil
}

// ImageExists checks if a Docker image exists locally
func (b *Builder) ImageExists(ctx context.Context, tag string) (bool, error) {
	slog.Debug("Checking if Docker image exists", "tag", tag)
	_, err := b.cli.ImageInspect(ctx, tag)
	if err != nil {
		if cerrdefs.IsNotFound(err) {
			slog.Debug("Docker image not found locally", "tag", tag)
			return false, nil
		}
		slog.Error("Failed to inspect Docker image",
			"tag", tag,
			"error", err)
		return false, fmt.Errorf("inspect image %s: %w", tag, err)
	}
	slog.Debug("Docker image exists locally", "tag", tag)
	return true, nil
}

// RunnerConfigs returns the default configuration for all runner images
func RunnerConfigs(baseDir string) []ImageConfig {
	return []ImageConfig{
		{
			Name:           "go-runner",
			DockerfilePath: filepath.Join(baseDir, "deployments/containers/runners/go/Dockerfile"),
			Tags:           []string{"archesai/runner-go:latest"},
		},
		{
			Name: "node-runner",
			DockerfilePath: filepath.Join(
				baseDir,
				"deployments/containers/runners/node/Dockerfile",
			),
			Tags: []string{"archesai/runner-node:latest"},
		},
		{
			Name: "python-runner",
			DockerfilePath: filepath.Join(
				baseDir,
				"deployments/containers/runners/python/Dockerfile",
			),
			Tags: []string{"archesai/runner-python:latest"},
		},
	}
}
