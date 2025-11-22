package codegen

// GenerateSQLC generates SQLC queries from SQL files

import (
	"bytes"
	"fmt"
	"path/filepath"

	sqlc "github.com/sqlc-dev/sqlc/pkg/cli"
)

// GenerateSQLC generates SQLC queries from SQL files
func (g *Generator) GenerateSQLC() error {

	// Generate the sqlc.yaml file with the correct output paths
	sqlcConfigPath := filepath.Join(
		g.outputDir,
		"generated",
		"infrastructure",
		"persistence",
		"sqlc.gen.yaml",
	)

	// Generate sqlc.yaml
	data := map[string]string{
		"OutputDir": g.outputDir,
	}

	var buf bytes.Buffer
	if err := g.renderer.Render(&buf, "sqlc.yaml.tmpl", data); err != nil {
		return fmt.Errorf("failed to render sqlc.yaml: %w", err)
	}

	if err := g.storage.WriteFile(sqlcConfigPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write sqlc.yaml: %w", err)
	}

	// Now run sqlc with the generated config file
	code := sqlc.Run(
		[]string{"generate", "--file", sqlcConfigPath},
	)
	if code != 0 {
		return fmt.Errorf("sqlc generation failed with code %d", code)
	}

	return nil
}
