package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime"

	"github.com/archesai/archesai/pkg/validation"
)

// DecodeJSON decodes the JSON request body into the target struct.
func DecodeJSON[T any](r *http.Request, v *T) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(v)
}

// DecodeAndValidate decodes the JSON request body and validates it if the type implements Validator.
// Returns validation.Errors if validation fails, or a regular error if decoding fails.
func DecodeAndValidate[T any](r *http.Request, v *T) error {
	if err := DecodeJSON(r, v); err != nil {
		return err
	}
	if errs := validation.ValidateStruct(v); errs.HasErrors() {
		return errs
	}
	return nil
}

// BindPathParamUUID extracts and parses a UUID path parameter.
func BindPathParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	value := r.PathValue(name)
	if value == "" {
		return uuid.Nil, fmt.Errorf("path parameter %s is required", name)
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID for parameter %s: %w", name, err)
	}
	return id, nil
}

// BindPathParamString extracts a string path parameter.
func BindPathParamString(r *http.Request, name string) (string, error) {
	value := r.PathValue(name)
	if value == "" {
		return "", fmt.Errorf("path parameter %s is required", name)
	}
	return value, nil
}

// BindPathParamInt extracts and parses an integer path parameter.
func BindPathParamInt(r *http.Request, name string) (int, error) {
	value := r.PathValue(name)
	if value == "" {
		return 0, fmt.Errorf("path parameter %s is required", name)
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for parameter %s: %w", name, err)
	}
	return i, nil
}

// BindPathParam uses the oapi-codegen runtime to bind a styled path parameter.
// This provides compatibility with the existing parameter binding approach.
func BindPathParam[T any](r *http.Request, name string, target *T) error {
	return runtime.BindStyledParameterWithOptions(
		"simple",
		name,
		r.PathValue(name),
		target,
		runtime.BindStyledParameterOptions{
			ParamLocation: runtime.ParamLocationPath,
			Explode:       false,
			Required:      true,
		},
	)
}

// BindQueryParam uses the oapi-codegen runtime to bind a query parameter with deepObject style.
func BindQueryParam[T any](r *http.Request, name string, target *T) error {
	return runtime.BindQueryParameter(
		"deepObject",
		true,  // explode
		false, // required
		name,
		r.URL.Query(),
		target,
	)
}

// RequiredHeader extracts a required header value.
func RequiredHeader(r *http.Request, name string) (string, error) {
	value := r.Header.Get(name)
	if value == "" {
		return "", fmt.Errorf("missing required header: %s", name)
	}
	return value, nil
}

// OptionalHeader extracts an optional header value.
func OptionalHeader(r *http.Request, name string) string {
	return r.Header.Get(name)
}
