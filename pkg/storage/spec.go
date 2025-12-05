//go:generate go run ../../cmd/archesai generate --spec ./spec/openapi.yaml --output . --only models,controllers,application,repositories,routes,bootstrap_handlers --pretty
package storage

import "embed"

// APISpec embeds the OpenAPI specification files for the storage package.
//
//go:embed spec
var APISpec embed.FS
