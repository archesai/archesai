package codegen

import (
	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/auth"
	"github.com/archesai/archesai/pkg/config"
	"github.com/archesai/archesai/pkg/executor"
	"github.com/archesai/archesai/pkg/pipelines"
	"github.com/archesai/archesai/pkg/server"
	"github.com/archesai/archesai/pkg/storage"
)

// NewDefaultIncludeMerger creates an IncludeMerger with all standard includes registered.
// This is defined in codegen (not parsers) to avoid circular dependencies between
// internal/parsers and pkg/* packages that use go:generate to run archesai.
func NewDefaultIncludeMerger() *parsers.IncludeMerger {
	merger := parsers.NewIncludeMerger()
	merger.RegisterInclude("auth", auth.APISpec)
	merger.RegisterInclude("config", config.APISpec)
	merger.RegisterInclude("server", server.APISpec)
	merger.RegisterInclude("storage", storage.APISpec)
	merger.RegisterInclude("pipelines", pipelines.APISpec)
	merger.RegisterInclude("executor", executor.APISpec)
	return merger
}
