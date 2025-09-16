package pipelines

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package pipelines --include-tags Pipelines ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package pipelines --include-tags Pipelines ../../api/openapi.bundled.yaml
