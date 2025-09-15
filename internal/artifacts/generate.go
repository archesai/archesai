package artifacts

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package artifacts --include-tags Artifacts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package artifacts --include-tags Artifacts ../../api/openapi.bundled.yaml
