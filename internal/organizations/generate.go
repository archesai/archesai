package organizations

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package organizations --include-tags Organizations,Members,Invitations ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package organizations --include-tags Organizations,Members,Invitations ../../api/openapi.bundled.yaml
