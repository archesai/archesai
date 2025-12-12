//go:generate go run ../../cmd/archesai generate
package auth

import "embed"

// API embeds the OpenAPI specification files for the auth package.
//
//go:embed api
var API embed.FS
