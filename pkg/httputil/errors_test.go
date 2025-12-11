package httputil

import (
	"net/http"
	"testing"
)

func TestBadRequestError(t *testing.T) {
	err := BadRequestError{Detail: "invalid input"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "invalid input" {
			t.Errorf("Error() = %q, want %q", err.Error(), "invalid input")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusBadRequest {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusBadRequest)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users")
		if pd.Status != http.StatusBadRequest {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusBadRequest)
		}
		if pd.Title != "Bad Request" {
			t.Errorf("Title = %q, want %q", pd.Title, "Bad Request")
		}
		if pd.Detail != "invalid input" {
			t.Errorf("Detail = %q, want %q", pd.Detail, "invalid input")
		}
		if pd.Instance != "/api/users" {
			t.Errorf("Instance = %q, want %q", pd.Instance, "/api/users")
		}
	})
}

func TestUnauthorizedError(t *testing.T) {
	err := UnauthorizedError{Detail: "missing token"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "missing token" {
			t.Errorf("Error() = %q, want %q", err.Error(), "missing token")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusUnauthorized {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusUnauthorized)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users")
		if pd.Status != http.StatusUnauthorized {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusUnauthorized)
		}
		if pd.Title != "Unauthorized" {
			t.Errorf("Title = %q, want %q", pd.Title, "Unauthorized")
		}
	})
}

func TestForbiddenError(t *testing.T) {
	err := ForbiddenError{Detail: "access denied"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "access denied" {
			t.Errorf("Error() = %q, want %q", err.Error(), "access denied")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusForbidden {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusForbidden)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/admin")
		if pd.Status != http.StatusForbidden {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusForbidden)
		}
		if pd.Title != "Forbidden" {
			t.Errorf("Title = %q, want %q", pd.Title, "Forbidden")
		}
	})
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError{Detail: "user not found"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "user not found" {
			t.Errorf("Error() = %q, want %q", err.Error(), "user not found")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusNotFound {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusNotFound)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users/123")
		if pd.Status != http.StatusNotFound {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusNotFound)
		}
		if pd.Title != "Not Found" {
			t.Errorf("Title = %q, want %q", pd.Title, "Not Found")
		}
	})
}

func TestConflictError(t *testing.T) {
	err := ConflictError{Detail: "email already exists"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "email already exists" {
			t.Errorf("Error() = %q, want %q", err.Error(), "email already exists")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusConflict {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusConflict)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users")
		if pd.Status != http.StatusConflict {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusConflict)
		}
		if pd.Title != "Conflict" {
			t.Errorf("Title = %q, want %q", pd.Title, "Conflict")
		}
	})
}

func TestUnprocessableEntityError(t *testing.T) {
	err := UnprocessableEntityError{Detail: "validation failed"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "validation failed" {
			t.Errorf("Error() = %q, want %q", err.Error(), "validation failed")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusUnprocessableEntity {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusUnprocessableEntity)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users")
		if pd.Status != http.StatusUnprocessableEntity {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusUnprocessableEntity)
		}
		if pd.Title != "Unprocessable Entity" {
			t.Errorf("Title = %q, want %q", pd.Title, "Unprocessable Entity")
		}
	})
}

func TestTooManyRequestsError(t *testing.T) {
	err := TooManyRequestsError{Detail: "rate limit exceeded"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "rate limit exceeded" {
			t.Errorf("Error() = %q, want %q", err.Error(), "rate limit exceeded")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusTooManyRequests {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusTooManyRequests)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users")
		if pd.Status != http.StatusTooManyRequests {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusTooManyRequests)
		}
		if pd.Title != "Too Many Requests" {
			t.Errorf("Title = %q, want %q", pd.Title, "Too Many Requests")
		}
	})
}

func TestInternalServerError(t *testing.T) {
	err := InternalServerError{Detail: "database connection failed"}

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "database connection failed" {
			t.Errorf("Error() = %q, want %q", err.Error(), "database connection failed")
		}
	})

	t.Run("StatusCode", func(t *testing.T) {
		if err.StatusCode() != http.StatusInternalServerError {
			t.Errorf("StatusCode() = %d, want %d", err.StatusCode(), http.StatusInternalServerError)
		}
	})

	t.Run("ProblemDetails", func(t *testing.T) {
		pd := err.ProblemDetails("/api/users")
		if pd.Status != http.StatusInternalServerError {
			t.Errorf("Status = %d, want %d", pd.Status, http.StatusInternalServerError)
		}
		if pd.Title != "Internal Server Error" {
			t.Errorf("Title = %q, want %q", pd.Title, "Internal Server Error")
		}
	})
}

func TestHTTPError_Interface(t *testing.T) {
	// Verify all error types implement HTTPError interface
	tests := []struct {
		name string
		err  HTTPError
	}{
		{"BadRequestError", BadRequestError{Detail: "test"}},
		{"UnauthorizedError", UnauthorizedError{Detail: "test"}},
		{"ForbiddenError", ForbiddenError{Detail: "test"}},
		{"NotFoundError", NotFoundError{Detail: "test"}},
		{"ConflictError", ConflictError{Detail: "test"}},
		{"UnprocessableEntityError", UnprocessableEntityError{Detail: "test"}},
		{"TooManyRequestsError", TooManyRequestsError{Detail: "test"}},
		{"InternalServerError", InternalServerError{Detail: "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify error interface
			if tt.err.Error() != "test" {
				t.Errorf("Error() = %q, want %q", tt.err.Error(), "test")
			}

			// Verify StatusCode returns valid HTTP status
			status := tt.err.StatusCode()
			if status < 400 || status >= 600 {
				t.Errorf("StatusCode() = %d, want 4xx or 5xx", status)
			}

			// Verify ProblemDetails returns valid structure
			pd := tt.err.ProblemDetails("/test")
			if pd.Status != status {
				t.Errorf("ProblemDetails.Status = %d, want %d", pd.Status, status)
			}
			if pd.Instance != "/test" {
				t.Errorf("ProblemDetails.Instance = %q, want %q", pd.Instance, "/test")
			}
			if pd.Type == "" {
				t.Error("ProblemDetails.Type should not be empty")
			}
			if pd.Title == "" {
				t.Error("ProblemDetails.Title should not be empty")
			}
		})
	}
}
