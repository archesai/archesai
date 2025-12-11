package validation

import (
	"testing"
)

func TestErrors_Error(t *testing.T) {
	tests := []struct {
		name   string
		errors Errors
		want   string
	}{
		{
			name:   "empty errors",
			errors: Errors{},
			want:   "",
		},
		{
			name: "single error",
			errors: Errors{
				{Field: "name", Message: "is required"},
			},
			want: "validation failed: name is required",
		},
		{
			name: "multiple errors",
			errors: Errors{
				{Field: "name", Message: "is required"},
				{Field: "email", Message: "must be a valid email address"},
			},
			want: "validation failed: name is required, email must be a valid email address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errors.Error(); got != tt.want {
				t.Errorf("Errors.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name   string
		errors Errors
		want   bool
	}{
		{
			name:   "empty errors",
			errors: Errors{},
			want:   false,
		},
		{
			name: "with errors",
			errors: Errors{
				{Field: "name", Message: "is required"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errors.HasErrors(); got != tt.want {
				t.Errorf("Errors.HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrors_Add(t *testing.T) {
	var errs Errors
	errs.Add("name", "is required")
	errs.Add("email", "must be valid")

	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}

	if errs[0].Field != "name" || errs[0].Message != "is required" {
		t.Errorf("first error mismatch: %+v", errs[0])
	}

	if errs[1].Field != "email" || errs[1].Message != "must be valid" {
		t.Errorf("second error mismatch: %+v", errs[1])
	}
}

func TestErrors_AddWithCode(t *testing.T) {
	var errs Errors
	errs.AddWithCode("name", "is required", "REQUIRED")

	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}

	if errs[0].Code != "REQUIRED" {
		t.Errorf("expected code REQUIRED, got %s", errs[0].Code)
	}
}

func TestErrors_Merge(t *testing.T) {
	var errs1 Errors
	errs1.Add("name", "is required")

	errs2 := Errors{
		{Field: "email", Message: "must be valid"},
		{Field: "age", Message: "must be positive"},
	}

	errs1.Merge(errs2)

	if len(errs1) != 3 {
		t.Errorf("expected 3 errors after merge, got %d", len(errs1))
	}
}
