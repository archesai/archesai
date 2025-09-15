package config

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package config --include-tags Config ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package config --include-tags Config ../../api/openapi.bundled.yaml
