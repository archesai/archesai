package organizations

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package organizations --include-tags Organizations ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package organizations --include-tags Organizations ../../api/openapi.bundled.yaml
