// Package invitations provides domain logic for invitation management
package invitations

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package invitations --include-tags Invitations ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package invitations --include-tags Invitations ../../api/openapi.bundled.yaml
