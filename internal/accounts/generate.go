// Package accounts provides account management functionality.
package accounts

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package accounts --include-tags Accounts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package accounts --include-tags Accounts ../../api/openapi.bundled.yaml
