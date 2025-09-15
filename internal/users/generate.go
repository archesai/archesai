package users

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package users --include-tags Users ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package users --include-tags Users ../../api/openapi.bundled.yaml
