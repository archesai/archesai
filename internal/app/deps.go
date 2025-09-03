// Package app provides application dependency injection and route registration.
package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/archesai/archesai/internal/domains/auth"
	"github.com/archesai/archesai/internal/domains/content"
	"github.com/archesai/archesai/internal/domains/organizations"
	"github.com/archesai/archesai/internal/domains/workflows"
	"github.com/archesai/archesai/internal/generated/api"
	postgresqlgen "github.com/archesai/archesai/internal/generated/database/postgresql"
	sqlitegen "github.com/archesai/archesai/internal/generated/database/sqlite"
	"github.com/archesai/archesai/internal/infrastructure/config"
	"github.com/archesai/archesai/internal/infrastructure/database"
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
}

// NewContainer creates and initializes all application dependencies
func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize slog logger
	var logger *slog.Logger
	var logLevel slog.Level

	// Parse log level
	switch cfg.Logging.Level {
	case api.LoggingConfigLevelDebug:
		logLevel = slog.LevelDebug
	case api.LoggingConfigLevelInfo:
		logLevel = slog.LevelInfo
	case api.LoggingConfigLevelWarn:
		logLevel = slog.LevelWarn
	case api.LoggingConfigLevelError:
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
	var authRepo auth.Repository
	var organizationsRepo organizations.Repository
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
			pgQueries = postgresqlgen.New(pool)
			authRepo = auth.NewPostgresRepository(pgQueries)
			organizationsRepo = organizations.NewPostgresRepository(pgQueries)
			workflowsRepo = workflows.NewPostgresRepository(pgQueries)
			contentRepo = content.NewPostgresRepository(pgQueries)
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

	authConfig := auth.Config{
		JWTSecret:          cfg.GetJWTSecret(),
		AccessTokenExpiry:  accessTokenTTL,
		RefreshTokenExpiry: refreshTokenTTL,
	}
	authService := auth.NewService(authRepo, authConfig, logger)
	authHandler := auth.NewHandler(authService, logger)

	// Initialize organizations domain
	organizationsService := organizations.NewService(organizationsRepo, logger)
	organizationsHandler := organizations.NewHandler(organizationsService, logger)

	// Initialize workflows domain
	workflowsService := workflows.NewService(workflowsRepo, logger)
	workflowsHandler := workflows.NewHandler(workflowsService, logger)

	// Initialize content domain
	contentService := content.NewService(contentRepo, logger)
	contentHandler := content.NewHandler(contentService, logger)

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

	// Setup API documentation if enabled
	if c.Config.Api.Docs {
		swagger, err := api.GetSwagger()
		if err != nil {
			c.Logger.Error("failed to load OpenAPI spec", "error", err)
		} else {
			if err := c.Server.SetupDocs(swagger); err != nil {
				c.Logger.Error("failed to setup API docs", "error", err)
			}
		}
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
