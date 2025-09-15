package sessions

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package sessions --include-tags Sessions ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package sessions --include-tags Sessions ../../api/openapi.bundled.yaml
