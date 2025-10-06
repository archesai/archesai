package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// SetupInfrastructureRoutes configures infrastructure routes only.
func (s *Server) SetupInfrastructureRoutes() {
	// Health check - simple liveness probe
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// API version endpoint
	s.echo.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version": "1.0.0",
			"build":   "development",
		})
	})

	// 404 handler - must be registered last (will be overridden when container registers routes)
	s.echo.RouteNotFound("/*", func(_ echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "route not found")
	})
}

// SetReadinessCheck allows the container to provide a readiness check function.
func (s *Server) SetReadinessCheck(checkFunc func(echo.Context) error) {
	s.echo.GET("/ready", checkFunc)
}
