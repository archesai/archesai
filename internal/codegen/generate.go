package codegen

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package codegen --generate skip-prune,models ../../api/components/schemas/XCodegenWrapper.yaml
