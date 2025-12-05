//go:generate go run ../../cmd/archesai generate --spec ./spec/openapi.yaml --output . --only models,controllers,application,repositories,routes,bootstrap_handlers --pretty
package config

import "embed"

// APISpec embeds the OpenAPI specification files for the config package.
// This allows the config schemas to be merged into user specs during bundling.
//
//go:embed spec
var APISpec embed.FS
