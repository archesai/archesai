package spec

import (
	"io/fs"

	authspec "github.com/archesai/archesai/pkg/auth/spec"
	configspec "github.com/archesai/archesai/pkg/config/spec"
	executorspec "github.com/archesai/archesai/pkg/executor/spec"
	pipelinesspec "github.com/archesai/archesai/pkg/pipelines/spec"
	serverspec "github.com/archesai/archesai/pkg/server/spec"
	"github.com/archesai/archesai/pkg/storage"
	storagespec "github.com/archesai/archesai/pkg/storage/spec"
)

// IncludeRegistry holds registered package filesystems for includes.
var IncludeRegistry = map[string]fs.FS{
	"server":    serverspec.FS,
	"auth":      authspec.FS,
	"config":    configspec.FS,
	"pipelines": pipelinesspec.FS,
	"executor":  executorspec.FS,
	"storage":   storagespec.FS,
}

// BuildIncludeFS creates a composite filesystem with include filesystems as base layers.
// The project filesystem is placed on top so its files override include files.
func BuildIncludeFS(projectFS fs.FS, includeNames []string) fs.FS {
	if len(includeNames) == 0 {
		return projectFS
	}

	var layers []fs.FS
	for _, name := range includeNames {
		if includeFS := IncludeRegistry[name]; includeFS != nil {
			layers = append(layers, includeFS)
		}
	}

	// Project FS is the top layer (overrides includes)
	layers = append(layers, projectFS)
	return storage.NewComposite(layers...)
}
