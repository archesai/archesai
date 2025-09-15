// Package oauth provides OAuth authentication functionality
package oauth

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package oauth --include-tags OAuth ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package oauth --include-tags OAuth ../../api/openapi.bundled.yaml
