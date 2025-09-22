package health

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package health --include-tags Health ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package health --include-tags Health ../../api/openapi.bundled.yaml
