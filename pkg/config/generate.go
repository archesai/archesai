//go:generate go run ../../cmd/archesai generate
package config

import "embed"

// API embeds the OpenAPI specification files for the config package.
// This allows the config schemas to be merged into user specs during bundling.
//
//go:embed api
var API embed.FS
