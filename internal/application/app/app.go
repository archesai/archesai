// Package app provides dependency injection and application container management
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/labstack/echo/v4"

	"github.com/archesai/archesai/internal/adapters/http/controllers"
	"github.com/archesai/archesai/internal/adapters/http/server"
	"github.com/archesai/archesai/internal/infrastructure/config"
	"github.com/archesai/archesai/internal/infrastructure/http/middleware"
	database "github.com/archesai/archesai/internal/infrastructure/persistence"
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
	// AuthMiddleware *middleware.AuthMiddleware  // TODO: Fix middleware type issue

	// HTTP handlers
	AccountsHandler      controllers.AccountsController
	UsersHandler         controllers.UsersController
	OrganizationsHandler controllers.OrganizationsController
	InvitationsHandler   controllers.InvitationsController
	ArtifactsHandler     controllers.ArtifactsController
	LabelsHandler        controllers.LabelsController
	MembersHandler       controllers.MembersController
	PipelinesHandler     controllers.PipelinesController
	RunsHandler          controllers.RunsController
	ToolsHandler         controllers.ToolsController
	HealthHandler        controllers.HealthController
	ConfigHandler        controllers.ConfigController
	AuthHandler          controllers.AuthController
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
		if err := database.RunMigrations(infra.Database.SQLDB(), infra.Database.TypeString(), log); err != nil {
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
		app.AccountsHandler = controllers.AccountsController(app.AccountsService)
		log.Info("accounts domain ready")
		return nil
	})

	// Initialize users domain
	g.Go(func() error {
		log.Info("initializing users domain")
		app.UsersHandler = controllers.UsersController(app.UsersService)
		log.Info("users domain ready")
		return nil
	})

	// Initialize organizations domain
	g.Go(func() error {
		log.Info("initializing organizations domain")
		app.OrganizationsHandler = controllers.OrganizationsController(
			app.OrganizationsService,
		)
		log.Info("organizations domain ready")
		return nil
	})

	// Initialize pipelines domain
	g.Go(func() error {
		log.Info("initializing pipelines domain")
		app.PipelinesHandler = *controllers.NewPipelinesController(app.PipelinesService)
		log.Info("pipelines domain ready")
		return nil
	})

	// Initialize runs domain
	g.Go(func() error {
		log.Info("initializing runs domain")
		app.RunsHandler = controllers.RunsController(app.RunsService)
		log.Info("runs domain ready")
		return nil
	})

	// Initialize tools domain
	g.Go(func() error {
		log.Info("initializing tools domain")
		app.ToolsHandler = controllers.ToolsController(app.ToolsService)
		log.Info("tools domain ready")
		return nil
	})

	// Initialize artifacts domain
	g.Go(func() error {
		log.Info("initializing artifacts domain")
		app.ArtifactsHandler = controllers.ArtifactsController(app.ArtifactsService)
		log.Info("artifacts domain ready")
		return nil
	})

	// Initialize labels domain
	g.Go(func() error {
		log.Info("initializing labels domain")
		app.LabelsHandler = controllers.LabelsController(app.LabelsService)
		log.Info("labels domain ready")
		return nil
	})

	// Initialize members domain
	g.Go(func() error {
		log.Info("initializing members domain")
		app.MembersHandler = controllers.MembersController(app.MembersService)
		log.Info("members domain ready")
		return nil
	})

	// Initialize invitations domain
	g.Go(func() error {
		log.Info("initializing invitations domain")
		app.InvitationsHandler = controllers.InvitationsController(app.InvitationsService)
		log.Info("invitations domain ready")
		return nil
	})

	// Initialize health domain
	g.Go(func() error {
		log.Info("initializing health domain")
		app.HealthHandler = controllers.HealthController(app.HealthService)
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
	tokenManager := tokens.NewManager(cfg.Auth.Local.JwtSecret)

	// Create magic link repository and store
	magicLinkRepo := stores.NewPostgresMagicLinkRepository(infra.Database.SQLDB())
	magicLinkStore := stores.NewMagicLinkStore(
		magicLinkRepo,
		15*time.Minute,
		5, // rate limit
	)

	// Create API token store
	// We need to create a postgresql.Queries instance for the API token repository
	if infra.Database.IsPostgreSQL() {
		pgQueries := database.NewDatabase(
			infra.Database.SQLDB(),
			infra.Database.PgxPool(),
			database.TypePostgreSQL,
		)
		apiKeyStore = stores.NewAPIKeyRepository(pgQueries)
	} else {
		// TODO: Implement SQLite API token store
		apiKeyStore = nil
	}
	apiKeyValidator := tokens.NewAPIKeyValidator()

	// Register OAuth providers
	if cfg.Auth.Google.Enabled {
		app.AuthService.RegisterProvider(dto.ProviderGoogle,
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
	app.AuthHandler = controllers.AuthController(app.AuthService)
	log.Info("unified auth handler ready")

	// Connect auth service to middleware
	app.AuthMiddleware.SetAuthService(app.AuthService)
	log.Info("connected auth service to middleware")

	// Initialize config handler
	log.Info("initializing config handler")
	configHandler := controllers.ConfigController(cfg)
	app.ConfigHandler = configHandler
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
	a.registerRoutes(e)
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
	controllers.RegisterAccountsRoutes(v1, &a.AccountsService)

	// Unified Auth routes - all authentication methods in one handler
	a.Logger.Info("registering unified auth routes (public)")
	controllers.RegisterAuthRoutes(v1, &a.AuthService)

	// Health routes - public endpoints for monitoring
	a.Logger.Info("registering health routes (public)")
	controllers.RegisterHealthRoutes(v1, &a.HealthService)

	// Config routes - public endpoints (sanitized configuration)
	a.Logger.Info("registering config routes (public)")
	controllers.RegisterConfigRoutes(v1, a.ConfigHandler)

	// ========================================
	// PROTECTED ROUTES (Authentication required)
	// ========================================

	// Create a protected group with authentication middleware
	a.Logger.Info("creating protected routes group")
	protected := v1.Group("")
	protected.Use(a.AuthMiddleware.RequireAuth())

	// Users routes - require authentication
	a.Logger.Info("registering users routes (protected)")
	controllers.RegisterUsersRoutes(protected, &a.UsersService)

	// Organizations routes - require authentication and organization membership
	// Create a separate group for organization routes with additional checks
	orgGroup := v1.Group("")
	orgGroup.Use(a.AuthMiddleware.RequireAuth(), a.AuthMiddleware.RequireOrganizationMember())
	controllers.RegisterOrganizationsRoutes(orgGroup, &a.OrganizationsService)

	// Pipelines routes - require authentication
	a.Logger.Info("registering pipelines routes (protected)")
	controllers.RegisterPipelinesRoutes(protected, &a.PipelinesService)

	// Runs routes - require authentication
	controllers.RegisterRunsRoutes(protected, &a.RunsService)

	// Tools routes - require authentication
	controllers.RegisterToolsRoutes(protected, &a.ToolsService)

	// Artifacts routes - require authentication
	controllers.RegisterArtifactsRoutes(protected, &a.ArtifactsService)

	// Labels routes - require authentication
	controllers.RegisterLabelsRoutes(protected, &a.LabelsService)

	// Members routes - require authentication and organization membership
	controllers.RegisterMembersRoutes(orgGroup, &a.MembersService)

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
