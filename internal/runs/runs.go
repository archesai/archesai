// Package runs provides the domain logic for managing runs of pipelines and tools.
package runs

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package runs --include-tags Runs ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package runs --include-tags Runs ../../api/openapi.bundled.yaml
