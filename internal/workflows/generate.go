package workflows

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package workflows --include-tags Workflows,Pipelines,Runs,Tools ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package workflows --include-tags Workflows,Pipelines,Runs,Tools ../../api/openapi.bundled.yaml
