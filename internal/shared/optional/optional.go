// Package optional provides a generic Optional type for handling nullable/optional fields in Go structs.
package optional

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// Optional is a generic type that represents a field that can be either:
// - absent (not provided in the request)
// - present (provided in the request, with any value including null)
//
// This is a two-state system, suitable for cases where you need to distinguish
// between "field not sent" vs "field sent" but don't need to separately track
// explicit null values.
//
// Use with Go 1.24+ and the `omitzero` JSON tag for automatic omission when absent.
//
// Example:
//
//	type UpdateUser struct {
//	    Email Optional[string] `json:"email,omitzero"`
//	}
type Optional[T any] struct {
	Val     T
	Present bool
}

// NewOptional creates an Optional with a value present
func NewOptional[T any](val T) Optional[T] {
	return Optional[T]{
		Val:     val,
		Present: true,
	}
}

// NewAbsent creates an Optional that is explicitly absent
func NewAbsent[T any]() Optional[T] {
	return Optional[T]{
		Val:     *new(T),
		Present: false,
	}
}

// Get retrieves the underlying value if present, otherwise returns an error
func (o Optional[T]) Get() (T, error) {
	if !o.Present {
		var zero T
		return zero, errors.New("optional value is not present")
	}
	return o.Val, nil
}

// MustGet retrieves the underlying value if present, otherwise panics
func (o Optional[T]) MustGet() T {
	if !o.Present {
		panic("optional value is not present")
	}
	return o.Val
}

// GetOr returns the value if present, otherwise returns the provided default
func (o Optional[T]) GetOr(defaultVal T) T {
	if !o.Present {
		return defaultVal
	}
	return o.Val
}

// IsPresent returns true if the field was provided in the request
func (o Optional[T]) IsPresent() bool {
	return o.Present
}

// IsAbsent returns true if the field was not provided in the request
func (o Optional[T]) IsAbsent() bool {
	return !o.Present
}

// Set sets the value and marks it as present
func (o *Optional[T]) Set(val T) {
	o.Val = val
	o.Present = true
}

// Clear marks the value as absent
func (o *Optional[T]) Clear() {
	var zero T
	o.Val = zero
	o.Present = false
}

// IsZero implements the interface used by omitzero
// When Present is false, this returns true, causing JSON marshaling to omit the field
func (o Optional[T]) IsZero() bool {
	return !o.Present
}

// MarshalJSON implements json.Marshaler
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Present {
		// This shouldn't normally be called when using omitzero,
		// but we handle it gracefully
		return []byte("null"), nil
	}

	data, err := json.Marshal(o.Val)
	if err != nil {
		return nil, fmt.Errorf("optional: couldn't marshal JSON: %w", err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	// If this method is called, the field was present in the JSON
	if err := json.Unmarshal(data, &o.Val); err != nil {
		return fmt.Errorf("optional: couldn't unmarshal JSON: %w", err)
	}

	o.Present = true
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (o *Optional[T]) UnmarshalYAML(node *yaml.Node) error {
	// If this method is called, the field was present in the YAML
	if err := node.Decode(&o.Val); err != nil {
		return fmt.Errorf("optional: couldn't unmarshal YAML: %w", err)
	}

	o.Present = true
	return nil
}

// MarshalYAML implements yaml.Marshaler
func (o Optional[T]) MarshalYAML() (interface{}, error) {
	if !o.Present {
		// Return nil to trigger omitempty behavior
		return nil, nil
	}
	return o.Val, nil
}
