// Package accounts provides account management functionality.
package accounts

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package accounts --include-tags Accounts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package accounts --include-tags Accounts ../../api/openapi.bundled.yaml
