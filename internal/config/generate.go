package config

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package config --include-tags Config ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package config --include-tags Config ../../api/openapi.bundled.yaml
