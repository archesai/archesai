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
	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/auth/deliverers"
	"github.com/archesai/archesai/internal/auth/providers"
	"github.com/archesai/archesai/internal/auth/stores"
	"github.com/archesai/archesai/internal/auth/tokens"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/database/postgresql"
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
	AuthService          *auth.Service

	// HTTP handlers
	AccountsHandler      accounts.ServerInterface
	UsersHandler         users.ServerInterface
	OrganizationsHandler organizations.ServerInterface
	InvitationsHandler   invitations.ServerInterface
	ArtifactsHandler     artifacts.ServerInterface
	LabelsHandler        labels.ServerInterface
	MembersHandler       members.ServerInterface
	PipelinesHandler     pipelines.ServerInterface
	RunsHandler          runs.ServerInterface
	ToolsHandler         tools.ServerInterface
	HealthHandler        health.ServerInterface
	ConfigHandler        config.ServerInterface

	// Unified authentication handler
	AuthHandler auth.ServerInterface
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
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.Local.JWTSecret, cfg, log)

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
		app.AccountsHandler = accounts.NewStrictHandler(app.AccountsService, nil)
		log.Info("accounts domain ready")
		return nil
	})

	// Initialize users domain
	g.Go(func() error {
		log.Info("initializing users domain")
		app.UsersService = users.NewService(repos.Users, nil, log)
		app.UsersHandler = users.NewStrictHandler(app.UsersService, nil)
		log.Info("users domain ready")
		return nil
	})

	// Initialize organizations domain
	g.Go(func() error {
		log.Info("initializing organizations domain")
		app.OrganizationsService = organizations.NewService(repos.Organizations, nil, log)
		app.OrganizationsHandler = organizations.NewStrictHandler(app.OrganizationsService, nil)
		log.Info("organizations domain ready")
		return nil
	})

	// Initialize pipelines domain
	g.Go(func() error {
		log.Info("initializing pipelines domain")
		app.PipelinesService = pipelines.NewService(repos.Pipelines, nil, log)
		app.PipelinesHandler = pipelines.NewStrictHandler(app.PipelinesService, nil)
		log.Info("pipelines domain ready")
		return nil
	})

	// Initialize runs domain
	g.Go(func() error {
		log.Info("initializing runs domain")
		app.RunsService = runs.NewService(repos.Runs, nil, log)
		app.RunsHandler = runs.NewStrictHandler(app.RunsService, nil)
		log.Info("runs domain ready")
		return nil
	})

	// Initialize tools domain
	g.Go(func() error {
		log.Info("initializing tools domain")
		app.ToolsService = tools.NewService(repos.Tools, nil, log)
		app.ToolsHandler = tools.NewStrictHandler(app.ToolsService, nil)
		log.Info("tools domain ready")
		return nil
	})

	// Initialize artifacts domain
	g.Go(func() error {
		log.Info("initializing artifacts domain")
		app.ArtifactsService = artifacts.NewService(repos.Artifacts, nil, log)
		app.ArtifactsHandler = artifacts.NewStrictHandler(app.ArtifactsService, nil)
		log.Info("artifacts domain ready")
		return nil
	})

	// Initialize labels domain
	g.Go(func() error {
		log.Info("initializing labels domain")
		app.LabelsService = labels.NewService(repos.Labels, nil, log)
		app.LabelsHandler = labels.NewStrictHandler(app.LabelsService, nil)
		log.Info("labels domain ready")
		return nil
	})

	// Initialize members domain
	g.Go(func() error {
		log.Info("initializing members domain")
		app.MembersService = members.NewService(repos.Members, nil, log)
		app.MembersHandler = members.NewStrictHandler(app.MembersService, nil)
		log.Info("members domain ready")
		return nil
	})

	// Initialize invitations domain
	g.Go(func() error {
		log.Info("initializing invitations domain")
		app.InvitationsService = invitations.NewService(repos.Invitations, nil, log)
		app.InvitationsHandler = invitations.NewStrictHandler(app.InvitationsService, nil)
		log.Info("invitations domain ready")
		return nil
	})

	// Initialize health domain
	g.Go(func() error {
		log.Info("initializing health domain")
		app.HealthService = health.NewService(repos.Health, log)
		healthHandler := health.NewHandler(app.HealthService, log)
		app.HealthHandler = health.NewStrictHandler(healthHandler, nil)
		log.Info("health domain ready")
		return nil
	})

	// Wait for all parallel initializations to complete
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to initialize domains: %w", err)
	}

	// Initialize unified auth service
	log.Info("initializing unified auth service")

	// Create session store with cache
	sessionStore := stores.NewSessionStore(
		repos.Sessions,
		nil, // Add cache if available
		30*24*time.Hour,
	)

	// Create token manager
	tokenManager := tokens.NewManager(cfg.Auth.Local.JWTSecret)

	// Create magic link repository and store
	magicLinkRepo := stores.NewPostgresMagicLinkRepository(infra.Database.SQLDB())
	magicLinkStore := stores.NewMagicLinkStore(
		magicLinkRepo,
		15*time.Minute,
		5, // rate limit
	)

	// Create API token store
	// We need to create a postgresql.Queries instance for the API token repository
	var apiKeyStore auth.APIKeyStore
	if infra.Database.IsPostgreSQL() {
		pgQueries := postgresql.New(infra.Database.PgxPool())
		apiKeyStore = stores.NewAPIKeyRepository(pgQueries)
	} else {
		// TODO: Implement SQLite API token store
		apiKeyStore = nil
	}
	apiKeyValidator := tokens.NewAPIKeyValidator()

	// Create unified auth service
	app.AuthService = auth.NewService(
		cfg,
		log,
		app.UsersService,
		repos.Accounts,
		sessionStore,
		tokenManager,
		magicLinkStore,
		apiKeyStore,
		apiKeyValidator,
	)

	// Register OAuth providers
	if cfg.Auth.Google.Enabled {
		app.AuthService.RegisterProvider(auth.ProviderGoogle,
			providers.NewGoogleProvider(
				*cfg.Auth.Google.ClientID,
				*cfg.Auth.Google.ClientSecret,
				*cfg.Auth.Google.RedirectURL,
			),
		)
	}

	if cfg.Auth.Github.Enabled {
		app.AuthService.RegisterProvider(auth.ProviderGitHub,
			providers.NewGitHubProvider(
				*cfg.Auth.Github.ClientID,
				*cfg.Auth.Github.ClientSecret,
				*cfg.Auth.Github.RedirectURL,
			),
		)
	}

	if cfg.Auth.Microsoft.Enabled {
		app.AuthService.RegisterProvider(auth.ProviderMicrosoft,
			providers.NewMicrosoftProvider(
				*cfg.Auth.Microsoft.ClientID,
				*cfg.Auth.Microsoft.ClientSecret,
				*cfg.Auth.Microsoft.RedirectURL,
			),
		)
	}

	// Register deliverers
	app.AuthService.RegisterDeliverer(auth.DeliveryConsole,
		deliverers.NewConsoleDeliverer(log))
	app.AuthService.RegisterDeliverer(auth.DeliveryOTP,
		deliverers.NewOTPDeliverer(log))

	// Create unified auth handler
	app.AuthHandler = auth.NewHandler(app.AuthService, log)
	log.Info("unified auth handler ready")

	// Connect auth service to middleware
	app.AuthMiddleware.SetAuthService(app.AuthService)
	log.Info("connected auth service to middleware")

	// Initialize config handler
	log.Info("initializing config handler")
	configHandler := config.NewHandler(cfg, log)
	app.ConfigHandler = config.NewStrictHandler(configHandler, nil)
	log.Info("config handler ready")

	// Create the HTTP server
	log.Info("creating HTTP server")
	app.Server = server.NewServer(cfg.API, log)

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
