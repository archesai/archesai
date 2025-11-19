package codegen

// GenerateSQLC generates SQLC queries from SQL files

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	sqlc "github.com/sqlc-dev/sqlc/pkg/cli"
)

// GenerateSQLC generates SQLC queries from SQL files
func (g *Generator) GenerateSQLC() error {

	// Second, generate the sqlc.yaml file with the correct output paths
	sqlcConfigPath := filepath.Join(g.outputDir, "sqlc.yaml")

	// Load the sqlc.yaml template
	tmplContent, err := GetTemplate("sqlc.yaml.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load sqlc.yaml template: %w", err)
	}

	// Parse and execute the template
	tmpl, err := template.New("sqlc.yaml").Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse sqlc.yaml template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{
		"OutputDir": g.outputDir,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute sqlc.yaml template: %w", err)
	}

	// Write the sqlc.yaml file
	if err := os.WriteFile(sqlcConfigPath, buf.Bytes(), 0644); err != nil {
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
