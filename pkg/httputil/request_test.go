package httputil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/archesai/archesai/pkg/validation"
)

func TestDecodeJSON(t *testing.T) {
	type Input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	t.Run("decodes valid JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"John","email":"john@example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")

		var input Input
		err := DecodeJSON(req, &input)
		if err != nil {
			t.Fatalf("DecodeJSON error: %v", err)
		}

		if input.Name != "John" {
			t.Errorf("Name = %q, want %q", input.Name, "John")
		}
		if input.Email != "john@example.com" {
			t.Errorf("Email = %q, want %q", input.Email, "john@example.com")
		}
	})

	t.Run("returns error for nil body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Body = nil

		var input Input
		err := DecodeJSON(req, &input)
		if err == nil {
			t.Error("expected error for nil body")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid json}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)

		var input Input
		err := DecodeJSON(req, &input)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("returns error for wrong type", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name": 123}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)

		var input Input
		err := DecodeJSON(req, &input)
		if err == nil {
			t.Error("expected error for wrong type")
		}
	})
}

// ValidatableInput implements validation.Validator
type ValidatableInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (v *ValidatableInput) Validate() validation.Errors {
	var errs validation.Errors
	validation.Required(v.Name, "name", &errs)
	validation.Required(v.Email, "email", &errs)
	if v.Email != nil {
		validation.Email(v.Email, "email", &errs)
	}
	return errs
}

func TestDecodeAndValidate(t *testing.T) {
	t.Run("decodes and validates valid input", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"John","email":"john@example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)

		var input ValidatableInput
		err := DecodeAndValidate(req, &input)
		if err != nil {
			t.Fatalf("DecodeAndValidate error: %v", err)
		}

		if input.Name == nil || *input.Name != "John" {
			t.Errorf("Name = %v, want pointer to 'John'", input.Name)
		}
	})

	t.Run("returns validation errors", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"","email":""}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)

		var input ValidatableInput
		err := DecodeAndValidate(req, &input)
		if err == nil {
			t.Fatal("expected validation error")
		}

		// Check if it's a validation.Errors type
		if _, ok := err.(validation.Errors); !ok {
			t.Errorf("expected validation.Errors, got %T", err)
		}
	})

	t.Run("returns decode error before validation", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)

		var input ValidatableInput
		err := DecodeAndValidate(req, &input)
		if err == nil {
			t.Fatal("expected decode error")
		}

		// Should not be a validation error
		if _, ok := err.(validation.Errors); ok {
			t.Error("expected decode error, not validation error")
		}
	})
}

func TestBindPathParamUUID(t *testing.T) {
	t.Run("parses valid UUID", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/users/550e8400-e29b-41d4-a716-446655440000",
			nil,
		)
		req.SetPathValue("id", "550e8400-e29b-41d4-a716-446655440000")

		id, err := BindPathParamUUID(req, "id")
		if err != nil {
			t.Fatalf("BindPathParamUUID error: %v", err)
		}

		expected := "550e8400-e29b-41d4-a716-446655440000"
		if id.String() != expected {
			t.Errorf("UUID = %q, want %q", id.String(), expected)
		}
	})

	t.Run("returns error for missing param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		// Don't set path value

		_, err := BindPathParamUUID(req, "id")
		if err == nil {
			t.Error("expected error for missing param")
		}
	})

	t.Run("returns error for invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/not-a-uuid", nil)
		req.SetPathValue("id", "not-a-uuid")

		_, err := BindPathParamUUID(req, "id")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestBindPathParamString(t *testing.T) {
	t.Run("returns value when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/john", nil)
		req.SetPathValue("username", "john")

		value, err := BindPathParamString(req, "username")
		if err != nil {
			t.Fatalf("BindPathParamString error: %v", err)
		}

		if value != "john" {
			t.Errorf("value = %q, want %q", value, "john")
		}
	})

	t.Run("returns error when missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		// Don't set path value

		_, err := BindPathParamString(req, "username")
		if err == nil {
			t.Error("expected error for missing param")
		}
	})
}

func TestBindPathParamInt(t *testing.T) {
	t.Run("parses valid integer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pages/42", nil)
		req.SetPathValue("page", "42")

		value, err := BindPathParamInt(req, "page")
		if err != nil {
			t.Fatalf("BindPathParamInt error: %v", err)
		}

		if value != 42 {
			t.Errorf("value = %d, want %d", value, 42)
		}
	})

	t.Run("parses negative integer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/offset/-10", nil)
		req.SetPathValue("offset", "-10")

		value, err := BindPathParamInt(req, "offset")
		if err != nil {
			t.Fatalf("BindPathParamInt error: %v", err)
		}

		if value != -10 {
			t.Errorf("value = %d, want %d", value, -10)
		}
	})

	t.Run("returns error for missing param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pages/", nil)
		// Don't set path value

		_, err := BindPathParamInt(req, "page")
		if err == nil {
			t.Error("expected error for missing param")
		}
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pages/abc", nil)
		req.SetPathValue("page", "abc")

		_, err := BindPathParamInt(req, "page")
		if err == nil {
			t.Error("expected error for invalid integer")
		}
	})

	t.Run("returns error for float", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pages/3.14", nil)
		req.SetPathValue("page", "3.14")

		_, err := BindPathParamInt(req, "page")
		if err == nil {
			t.Error("expected error for float")
		}
	})
}

func TestRequiredHeader(t *testing.T) {
	t.Run("returns value when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token123")

		value, err := RequiredHeader(req, "Authorization")
		if err != nil {
			t.Fatalf("RequiredHeader error: %v", err)
		}

		if value != "Bearer token123" {
			t.Errorf("value = %q, want %q", value, "Bearer token123")
		}
	})

	t.Run("returns error when missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Don't set header

		_, err := RequiredHeader(req, "Authorization")
		if err == nil {
			t.Error("expected error for missing header")
		}
	})

	t.Run("returns error for empty header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "")

		_, err := RequiredHeader(req, "Authorization")
		if err == nil {
			t.Error("expected error for empty header")
		}
	})
}

func TestOptionalHeader(t *testing.T) {
	t.Run("returns value when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", "abc123")

		value := OptionalHeader(req, "X-Request-ID")
		if value != "abc123" {
			t.Errorf("value = %q, want %q", value, "abc123")
		}
	})

	t.Run("returns empty string when missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Don't set header

		value := OptionalHeader(req, "X-Request-ID")
		if value != "" {
			t.Errorf("value = %q, want empty string", value)
		}
	})
}
