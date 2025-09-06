package app

import (
	authhttp "github.com/archesai/archesai/internal/auth/adapters/http"
	confighttp "github.com/archesai/archesai/internal/config/adapters/http"
	contenthttp "github.com/archesai/archesai/internal/content/adapters/http"
	healthhttp "github.com/archesai/archesai/internal/health/adapters/http"
	organizationshttp "github.com/archesai/archesai/internal/organizations/adapters/http"
	workflowshttp "github.com/archesai/archesai/internal/workflows/adapters/http"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func (a *App) RegisterRoutes(e *echo.Echo) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes using StrictHandler pattern
	strictAuthHandler := authhttp.NewAuthStrictHandlerWithMiddleware(a.AuthHandler)
	authhttp.RegisterHandlers(v1, strictAuthHandler)

	// Register organizations routes using StrictHandler pattern
	strictOrganizationsHandler := organizationshttp.NewOrganizationStrictHandler(a.OrganizationsHandler)
	organizationshttp.RegisterHandlers(v1, strictOrganizationsHandler)

	// Register workflows routes using StrictHandler pattern
	strictWorkflowsHandler := workflowshttp.NewWorkflowStrictHandler(a.WorkflowsHandler)
	workflowshttp.RegisterHandlers(v1, strictWorkflowsHandler)

	// Register content routes using StrictHandler pattern
	strictContentHandler := contenthttp.NewContentStrictHandler(a.ContentHandler)
	contenthttp.RegisterHandlers(v1, strictContentHandler)

	// Register health routes using StrictHandler pattern
	strictHealthHandler := healthhttp.NewStrictHandler(a.HealthHandler, nil)
	healthhttp.RegisterHandlers(v1, strictHealthHandler)

	// Register config routes using StrictHandler pattern
	strictConfigHandler := confighttp.NewStrictHandler(a.ConfigHandler, nil)
	confighttp.RegisterHandlers(v1, strictConfigHandler)
}
