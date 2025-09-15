// Package members provides domain logic for member management
package members

import "errors"

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package members --include-tags Members ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package members --include-tags Members ../../api/openapi.bundled.yaml

// Member errors
var (
	ErrMemberNotFound = errors.New("member not found")
	// ErrDuplicateMember is returned when a member already exists
	ErrDuplicateMember = errors.New("member already exists")
)
