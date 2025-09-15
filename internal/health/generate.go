package health

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package health --include-tags Health ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package health --include-tags Health ../../api/openapi.bundled.yaml
