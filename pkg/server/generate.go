//go:generate go run ../../cmd/archesai generate
package server

import "embed"

// API embeds the OpenAPI specification files for the server package.
//
//go:embed api
var API embed.FS
