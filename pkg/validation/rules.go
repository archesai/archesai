package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

// Required validates that a string pointer is not nil and not empty.
func Required(value *string, field string, errs *Errors) {
	if value == nil || *value == "" {
		errs.Add(field, "is required")
	}
}

// RequiredString validates that a string is not empty (for non-pointer fields).
func RequiredString(value string, field string, errs *Errors) {
	if value == "" {
		errs.Add(field, "is required")
	}
}

// RequiredUUID validates that a UUID is not the zero value.
func RequiredUUID(value uuid.UUID, field string, errs *Errors) {
	if value == uuid.Nil {
		errs.Add(field, "is required")
	}
}

// RequiredInt validates that an int pointer is not nil.
func RequiredInt[T constraints.Integer](value *T, field string, errs *Errors) {
	if value == nil {
		errs.Add(field, "is required")
	}
}

// MinLength validates that a string pointer has at least minLen characters.
func MinLength(value *string, minLen int, field string, errs *Errors) {
	if value == nil {
		return
	}
	if utf8.RuneCountInString(*value) < minLen {
		errs.Add(field, fmt.Sprintf("must be at least %d characters", minLen))
	}
}

// MinLengthString validates that a string has at least minLen characters (for non-pointer fields).
func MinLengthString(value string, minLen int, field string, errs *Errors) {
	if utf8.RuneCountInString(value) < minLen {
		errs.Add(field, fmt.Sprintf("must be at least %d characters", minLen))
	}
}

// MaxLength validates that a string pointer has at most maxLen characters.
func MaxLength(value *string, maxLen int, field string, errs *Errors) {
	if value == nil {
		return
	}
	if utf8.RuneCountInString(*value) > maxLen {
		errs.Add(field, fmt.Sprintf("must be at most %d characters", maxLen))
	}
}

// MaxLengthString validates that a string has at most maxLen characters (for non-pointer fields).
func MaxLengthString(value string, maxLen int, field string, errs *Errors) {
	if utf8.RuneCountInString(value) > maxLen {
		errs.Add(field, fmt.Sprintf("must be at most %d characters", maxLen))
	}
}

// Email validates that a string pointer is a valid email address.
func Email(value *string, field string, errs *Errors) {
	if value == nil || *value == "" {
		return
	}
	_, err := mail.ParseAddress(*value)
	if err != nil {
		errs.Add(field, "must be a valid email address")
	}
}

// EmailString validates that a string is a valid email address (for non-pointer fields).
func EmailString(value string, field string, errs *Errors) {
	if value == "" {
		return
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		errs.Add(field, "must be a valid email address")
	}
}

// UUID validates that a string pointer is a valid UUID.
func UUID(value *string, field string, errs *Errors) {
	if value == nil || *value == "" {
		return
	}
	_, err := uuid.Parse(*value)
	if err != nil {
		errs.Add(field, "must be a valid UUID")
	}
}

// UUIDString validates that a string is a valid UUID (for non-pointer fields).
func UUIDString(value string, field string, errs *Errors) {
	if value == "" {
		return
	}
	_, err := uuid.Parse(value)
	if err != nil {
		errs.Add(field, "must be a valid UUID")
	}
}

// Min validates that a numeric value is at least minVal.
func Min[T constraints.Ordered](value *T, minVal T, field string, errs *Errors) {
	if value == nil {
		return
	}
	if *value < minVal {
		errs.Add(field, fmt.Sprintf("must be at least %v", minVal))
	}
}

// Max validates that a numeric value is at most maxVal.
func Max[T constraints.Ordered](value *T, maxVal T, field string, errs *Errors) {
	if value == nil {
		return
	}
	if *value > maxVal {
		errs.Add(field, fmt.Sprintf("must be at most %v", maxVal))
	}
}

// Pattern validates that a string matches the given regular expression pattern.
func Pattern(value *string, pattern, field string, errs *Errors) {
	if value == nil || *value == "" {
		return
	}
	matched, err := regexp.MatchString(pattern, *value)
	if err != nil || !matched {
		errs.Add(field, fmt.Sprintf("must match pattern %s", pattern))
	}
}

// OneOf validates that a string value is one of the allowed values.
func OneOf[T comparable](value *T, allowed []T, field string, errs *Errors) {
	if value == nil {
		return
	}
	for _, a := range allowed {
		if *value == a {
			return
		}
	}
	errs.Add(field, "must be one of the allowed values")
}

// NotEmpty validates that a slice is not empty.
func NotEmpty[T any](value []T, field string, errs *Errors) {
	if len(value) == 0 {
		errs.Add(field, "must not be empty")
	}
}

// MinItems validates that a slice has at least minCount items.
func MinItems[T any](value []T, minCount int, field string, errs *Errors) {
	if len(value) < minCount {
		errs.Add(field, fmt.Sprintf("must have at least %d items", minCount))
	}
}

// MaxItems validates that a slice has at most maxCount items.
func MaxItems[T any](value []T, maxCount int, field string, errs *Errors) {
	if len(value) > maxCount {
		errs.Add(field, fmt.Sprintf("must have at most %d items", maxCount))
	}
}
