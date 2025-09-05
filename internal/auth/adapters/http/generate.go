// Package http contains HTTP handlers for the auth domain.
package http

// Generate server interfaces for auth domain (Auth, Users, Sessions, Accounts tags)
//go:generate go tool oapi-codegen --config=oapi-codegen.yaml ../../../../api/openapi.bundled.yaml
