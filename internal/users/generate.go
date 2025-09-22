package users

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package users --include-tags Users ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package users --include-tags Users ../../api/openapi.bundled.yaml
