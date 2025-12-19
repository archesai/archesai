package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScalarHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name             string
		title            string
		specURL          string
		wantStatus       int
		wantTitle        string
		wantSpecURL      string
		wantScalarScript bool
	}{
		{
			name:             "serves HTML with correct title and spec URL",
			title:            "My API",
			specURL:          "/openapi.yaml",
			wantStatus:       http.StatusOK,
			wantTitle:        "My API - API Reference",
			wantSpecURL:      "/openapi.yaml",
			wantScalarScript: true,
		},
		{
			name:             "handles custom spec URL",
			title:            "Custom API",
			specURL:          "/api/v1/spec.json",
			wantStatus:       http.StatusOK,
			wantTitle:        "Custom API - API Reference",
			wantSpecURL:      "/api/v1/spec.json",
			wantScalarScript: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewScalarHandler(tt.title, tt.specURL)
			req := httptest.NewRequest(http.MethodGet, "/docs", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			bodyStr := string(body)
			assert.Contains(t, bodyStr, tt.wantTitle)
			assert.Contains(t, bodyStr, tt.wantSpecURL)

			if tt.wantScalarScript {
				assert.Contains(t, bodyStr, "cdn.jsdelivr.net/npm/@scalar/api-reference")
				assert.Contains(t, bodyStr, "Scalar.createApiReference")
			}
		})
	}
}

func TestOpenAPISpecHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name            string
		spec            []byte
		contentType     string
		wantStatus      int
		wantContentType string
	}{
		{
			name:            "serves YAML spec",
			spec:            []byte("openapi: 3.1.0\ninfo:\n  title: Test API"),
			contentType:     "application/yaml",
			wantStatus:      http.StatusOK,
			wantContentType: "application/yaml",
		},
		{
			name:            "serves JSON spec",
			spec:            []byte(`{"openapi": "3.1.0", "info": {"title": "Test API"}}`),
			contentType:     "application/json",
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewOpenAPISpecHandler(tt.spec, tt.contentType)
			req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, tt.wantContentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, "public, max-age=3600", resp.Header.Get("Cache-Control"))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.spec, body)
		})
	}
}

func TestRegisterScalarRoutes(t *testing.T) {
	spec := []byte("openapi: 3.1.0\ninfo:\n  title: Test API")
	mux := http.NewServeMux()

	RegisterScalarRoutes(mux, spec, "Test API")

	t.Run("docs endpoint works", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), "Test API - API Reference")
	})

	t.Run("openapi.yaml endpoint works", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))
		assert.Equal(t, string(spec), rec.Body.String())
	})
}

func TestSprintf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []string
		want   string
	}{
		{
			name:   "single substitution",
			format: "Hello %s!",
			args:   []string{"World"},
			want:   "Hello World!",
		},
		{
			name:   "multiple substitutions",
			format: "%s says %s",
			args:   []string{"Alice", "hello"},
			want:   "Alice says hello",
		},
		{
			name:   "no substitutions",
			format: "No placeholders",
			args:   []string{},
			want:   "No placeholders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sprintf(tt.format, tt.args...)
			assert.Equal(t, tt.want, got)
		})
	}
}
