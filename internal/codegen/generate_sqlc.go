package codegen

// GenerateSQLC generates SQLC queries from SQL files

import (
	"fmt"

	sqlc "github.com/sqlc-dev/sqlc/pkg/cli"
)

// GenerateSQLC generates SQLC queries from SQL files
func (g *Generator) GenerateSQLC() error {

	code := sqlc.Run(
		[]string{"generate", "--file", "internal/infrastructure/persistence/sqlc.yaml"},
	)
	if code != 0 {
		return fmt.Errorf("sqlc generation failed with code %d", code)
	}

	return nil
}
