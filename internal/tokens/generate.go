// Package tokens provides API token management functionality for programmatic access.
package tokens

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package tokens --include-tags Tokens ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package tokens --include-tags Tokens ../../api/openapi.bundled.yaml
