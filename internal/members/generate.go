// Package members provides domain logic for member management
package members

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package members --include-tags Members ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package members --include-tags Members ../../api/openapi.bundled.yaml
