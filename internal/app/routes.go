package app

import (
	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/labels"
	"github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/oauth"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/sessions"
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
	strictAccountsHandler := accounts.NewStrictHandler(a.AccountsHandler, nil)
	accounts.RegisterHandlers(v1, strictAccountsHandler)

	// Sessions routes - public endpoints for authentication
	// These endpoints handle login/logout
	a.Logger.Info("registering session routes (public)")
	strictSessionsHandler := sessions.NewStrictHandler(a.SessionsHandler, nil)
	sessions.RegisterHandlers(v1, strictSessionsHandler)

	// Magic Link routes - public endpoints for passwordless authentication
	a.Logger.Info("registering magic link routes (public)")
	v1.POST("/auth/magic-links/request", a.MagicLinkHandler.RequestMagicLink)
	v1.POST("/auth/magic-links/verify", a.MagicLinkHandler.VerifyMagicLink)

	// OAuth routes - public endpoints for OAuth authentication
	// Note: OAuth routes don't use the v1 prefix as per the OpenAPI spec
	a.Logger.Info("registering OAuth routes (public)")
	strictOAuthHandler := oauth.NewStrictHandler(a.OAuthHandler, nil)
	oauth.RegisterHandlers(e, strictOAuthHandler)

	// Health routes - public endpoints for monitoring
	a.Logger.Info("registering health routes (public)")
	strictHealthHandler := health.NewStrictHandler(a.HealthHandler, nil)
	health.RegisterHandlers(v1, strictHealthHandler)

	// ========================================
	// PROTECTED ROUTES (Authentication required)
	// ========================================

	// Create a protected group with authentication middleware
	a.Logger.Info("creating protected routes group")
	protected := v1.Group("")
	protected.Use(a.AuthMiddleware.RequireAuth())

	// Users routes - require authentication
	a.Logger.Info("registering users routes (protected)")
	strictUsersHandler := users.NewStrictHandler(a.UsersHandler, nil)
	users.RegisterHandlers(protected, strictUsersHandler)

	// Organizations routes - require authentication and organization membership
	strictOrganizationsHandler := organizations.NewStrictHandler(a.OrganizationsHandler, nil)
	// Create a separate group for organization routes with additional checks
	orgGroup := v1.Group("")
	orgGroup.Use(a.AuthMiddleware.RequireAuth(), a.AuthMiddleware.RequireOrganizationMember())
	organizations.RegisterHandlers(orgGroup, strictOrganizationsHandler)

	// Pipelines routes - require authentication
	a.Logger.Info("registering pipelines routes (protected)")
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

	// Config routes - require authentication
	a.Logger.Info("registering config routes (protected)")
	strictConfigHandler := config.NewStrictHandler(a.ConfigHandler, nil)
	config.RegisterHandlers(protected, strictConfigHandler)

	a.Logger.Info("all routes registered successfully")
}
