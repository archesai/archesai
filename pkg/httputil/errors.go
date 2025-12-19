package httputil

import "net/http"

// HTTPError is an interface for errors that carry HTTP status information.
// Handlers can return these errors and the route layer will respond appropriately.
type HTTPError interface {
	error
	StatusCode() int
	ProblemDetails(instance string) ProblemDetails
}

// BadRequestError represents a 400 Bad Request error.
type BadRequestError struct {
	Detail string
}

func (e BadRequestError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for BadRequestError.
func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

// ProblemDetails returns RFC 7807 problem details for BadRequestError.
func (e BadRequestError) ProblemDetails(instance string) ProblemDetails {
	return NewBadRequestResponse(e.Detail, instance)
}

// UnauthorizedError represents a 401 Unauthorized error.
type UnauthorizedError struct {
	Detail string
}

func (e UnauthorizedError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for UnauthorizedError.
func (e UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

// ProblemDetails returns RFC 7807 problem details for UnauthorizedError.
func (e UnauthorizedError) ProblemDetails(instance string) ProblemDetails {
	return NewUnauthorizedResponse(e.Detail, instance)
}

// ForbiddenError represents a 403 Forbidden error.
type ForbiddenError struct {
	Detail string
}

func (e ForbiddenError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for ForbiddenError.
func (e ForbiddenError) StatusCode() int {
	return http.StatusForbidden
}

// ProblemDetails returns RFC 7807 problem details for ForbiddenError.
func (e ForbiddenError) ProblemDetails(instance string) ProblemDetails {
	return NewForbiddenResponse(e.Detail, instance)
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	Detail string
}

func (e NotFoundError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for NotFoundError.
func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// ProblemDetails returns RFC 7807 problem details for NotFoundError.
func (e NotFoundError) ProblemDetails(instance string) ProblemDetails {
	return NewNotFoundResponse(e.Detail, instance)
}

// ConflictError represents a 409 Conflict error.
type ConflictError struct {
	Detail string
}

func (e ConflictError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for ConflictError.
func (e ConflictError) StatusCode() int {
	return http.StatusConflict
}

// ProblemDetails returns RFC 7807 problem details for ConflictError.
func (e ConflictError) ProblemDetails(instance string) ProblemDetails {
	return NewConflictResponse(e.Detail, instance)
}

// UnprocessableEntityError represents a 422 Unprocessable Entity error.
type UnprocessableEntityError struct {
	Detail string
}

func (e UnprocessableEntityError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for UnprocessableEntityError.
func (e UnprocessableEntityError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

// ProblemDetails returns RFC 7807 problem details for UnprocessableEntityError.
func (e UnprocessableEntityError) ProblemDetails(instance string) ProblemDetails {
	return NewUnprocessableEntityResponse(e.Detail, instance)
}

// TooManyRequestsError represents a 429 Too Many Requests error.
type TooManyRequestsError struct {
	Detail string
}

func (e TooManyRequestsError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for TooManyRequestsError.
func (e TooManyRequestsError) StatusCode() int {
	return http.StatusTooManyRequests
}

// ProblemDetails returns RFC 7807 problem details for TooManyRequestsError.
func (e TooManyRequestsError) ProblemDetails(instance string) ProblemDetails {
	return NewTooManyRequestsResponse(e.Detail, instance)
}

// InternalServerError represents a 500 Internal Server Error.
type InternalServerError struct {
	Detail string
}

func (e InternalServerError) Error() string {
	return e.Detail
}

// StatusCode returns the HTTP status code for InternalServerError.
func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

// ProblemDetails returns RFC 7807 problem details for InternalServerError.
func (e InternalServerError) ProblemDetails(instance string) ProblemDetails {
	return NewInternalServerErrorResponse(e.Detail, instance)
}
