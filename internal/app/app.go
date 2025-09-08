// Package app provides dependency injection and application container management
package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/config"
	"github.com/archesai/archesai/internal/content"
	"github.com/archesai/archesai/internal/database"
	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/database/sqlite"
	"github.com/archesai/archesai/internal/health"
	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/migrations"
	"github.com/archesai/archesai/internal/organizations"
	"github.com/archesai/archesai/internal/server"
	"github.com/archesai/archesai/internal/workflows"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// App holds all application dependencies
type App struct {
	// Infrastructure
	DB            database.Database
	PgQueries     *postgresql.Queries // PostgreSQL queries (if using PostgreSQL)
	SqliteQueries *sqlite.Queries     // SQLite queries (if using SQLite)
	Logger        *slog.Logger
	Config        *config.Config
	Server        *server.Server // The HTTP server

	// Auth domain
	AuthRepository auth.Repository
	AuthService    *auth.Service
	AuthHandler    *auth.Handler

	// Organizations domain
	OrganizationsRepository organizations.Repository
	OrganizationsService    *organizations.Service
	OrganizationsHandler    *organizations.Handler

	// Workflows domain
	WorkflowsRepository workflows.Repository
	WorkflowsService    *workflows.Service
	WorkflowsHandler    *workflows.Handler

	// Content domain
	ContentRepository content.Repository
	ContentService    *content.Service
	ContentHandler    *content.Handler

	// Health domain
	HealthService *health.Service
	HealthHandler *health.Handler

	// Config handler
	// TODO: Implement config handler when needed
	// ConfigHandler *config.Handler
}

// NewApp creates and initializes all application dependencies
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize logger using the logger package
	loggerCfg := logger.Config{
		Level:  string(cfg.Logging.Level),
		Pretty: cfg.Logging.Pretty,
	}
	log := logger.New(loggerCfg)

	// Set as default logger
	slog.SetDefault(log)

	// Initialize database
	dbFactory := database.NewFactory(log)
	db, err := dbFactory.Create(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

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

	// Create queries based on database type
	var pgQueries *postgresql.Queries
	var sqliteQueries *sqlite.Queries
	var authRepo auth.ExtendedRepository
	var organizationsRepo organizations.ExtendedRepository
	var workflowsRepo workflows.Repository
	var contentRepo content.Repository

	// Determine actual database type
	var dbType database.Type
	if cfg.Database.Type != "" {
		dbType = database.ParseTypeFromString(string(cfg.Database.Type))
	} else {
		dbType = database.DetectTypeFromURL(cfg.Database.Url)
	}

	switch dbType {
	case database.TypePostgreSQL:
		// Get the underlying pgxpool for PostgreSQL
		if pool, ok := db.Underlying().(*pgxpool.Pool); ok && pool != nil {
			pgQueries = postgresql.New(pool)
			authRepo = auth.NewPostgresRepository(pgQueries)
			organizationsRepo = organizations.NewPostgresRepository(pgQueries)
			workflowsRepo = workflows.NewPostgresRepository(pgQueries)
			contentRepo = content.NewPostgresRepository(pgQueries)
		} else {
			return nil, fmt.Errorf("failed to get PostgreSQL connection pool")
		}
	case database.TypeSQLite:
		if sqlDB, ok := db.Underlying().(*sql.DB); ok && sqlDB != nil {
			sqliteQueries = sqlite.New(sqlDB)
			// Use SQLite repositories
			authRepo = auth.NewSQLiteRepository(sqliteQueries)
			organizationsRepo = organizations.NewSQLiteRepository(sqliteQueries)
			workflowsRepo = workflows.NewSQLiteRepository(sqliteQueries)
			contentRepo = content.NewSQLiteRepository(sqliteQueries)
			log.Info("Using SQLite repositories")
		}
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

	authConfig := auth.Config{
		JWTSecret:          cfg.GetJWTSecret(),
		AccessTokenExpiry:  accessTokenTTL,
		RefreshTokenExpiry: refreshTokenTTL,
	}
	authService := auth.NewService(authRepo, authConfig, log)
	authHandler := auth.NewHandler(authService, log)

	// Initialize organizations domain
	organizationsService := organizations.NewService(organizationsRepo, log)
	organizationsHandler := organizations.NewHandler(organizationsService, log)

	// Initialize workflows domain
	workflowsService := workflows.NewService(workflowsRepo, log)
	workflowsHandler := workflows.NewHandler(workflowsService, log)

	// Initialize content domain
	contentService := content.NewService(contentRepo, log)
	contentHandler := content.NewHandler(contentService, log)

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
		DB:            db,
		PgQueries:     pgQueries,
		SqliteQueries: sqliteQueries,
		Logger:        log,
		Config:        cfg,
		Server:        httpServer,

		// Auth domain
		AuthRepository: authRepo,
		AuthService:    authService,
		AuthHandler:    authHandler,

		// Organizations domain
		OrganizationsRepository: organizationsRepo,
		OrganizationsService:    organizationsService,
		OrganizationsHandler:    organizationsHandler,

		// Workflows domain
		WorkflowsRepository: workflowsRepo,
		WorkflowsService:    workflowsService,
		WorkflowsHandler:    workflowsHandler,

		// Content domain
		ContentRepository: contentRepo,
		ContentService:    contentService,
		ContentHandler:    contentHandler,

		// Health domain
		HealthService: healthService,
		HealthHandler: healthHandler,

		// Config handler
		// ConfigHandler: configHandler,
	}

	// Register all application routes
	app.registerRoutes()

	return app, nil
}

// Close cleans up all resources
func (a *App) Close() error {
	// Close database connection
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			a.Logger.Error("failed to close database connection", "error", err)
		}
	}

	// Logger cleanup not needed for slog (it uses os.Stdout/Stderr)

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
	if err := a.DB.Ping(ctx.Request().Context()); err != nil {
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
