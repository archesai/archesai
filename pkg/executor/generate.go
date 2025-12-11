//go:generate go run ../../cmd/archesai generate --spec ./api/openapi.yaml --output . --only models,routes,handlers,repositories,bootstrap_handlers,bootstrap_routes --pretty
package executor

import "embed"

// APISpec embeds the OpenAPI specification files for the executor package.
//
//go:embed api
var APISpec embed.FS
