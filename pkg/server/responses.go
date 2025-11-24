package server

import (
	"net/http"
	"time"
)

// Standard response types that match the OpenAPI spec

// ProblemDetails represents an RFC 7807 problem details response.
type ProblemDetails struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Detail    string    `json:"detail,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// NewBadRequestResponse creates a new 400 Bad Request response
func NewBadRequestResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc7231#section-6.5.1",
		Title:     "Bad Request",
		Status:    http.StatusBadRequest,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}

// NewUnauthorizedResponse creates a new 401 Unauthorized response
func NewUnauthorizedResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc7235#section-3.1",
		Title:     "Unauthorized",
		Status:    http.StatusUnauthorized,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}

// NewForbiddenResponse creates a new 403 Forbidden response
func NewForbiddenResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc7231#section-6.5.3",
		Title:     "Forbidden",
		Status:    http.StatusForbidden,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}

// NewNotFoundResponse creates a new 404 Not Found response
func NewNotFoundResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc7231#section-6.5.4",
		Title:     "Not Found",
		Status:    http.StatusNotFound,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}

// NewConflictResponse creates a new 409 Conflict response
func NewConflictResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc7231#section-6.5.8",
		Title:     "Conflict",
		Status:    http.StatusConflict,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}

// NewInternalServerErrorResponse creates a new 500 Internal Server Error response
func NewInternalServerErrorResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc7231#section-6.6.1",
		Title:     "Internal Server Error",
		Status:    http.StatusInternalServerError,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}

// NewTooManyRequestsResponse creates a new 429 Too Many Requests response
func NewTooManyRequestsResponse(detail, instance string) ProblemDetails {
	return ProblemDetails{
		Type:      "https://tools.ietf.org/html/rfc6585#section-4",
		Title:     "Too Many Requests",
		Status:    http.StatusTooManyRequests,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
}
