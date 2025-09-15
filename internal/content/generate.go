package content

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package content --include-tags Content,Artifacts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package content --include-tags Content,Artifacts ../../api/openapi.bundled.yaml
