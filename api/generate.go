//go:generate oapi-codegen -config oapi-codegen.yaml -package common -o ../gen/api/common/errors.gen.go   -generate types,skip-prune ./common/errors.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package common -o ../gen/api/common/filters.gen.go  -generate types,skip-prune ./common/filters.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package common -o ../gen/api/common/security.gen.go -generate types,skip-prune ./common/security.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package config -o ../gen/api/infrastructure/config/schemas.gen.go -generate types,skip-prune                 ./infrastructure/config/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package config -o ../gen/api/infrastructure/config/paths.gen.go   -generate server,strict-server,skip-prune,types  ./infrastructure/config/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package health -o ../gen/api/infrastructure/health/schemas.gen.go -generate types,skip-prune		           ./infrastructure/health/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package health -o ../gen/api/infrastructure/health/paths.gen.go   -generate server,strict-server,skip-prune,types  ./infrastructure/health/paths.yaml

// Auth features
//go:generate oapi-codegen -config oapi-codegen.yaml -package accounts         -o ../gen/api/features/auth/accounts/schemas.gen.go         -generate types,skip-prune                 ./features/auth/accounts/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package accounts         -o ../gen/api/features/auth/accounts/paths.gen.go           -generate server,strict-server,skip-prune,types  ./features/auth/accounts/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package apitokens        -o ../gen/api/features/auth/api-tokens/schemas.gen.go       -generate types,skip-prune                 ./features/auth/api-tokens/schemas.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package invitations      -o ../gen/api/features/auth/invitations/schemas.gen.go      -generate types,skip-prune                 ./features/auth/invitations/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package invitations      -o ../gen/api/features/auth/invitations/paths.gen.go        -generate server,strict-server,skip-prune,types  ./features/auth/invitations/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package members          -o ../gen/api/features/auth/members/schemas.gen.go          -generate types,skip-prune                 ./features/auth/members/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package members          -o ../gen/api/features/auth/members/paths.gen.go            -generate server,strict-server,skip-prune,types  ./features/auth/members/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package organizations    -o ../gen/api/features/auth/organizations/schemas.gen.go    -generate types,skip-prune                 ./features/auth/organizations/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package organizations    -o ../gen/api/features/auth/organizations/paths.gen.go      -generate server,strict-server,skip-prune,types  ./features/auth/organizations/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package sessions         -o ../gen/api/features/auth/sessions/schemas.gen.go         -generate types,skip-prune                 ./features/auth/sessions/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package sessions         -o ../gen/api/features/auth/sessions/paths.gen.go           -generate server,strict-server,skip-prune,types  ./features/auth/sessions/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package users            -o ../gen/api/features/auth/users/schemas.gen.go            -generate types,skip-prune                 ./features/auth/users/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package users            -o ../gen/api/features/auth/users/paths.gen.go              -generate server,strict-server,skip-prune,types  ./features/auth/users/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package verificationtokens -o ../gen/api/features/auth/verification-tokens/schemas.gen.go -generate types,skip-prune               ./features/auth/verification-tokens/schemas.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package auth             -o ../gen/api/features/auth/auth.gen.go                     -generate types,skip-prune                 ./features/auth/auth/paths.yaml

// Intelligence features
//go:generate oapi-codegen -config oapi-codegen.yaml -package artifacts        -o ../gen/api/features/intelligence/artifacts/schemas.gen.go -generate types,skip-prune                 ./features/intelligence/artifacts/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package artifacts        -o ../gen/api/features/intelligence/artifacts/paths.gen.go   -generate server,strict-server,skip-prune,types  ./features/intelligence/artifacts/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package labels           -o ../gen/api/features/intelligence/labels/schemas.gen.go    -generate types,skip-prune                 ./features/intelligence/labels/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package labels           -o ../gen/api/features/intelligence/labels/paths.gen.go      -generate server,strict-server,skip-prune,types  ./features/intelligence/labels/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package pipelines        -o ../gen/api/features/intelligence/pipelines/schemas.gen.go -generate types,skip-prune                 ./features/intelligence/pipelines/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package pipelines        -o ../gen/api/features/intelligence/pipelines/paths.gen.go   -generate server,strict-server,skip-prune,types  ./features/intelligence/pipelines/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package runs             -o ../gen/api/features/intelligence/runs/schemas.gen.go      -generate types,skip-prune                 ./features/intelligence/runs/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package runs             -o ../gen/api/features/intelligence/runs/paths.gen.go        -generate server,strict-server,skip-prune,types  ./features/intelligence/runs/paths.yaml

//go:generate oapi-codegen -config oapi-codegen.yaml -package tools            -o ../gen/api/features/intelligence/tools/schemas.gen.go     -generate types,skip-prune                 ./features/intelligence/tools/schemas.yaml
//go:generate oapi-codegen -config oapi-codegen.yaml -package tools            -o ../gen/api/features/intelligence/tools/paths.gen.go       -generate server,strict-server,skip-prune,types  ./features/intelligence/tools/paths.yaml

//go:generate oapi-codegen -import-mapping "./schemas.yaml:-,../../../common/filters.yaml:-,../../../common/errors.yaml:-" -package api 				-o ../gen/api/api.gen.go									-generate spec								./openapi.yaml

package api
