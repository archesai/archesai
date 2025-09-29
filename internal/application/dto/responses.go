package dto

import (
	"net/http"
	"time"
)

// Standard response types that match the OpenAPI spec

// BadRequestResponse represents a 400 Bad Request error
type BadRequestResponse struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Detail    string    `json:"detail,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewBadRequestResponse(detail string) BadRequestResponse {
	return BadRequestResponse{
		Type:      "https://example.com/probs/bad-request",
		Title:     "Bad Request",
		Status:    http.StatusBadRequest,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// UnauthorizedResponse represents a 401 Unauthorized error
type UnauthorizedResponse struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Detail    string    `json:"detail,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewUnauthorizedResponse(detail string) UnauthorizedResponse {
	return UnauthorizedResponse{
		Type:      "https://example.com/probs/unauthorized",
		Title:     "Unauthorized",
		Status:    http.StatusUnauthorized,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// NotFoundResponse represents a 404 Not Found error
type NotFoundResponse struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Detail    string    `json:"detail,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewNotFoundResponse(resource, id string) NotFoundResponse {
	detail := resource + " not found"
	if id != "" {
		detail = resource + " with id " + id + " not found"
	}
	return NotFoundResponse{
		Type:      "https://example.com/probs/not-found",
		Title:     "Not Found",
		Status:    http.StatusNotFound,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// ConflictResponse represents a 409 Conflict error
type ConflictResponse struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Detail    string    `json:"detail,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewConflictResponse(detail string) ConflictResponse {
	return ConflictResponse{
		Type:      "https://example.com/probs/conflict",
		Title:     "Conflict",
		Status:    http.StatusConflict,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// InternalServerErrorResponse represents a 500 Internal Server Error
type InternalServerErrorResponse struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Status    int       `json:"status"`
	Detail    string    `json:"detail,omitempty"`
	Instance  string    `json:"instance,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewInternalServerErrorResponse(detail string) InternalServerErrorResponse {
	return InternalServerErrorResponse{
		Type:      "https://example.com/probs/internal-server-error",
		Title:     "Internal Server Error",
		Status:    http.StatusInternalServerError,
		Detail:    detail,
		Timestamp: time.Now(),
	}
}

// TooManyRequestsResponse represents a 429 Too Many Requests error
type TooManyRequestsResponse struct {
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Status     int       `json:"status"`
	Detail     string    `json:"detail,omitempty"`
	Instance   string    `json:"instance,omitempty"`
	RetryAfter int       `json:"retryAfter,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewTooManyRequestsResponse(retryAfter int) TooManyRequestsResponse {
	return TooManyRequestsResponse{
		Type:       "https://example.com/probs/too-many-requests",
		Title:      "Too Many Requests",
		Status:     http.StatusTooManyRequests,
		Detail:     "Rate limit exceeded",
		RetryAfter: retryAfter,
		Timestamp:  time.Now(),
	}
}

// NoContentResponse represents a 204 No Content response
type NoContentResponse struct{}
