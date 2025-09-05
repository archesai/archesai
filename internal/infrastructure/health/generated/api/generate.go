// Package api contains generated types and server interfaces for health infrastructure.
package api

// Generate types and server interfaces for health infrastructure (Health tag)
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../../oapi-codegen.yaml -o types.gen.go  -generate types,skip-prune       --include-tags=Health ../../../../../api/openapi.bundled.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../../oapi-codegen.yaml -o server.gen.go -generate echo-server,skip-prune --include-tags=Health ../../../../../api/openapi.bundled.yaml
