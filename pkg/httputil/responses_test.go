package httputil

import (
	"net/http"
	"testing"
)

func TestNewBadRequestResponse(t *testing.T) {
	pd := NewBadRequestResponse("invalid input", "/api/users")

	if pd.Type != "https://tools.ietf.org/html/rfc7231#section-6.5.1" {
		t.Errorf("Type = %q, want RFC 7231 section 6.5.1", pd.Type)
	}
	if pd.Title != "Bad Request" {
		t.Errorf("Title = %q, want %q", pd.Title, "Bad Request")
	}
	if pd.Status != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusBadRequest)
	}
	if pd.Detail != "invalid input" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "invalid input")
	}
	if pd.Instance != "/api/users" {
		t.Errorf("Instance = %q, want %q", pd.Instance, "/api/users")
	}
	if pd.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestNewUnauthorizedResponse(t *testing.T) {
	pd := NewUnauthorizedResponse("missing token", "/api/protected")

	if pd.Type != "https://tools.ietf.org/html/rfc7235#section-3.1" {
		t.Errorf("Type = %q, want RFC 7235 section 3.1", pd.Type)
	}
	if pd.Title != "Unauthorized" {
		t.Errorf("Title = %q, want %q", pd.Title, "Unauthorized")
	}
	if pd.Status != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusUnauthorized)
	}
	if pd.Detail != "missing token" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "missing token")
	}
	if pd.Instance != "/api/protected" {
		t.Errorf("Instance = %q, want %q", pd.Instance, "/api/protected")
	}
}

func TestNewForbiddenResponse(t *testing.T) {
	pd := NewForbiddenResponse("access denied", "/api/admin")

	if pd.Type != "https://tools.ietf.org/html/rfc7231#section-6.5.3" {
		t.Errorf("Type = %q, want RFC 7231 section 6.5.3", pd.Type)
	}
	if pd.Title != "Forbidden" {
		t.Errorf("Title = %q, want %q", pd.Title, "Forbidden")
	}
	if pd.Status != http.StatusForbidden {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusForbidden)
	}
	if pd.Detail != "access denied" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "access denied")
	}
}

func TestNewNotFoundResponse(t *testing.T) {
	pd := NewNotFoundResponse("user not found", "/api/users/123")

	if pd.Type != "https://tools.ietf.org/html/rfc7231#section-6.5.4" {
		t.Errorf("Type = %q, want RFC 7231 section 6.5.4", pd.Type)
	}
	if pd.Title != "Not Found" {
		t.Errorf("Title = %q, want %q", pd.Title, "Not Found")
	}
	if pd.Status != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusNotFound)
	}
	if pd.Detail != "user not found" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "user not found")
	}
}

func TestNewConflictResponse(t *testing.T) {
	pd := NewConflictResponse("email already exists", "/api/users")

	if pd.Type != "https://tools.ietf.org/html/rfc7231#section-6.5.8" {
		t.Errorf("Type = %q, want RFC 7231 section 6.5.8", pd.Type)
	}
	if pd.Title != "Conflict" {
		t.Errorf("Title = %q, want %q", pd.Title, "Conflict")
	}
	if pd.Status != http.StatusConflict {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusConflict)
	}
	if pd.Detail != "email already exists" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "email already exists")
	}
}

func TestNewInternalServerErrorResponse(t *testing.T) {
	pd := NewInternalServerErrorResponse("database error", "/api/users")

	if pd.Type != "https://tools.ietf.org/html/rfc7231#section-6.6.1" {
		t.Errorf("Type = %q, want RFC 7231 section 6.6.1", pd.Type)
	}
	if pd.Title != "Internal Server Error" {
		t.Errorf("Title = %q, want %q", pd.Title, "Internal Server Error")
	}
	if pd.Status != http.StatusInternalServerError {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusInternalServerError)
	}
	if pd.Detail != "database error" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "database error")
	}
}

func TestNewTooManyRequestsResponse(t *testing.T) {
	pd := NewTooManyRequestsResponse("rate limit exceeded", "/api/users")

	if pd.Type != "https://tools.ietf.org/html/rfc6585#section-4" {
		t.Errorf("Type = %q, want RFC 6585 section 4", pd.Type)
	}
	if pd.Title != "Too Many Requests" {
		t.Errorf("Title = %q, want %q", pd.Title, "Too Many Requests")
	}
	if pd.Status != http.StatusTooManyRequests {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusTooManyRequests)
	}
	if pd.Detail != "rate limit exceeded" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "rate limit exceeded")
	}
}

func TestNewUnprocessableEntityResponse(t *testing.T) {
	pd := NewUnprocessableEntityResponse("validation failed", "/api/users")

	if pd.Type != "https://tools.ietf.org/html/rfc4918#section-11.2" {
		t.Errorf("Type = %q, want RFC 4918 section 11.2", pd.Type)
	}
	if pd.Title != "Unprocessable Entity" {
		t.Errorf("Title = %q, want %q", pd.Title, "Unprocessable Entity")
	}
	if pd.Status != http.StatusUnprocessableEntity {
		t.Errorf("Status = %d, want %d", pd.Status, http.StatusUnprocessableEntity)
	}
	if pd.Detail != "validation failed" {
		t.Errorf("Detail = %q, want %q", pd.Detail, "validation failed")
	}
}

func TestProblemDetails_AllResponses(t *testing.T) {
	// Table-driven test to verify all responses follow RFC 7807 structure
	tests := []struct {
		name           string
		createResponse func(detail, instance string) ProblemDetails
		expectedStatus int
		expectedTitle  string
	}{
		{
			name:           "BadRequest",
			createResponse: NewBadRequestResponse,
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
		},
		{
			name:           "Unauthorized",
			createResponse: NewUnauthorizedResponse,
			expectedStatus: http.StatusUnauthorized,
			expectedTitle:  "Unauthorized",
		},
		{
			name:           "Forbidden",
			createResponse: NewForbiddenResponse,
			expectedStatus: http.StatusForbidden,
			expectedTitle:  "Forbidden",
		},
		{
			name:           "NotFound",
			createResponse: NewNotFoundResponse,
			expectedStatus: http.StatusNotFound,
			expectedTitle:  "Not Found",
		},
		{
			name:           "Conflict",
			createResponse: NewConflictResponse,
			expectedStatus: http.StatusConflict,
			expectedTitle:  "Conflict",
		},
		{
			name:           "InternalServerError",
			createResponse: NewInternalServerErrorResponse,
			expectedStatus: http.StatusInternalServerError,
			expectedTitle:  "Internal Server Error",
		},
		{
			name:           "TooManyRequests",
			createResponse: NewTooManyRequestsResponse,
			expectedStatus: http.StatusTooManyRequests,
			expectedTitle:  "Too Many Requests",
		},
		{
			name:           "UnprocessableEntity",
			createResponse: NewUnprocessableEntityResponse,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedTitle:  "Unprocessable Entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd := tt.createResponse("test detail", "/test/instance")

			// Verify RFC 7807 required fields
			if pd.Type == "" {
				t.Error("Type should not be empty (RFC 7807)")
			}
			if pd.Title != tt.expectedTitle {
				t.Errorf("Title = %q, want %q", pd.Title, tt.expectedTitle)
			}
			if pd.Status != tt.expectedStatus {
				t.Errorf("Status = %d, want %d", pd.Status, tt.expectedStatus)
			}
			if pd.Detail != "test detail" {
				t.Errorf("Detail = %q, want %q", pd.Detail, "test detail")
			}
			if pd.Instance != "/test/instance" {
				t.Errorf("Instance = %q, want %q", pd.Instance, "/test/instance")
			}
			if pd.Timestamp.IsZero() {
				t.Error("Timestamp should be set")
			}
		})
	}
}
