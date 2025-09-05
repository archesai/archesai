// Package container provides dependency injection and application container management
package container

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	authhandlers "github.com/archesai/archesai/internal/auth/adapters/http"
	authinfra "github.com/archesai/archesai/internal/auth/adapters/postgres"
	authsqlite "github.com/archesai/archesai/internal/auth/adapters/sqlite"
	authcore "github.com/archesai/archesai/internal/auth/domain"
	"github.com/archesai/archesai/internal/config"
	contenthandlers "github.com/archesai/archesai/internal/content/adapters/http"
	contentinfra "github.com/archesai/archesai/internal/content/adapters/postgres"
	contentsqlite "github.com/archesai/archesai/internal/content/adapters/sqlite"
	contentcore "github.com/archesai/archesai/internal/content/domain"
	"github.com/archesai/archesai/internal/database"
	postgresqlgen "github.com/archesai/archesai/internal/database/postgresql"
	sqlitegen "github.com/archesai/archesai/internal/database/sqlite"
	orghandlers "github.com/archesai/archesai/internal/organizations/adapters/http"
	orginfra "github.com/archesai/archesai/internal/organizations/adapters/postgres"
	orgsqlite "github.com/archesai/archesai/internal/organizations/adapters/sqlite"
	orgcore "github.com/archesai/archesai/internal/organizations/domain"
	serverhttp "github.com/archesai/archesai/internal/server/http"
	workflowhandlers "github.com/archesai/archesai/internal/workflows/adapters/http"
	workflowinfra "github.com/archesai/archesai/internal/workflows/adapters/postgres"
	workflowsqlite "github.com/archesai/archesai/internal/workflows/adapters/sqlite"
	workflowcore "github.com/archesai/archesai/internal/workflows/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// Container holds all application dependencies
type Container struct {
	// Infrastructure
	DB            database.Database
	PgQueries     *postgresqlgen.Queries // PostgreSQL queries (if using PostgreSQL)
	SqliteQueries *sqlitegen.Queries     // SQLite queries (if using SQLite)
	Logger        *slog.Logger
	Config        *config.Config
	Server        *serverhttp.Server // The HTTP server

	// Auth domain
	AuthRepository authcore.AuthRepository
	AuthService    *authcore.AuthService
	AuthHandler    *authhandlers.AuthHandler

	// Organizations domain
	OrganizationsRepository orgcore.OrganizationRepository
	OrganizationsService    *orgcore.OrganizationService
	OrganizationsHandler    *orghandlers.OrganizationHandler

	// Workflows domain
	WorkflowsRepository workflowcore.WorkflowRepository
	WorkflowsService    *workflowcore.WorkflowService
	WorkflowsHandler    *workflowhandlers.WorkflowHandler

	// Content domain
	ContentRepository contentcore.ContentRepository
	ContentService    *contentcore.ContentService
	ContentHandler    *contenthandlers.ContentHandler
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
	var authRepo authcore.AuthRepository
	var organizationsRepo orgcore.OrganizationRepository
	var workflowsRepo workflowcore.WorkflowRepository
	var contentRepo contentcore.ContentRepository

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
			authRepo = authinfra.NewAuthPostgresRepository(pgQueries)
			organizationsRepo = orginfra.NewOrganizationPostgresRepository(pgQueries)
			workflowsRepo = workflowinfra.NewWorkflowPostgresRepository(pgQueries)
			contentRepo = contentinfra.NewContentPostgresRepository(pgQueries)
		} else {
			return nil, fmt.Errorf("failed to get PostgreSQL connection pool")
		}
	case database.TypeSQLite:
		if sqlDB, ok := db.Underlying().(*sql.DB); ok && sqlDB != nil {
			sqliteQueries = sqlitegen.New(sqlDB)
			// Use SQLite repositories
			authRepo = authsqlite.NewAuthSQLiteRepository(sqliteQueries)
			organizationsRepo = orgsqlite.NewOrganizationSQLiteRepository(sqliteQueries)
			workflowsRepo = workflowsqlite.NewWorkflowSQLiteRepository(sqliteQueries)
			contentRepo = contentsqlite.NewContentSQLiteRepository(sqliteQueries)
			logger.Info("Using SQLite repositories")
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
	authService := authcore.NewAuthService(authRepo, authConfig, logger)
	authHandler := authhandlers.NewAuthHandler(authService, logger)

	// Initialize organizations domain
	organizationsService := orgcore.NewOrganizationService(organizationsRepo, logger)
	organizationsHandler := orghandlers.NewOrganizationHandler(organizationsService, logger)

	// Initialize workflows domain
	workflowsService := workflowcore.NewWorkflowService(workflowsRepo, logger)
	workflowsHandler := workflowhandlers.NewWorkflowHandler(workflowsService, logger)

	// Initialize content domain
	contentService := contentcore.NewContentService(contentRepo, logger)
	contentHandler := contenthandlers.NewContentHandler(contentService, logger)

	// Create the HTTP server
	serverConfig := &serverhttp.Config{
		Port:           fmt.Sprintf("%d", int(cfg.Api.Port)),
		AllowedOrigins: cfg.GetAllowedOrigins(),
		DocsEnabled:    cfg.Api.Docs,
	}
	httpServer := serverhttp.NewServer(serverConfig, logger)

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
		if err := c.Server.SetupDocs(); err != nil {
			c.Logger.Error("failed to setup API docs", "error", err)
		}
	}

	// Register all application routes
	c.RegisterRoutes(e)
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
