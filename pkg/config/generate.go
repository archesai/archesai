//go:generate go run ../../cmd/archesai generate --spec ./api/openapi.yaml --output . --only models,routes,handlers,repositories,bootstrap_handlers,bootstrap_routes --pretty
package config

import "embed"

// APISpec embeds the OpenAPI specification files for the config package.
// This allows the config schemas to be merged into user specs during bundling.
//
//go:embed api
var APISpec embed.FS
