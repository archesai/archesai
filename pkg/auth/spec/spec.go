// Package spec embeds the OpenAPI specification files for the auth package.
package spec

import "embed"

// FS embeds the OpenAPI specification files.
//
//go:embed *
var FS embed.FS
