package server

import (
	"net/http"
)

const scalarHTMLTemplate = `<!doctype html>
<html>
  <head>
    <title>%s - API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '%s',
      })
    </script>
  </body>
</html>
`

// ScalarHandler serves the Scalar API Reference UI.
type ScalarHandler struct {
	title   string
	specURL string
}

// NewScalarHandler creates a new ScalarHandler.
// title is used in the page title.
// specURL is the URL where the OpenAPI spec is served (e.g., "/openapi.yaml").
func NewScalarHandler(title, specURL string) *ScalarHandler {
	return &ScalarHandler{
		title:   title,
		specURL: specURL,
	}
}

// ServeHTTP serves the Scalar API Reference HTML page.
func (h *ScalarHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Override CSP to allow Scalar CDN and inline scripts
	w.Header().
		Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self' https:")
	w.WriteHeader(http.StatusOK)

	// Format the HTML with title and spec URL
	html := []byte(sprintf(scalarHTMLTemplate, h.title, h.specURL))
	_, _ = w.Write(html)
}

// sprintf is a local helper to avoid importing fmt just for Sprintf.
func sprintf(format string, args ...string) string {
	result := format
	for _, arg := range args {
		// Replace first %s with the argument
		for i := 0; i < len(result)-1; i++ {
			if result[i] == '%' && result[i+1] == 's' {
				result = result[:i] + arg + result[i+2:]
				break
			}
		}
	}
	return result
}

// OpenAPISpecHandler serves the bundled OpenAPI specification.
type OpenAPISpecHandler struct {
	spec        []byte
	contentType string
}

// NewOpenAPISpecHandler creates a new handler for serving the OpenAPI spec.
// spec is the pre-bundled OpenAPI specification bytes.
// contentType should be "application/yaml" or "application/json".
func NewOpenAPISpecHandler(spec []byte, contentType string) *OpenAPISpecHandler {
	return &OpenAPISpecHandler{
		spec:        spec,
		contentType: contentType,
	}
}

// ServeHTTP serves the OpenAPI specification.
func (h *OpenAPISpecHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", h.contentType)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.spec)
}

// RegisterScalarRoutes registers the Scalar API reference routes.
// It registers:
//   - GET /docs - The Scalar API Reference UI
//   - GET /openapi.yaml - The bundled OpenAPI specification
//
// spec is the pre-bundled OpenAPI specification in YAML format.
// title is used in the page title.
func RegisterScalarRoutes(mux *http.ServeMux, spec []byte, title string) {
	specHandler := NewOpenAPISpecHandler(spec, "application/yaml")
	mux.Handle("GET /openapi.yaml", specHandler)

	scalarHandler := NewScalarHandler(title, "/openapi.yaml")
	mux.Handle("GET /docs", scalarHandler)
}
