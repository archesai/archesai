// Package http provides HTTP server implementation and middleware
package http

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed docs.html
var docsHTML []byte

// SetupDocs configures API documentation endpoints with the provided OpenAPI spec
func (s *Server) SetupDocs(openapiSpec interface{}) error {
	if !s.config.DocsEnabled {
		s.logger.Info("API documentation disabled")
		return nil
	}

	// Validate that we can marshal the spec
	specJSON, err := json.Marshal(openapiSpec)
	if err != nil {
		return err
	}

	// Serve OpenAPI spec as JSON
	s.echo.GET("/openapi.json", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, specJSON)
	})

	// Serve OpenAPI spec as YAML (Scalar handles JSON fine even on .yaml endpoint)
	s.echo.GET("/openapi.yaml", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, specJSON)
	})

	// Serve the Scalar documentation UI
	s.echo.GET("/docs", func(c echo.Context) error {
		return c.HTMLBlob(http.StatusOK, docsHTML)
	})

	// Redirect root /docs/ to /docs
	s.echo.GET("/docs/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/docs")
	})

	s.logger.Info("API documentation enabled",
		"docs_url", "/docs",
		"spec_url", "/openapi.yaml",
	)

	return nil
}
