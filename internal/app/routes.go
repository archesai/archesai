package app

import (
	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/labels"
	"github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/tools"
	"github.com/archesai/archesai/internal/users"
)

// RegisterRoutes registers all application routes with the Echo server.
func (a *App) RegisterRoutes(e *echo.Echo) {
	a.Logger.Info("registering API routes...")

	// API v1 group
	v1 := e.Group("/api/v1")

	// ========================================
	// PUBLIC ROUTES (No authentication required)
	// ========================================

	// Account routes - public endpoints
	// These endpoints handle registration, password reset, email verification, etc.
	a.Logger.Info("registering account routes (public)")
	accounts.RegisterHandlers(v1, a.AccountsHandler)

	// Unified Auth routes - all authentication methods in one handler
	a.Logger.Info("registering unified auth routes (public)")
	auth.RegisterHandlers(v1, a.AuthHandler)

	// Health routes - public endpoints for monitoring
	a.Logger.Info("registering health routes (public)")
	health.RegisterHandlers(v1, a.HealthHandler)

	// Config routes - public endpoints (sanitized configuration)
	a.Logger.Info("registering config routes (public)")
	config.RegisterHandlers(v1, a.ConfigHandler)

	// ========================================
	// PROTECTED ROUTES (Authentication required)
	// ========================================

	// Create a protected group with authentication middleware
	a.Logger.Info("creating protected routes group")
	protected := v1.Group("")
	protected.Use(a.AuthMiddleware.RequireAuth())

	// Users routes - require authentication
	a.Logger.Info("registering users routes (protected)")
	users.RegisterHandlers(protected, a.UsersHandler)

	// Organizations routes - require authentication and organization membership
	// Create a separate group for organization routes with additional checks
	orgGroup := v1.Group("")
	orgGroup.Use(a.AuthMiddleware.RequireAuth(), a.AuthMiddleware.RequireOrganizationMember())
	organizations.RegisterHandlers(orgGroup, a.OrganizationsHandler)

	// Pipelines routes - require authentication
	a.Logger.Info("registering pipelines routes (protected)")
	pipelines.RegisterHandlers(protected, a.PipelinesHandler)

	// Runs routes - require authentication
	runs.RegisterHandlers(protected, a.RunsHandler)

	// Tools routes - require authentication
	tools.RegisterHandlers(protected, a.ToolsHandler)

	// Artifacts routes - require authentication
	artifacts.RegisterHandlers(protected, a.ArtifactsHandler)

	// Labels routes - require authentication
	labels.RegisterHandlers(protected, a.LabelsHandler)

	// Members routes - require authentication and organization membership
	members.RegisterHandlers(orgGroup, a.MembersHandler)

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

	a.Logger.Info("all routes registered successfully")
}
