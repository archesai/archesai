package labels

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package labels --include-tags Labels ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package labels --include-tags Labels ../../api/openapi.bundled.yaml
