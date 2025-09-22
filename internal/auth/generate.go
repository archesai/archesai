package auth

// Sessions API generation
//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package auth --include-tags Sessions -o sessions_types.gen.go ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package auth --include-tags Sessions -o sessions_api.gen.go ../../api/openapi.bundled.yaml

// Tokens API generation
//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package auth --include-tags Tokens -o tokens_types.gen.go ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package auth --include-tags Tokens -o tokens_api.gen.go ../../api/openapi.bundled.yaml
