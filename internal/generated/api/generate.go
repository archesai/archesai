// Generate all types and server interfaces from bundled OpenAPI spec
//go:generate go tool oapi-codegen -config oapi-codegen.yaml -package api -o types.gen.go  -generate types,skip-prune                ../../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen -config oapi-codegen.yaml -package api -o server.gen.go -generate server,strict-server,skip-prune ../../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen -config oapi-codegen.yaml -package api -o spec.gen.go   -generate spec                            ../../../api/openapi.bundled.yaml

// Package api contains generated OpenAPI client and server code.
package api
