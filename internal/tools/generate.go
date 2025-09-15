package tools

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package tools --include-tags Tools ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package tools --include-tags Tools ../../api/openapi.bundled.yaml
