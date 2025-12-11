//go:generate go run ../../cmd/archesai generate
package storage

import "embed"

// API embeds the OpenAPI specification files for the storage package.
//
//go:embed api
var API embed.FS
