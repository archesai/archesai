package pipelines

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package pipelines --include-tags Pipelines ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package pipelines --include-tags Pipelines ../../api/openapi.bundled.yaml
