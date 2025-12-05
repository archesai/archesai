//go:generate go run ../../cmd/archesai generate --spec ./spec/openapi.yaml --output ./gen --only models,controllers,application,repositories,routes,bootstrap_handlers --pretty
package server

import "embed"

// APISpec embeds the OpenAPI specification files for the server package.
//
//go:embed spec
var APISpec embed.FS
