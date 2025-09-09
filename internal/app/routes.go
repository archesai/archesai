package app

import (
	"github.com/archesai/archesai/internal/auth"
	// "github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/content"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/users"
	"github.com/archesai/archesai/internal/workflows"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func (a *App) RegisterRoutes(e *echo.Echo) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes using StrictHandler pattern with rate limiting
	strictAuthHandler := auth.NewAuthStrictHandler(a.AuthHandler)

	// Apply rate limiting to the entire auth group
	// This will limit registration, login, and password reset endpoints
	authGroup := v1.Group("")
	authGroup.Use(auth.RateLimitMiddleware(10, 60)) // 10 requests per minute
	auth.RegisterHandlers(authGroup, strictAuthHandler)

	// Register users routes using StrictHandler pattern
	strictUsersHandler := users.NewUserStrictHandler(a.UsersHandler)
	users.RegisterHandlers(v1, strictUsersHandler)

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
	// TODO: Implement config handler when needed
	// strictConfigHandler := config.NewStrictHandler(a.ConfigHandler, nil)
	// config.RegisterHandlers(v1, strictConfigHandler)
}
