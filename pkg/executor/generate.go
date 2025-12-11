//go:generate go run ../../cmd/archesai generate --spec ./api/openapi.yaml --output ./gen --only models,controllers,application,repositories,routes,bootstrap_handlers --pretty
package executor

import "embed"

// APISpec embeds the OpenAPI specification files for the executor package.
//
//go:embed api
var APISpec embed.FS
