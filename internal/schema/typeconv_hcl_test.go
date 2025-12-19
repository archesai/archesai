package schema

import "testing"

func TestSchema_HCLType(t *testing.T) {
	tests := []struct {
		name   string
		schema Schema
		want   string
	}{
		// Format-based types
		{"uuid format", Schema{Format: FormatUUID}, `sql("uuid")`},
		{"datetime format", Schema{Format: FormatDateTime}, `sql("timestamptz")`},
		{"date format", Schema{Format: FormatDate}, `sql("date")`},
		{"time format", Schema{Format: "time"}, `sql("time")`},
		{"email format", Schema{Format: FormatEmail}, `sql("text")`},
		{"uri format", Schema{Format: FormatURI}, `sql("text")`},
		{"hostname format", Schema{Format: FormatHostname}, `sql("text")`},
		{"ipv4 format", Schema{Format: "ipv4"}, `sql("text")`},
		{"binary format", Schema{Format: "binary"}, `sql("bytea")`},
		{"int32 format", Schema{Format: FormatInt32}, `sql("integer")`},
		{"int64 format", Schema{Format: FormatInt64}, `sql("bigint")`},
		{"float format", Schema{Format: FormatFloat}, `sql("real")`},
		{"double format", Schema{Format: FormatDouble}, `sql("double precision")`},

		// Type-based (no format)
		{"string type", Schema{Type: PropertyType{Types: []string{TypeString}}}, `sql("text")`},
		{
			"integer type",
			Schema{Type: PropertyType{Types: []string{TypeInteger}}},
			`sql("integer")`,
		},
		{
			"integer int64",
			Schema{Type: PropertyType{Types: []string{TypeInteger}}, Format: FormatInt64},
			`sql("bigint")`,
		},
		{"number type", Schema{Type: PropertyType{Types: []string{TypeNumber}}}, `sql("numeric")`},
		{
			"number float",
			Schema{Type: PropertyType{Types: []string{TypeNumber}}, Format: FormatFloat},
			`sql("real")`,
		},
		{
			"boolean type",
			Schema{Type: PropertyType{Types: []string{TypeBoolean}}},
			`sql("boolean")`,
		},
		{"object type", Schema{Type: PropertyType{Types: []string{TypeObject}}}, `sql("jsonb")`},

		// String with enum
		{
			"string with enum",
			Schema{Type: PropertyType{Types: []string{TypeString}}, Enum: []string{"a", "b"}},
			`sql("text")`,
		},

		// Unknown/default
		{"unknown type", Schema{Type: PropertyType{Types: []string{"unknown"}}}, `sql("text")`},
		{"empty type", Schema{}, `sql("text")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schema.HCLType()
			if got != tt.want {
				t.Errorf("HCLType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSchema_SQLiteHCLType(t *testing.T) {
	tests := []struct {
		name   string
		schema Schema
		want   string
	}{
		// Format-based types
		{"uuid format", Schema{Format: FormatUUID}, `sql("TEXT")`},
		{"datetime format", Schema{Format: FormatDateTime}, `sql("TEXT")`},
		{"date format", Schema{Format: FormatDate}, `sql("TEXT")`},
		{"time format", Schema{Format: "time"}, `sql("TEXT")`},
		{"email format", Schema{Format: FormatEmail}, `sql("TEXT")`},
		{"binary format", Schema{Format: "binary"}, `sql("BLOB")`},
		{"int32 format", Schema{Format: FormatInt32}, `sql("INTEGER")`},
		{"int64 format", Schema{Format: FormatInt64}, `sql("INTEGER")`},
		{"float format", Schema{Format: FormatFloat}, `sql("REAL")`},
		{"double format", Schema{Format: FormatDouble}, `sql("REAL")`},

		// Type-based
		{"string type", Schema{Type: PropertyType{Types: []string{TypeString}}}, `sql("TEXT")`},
		{
			"integer type",
			Schema{Type: PropertyType{Types: []string{TypeInteger}}},
			`sql("INTEGER")`,
		},
		{"number type", Schema{Type: PropertyType{Types: []string{TypeNumber}}}, `sql("REAL")`},
		{
			"boolean type",
			Schema{Type: PropertyType{Types: []string{TypeBoolean}}},
			`sql("INTEGER")`,
		},
		{"object type", Schema{Type: PropertyType{Types: []string{TypeObject}}}, `sql("TEXT")`},
		{"array type", Schema{Type: PropertyType{Types: []string{TypeArray}}}, `sql("TEXT")`},

		// Default
		{"empty type", Schema{}, `sql("TEXT")`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schema.SQLiteHCLType()
			if got != tt.want {
				t.Errorf("SQLiteHCLType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSchema_HCLDefault(t *testing.T) {
	tests := []struct {
		name   string
		schema Schema
		want   string
	}{
		// Nil default
		{"nil default", Schema{}, ""},

		// String defaults
		{"string default", Schema{Default: "hello"}, `"hello"`},
		{"sql expression", Schema{Default: `sql("NOW()")`}, `sql("NOW()")`},

		// Numeric defaults
		{"int default", Schema{Default: 42}, "42"},
		{"int32 default", Schema{Default: int32(42)}, "42"},
		{"int64 default", Schema{Default: int64(42)}, "42"},
		{"float32 default", Schema{Default: float32(3.14)}, "3.14"},
		{"float64 default", Schema{Default: 3.14159}, "3.14159"},

		// Boolean defaults
		{"bool true", Schema{Default: true}, "true"},
		{"bool false", Schema{Default: false}, "false"},

		// Special formats
		{
			"uuid format with default",
			Schema{Format: FormatUUID, Default: ""},
			`sql("gen_random_uuid()")`,
		},
		{
			"datetime format with default",
			Schema{Format: FormatDateTime, Default: ""},
			`sql("CURRENT_TIMESTAMP")`,
		},

		// Array defaults
		{
			"empty array",
			Schema{Type: PropertyType{Types: []string{TypeArray}}, Default: []any{}},
			`"{}"`,
		},
		{
			"array string default",
			Schema{Type: PropertyType{Types: []string{TypeArray}}, Default: "{}"},
			`"{}"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schema.HCLDefault()
			if got != tt.want {
				t.Errorf("HCLDefault() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSchema_SQLiteHCLDefault(t *testing.T) {
	tests := []struct {
		name   string
		schema Schema
		want   string
	}{
		// Nil default
		{"nil default", Schema{}, ""},

		// String defaults
		{"string default", Schema{Default: "hello"}, `"hello"`},
		{"sql expression", Schema{Default: `sql("NOW()")`}, `sql("NOW()")`},

		// Numeric defaults
		{"int default", Schema{Default: 42}, "42"},
		{"float64 default", Schema{Default: 3.14159}, "3.14159"},

		// Boolean defaults (SQLite uses 0/1)
		{"bool true", Schema{Default: true}, "1"},
		{"bool false", Schema{Default: false}, "0"},

		// Special formats
		{
			"uuid format with default",
			Schema{Format: FormatUUID, Default: ""},
			`sql("lower(hex(randomblob(16)))")`,
		},
		{
			"datetime format with default",
			Schema{Format: FormatDateTime, Default: ""},
			`sql("CURRENT_TIMESTAMP")`,
		},

		// Array defaults (SQLite uses JSON)
		{
			"empty array",
			Schema{Type: PropertyType{Types: []string{TypeArray}}, Default: []any{}},
			`"[]"`,
		},
		{
			"array string default",
			Schema{Type: PropertyType{Types: []string{TypeArray}}, Default: "{}"},
			`"[]"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schema.SQLiteHCLDefault()
			if got != tt.want {
				t.Errorf("SQLiteHCLDefault() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractYAMLValue(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		marker string
		want   string
	}{
		{"extract string value", "&{0 0 !!str hello }", "!!str ", "hello"},
		{"extract int value", "&{0 0 !!int 42 }", "!!int ", "42"},
		{"no marker", "no marker here", "!!str ", ""},
		{"empty after marker", "!!str ", "!!str ", ""},
		{"multiple words", "&{0 0 !!str hello world}", "!!str ", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractYAMLValue(tt.str, tt.marker)
			if got != tt.want {
				t.Errorf("extractYAMLValue(%q, %q) = %q, want %q", tt.str, tt.marker, got, tt.want)
			}
		})
	}
}

func TestHCLTypeConstants(t *testing.T) {
	// Verify PostgreSQL type constants
	pgTests := []struct {
		constant string
		expected string
	}{
		{hclTypeText, `sql("text")`},
		{hclTypeInteger, `sql("integer")`},
		{hclTypeBigint, `sql("bigint")`},
		{hclTypeBoolean, `sql("boolean")`},
		{hclTypeUUID, `sql("uuid")`},
		{hclTypeTimestamptz, `sql("timestamptz")`},
		{hclTypeDate, `sql("date")`},
		{hclTypeJSONB, `sql("jsonb")`},
	}

	for _, tt := range pgTests {
		if tt.constant != tt.expected {
			t.Errorf("constant = %q, want %q", tt.constant, tt.expected)
		}
	}

	// Verify SQLite type constants
	sqliteTests := []struct {
		constant string
		expected string
	}{
		{hclTypeSQLiteText, `sql("TEXT")`},
		{hclTypeSQLiteInteger, `sql("INTEGER")`},
		{hclTypeSQLiteReal, `sql("REAL")`},
		{hclTypeSQLiteBlob, `sql("BLOB")`},
	}

	for _, tt := range sqliteTests {
		if tt.constant != tt.expected {
			t.Errorf("constant = %q, want %q", tt.constant, tt.expected)
		}
	}
}
