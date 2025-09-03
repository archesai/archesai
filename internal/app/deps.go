// Package app provides application dependency injection and route registration.
package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/domains/auth/adapters/postgres"
	"github.com/archesai/archesai/internal/domains/auth/handlers"
	"github.com/archesai/archesai/internal/domains/auth/repositories"
	"github.com/archesai/archesai/internal/domains/auth/services"
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

	// TODO: Register other feature routes as they are implemented
	// Example for intelligence feature:
	// intelligenceGroup := v1.Group("/intelligence")
	// artifacts.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// labels.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// pipelines.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// runs.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// tools.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
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

	// Auth feature
	AuthRepository repositories.Repository
	AuthService    *services.Service
	AuthHandler    *handlers.Handler

	// TODO: Add other features as they are implemented
	// IntelligenceRepository intelligence.Repository
	// IntelligenceService    intelligence.Service
	// IntelligenceHandler    *intelligencehttp.Handler
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
	var dbType database.Type
	if cfg.Database.Type != nil {
		dbType = database.ParseTypeString(string(*cfg.Database.Type))
	} else {
		// Auto-detect from URL
		url := cfg.Database.Url
		if strings.HasPrefix(url, "postgresql://") || strings.HasPrefix(url, "postgres://") {
			dbType = database.TypePostgreSQL
		} else if strings.HasPrefix(url, "sqlite://") || strings.Contains(url, ".db") || url == ":memory:" {
			dbType = database.TypeSQLite
		} else {
			dbType = database.TypePostgreSQL // default
		}
	}

	// Create database config
	dbConfig := &database.Config{
		URL:  cfg.Database.Url,
		Type: dbType,
	}

	// Apply optional configuration
	if cfg.Database.MaxConns != nil {
		dbConfig.PgMaxConns = int32(*cfg.Database.MaxConns)
	}
	if cfg.Database.MinConns != nil {
		dbConfig.PgMinConns = int32(*cfg.Database.MinConns)
	}
	if cfg.Database.ConnMaxLifetime != nil {
		dbConfig.ConnMaxLifetime = *cfg.Database.ConnMaxLifetime
	}
	if cfg.Database.ConnMaxIdleTime != nil {
		dbConfig.ConnMaxIdleTime = *cfg.Database.ConnMaxIdleTime
	}
	if cfg.Database.HealthCheckPeriod != nil {
		dbConfig.PgHealthCheckPeriod = *cfg.Database.HealthCheckPeriod
	}
	if cfg.Database.RunMigrations != nil {
		dbConfig.RunMigrations = *cfg.Database.RunMigrations
	}

	dbFactory := database.NewFactory(logger)
	db, err := dbFactory.Create(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations if enabled
	if dbConfig.RunMigrations {
		if err := database.RunMigrations(db, logger); err != nil {
			logger.Error("failed to run migrations", "error", err)
			if cfg.Server.Environment != "development" {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		}
	}

	// Create queries based on database type
	var pgQueries *postgresqlgen.Queries
	var sqliteQueries *sqlitegen.Queries
	var authRepo repositories.Repository

	switch dbConfig.Type {
	case database.TypePostgreSQL:
		// Get the underlying pgxpool for PostgreSQL
		if pool, ok := db.Underlying().(*pgxpool.Pool); ok && pool != nil {
			pgQueries = postgresqlgen.New(pool)
			authRepo = postgres.NewRepository(pgQueries)
		} else {
			return nil, fmt.Errorf("failed to get PostgreSQL connection pool")
		}
	case database.TypeSQLite:
		if sqlDB, ok := db.Underlying().(*sql.DB); ok && sqlDB != nil {
			sqliteQueries = sqlitegen.New(sqlDB)
			// TODO: Create SQLite auth repository when available
			logger.Warn("SQLite repositories not yet implemented, auth features may not work")
		}
	}

	// Initialize auth feature
	authConfig := services.Config{
		JWTSecret:          cfg.Auth.JWTSecret,
		AccessTokenExpiry:  cfg.Auth.AccessTokenTTL,
		RefreshTokenExpiry: cfg.Auth.RefreshTokenTTL,
	}
	authService := services.NewService(authRepo, authConfig, logger)
	authHandler := handlers.NewHandler(authService, logger)

	// TODO: Initialize other features
	// intelligenceRepo := intelligencepostgresql.NewRepository(db)
	// intelligenceService := intelligenceusecase.NewService(intelligenceRepo, logger)
	// intelligenceHandler := intelligencehttp.NewHandler(intelligenceService, logger)

	// Create the HTTP server
	serverConfig := &server.Config{
		Port:           fmt.Sprintf("%d", cfg.Server.Port),
		AllowedOrigins: cfg.GetAllowedOrigins(),
		DocsEnabled:    cfg.Server.DocsEnabled,
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

		// Auth feature
		AuthRepository: authRepo,
		AuthService:    authService,
		AuthHandler:    authHandler,

		// TODO: Add other features
		// IntelligenceRepository: intelligenceRepo,
		// IntelligenceService:    intelligenceService,
		// IntelligenceHandler:    intelligenceHandler,
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
	if c.Config.Server.DocsEnabled {
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
