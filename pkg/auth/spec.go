//go:generate go run ../../cmd/archesai/main.go generate --spec ./spec/openapi.yaml --output . --only bundle,models,controllers,application,repositories,routes,bootstrap_handlers --pretty
package auth

import "embed"

// APISpec embeds the OpenAPI specification files for the auth package.
//
//go:embed spec
var APISpec embed.FS
