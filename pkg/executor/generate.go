//go:generate go run ../../cmd/archesai generate
package executor

import "embed"

// API embeds the OpenAPI specification files for the executor package.
//
//go:embed api
var API embed.FS
