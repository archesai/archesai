package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ProblemDetails represents an RFC 7807 problem details response.
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// Common problem types
const (
	ProblemTypeNotFound     = "https://tools.ietf.org/html/rfc7231#section-6.5.4"
	ProblemTypeUnauthorized = "https://tools.ietf.org/html/rfc7235#section-3.1"
	ProblemTypeForbidden    = "https://tools.ietf.org/html/rfc7231#section-6.5.3"
	ProblemTypeBadRequest   = "https://tools.ietf.org/html/rfc7231#section-6.5.1"
	ProblemTypeInternal     = "https://tools.ietf.org/html/rfc7231#section-6.6.1"
)

// NewProblemDetails creates a new RFC 7807 problem details response.
func NewProblemDetails(status int, title, detail, instance string) *ProblemDetails {
	var problemType string

	switch status {
	case http.StatusNotFound:
		problemType = ProblemTypeNotFound
	case http.StatusUnauthorized:
		problemType = ProblemTypeUnauthorized
	case http.StatusForbidden:
		problemType = ProblemTypeForbidden
	case http.StatusBadRequest:
		problemType = ProblemTypeBadRequest
	case http.StatusInternalServerError:
		problemType = ProblemTypeInternal
	default:
		problemType = "about:blank"
	}

	return &ProblemDetails{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

// CustomErrorHandler handles HTTP errors and returns RFC 7807 problem details.
func (s *Server) CustomErrorHandler(err error, c echo.Context) {
	var (
		code   = http.StatusInternalServerError
		title  = "Internal Server Error"
		detail = err.Error()
	)

	// Extract error details from Echo HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			detail = msg
		}
	}

	// Set appropriate title based on status code
	switch code {
	case http.StatusNotFound:
		title = "Not Found"
		detail = "The requested resource could not be found"
	case http.StatusUnauthorized:
		title = "Unauthorized"
		detail = "Authentication is required to access this resource"
	case http.StatusForbidden:
		title = "Forbidden"
		detail = "You do not have permission to access this resource"
	case http.StatusBadRequest:
		title = "Bad Request"
		if detail == "" {
			detail = "The request could not be understood by the server"
		}
	case http.StatusMethodNotAllowed:
		title = "Method Not Allowed"
		detail = "The request method is not supported for this resource"
	case http.StatusTooManyRequests:
		title = "Too Many Requests"
		detail = "Rate limit exceeded"
	case http.StatusInternalServerError:
		title = "Internal Server Error"
		// Don't expose internal error details in production
		if s.config.Environment != "development" {
			detail = "An unexpected error occurred"
		}
	}

	// Create RFC 7807 problem details
	problem := NewProblemDetails(code, title, detail, c.Request().URL.Path)

	// Log the error
	s.logger.Error("HTTP error",
		"status", code,
		"path", c.Request().URL.Path,
		"method", c.Request().Method,
		"error", err,
		"request_id", c.Response().Header().Get(echo.HeaderXRequestID),
	)

	// Set content type for problem details
	c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")

	// Send the response
	if err := c.JSON(code, problem); err != nil {
		s.logger.Error("Failed to send error response", "error", err)
	}
}
