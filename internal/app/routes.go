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

	// ========================================
	// PUBLIC ROUTES (No authentication required)
	// ========================================

	// Auth routes - public endpoints with rate limiting
	// These endpoints handle login, registration, password reset, etc.
	strictAuthHandler := auth.NewAuthStrictHandler(a.AuthHandler)
	publicAuthGroup := v1.Group("")
	publicAuthGroup.Use(auth.RateLimitMiddleware(10, 60)) // 10 requests per minute
	auth.RegisterHandlers(publicAuthGroup, strictAuthHandler)

	// Health routes - public endpoints for monitoring
	strictHealthHandler := health.NewStrictHandler(a.HealthHandler, nil)
	health.RegisterHandlers(v1, strictHealthHandler)

	// ========================================
	// PROTECTED ROUTES (Authentication required)
	// ========================================

	// Create a protected group with authentication middleware
	protected := v1.Group("")
	protected.Use(auth.RequireAuth(a.AuthService, auth.MiddlewarePresets.Authenticated, a.Logger))

	// Users routes - require authentication
	strictUsersHandler := users.NewUserStrictHandler(a.UsersHandler)
	users.RegisterHandlers(protected, strictUsersHandler)

	// Organizations routes - require authentication and organization membership
	strictOrganizationsHandler := organizations.NewOrganizationStrictHandler(a.OrganizationsHandler)
	// Create a separate group for organization routes with additional checks
	orgGroup := v1.Group("")
	orgGroup.Use(auth.RequireAuth(a.AuthService, auth.MiddlewarePresets.OrganizationMember, a.Logger))
	organizations.RegisterHandlers(orgGroup, strictOrganizationsHandler)

	// Workflows routes - require authentication
	strictWorkflowsHandler := workflows.NewWorkflowStrictHandler(a.WorkflowsHandler)
	workflows.RegisterHandlers(protected, strictWorkflowsHandler)

	// Content routes - require authentication
	strictContentHandler := content.NewContentStrictHandler(a.ContentHandler)
	content.RegisterHandlers(protected, strictContentHandler)

	// ========================================
	// API-ONLY ROUTES (API key authentication)
	// ========================================

	// Create an API-only group for endpoints that should use API keys
	// (Currently none, but structure is in place for future use)
	// apiOnly := v1.Group("")
	// apiOnly.Use(auth.RequireAuth(a.AuthService, auth.MiddlewarePresets.APIOnly, a.Logger))

	// ========================================
	// ADMIN ROUTES (Admin role required)
	// ========================================

	// Create an admin group for administrative endpoints
	// (Currently none, but structure is in place for future use)
	// admin := v1.Group("/admin")
	// admin.Use(auth.RequireAuth(a.AuthService, auth.MiddlewarePresets.AdminOnly, a.Logger))

	// Register config routes using StrictHandler pattern
	// TODO: Implement config handler when needed
	// strictConfigHandler := config.NewStrictHandler(a.ConfigHandler, nil)
	// config.RegisterHandlers(protected, strictConfigHandler)
}
