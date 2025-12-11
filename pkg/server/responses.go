package server

import (
	"github.com/archesai/archesai/pkg/httputil"
)

// ProblemDetails is an alias for httputil.ProblemDetails for backwards compatibility.
type ProblemDetails = httputil.ProblemDetails

// NewBadRequestResponse creates a new 400 Bad Request response
func NewBadRequestResponse(detail, instance string) ProblemDetails {
	return httputil.NewBadRequestResponse(detail, instance)
}

// NewUnauthorizedResponse creates a new 401 Unauthorized response
func NewUnauthorizedResponse(detail, instance string) ProblemDetails {
	return httputil.NewUnauthorizedResponse(detail, instance)
}

// NewForbiddenResponse creates a new 403 Forbidden response
func NewForbiddenResponse(detail, instance string) ProblemDetails {
	return httputil.NewForbiddenResponse(detail, instance)
}

// NewNotFoundResponse creates a new 404 Not Found response
func NewNotFoundResponse(detail, instance string) ProblemDetails {
	return httputil.NewNotFoundResponse(detail, instance)
}

// NewConflictResponse creates a new 409 Conflict response
func NewConflictResponse(detail, instance string) ProblemDetails {
	return httputil.NewConflictResponse(detail, instance)
}

// NewInternalServerErrorResponse creates a new 500 Internal Server Error response
func NewInternalServerErrorResponse(detail, instance string) ProblemDetails {
	return httputil.NewInternalServerErrorResponse(detail, instance)
}

// NewTooManyRequestsResponse creates a new 429 Too Many Requests response
func NewTooManyRequestsResponse(detail, instance string) ProblemDetails {
	return httputil.NewTooManyRequestsResponse(detail, instance)
}
