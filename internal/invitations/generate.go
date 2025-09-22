package invitations

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package invitations --include-tags Invitations ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package invitations --include-tags Invitations ../../api/openapi.bundled.yaml
