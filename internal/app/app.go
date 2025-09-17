// Package app provides dependency injection and application container management
package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/invitations"
	"github.com/archesai/archesai/internal/labels"
	"github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/middleware"
	"github.com/archesai/archesai/internal/migrations"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/server"
	"github.com/archesai/archesai/internal/tools"
	"github.com/archesai/archesai/internal/users"
)

// App holds all application dependencies.
type App struct {
	// Core infrastructure
	infra *Infrastructure // Private infrastructure holder

	// Public infrastructure access
	Logger *slog.Logger
	Config *config.Config
	Server *server.Server

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Domain services (public for handler access)
	AccountsService      *accounts.Service
	UsersService         *users.Service
	OrganizationsService *organizations.Service
	InvitationsService   *invitations.Service
	ArtifactsService     *artifacts.Service
	LabelsService        *labels.Service
	MembersService       *members.Service
	PipelinesService     *pipelines.Service
	RunsService          *runs.Service
	ToolsService         *tools.Service
	HealthService        *health.Service

	// HTTP handlers
	AccountsHandler      accounts.StrictServerInterface
	UsersHandler         users.StrictServerInterface
	OrganizationsHandler organizations.StrictServerInterface
	InvitationsHandler   invitations.StrictServerInterface
	ArtifactsHandler     artifacts.StrictServerInterface
	LabelsHandler        labels.StrictServerInterface
	MembersHandler       members.StrictServerInterface
	PipelinesHandler     pipelines.StrictServerInterface
	RunsHandler          runs.StrictServerInterface
	ToolsHandler         tools.StrictServerInterface
	HealthHandler        *health.Handler
	// ConfigHandler        *config.Handler
}

// NewApp creates and initializes all application dependencies.
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
			isProduction := cfg.API.Environment == "production"
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

	// Initialize authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.Local.JwtSecret, log)

	// Initialize accounts domain
	// accountsEvents := accounts.NewEventPublisher(infra.EventPublisher)
	accountsService := accounts.NewService(repos.Accounts, nil, log)
	accountsHandler := accounts.NewStrictServer(accountsService, log)

	// Initialize users domain
	// usersEvents := users.NewEventPublisher(infra.EventPublisher)
	usersService := users.NewService(repos.Users, nil, log)
	usersHandler := users.NewStrictServer(usersService, log)

	// Initialize organizations domain
	// organizationsEvents := organizations.NewEventPublisher(infra.EventPublisher)
	organizationsService := organizations.NewService(repos.Organizations, nil, log)
	organizationsHandler := organizations.NewStrictServer(organizationsService, log)

	// Initialize pipelines domain
	pipelinesService := pipelines.NewService(repos.Pipelines, nil, log)
	pipelinesHandler := pipelines.NewHandler(pipelinesService, log)

	// Initialize runs domain
	runsService := runs.NewService(repos.Runs, nil, log)
	runsHandler := runs.NewStrictServer(runsService, log)

	// Initialize tools domain
	toolsService := tools.NewService(repos.Tools, nil, log)
	toolsHandler := tools.NewStrictServer(toolsService, log)

	// Initialize artifacts domain
	artifactsService := artifacts.NewService(repos.Artifacts, nil, log)
	artifactsHandler := artifacts.NewStrictServer(artifactsService, log)

	// Initialize labels domain
	labelsService := labels.NewService(repos.Labels, nil, log)
	labelsHandler := labels.NewStrictServer(labelsService, log)

	// Initialize members domain
	membersService := members.NewService(repos.Members, nil, log)
	membersHandler := members.NewStrictServer(membersService, log)

	// Initialize invitations domain
	invitationsService := invitations.NewService(repos.Invitations, nil, log)
	invitationsHandler := invitations.NewStrictServer(invitationsService, log)

	// Initialize health domain
	healthService := health.NewService(&infra.Database, infra.Logger)
	healthHandler := health.NewHandler(healthService, log)

	// Initialize config handler
	// TODO: Implement config handler when needed
	// configHandler := config.NewHandler(cfg, log)

	// Create the HTTP server
	httpServer := server.NewServer(&cfg.API, log)

	// Create app with all dependencies
	app := &App{
		// Infrastructure
		infra:  infra,
		Logger: log,
		Config: cfg,
		Server: httpServer,

		// Middleware
		AuthMiddleware: authMiddleware,

		// Domain services
		AccountsService:      accountsService,
		UsersService:         usersService,
		OrganizationsService: organizationsService,
		PipelinesService:     pipelinesService,
		RunsService:          runsService,
		ToolsService:         toolsService,
		ArtifactsService:     artifactsService,
		LabelsService:        labelsService,
		MembersService:       membersService,
		InvitationsService:   invitationsService,
		HealthService:        healthService,
		// ConfigService:     configService,

		// HTTP handlers
		AccountsHandler:      accountsHandler,
		UsersHandler:         usersHandler,
		OrganizationsHandler: organizationsHandler,
		PipelinesHandler:     pipelinesHandler,
		RunsHandler:          runsHandler,
		ToolsHandler:         toolsHandler,
		ArtifactsHandler:     artifactsHandler,
		LabelsHandler:        labelsHandler,
		MembersHandler:       membersHandler,
		InvitationsHandler:   invitationsHandler,
		HealthHandler:        healthHandler,
		// ConfigHandler:        configHandler,
	}

	// Register all application routes
	app.registerRoutes()

	return app, nil
}

// Close cleans up all resources.
func (a *App) Close() error {
	if a.infra != nil {
		return a.infra.Close()
	}
	return nil
}

// registerRoutes registers all application routes with the server.
func (a *App) registerRoutes() {
	// Get the echo instance from the server
	e := a.Server.Echo()

	// Register readiness check that can access the database
	a.Server.SetReadinessCheck(a.readinessCheck)

	// Setup API documentation if enabled
	if a.Config.API.Docs {
		if err := a.Server.SetupDocs(); err != nil {
			a.Logger.Error("failed to setup API docs", "error", err)
		}
	}

	// Register all application routes
	a.RegisterRoutes(e)
}

// readinessCheck checks if the service is ready to handle requests.
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
