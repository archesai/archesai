package auth

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package auth --include-tags Auth,Sessions,Accounts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package auth --include-tags Auth,Sessions,Accounts ../../api/openapi.bundled.yaml
