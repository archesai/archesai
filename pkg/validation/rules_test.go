package validation

import (
	"testing"

	"github.com/archesai/archesai/internal/ptr"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		wantError bool
	}{
		{"nil value", nil, true},
		{"empty string", ptr.To(""), true},
		{"valid string", ptr.To("hello"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			Required(tt.value, "field", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		min       int
		wantError bool
	}{
		{"nil value", nil, 3, false},
		{"too short", ptr.To("ab"), 3, true},
		{"exact length", ptr.To("abc"), 3, false},
		{"longer", ptr.To("abcdef"), 3, false},
		{"unicode", ptr.To("日本"), 3, true},
		{"unicode valid", ptr.To("日本語"), 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			MinLength(tt.value, tt.min, "field", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		max       int
		wantError bool
	}{
		{"nil value", nil, 3, false},
		{"within limit", ptr.To("ab"), 3, false},
		{"exact length", ptr.To("abc"), 3, false},
		{"too long", ptr.To("abcdef"), 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			MaxLength(tt.value, tt.max, "field", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		wantError bool
	}{
		{"nil value", nil, false},
		{"empty string", ptr.To(""), false},
		{"valid email", ptr.To("test@example.com"), false},
		{"valid email with name", ptr.To("Test User <test@example.com>"), false},
		{"invalid email", ptr.To("not-an-email"), true},
		{"missing domain", ptr.To("test@"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			Email(tt.value, "email", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		wantError bool
	}{
		{"nil value", nil, false},
		{"empty string", ptr.To(""), false},
		{"valid uuid", ptr.To("550e8400-e29b-41d4-a716-446655440000"), false},
		{"invalid uuid", ptr.To("not-a-uuid"), true},
		{"invalid format", ptr.To("550e8400-e29b-41d4-a716"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			UUID(tt.value, "id", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name      string
		value     *int
		min       int
		wantError bool
	}{
		{"nil value", nil, 0, false},
		{"below min", ptr.To(-1), 0, true},
		{"at min", ptr.To(0), 0, false},
		{"above min", ptr.To(10), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			Min(tt.value, tt.min, "count", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name      string
		value     *int
		max       int
		wantError bool
	}{
		{"nil value", nil, 100, false},
		{"below max", ptr.To(50), 100, false},
		{"at max", ptr.To(100), 100, false},
		{"above max", ptr.To(150), 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			Max(tt.value, tt.max, "count", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestPattern(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		pattern   string
		wantError bool
	}{
		{"nil value", nil, "^[a-z]+$", false},
		{"empty string", ptr.To(""), "^[a-z]+$", false},
		{"matching pattern", ptr.To("abc"), "^[a-z]+$", false},
		{"not matching", ptr.To("ABC"), "^[a-z]+$", true},
		{"partial match not allowed", ptr.To("abc123"), "^[a-z]+$", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			Pattern(tt.value, tt.pattern, "field", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	tests := []struct {
		name      string
		value     *string
		allowed   []string
		wantError bool
	}{
		{"nil value", nil, []string{"a", "b", "c"}, false},
		{"valid value", ptr.To("b"), []string{"a", "b", "c"}, false},
		{"invalid value", ptr.To("d"), []string{"a", "b", "c"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			OneOf(tt.value, tt.allowed, "field", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		value     []string
		wantError bool
	}{
		{"nil slice", nil, true},
		{"empty slice", []string{}, true},
		{"non-empty slice", []string{"a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			NotEmpty(tt.value, "items", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestMinItems(t *testing.T) {
	tests := []struct {
		name      string
		value     []string
		min       int
		wantError bool
	}{
		{"below min", []string{"a"}, 2, true},
		{"at min", []string{"a", "b"}, 2, false},
		{"above min", []string{"a", "b", "c"}, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			MinItems(tt.value, tt.min, "items", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestMaxItems(t *testing.T) {
	tests := []struct {
		name      string
		value     []string
		max       int
		wantError bool
	}{
		{"below max", []string{"a"}, 2, false},
		{"at max", []string{"a", "b"}, 2, false},
		{"above max", []string{"a", "b", "c"}, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errs Errors
			MaxItems(tt.value, tt.max, "items", &errs)
			if tt.wantError && !errs.HasErrors() {
				t.Error("expected error, got none")
			}
			if !tt.wantError && errs.HasErrors() {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}
