package app

import (
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/auth"

	// "github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/labels"
	"github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/tools"
	"github.com/archesai/archesai/internal/users"

	// "github.com/archesai/archesai/internal/workflows"
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
	strictOrganizationsHandler := organizations.NewStrictHandler(a.OrganizationsHandler, nil)
	// Create a separate group for organization routes with additional checks
	orgGroup := v1.Group("")
	orgGroup.Use(auth.RequireAuth(a.AuthService, auth.MiddlewarePresets.OrganizationMember, a.Logger))
	organizations.RegisterHandlers(orgGroup, strictOrganizationsHandler)

	// Pipelines routes - require authentication
	strictPipelinesHandler := pipelines.NewStrictHandler(a.PipelinesHandler, nil)
	pipelines.RegisterHandlers(protected, strictPipelinesHandler)

	// Runs routes - require authentication
	strictRunsHandler := runs.NewStrictHandler(a.RunsHandler, nil)
	runs.RegisterHandlers(protected, strictRunsHandler)

	// Tools routes - require authentication
	strictToolsHandler := tools.NewStrictHandler(a.ToolsHandler, nil)
	tools.RegisterHandlers(protected, strictToolsHandler)

	// Artifacts routes - require authentication
	strictArtifactsHandler := artifacts.NewStrictHandler(a.ArtifactsHandler, nil)
	artifacts.RegisterHandlers(protected, strictArtifactsHandler)

	// Labels routes - require authentication
	strictLabelsHandler := labels.NewStrictHandler(a.LabelsHandler, nil)
	labels.RegisterHandlers(protected, strictLabelsHandler)

	// Members routes - require authentication and organization membership
	strictMembersHandler := members.NewStrictHandler(a.MembersHandler, nil)
	members.RegisterHandlers(orgGroup, strictMembersHandler)

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
