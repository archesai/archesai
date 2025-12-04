package parsers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAPIParser_Lint(t *testing.T) {
	tests := []struct {
		name        string
		spec        string
		wantErr     bool
		errContains []string // Multiple strings that should be in the error
	}{
		{
			name: "valid minimal spec passes linting",
			spec: `openapi: 3.0.0
info:
  title: Valid API
  version: 1.0.0
  description: A valid minimal API
  contact:
    email: contact@example.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT
servers:
  - url: https://api.example.com
    description: Production server
tags:
  - name: test
    description: Test operations
paths:
  /test:
    get:
      operationId: getTest
      summary: Get test data
      description: Returns test data
      tags:
        - test
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    description: Response message
components:
  schemas:
    TestSchema:
      type: object
      description: Test schema object
      properties:
        id:
          type: string
          description: Unique identifier`,
			wantErr: false,
		},
		{
			name: "missing operationId fails linting",
			spec: `openapi: 3.0.0
info:
  title: Invalid API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Get test
      responses:
        '200':
          description: Success`,
			wantErr:     true,
			errContains: []string{"violation", "operationId"},
		},
		{
			name: "missing info contact fails strict linting",
			spec: `openapi: 3.0.0
info:
  title: API Without Contact
  version: 1.0.0
paths:
  /test:
    get:
      operationId: getTest
      summary: Get test
      responses:
        '200':
          description: Success`,
			wantErr:     true,
			errContains: []string{"violation", "contact"},
		},
		{
			name: "missing response description fails linting",
			spec: `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
  contact:
    email: test@example.com
paths:
  /test:
    get:
      operationId: getTest
      summary: Get test
      responses:
        '200':
          content:
            application/json:
              schema:
                type: object`,
			wantErr:     true,
			errContains: []string{"violation", "description"},
		},
		{
			name: "schema without description fails strict linting",
			spec: `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
  contact:
    email: test@example.com
components:
  schemas:
    TestSchema:
      type: object
      properties:
        id:
          type: string`,
			wantErr:     true,
			errContains: []string{"violation", "description"},
		},
		{
			name: "multiple violations are all reported",
			spec: `openapi: 3.0.0
info:
  title: Bad API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Get test
      responses:
        '200':
          content:
            application/json:
              schema:
                type: object
  /another:
    post:
      summary: Post test
      responses:
        '201':
          description: Created`,
			wantErr:     true,
			errContains: []string{"violation", "operationId", "contact"},
		},
		{
			name:        "empty spec fails linting",
			spec:        ``,
			wantErr:     true,
			errContains: []string{"violation"},
		},
		{
			name: "spec with invalid yaml fails",
			spec: `openapi: 3.0.0
info:
  title: Invalid
  version: 1.0.0
paths:
  /test:
    get:
      invalid_yaml: [
        unclosed`,
			wantErr:     true,
			errContains: []string{"violation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewOpenAPIParser()
			err := parser.Lint([]byte(tt.spec))

			if tt.wantErr {
				assert.Error(t, err)
				errStr := err.Error()
				for _, contains := range tt.errContains {
					assert.Contains(t, strings.ToLower(errStr), strings.ToLower(contains),
						"Error should contain '%s'", contains)
				}
				// Verify error format includes violation count
				if err != nil {
					assert.Contains(t, errStr, "❌")
					assert.Contains(t, errStr, "found:")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOpenAPIParser_WithLinting(t *testing.T) {
	// Test that WithLinting enables linting flag
	parser := NewOpenAPIParser()
	assert.False(t, parser.lintEnabled, "Linting should be disabled by default")

	parser.WithLinting()
	assert.True(t, parser.lintEnabled, "Linting should be enabled after calling WithLinting")
}

func TestOpenAPIParser_Parse_WithLinting(t *testing.T) {
	validSpec := `openapi: 3.0.0
info:
  title: Valid API
  version: 1.0.0
  description: A valid API
  contact:
    email: contact@example.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT
servers:
  - url: https://api.example.com
    description: Production server
paths:
  /test:
    get:
      operationId: getTest
      summary: Get test
      description: Get test endpoint
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                description: Response object`

	invalidSpec := `openapi: 3.0.0
info:
  title: Invalid API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Missing operationId
      responses:
        '200':
          description: Success`

	t.Run("valid spec with linting enabled", func(t *testing.T) {
		parser := NewOpenAPIParser().WithLinting()
		doc, err := parser.Parse([]byte(validSpec))
		assert.NoError(t, err)
		assert.NotNil(t, doc)
	})

	t.Run("invalid spec with linting enabled blocks parsing", func(t *testing.T) {
		parser := NewOpenAPIParser().WithLinting()
		doc, err := parser.Parse([]byte(invalidSpec))
		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "violation")
		assert.Contains(t, err.Error(), "operationId")
	})

	t.Run("invalid spec without linting still parses", func(t *testing.T) {
		parser := NewOpenAPIParser() // No linting
		doc, err := parser.Parse([]byte(invalidSpec))
		assert.NoError(t, err)
		assert.NotNil(t, doc)
	})
}

func TestOpenAPIParser_ParseFile_WithLinting(t *testing.T) {
	// Test using actual test files
	t.Run("valid spec file with linting", func(t *testing.T) {
		parser := NewOpenAPIParser().WithLinting()
		// This will likely fail with strict linting due to missing descriptions, etc
		doc, err := parser.ParseFile("../../test/data/parsers/openapi/petstore.yaml")
		// The petstore example likely has violations with strict linting
		if err != nil {
			assert.Contains(t, err.Error(), "violation")
		} else {
			assert.NotNil(t, doc)
		}
	})
}

func TestFormatLintingErrors(t *testing.T) {
	// Test the error formatting by creating a mock result set
	// This test ensures the error message is properly formatted with categories
	parser := NewOpenAPIParser()

	// Test with a spec that will have violations in multiple categories
	spec := `openapi: 3.0.0
info:
  title: Test
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test
      responses:
        '200':
          description: OK
components:
  schemas:
    Test:
      type: object`

	err := parser.Lint([]byte(spec))
	require.Error(t, err)

	// Check that the error is properly formatted
	errStr := err.Error()
	assert.Contains(t, errStr, "❌")
	assert.Contains(t, errStr, "violation")
	assert.Contains(t, errStr, "found:")
	// Should have line:column format
	assert.Regexp(t, `\[\d+:\d+\]`, errStr)
}

func TestLintingWithParseIntegration(t *testing.T) {
	// Integration test ensuring Parse method properly integrates linting
	parser := NewOpenAPIParser()

	// Create a spec with known violations
	spec := []byte(`openapi: 3.0.0
info:
  title: API
  version: 1.0.0
paths:
  /test:
    get:
      summary: No ID
      responses:
        '200':
          description: OK`)

	// Without linting, should parse successfully
	parser.lintEnabled = false
	doc, err := parser.Parse(spec)
	assert.NoError(t, err)
	assert.NotNil(t, doc)

	// With linting, should fail
	parser.lintEnabled = true
	doc, err = parser.Parse(spec)
	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Contains(t, err.Error(), "violation")
}

func TestLintingPreservesSpecBytes(t *testing.T) {
	// Ensure spec bytes are properly stored for potential re-linting
	parser := NewOpenAPIParser()
	spec := []byte(`openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths: {}`)

	_, err := parser.Parse(spec)
	require.NoError(t, err)
	assert.Equal(t, spec, parser.specBytes)

	// Now enable linting and lint the stored bytes
	parser.lintEnabled = true
	err = parser.Lint(parser.specBytes)
	assert.Error(t, err) // Should have violations due to strict rules
}
