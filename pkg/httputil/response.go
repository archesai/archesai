package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/archesai/archesai/pkg/validation"
)

// Responder is implemented by response types that know their HTTP status code.
type Responder interface {
	StatusCode() int
}

// WriteResponse writes a response that implements Responder.
// It automatically uses the correct status code and content type.
func WriteResponse(w http.ResponseWriter, response Responder) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode())
	return json.NewEncoder(w).Encode(response)
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON[T any](w http.ResponseWriter, status int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteNoContent writes a 204 No Content response.
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// WriteError writes a ProblemDetails error response with the given status code.
func WriteError(w http.ResponseWriter, status int, detail, instance string) error {
	problem := ProblemDetails{
		Type:      typeForStatus(status),
		Title:     http.StatusText(status),
		Status:    status,
		Detail:    detail,
		Instance:  instance,
		Timestamp: time.Now(),
	}
	return writeProblemDetails(w, status, problem)
}

// WriteHTTPError writes an HTTPError as a ProblemDetails response.
func WriteHTTPError(w http.ResponseWriter, err HTTPError, instance string) error {
	problem := err.ProblemDetails(instance)
	return writeProblemDetails(w, err.StatusCode(), problem)
}

// ValidationProblemDetails extends ProblemDetails with field-level validation errors.
type ValidationProblemDetails struct {
	ProblemDetails
	Errors []validation.FieldError `json:"errors,omitempty"`
}

// WriteValidationErrors writes validation errors as a 422 Unprocessable Entity response.
func WriteValidationErrors(w http.ResponseWriter, errs validation.Errors, instance string) error {
	problem := ValidationProblemDetails{
		ProblemDetails: ProblemDetails{
			Type:      "https://tools.ietf.org/html/rfc4918#section-11.2",
			Title:     "Validation Failed",
			Status:    http.StatusUnprocessableEntity,
			Detail:    fmt.Sprintf("Request validation failed with %d error(s)", len(errs)),
			Instance:  instance,
			Timestamp: time.Now(),
		},
		Errors: errs,
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	return json.NewEncoder(w).Encode(problem)
}

// writeProblemDetails writes a ProblemDetails response.
func writeProblemDetails(w http.ResponseWriter, status int, problem ProblemDetails) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(problem)
}

// typeForStatus returns the RFC reference URL for a given HTTP status code.
func typeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "https://tools.ietf.org/html/rfc7231#section-6.5.1"
	case http.StatusUnauthorized:
		return "https://tools.ietf.org/html/rfc7235#section-3.1"
	case http.StatusForbidden:
		return "https://tools.ietf.org/html/rfc7231#section-6.5.3"
	case http.StatusNotFound:
		return "https://tools.ietf.org/html/rfc7231#section-6.5.4"
	case http.StatusConflict:
		return "https://tools.ietf.org/html/rfc7231#section-6.5.8"
	case http.StatusUnprocessableEntity:
		return "https://tools.ietf.org/html/rfc4918#section-11.2"
	case http.StatusTooManyRequests:
		return "https://tools.ietf.org/html/rfc6585#section-4"
	case http.StatusInternalServerError:
		return "https://tools.ietf.org/html/rfc7231#section-6.6.1"
	default:
		return fmt.Sprintf("https://httpstatuses.com/%d", status)
	}
}
