// Package validation provides validation utilities for struct validation.
package validation

import (
	"fmt"
	"strings"
)

// FieldError represents a single field validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Errors is a collection of validation errors.
type Errors []FieldError

// Error implements the error interface.
func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return fmt.Sprintf("validation failed: %s %s", e[0].Field, e[0].Message)
	}

	var msgs []string
	for _, err := range e {
		msgs = append(msgs, fmt.Sprintf("%s %s", err.Field, err.Message))
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, ", "))
}

// HasErrors returns true if there are any validation errors.
func (e Errors) HasErrors() bool {
	return len(e) > 0
}

// Add adds a new field error with the given field name and message.
func (e *Errors) Add(field, message string) {
	*e = append(*e, FieldError{Field: field, Message: message})
}

// AddWithCode adds a new field error with the given field name, message, and error code.
func (e *Errors) AddWithCode(field, message, code string) {
	*e = append(*e, FieldError{Field: field, Message: message, Code: code})
}

// Merge combines another Errors collection into this one.
func (e *Errors) Merge(other Errors) {
	*e = append(*e, other...)
}
