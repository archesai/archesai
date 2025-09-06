// Package server provides HTTP server implementation and middleware
package server

import (
	_ "embed"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

//go:embed assets/docs.html
var docsHTML []byte

// SetupDocs configures API documentation endpoints
func (s *Server) SetupDocs() error {
	if !s.config.Docs {
		s.logger.Info("API documentation disabled")
		return nil
	}

	// Determine the OpenAPI spec path
	// Try relative to current working directory first
	specPath := "api/openapi.bundled.yaml"
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		// Try relative to executable
		exePath, _ := os.Executable()
		specPath = filepath.Join(filepath.Dir(exePath), "../api/openapi.bundled.yaml")
		if _, err := os.Stat(specPath); os.IsNotExist(err) {
			// Fallback to absolute path
			specPath = "/home/jonathan/Projects/archesai/api/openapi.bundled.yaml"
		}
	}

	// Read the OpenAPI spec file
	openapiYAML, err := os.ReadFile(specPath)
	if err != nil {
		s.logger.Error("Failed to read OpenAPI spec", "error", err, "path", specPath)
		return err
	}

	// Serve OpenAPI spec as YAML
	s.echo.GET("/openapi.yaml", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/x-yaml")
		return c.Blob(http.StatusOK, "application/x-yaml", openapiYAML)
	})

	// Serve OpenAPI spec as JSON (convert YAML to JSON)
	s.echo.GET("/openapi.json", func(c echo.Context) error {
		var spec map[string]interface{}
		if err := yaml.Unmarshal(openapiYAML, &spec); err != nil {
			s.logger.Error("Failed to parse OpenAPI YAML", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to parse OpenAPI specification",
			})
		}
		return c.JSON(http.StatusOK, spec)
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
		"openapi_yaml", "/openapi.yaml",
		"openapi_json", "/openapi.json",
	)

	return nil
}
