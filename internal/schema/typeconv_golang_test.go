package schema

import "testing"

func TestToGoType(t *testing.T) {
	tests := []struct {
		name       string
		schemaType string
		format     string
		want       string
	}{
		// String types
		{"string no format", TypeString, "", GoTypeString},
		{"string with email", TypeString, FormatEmail, GoTypeString},
		{"string with uri", TypeString, FormatURI, GoTypeString},
		{"string with hostname", TypeString, FormatHostname, GoTypeString},
		{"string with password", TypeString, FormatPassword, GoTypeString},
		{"string with uuid", TypeString, FormatUUID, GoTypeUUID},
		{"string with date-time", TypeString, FormatDateTime, GoTypeTime},
		{"string with date", TypeString, FormatDate, GoTypeTime},

		// Integer types
		{"integer no format", TypeInteger, "", GoTypeInt},
		{"integer with int32", TypeInteger, FormatInt32, GoTypeInt32},
		{"integer with int64", TypeInteger, FormatInt64, GoTypeInt64},

		// Number types
		{"number no format", TypeNumber, "", GoTypeFloat64},
		{"number with float", TypeNumber, FormatFloat, GoTypeFloat32},
		{"number with double", TypeNumber, FormatDouble, GoTypeFloat64},

		// Other types
		{"boolean", TypeBoolean, "", GoTypeBool},
		{"array", TypeArray, "", GoTypeSliceAny},
		{"object", TypeObject, "", GoTypeMapString},

		// Unknown type
		{"unknown type", "unknown", "", GoTypeInterface},
		{"empty type", "", "", GoTypeInterface},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToGoType(tt.schemaType, tt.format)
			if got != tt.want {
				t.Errorf(
					"ToGoType(%q, %q) = %q, want %q",
					tt.schemaType,
					tt.format,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestStringFormatToGoType(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{FormatDateTime, GoTypeTime},
		{FormatDate, GoTypeTime},
		{FormatUUID, GoTypeUUID},
		{FormatEmail, GoTypeString},
		{FormatURI, GoTypeString},
		{FormatHostname, GoTypeString},
		{FormatPassword, GoTypeString},
		{"", GoTypeString},
		{"custom", GoTypeString},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got := stringFormatToGoType(tt.format)
			if got != tt.want {
				t.Errorf("stringFormatToGoType(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

func TestIntegerFormatToGoType(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{FormatInt32, GoTypeInt32},
		{FormatInt64, GoTypeInt64},
		{"", GoTypeInt},
		{"custom", GoTypeInt},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got := integerFormatToGoType(tt.format)
			if got != tt.want {
				t.Errorf("integerFormatToGoType(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

func TestNumberFormatToGoType(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{FormatFloat, GoTypeFloat32},
		{FormatDouble, GoTypeFloat64},
		{"", GoTypeFloat64},
		{"custom", GoTypeFloat64},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got := numberFormatToGoType(tt.format)
			if got != tt.want {
				t.Errorf("numberFormatToGoType(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

func TestGoTypeConstants(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{GoTypeInterface, "any"},
		{GoTypeString, "string"},
		{GoTypeInt, "int"},
		{GoTypeInt32, "int32"},
		{GoTypeInt64, "int64"},
		{GoTypeFloat32, "float32"},
		{GoTypeFloat64, "float64"},
		{GoTypeBool, "bool"},
		{GoTypeTime, "time.Time"},
		{GoTypeUUID, "uuid.UUID"},
		{GoTypeMapString, "map[string]any"},
		{GoTypeSliceAny, "[]any"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("Go type constant = %q, want %q", tt.constant, tt.expected)
		}
	}
}

func TestDefaultConstraints(t *testing.T) {
	// Verify default constraints are sensible values
	if DefaultStringMinLength != 0 {
		t.Errorf("DefaultStringMinLength = %d, want 0", DefaultStringMinLength)
	}
	if DefaultStringMaxLength <= 0 {
		t.Errorf("DefaultStringMaxLength = %d, should be positive", DefaultStringMaxLength)
	}
	if DefaultUUIDLength != 36 {
		t.Errorf("DefaultUUIDLength = %d, want 36", DefaultUUIDLength)
	}
	if DefaultIntegerMinimum != 0 {
		t.Errorf("DefaultIntegerMinimum = %d, want 0", DefaultIntegerMinimum)
	}
	if DefaultArrayMaxItems <= 0 {
		t.Errorf("DefaultArrayMaxItems = %d, should be positive", DefaultArrayMaxItems)
	}
}
