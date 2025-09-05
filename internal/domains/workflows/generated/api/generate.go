// Package api contains generated types and server interfaces for the workflows domain.
package api

// Generate types and server interfaces for workflows domain (Workflows tag)
//go:generate go tool oapi-codegen --config=../../../../../oapi-codegen.yaml -o types.gen.go  -generate types,skip-prune       --include-tags=Workflows ../../../../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../../../../oapi-codegen.yaml -o server.gen.go -generate echo-server,skip-prune --include-tags=Workflows ../../../../../api/openapi.bundled.yaml
