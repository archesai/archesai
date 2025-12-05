package generators

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

// GoModTemplateData holds the data for rendering the go.mod template.
type GoModTemplateData struct {
	ProjectName string
}

// GoModGenerator generates the go.mod file for new projects.
type GoModGenerator struct{}

// Name returns the generator name.
func (g *GoModGenerator) Name() string { return "go.mod" }

// Priority returns the generator priority.
func (g *GoModGenerator) Priority() int { return PriorityFirst }

// Generate creates the go.mod file if it doesn't exist.
func (g *GoModGenerator) Generate(ctx *GeneratorContext) error {
	if ctx.ProjectName == "" {
		return nil
	}

	data := &GoModTemplateData{ProjectName: ctx.ProjectName}

	if err := ctx.RenderToFileIfNotExists("go.mod.tmpl", "go.mod", data); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	// slog.Info("generated go.mod", "project", ctx.ProjectName) FIXME: fix logging consistency
	return nil
}

// RunGoModTidy runs go mod tidy in the output directory.
func RunGoModTidy(baseDir string) error {
	goModPath := filepath.Join(baseDir, "go.mod")

	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = baseDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("running go mod tidy", "dir", baseDir)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}
