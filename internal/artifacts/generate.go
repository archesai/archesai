package artifacts

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package artifacts --include-tags Artifacts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package artifacts --include-tags Artifacts ../../api/openapi.bundled.yaml
