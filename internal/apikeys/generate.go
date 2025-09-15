// Package apikeys provides API key management functionality
package apikeys

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package apikeys --include-tags APIKeys ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package apikeys --include-tags APIKeys ../../api/openapi.bundled.yaml
