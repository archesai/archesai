package runs

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package runs --include-tags Runs ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package runs --include-tags Runs ../../api/openapi.bundled.yaml
