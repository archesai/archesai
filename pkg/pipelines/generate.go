//go:generate go run ../../cmd/archesai generate
package pipelines

import "embed"

// API embeds the OpenAPI specification files for the pipelines package.
//
//go:embed api
var API embed.FS
