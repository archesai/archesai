// Package app provides dependency injection and application container management
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/artifacts"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/invitations"
	"github.com/archesai/archesai/internal/labels"
	"github.com/archesai/archesai/internal/magiclink"
	"github.com/archesai/archesai/internal/members"
	"github.com/archesai/archesai/internal/middleware"
	"github.com/archesai/archesai/internal/migrations"
	"github.com/archesai/archesai/internal/oauth"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/pipelines"
	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/server"
	"github.com/archesai/archesai/internal/sessions"
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
	SessionsService      *sessions.Service
	OAuthService         *oauth.Service
	MagicLinkService     *magiclink.Service

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
	SessionsHandler      sessions.StrictServerInterface
	OAuthHandler         oauth.StrictServerInterface
	ConfigHandler        config.StrictServerInterface
	MagicLinkHandler     *magiclink.Handler
}

// NewApp creates and initializes all application dependencies.
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize infrastructure
	infra, err := NewInfrastructure(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	log := infra.Logger

	// Run migrations if enabled
	if cfg.Database.RunMigrations {
		log.Info("running database migrations")
		if err := migrations.RunMigrations(infra.Database.SQLDB(), infra.Database.TypeString(), log); err != nil {
			log.Error("failed to run migrations", "error", err)
			isProduction := cfg.API.Environment == "production"
			if isProduction {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		}
		log.Info("database migrations completed")
	}

	// Create repositories
	log.Info("creating repositories")
	repos, err := NewRepositories(infra)
	if err != nil {
		return nil, fmt.Errorf("failed to create repositries: %w", err)
	}

	// Initialize authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.Local.JwtSecret, log)

	// Create app instance to populate
	app := &App{
		// Infrastructure
		infra:  infra,
		Logger: log,
		Config: cfg,

		// Middleware
		AuthMiddleware: authMiddleware,
	}

	// Initialize all domains in parallel where possible
	// Group 1: Independent domains (no dependencies on other services)
	g, _ := errgroup.WithContext(context.Background())

	// Initialize accounts domain
	g.Go(func() error {
		log.Info("initializing accounts domain")
		app.AccountsService = accounts.NewService(repos.Accounts, nil, log)
		app.AccountsHandler = accounts.NewStrictServer(app.AccountsService, log)
		log.Info("accounts domain ready")
		return nil
	})

	// Initialize users domain
	g.Go(func() error {
		log.Info("initializing users domain")
		app.UsersService = users.NewService(repos.Users, nil, log)
		app.UsersHandler = users.NewStrictServer(app.UsersService, log)
		log.Info("users domain ready")
		return nil
	})

	// Initialize organizations domain
	g.Go(func() error {
		log.Info("initializing organizations domain")
		app.OrganizationsService = organizations.NewService(repos.Organizations, nil, log)
		app.OrganizationsHandler = organizations.NewStrictServer(app.OrganizationsService, log)
		log.Info("organizations domain ready")
		return nil
	})

	// Initialize pipelines domain
	g.Go(func() error {
		log.Info("initializing pipelines domain")
		app.PipelinesService = pipelines.NewService(repos.Pipelines, nil, log)
		app.PipelinesHandler = pipelines.NewHandler(app.PipelinesService, log)
		log.Info("pipelines domain ready")
		return nil
	})

	// Initialize runs domain
	g.Go(func() error {
		log.Info("initializing runs domain")
		app.RunsService = runs.NewService(repos.Runs, nil, log)
		app.RunsHandler = runs.NewStrictServer(app.RunsService, log)
		log.Info("runs domain ready")
		return nil
	})

	// Initialize tools domain
	g.Go(func() error {
		log.Info("initializing tools domain")
		app.ToolsService = tools.NewService(repos.Tools, nil, log)
		app.ToolsHandler = tools.NewStrictServer(app.ToolsService, log)
		log.Info("tools domain ready")
		return nil
	})

	// Initialize artifacts domain
	g.Go(func() error {
		log.Info("initializing artifacts domain")
		app.ArtifactsService = artifacts.NewService(repos.Artifacts, nil, log)
		app.ArtifactsHandler = artifacts.NewStrictServer(app.ArtifactsService, log)
		log.Info("artifacts domain ready")
		return nil
	})

	// Initialize labels domain
	g.Go(func() error {
		log.Info("initializing labels domain")
		app.LabelsService = labels.NewService(repos.Labels, nil, log)
		app.LabelsHandler = labels.NewStrictServer(app.LabelsService, log)
		log.Info("labels domain ready")
		return nil
	})

	// Initialize members domain
	g.Go(func() error {
		log.Info("initializing members domain")
		app.MembersService = members.NewService(repos.Members, nil, log)
		app.MembersHandler = members.NewStrictServer(app.MembersService, log)
		log.Info("members domain ready")
		return nil
	})

	// Initialize invitations domain
	g.Go(func() error {
		log.Info("initializing invitations domain")
		app.InvitationsService = invitations.NewService(repos.Invitations, nil, log)
		app.InvitationsHandler = invitations.NewStrictServer(app.InvitationsService, log)
		log.Info("invitations domain ready")
		return nil
	})

	// Initialize health domain
	g.Go(func() error {
		log.Info("initializing health domain")
		app.HealthService = health.NewService(repos.Health, log)
		app.HealthHandler = health.NewHandler(app.HealthService, log)
		log.Info("health domain ready")
		return nil
	})

	// Initialize sessions domain
	g.Go(func() error {
		log.Info("initializing sessions domain")
		app.SessionsService = sessions.NewService(repos.Sessions, cfg.Auth.Local.JwtSecret, log)
		app.SessionsHandler = sessions.NewStrictServer(app.SessionsService, log)
		log.Info("sessions domain ready")
		return nil
	})

	// Initialize OAuth domain (needs to be after sessions and users)
	// So we do it after the wait

	// Wait for all parallel initializations to complete
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to initialize domains: %w", err)
	}

	// Initialize OAuth service (depends on sessions and users services)
	log.Info("initializing oauth domain")
	app.OAuthService = oauth.NewService(cfg, app.SessionsService, app.UsersService, log)
	app.OAuthHandler = oauth.NewStrictServer(app.OAuthService, log)
	log.Info("oauth domain ready")

	// Initialize MagicLink service (depends on sessions and users services)
	log.Info("initializing magic link domain")
	// Determine protocol based on environment
	protocol := "http"
	if cfg.API.Environment == "production" {
		protocol = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", protocol, cfg.API.Host)
	if cfg.API.Port != 0 && cfg.API.Port != 80 && cfg.API.Port != 443 {
		baseURL = fmt.Sprintf("%s:%d", baseURL, int(cfg.API.Port))
	}
	magicLinkRepo := magiclink.NewPostgresRepository(infra.Database.SQLDB())
	app.MagicLinkService = magiclink.NewService(magicLinkRepo, log, baseURL)
	app.MagicLinkHandler = magiclink.NewHandler(
		app.MagicLinkService,
		app.SessionsService,
		app.UsersService,
		log,
	)
	log.Info("magic link domain ready")

	// Initialize config handler
	log.Info("initializing config handler")
	app.ConfigHandler = config.NewStrictServer(cfg, log)
	log.Info("config handler ready")

	// Create the HTTP server
	log.Info("creating HTTP server")
	app.Server = server.NewServer(&cfg.API, log)

	// Register all application routes
	app.registerRoutes()

	log.Info("application initialized successfully")
	return app, nil
}

// Close cleans up all resources.
func (a *App) Close() error {
	a.Logger.Info("shutting down application")
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
		a.Logger.Info("enabling API documentation")
		if err := a.Server.SetupDocs(); err != nil {
			a.Logger.Error("failed to setup API docs", "error", err)
		}
	}

	// Register all application routes
	a.RegisterRoutes(e)
	a.Logger.Info("routes registered")
}

// readinessCheck checks if the service is ready to handle requests.
func (a *App) readinessCheck(ctx echo.Context) error {
	// Check database connection
	if err := a.infra.Database.SQLDB().PingContext(ctx.Request().Context()); err != nil {
		a.Logger.Error("database health check failed", "error", err)
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
