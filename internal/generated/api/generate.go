//go:generate oapi-codegen -config oapi-codegen.yaml -package common -o common/responses.gen.go   -generate types,skip-prune ../../../api/specifications/common/responses/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package common -o common/filters.gen.go  -generate types,skip-prune ../../../api/specifications/common/filters/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package common -o common/security.gen.go -generate types,skip-prune ../../../api/specifications/common/security/schemas.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package config -o config/schemas.gen.go -generate types,skip-prune                 ../../../api/specifications/admin/config/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package config -o config/paths.gen.go   -generate server,strict-server,skip-prune,types  ../../../api/specifications/admin/config/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package health -o health/schemas.gen.go -generate types,skip-prune		           ../../../api/specifications/admin/health/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package health -o health/paths.gen.go   -generate server,strict-server,skip-prune,types  ../../../api/specifications/admin/health/paths.yaml

// Auth features
//go:generate oapi-codegen -config oapi-codegen.yaml -package accounts         -o auth/accounts/schemas.gen.go         -generate types,skip-prune                 ../../../api/specifications/auth/accounts/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package accounts         -o auth/accounts/parameters.gen.go      -generate types,skip-prune                 ../../../api/specifications/auth/accounts/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package accounts         -o auth/accounts/paths.gen.go           -generate server,strict-server,skip-prune,types  ../../../api/specifications/auth/accounts/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package apitokens        -o auth/api-tokens/schemas.gen.go       -generate types,skip-prune                 ../../../api/specifications/auth/api-tokens/schemas.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package invitations      -o auth/invitations/schemas.gen.go      -generate types,skip-prune                 ../../../api/specifications/auth/invitations/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package invitations      -o auth/invitations/parameters.gen.go   -generate types,skip-prune                 ../../../api/specifications/auth/invitations/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package invitations      -o auth/invitations/paths.gen.go        -generate server,strict-server,skip-prune,types  ../../../api/specifications/auth/invitations/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package members          -o auth/members/schemas.gen.go          -generate types,skip-prune                 ../../../api/specifications/auth/members/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package members          -o auth/members/parameters.gen.go       -generate types,skip-prune                 ../../../api/specifications/auth/members/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package members          -o auth/members/paths.gen.go            -generate server,strict-server,skip-prune,types  ../../../api/specifications/auth/members/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package organizations    -o auth/organizations/schemas.gen.go    -generate types,skip-prune                 ../../../api/specifications/auth/organizations/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package organizations    -o auth/organizations/parameters.gen.go -generate types,skip-prune                 ../../../api/specifications/auth/organizations/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package organizations    -o auth/organizations/paths.gen.go      -generate server,strict-server,skip-prune,types  ../../../api/specifications/auth/organizations/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package sessions         -o auth/sessions/schemas.gen.go         -generate types,skip-prune                 ../../../api/specifications/auth/sessions/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package sessions         -o auth/sessions/parameters.gen.go      -generate types,skip-prune                 ../../../api/specifications/auth/sessions/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package sessions         -o auth/sessions/paths.gen.go           -generate server,strict-server,skip-prune,types  ../../../api/specifications/auth/sessions/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package users            -o auth/users/schemas.gen.go            -generate types,skip-prune                 ../../../api/specifications/auth/users/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package users            -o auth/users/parameters.gen.go         -generate types,skip-prune                 ../../../api/specifications/auth/users/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package users            -o auth/users/paths.gen.go              -generate server,strict-server,skip-prune,types  ../../../api/specifications/auth/users/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package verificationtokens -o auth/verification-tokens/schemas.gen.go -generate types,skip-prune               ../../../api/specifications/auth/verification-tokens/schemas.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package auth             -o auth/auth.gen.go                     -generate types,skip-prune                 ../../../api/specifications/auth/auth/paths.yaml

// Intelligence features
//go:generate oapi-codegen -config oapi-codegen.yaml -package artifacts        -o intelligence/artifacts/schemas.gen.go    -generate types,skip-prune                 ../../../api/specifications/intelligence/artifacts/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package artifacts        -o intelligence/artifacts/parameters.gen.go -generate types,skip-prune                 ../../../api/specifications/intelligence/artifacts/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package artifacts        -o intelligence/artifacts/paths.gen.go      -generate server,strict-server,skip-prune,types  ../../../api/specifications/intelligence/artifacts/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package labels           -o intelligence/labels/schemas.gen.go       -generate types,skip-prune                 ../../../api/specifications/intelligence/labels/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package labels           -o intelligence/labels/parameters.gen.go    -generate types,skip-prune                 ../../../api/specifications/intelligence/labels/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package labels           -o intelligence/labels/paths.gen.go         -generate server,strict-server,skip-prune,types  ../../../api/specifications/intelligence/labels/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package pipelines        -o intelligence/pipelines/schemas.gen.go    -generate types,skip-prune                 ../../../api/specifications/intelligence/pipelines/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package pipelines        -o intelligence/pipelines/parameters.gen.go -generate types,skip-prune                 ../../../api/specifications/intelligence/pipelines/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package pipelines        -o intelligence/pipelines/paths.gen.go      -generate server,strict-server,skip-prune,types  ../../../api/specifications/intelligence/pipelines/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package runs             -o intelligence/runs/schemas.gen.go         -generate types,skip-prune                 ../../../api/specifications/intelligence/runs/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package runs             -o intelligence/runs/parameters.gen.go      -generate types,skip-prune                 ../../../api/specifications/intelligence/runs/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package runs             -o intelligence/runs/paths.gen.go           -generate server,strict-server,skip-prune,types  ../../../api/specifications/intelligence/runs/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package tools            -o intelligence/tools/schemas.gen.go        -generate types,skip-prune                 ../../../api/specifications/intelligence/tools/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package tools            -o intelligence/tools/parameters.gen.go     -generate types,skip-prune                 ../../../api/specifications/intelligence/tools/parameters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package tools            -o intelligence/tools/paths.gen.go          -generate server,strict-server,skip-prune,types  ../../../api/specifications/intelligence/tools/paths.yaml

//go:generate oapi-codegen -import-mapping "./schemas.yaml:-,./parameters.yaml:-,../../common/filters/schemas.yaml:-,../../common/responses/schemas.yaml:-,../../common/security/schemas.yaml:-,../../common/parameters/schemas.yaml:-,../filters/schemas.yaml:-" -package api 				-o api.gen.go									-generate spec		../../../api/openapi.yaml

package api
