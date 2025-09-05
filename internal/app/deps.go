// Package app provides application dependency injection and route registration.
package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	authcore "github.com/archesai/archesai/internal/domains/auth/core"
	authhandlers "github.com/archesai/archesai/internal/domains/auth/handlers/http"
	authinfra "github.com/archesai/archesai/internal/domains/auth/infrastructure"
	contentcore "github.com/archesai/archesai/internal/domains/content/core"
	contenthandlers "github.com/archesai/archesai/internal/domains/content/handlers/http"
	contentinfra "github.com/archesai/archesai/internal/domains/content/infrastructure"
	orgcore "github.com/archesai/archesai/internal/domains/organizations/core"
	orghandlers "github.com/archesai/archesai/internal/domains/organizations/handlers/http"
	orginfra "github.com/archesai/archesai/internal/domains/organizations/infrastructure"
	workflowcore "github.com/archesai/archesai/internal/domains/workflows/core"
	workflowhandlers "github.com/archesai/archesai/internal/domains/workflows/handlers/http"
	workflowinfra "github.com/archesai/archesai/internal/domains/workflows/infrastructure"
	"github.com/archesai/archesai/internal/infrastructure/config"
	"github.com/archesai/archesai/internal/infrastructure/database"
	postgresqlgen "github.com/archesai/archesai/internal/infrastructure/database/generated/postgresql"
	sqlitegen "github.com/archesai/archesai/internal/infrastructure/database/generated/sqlite"
	"github.com/archesai/archesai/internal/infrastructure/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func RegisterRoutes(e *echo.Echo, container *Container) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes
	container.AuthHandler.RegisterRoutes(v1)

	// Register organizations routes
	orgGroup := v1.Group("/organizations")
	container.OrganizationsHandler.RegisterRoutes(orgGroup)

	// Register workflows routes
	workflowsGroup := v1.Group("/workflows")
	container.WorkflowsHandler.RegisterRoutes(workflowsGroup)

	// Register content routes
	contentGroup := v1.Group("/content")
	container.ContentHandler.RegisterRoutes(contentGroup)
}

// Container holds all application dependencies
type Container struct {
	// Infrastructure
	DB            database.Database
	PgQueries     *postgresqlgen.Queries // PostgreSQL queries (if using PostgreSQL)
	SqliteQueries *sqlitegen.Queries     // SQLite queries (if using SQLite)
	Logger        *slog.Logger
	Config        *config.Config
	Server        *server.Server // The HTTP server

	// Auth domain
	AuthRepository authcore.Repository
	AuthService    *authcore.Service
	AuthHandler    *authhandlers.Handler

	// Organizations domain
	OrganizationsRepository orgcore.Repository
	OrganizationsService    *orgcore.Service
	OrganizationsHandler    *orghandlers.Handler

	// Workflows domain
	WorkflowsRepository workflowcore.Repository
	WorkflowsService    *workflowcore.Service
	WorkflowsHandler    *workflowhandlers.Handler

	// Content domain
	ContentRepository contentcore.Repository
	ContentService    *contentcore.Service
	ContentHandler    *contenthandlers.Handler
}

// NewContainer creates and initializes all application dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize slog logger
	var logger *slog.Logger
	var logLevel slog.Level

	// Parse log level
	switch cfg.Logging.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Configure handler based on format preference
	if cfg.Logging.Pretty {
		// Use text handler for pretty output
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: logLevel == slog.LevelDebug,
		}))
	} else {
		// Use JSON handler for production
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))
	}

	// Set as default logger
	slog.SetDefault(logger)

	// Initialize database
	dbFactory := database.NewFactory(logger)
	db, err := dbFactory.Create(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations if enabled
	if cfg.Database.RunMigrations {
		if err := database.RunMigrations(db, logger); err != nil {
			logger.Error("failed to run migrations", "error", err)
			isProduction := cfg.Api.Environment == "production"
			if isProduction {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		}
	}

	// Create queries based on database type
	var pgQueries *postgresqlgen.Queries
	var sqliteQueries *sqlitegen.Queries
	var authRepo authcore.Repository
	var organizationsRepo orgcore.Repository
	var workflowsRepo workflowcore.Repository
	var contentRepo contentcore.Repository

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
			pgQueries = postgresqlgen.New(pool)
			authRepo = authinfra.NewPostgresRepository(pgQueries)
			organizationsRepo = orginfra.NewPostgresRepository(pgQueries)
			workflowsRepo = workflowinfra.NewPostgresRepository(pgQueries)
			contentRepo = contentinfra.NewPostgresRepository(pgQueries)
		} else {
			return nil, fmt.Errorf("failed to get PostgreSQL connection pool")
		}
	case database.TypeSQLite:
		if sqlDB, ok := db.Underlying().(*sql.DB); ok && sqlDB != nil {
			sqliteQueries = sqlitegen.New(sqlDB)
			// TODO: Create SQLite repositories when available
			logger.Warn("SQLite repositories not yet implemented, features may not work")
			// authRepo = auth.NewSQLiteRepository(sqliteQueries)
			// organizationsRepo = organizations.NewSQLiteRepository(sqliteQueries)
			// workflowsRepo = workflows.NewSQLiteRepository(sqliteQueries)
			// contentRepo = content.NewSQLiteRepository(sqliteQueries)
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

	authConfig := authcore.Config{
		JWTSecret:          cfg.GetJWTSecret(),
		AccessTokenExpiry:  accessTokenTTL,
		RefreshTokenExpiry: refreshTokenTTL,
	}
	authService := authcore.NewService(authRepo, authConfig, logger)
	authHandler := authhandlers.NewHandler(authService, logger)

	// Initialize organizations domain
	organizationsService := orgcore.NewService(organizationsRepo, logger)
	organizationsHandler := orghandlers.NewHandler(organizationsService, logger)

	// Initialize workflows domain
	workflowsService := workflowcore.NewService(workflowsRepo, logger)
	workflowsHandler := workflowhandlers.NewHandler(workflowsService, logger)

	// Initialize content domain
	contentService := contentcore.NewService(contentRepo, logger)
	contentHandler := contenthandlers.NewHandler(contentService, logger)

	// Create the HTTP server
	serverConfig := &server.Config{
		Port:           fmt.Sprintf("%d", int(cfg.Api.Port)),
		AllowedOrigins: cfg.GetAllowedOrigins(),
		DocsEnabled:    cfg.Api.Docs,
	}
	httpServer := server.NewServer(serverConfig, logger)

	// Create container with all dependencies
	container := &Container{
		// Infrastructure
		DB:            db,
		PgQueries:     pgQueries,
		SqliteQueries: sqliteQueries,
		Logger:        logger,
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
	}

	// Register all application routes
	container.registerRoutes()

	return container, nil
}

// Close cleans up all resources
func (c *Container) Close() error {
	// Close database connection
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			c.Logger.Error("failed to close database connection", "error", err)
		}
	}

	// Logger cleanup not needed for slog (it uses os.Stdout/Stderr)

	return nil
}

// registerRoutes registers all application routes with the server
func (c *Container) registerRoutes() {
	// Get the echo instance from the server
	e := c.Server.Echo()

	// Register readiness check that can access the database
	c.Server.SetReadinessCheck(c.readinessCheck)

	// TODO: Setup domain-scoped API documentation
	// Setup API documentation if enabled
	if c.Config.Api.Docs {
		c.Logger.Info("API documentation setup temporarily disabled - needs domain-scoped implementation")
		// swagger, err := api.GetSwagger()
		// if err != nil {
		// 	c.Logger.Error("failed to load OpenAPI spec", "error", err)
		// } else {
		// 	if err := c.Server.SetupDocs(swagger); err != nil {
		// 		c.Logger.Error("failed to setup API docs", "error", err)
		// 	}
		// }
	}

	// Register all application routes
	RegisterRoutes(e, c)
}

// readinessCheck checks if the service is ready to handle requests
func (c *Container) readinessCheck(ctx echo.Context) error {
	// Check database connection
	if err := c.DB.Ping(ctx.Request().Context()); err != nil {
		c.Logger.Error("database health check failed")
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
