package members

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package members --include-tags Members ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package members --include-tags Members ../../api/openapi.bundled.yaml
