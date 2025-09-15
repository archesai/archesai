// Package app provides dependency injection and application container management
package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	// "github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/health"

	// "github.com/archesai/archesai/internal/invitations"
	// "github.com/archesai/archesai/internal/labels"
	// "github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/migrations"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/server"
	"github.com/archesai/archesai/internal/tools"
	"github.com/archesai/archesai/internal/users"
	"github.com/labstack/echo/v4"
)

// App holds all application dependencies
type App struct {
	// Core infrastructure
	infra *Infrastructure // Private infrastructure holder

	// Public infrastructure access
	Logger *slog.Logger
	Config *config.Config
	Server *server.Server

	// Domain services (public for handler access)
	AuthService          *auth.Service
	UsersService         *users.Service
	OrganizationsService *organizations.Service
	// TODO: Implement these services
	// ArtifactsService     *artifacts.Service
	// LabelsService        *labels.Service
	// MembersService       *members.Service
	// InvitationsService   *invitations.Service
	PipelinesService *pipelines.Service
	RunsService      *runs.Service
	ToolsService     *tools.Service
	HealthService    *health.Service

	// HTTP handlers
	AuthHandler          *auth.Handler
	UsersHandler         *users.Handler
	OrganizationsHandler *organizations.Handler
	// TODO: Implement these handlers
	// ArtifactsHandler     *artifacts.Handler
	// LabelsHandler        *labels.Handler
	// MembersHandler       *members.Handler
	// InvitationsHandler   *invitations.Handler
	PipelinesHandler *pipelines.Handler
	RunsHandler      *runs.Handler
	ToolsHandler     *tools.Handler
	HealthHandler    *health.Handler
}

// NewApp creates and initializes all application dependencies
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize infrastructure
	infra, err := NewInfrastructure(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	log := infra.Logger
	db := infra.Database

	// Run migrations if enabled
	if cfg.Database.RunMigrations {
		if err := migrations.RunMigrations(db, log); err != nil {
			log.Error("failed to run migrations", "error", err)
			isProduction := cfg.Api.Environment == "production"
			if isProduction {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		}
	}

	// Create repositories
	repos, err := NewRepositories(db, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create repositories: %w", err)
	}

	// Initialize auth feature
	accessTokenTTL, err := cfg.GetAccessTokenTTL()
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token TTL: %w", err)
	}
	refreshTokenTTL, err := cfg.GetRefreshTokenTTL()
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token TTL: %w", err)
	}

	// Initialize auth service
	authConfig := auth.Config{
		JWTSecret:          cfg.GetJWTSecret(),
		AccessTokenExpiry:  accessTokenTTL,
		RefreshTokenExpiry: refreshTokenTTL,
	}

	// Always pass cache to auth service (infra.AuthCache is never nil)
	authService := auth.NewService(repos.Accounts, repos.Sessions, repos.Users, infra.AuthCache, authConfig, log)
	authHandler := auth.NewHandler(authService, log)

	// Initialize users domain
	usersEvents := users.NewEventPublisher(infra.EventPublisher)
	usersService := users.NewService(repos.Users, infra.UsersCache, usersEvents, log)
	usersHandler := users.NewHandler(usersService, log)

	// Initialize organizations domain
	organizationsService := organizations.NewService(repos.Organizations, log)
	organizationsHandler := organizations.NewHandler(organizationsService, log)

	// Initialize pipelines domain
	pipelinesService := pipelines.NewService(repos.Pipelines, log)
	pipelinesHandler := pipelines.NewHandler(pipelinesService, log)

	// Initialize runs domain
	runsService := runs.NewService(repos.Runs, log)
	runsHandler := runs.NewHandler(runsService, log)

	// Initialize tools domain
	toolsService := tools.NewService(repos.Tools, log)
	toolsHandler := tools.NewHandler(toolsService, log)

	// TODO: Initialize these domains when services/handlers are implemented
	// Initialize artifacts domain
	// artifactsService := artifacts.NewArtifactsService(repos.Artifacts, repos.Labels, log)
	// artifactsHandler := artifacts.NewHandler(artifactsService, log)

	// Initialize labels domain
	// labelsService := labels.NewService(repos.Labels, log)
	// labelsHandler := labels.NewHandler(labelsService, log)

	// Initialize members domain
	// membersService := members.NewService(repos.Members, log)
	// membersHandler := members.NewHandler(membersService, log)

	// Initialize invitations domain
	// invitationsService := invitations.NewService(repos.Invitations, log)
	// invitationsHandler := invitations.NewHandler(invitationsService, log)

	// Initialize health domain
	healthService := health.NewService(log)
	healthHandler := health.NewHandler(healthService, log)

	// Initialize config handler
	// TODO: Implement config handler when needed
	// configHandler := config.NewHandler(cfg, log)

	// Create the HTTP server
	httpServer := server.NewServer(&cfg.Api, log)

	// Create app with all dependencies
	app := &App{
		// Infrastructure
		infra:  infra,
		Logger: log,
		Config: cfg,
		Server: httpServer,

		// Domain services
		AuthService:          authService,
		UsersService:         usersService,
		OrganizationsService: organizationsService,
		PipelinesService:     pipelinesService,
		RunsService:          runsService,
		ToolsService:         toolsService,
		// ArtifactsService:     artifactsService,
		// LabelsService:        labelsService,
		// MembersService:       membersService,
		// InvitationsService:   invitationsService,
		HealthService: healthService,

		// HTTP handlers
		AuthHandler:          authHandler,
		UsersHandler:         usersHandler,
		OrganizationsHandler: organizationsHandler,
		PipelinesHandler:     pipelinesHandler,
		RunsHandler:          runsHandler,
		ToolsHandler:         toolsHandler,
		// ArtifactsHandler:     artifactsHandler,
		// LabelsHandler:        labelsHandler,
		// MembersHandler:       membersHandler,
		// InvitationsHandler:   invitationsHandler,
		HealthHandler: healthHandler,
	}

	// Register all application routes
	app.registerRoutes()

	return app, nil
}

// Close cleans up all resources
func (a *App) Close() error {
	if a.infra != nil {
		return a.infra.Close()
	}
	return nil
}

// registerRoutes registers all application routes with the server
func (a *App) registerRoutes() {
	// Get the echo instance from the server
	e := a.Server.Echo()

	// Register readiness check that can access the database
	a.Server.SetReadinessCheck(a.readinessCheck)

	// Setup API documentation if enabled
	if a.Config.Api.Docs {
		if err := a.Server.SetupDocs(); err != nil {
			a.Logger.Error("failed to setup API docs", "error", err)
		}
	}

	// Register all application routes
	a.RegisterRoutes(e)
}

// readinessCheck checks if the service is ready to handle requests
func (a *App) readinessCheck(ctx echo.Context) error {
	// Check database connection
	if err := a.infra.Database.Ping(ctx.Request().Context()); err != nil {
		a.Logger.Error("database health check failed")
		return ctx.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().Unix(),
	})
}
