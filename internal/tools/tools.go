// Package tools provides tool management functionality.
package tools

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package tools --include-tags Tools ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package tools --include-tags Tools ../../api/openapi.bundled.yaml
