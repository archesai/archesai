package app

import (
	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/content"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/workflows"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func (a *App) RegisterRoutes(e *echo.Echo) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes using StrictHandler pattern
	strictAuthHandler := auth.NewAuthStrictHandlerWithMiddleware(a.AuthHandler)
	auth.RegisterHandlers(v1, strictAuthHandler)

	// Register organizations routes using StrictHandler pattern
	strictOrganizationsHandler := organizations.NewOrganizationStrictHandler(a.OrganizationsHandler)
	organizations.RegisterHandlers(v1, strictOrganizationsHandler)

	// Register workflows routes using StrictHandler pattern
	strictWorkflowsHandler := workflows.NewWorkflowStrictHandler(a.WorkflowsHandler)
	workflows.RegisterHandlers(v1, strictWorkflowsHandler)

	// Register content routes using StrictHandler pattern
	strictContentHandler := content.NewContentStrictHandler(a.ContentHandler)
	content.RegisterHandlers(v1, strictContentHandler)

	// Register health routes using StrictHandler pattern
	strictHealthHandler := health.NewStrictHandler(a.HealthHandler, nil)
	health.RegisterHandlers(v1, strictHealthHandler)

	// Register config routes using StrictHandler pattern
	strictConfigHandler := config.NewStrictHandler(a.ConfigHandler, nil)
	config.RegisterHandlers(v1, strictConfigHandler)
}
