// Package sessions provides session management functionality
package sessions

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package sessions --include-tags Sessions ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package sessions --include-tags Sessions ../../api/openapi.bundled.yaml
