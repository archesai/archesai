package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	postgresqlgen "github.com/archesai/archesai/gen/db/postgresql"
	"github.com/labstack/echo/v4"
	authhttp "github.com/archesai/archesai/internal/features/auth/adapters/http"
	"github.com/archesai/archesai/internal/features/auth/adapters/postgresql"
	"github.com/archesai/archesai/internal/features/auth/ports"
	"github.com/archesai/archesai/internal/features/auth/usecase"
	"github.com/archesai/archesai/internal/infrastructure/config"
	"github.com/archesai/archesai/internal/infrastructure/database"
	"github.com/archesai/archesai/internal/infrastructure/server"
)

// Container holds all application dependencies
type Container struct {
	// Infrastructure
	DB      *database.DB
	Queries *postgresqlgen.Queries
	Logger  *slog.Logger
	Config  *config.Config
	Server  *server.Server  // The HTTP server

	// Auth feature
	AuthRepository ports.Repository
	AuthService    ports.Service
	AuthHandler    *authhttp.Handler

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
			Level: logLevel,
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
	dbConfig := database.Config{
		URL: cfg.Database.URL,
	}
	db, err := database.NewConnection(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create queries instance for sqlc
	queries := postgresqlgen.New(db)

	// Initialize auth feature
	authRepo := postgresql.NewRepository(queries)
	authConfig := usecase.Config{
		JWTSecret:          cfg.Auth.JWTSecret,
		AccessTokenExpiry:  cfg.Auth.AccessTokenTTL,
		RefreshTokenExpiry: cfg.Auth.RefreshTokenTTL,
	}
	authService := usecase.NewService(authRepo, authConfig, logger)
	authHandler := authhttp.NewHandler(authService, logger)

	// TODO: Initialize other features
	// intelligenceRepo := intelligencepostgresql.NewRepository(db)
	// intelligenceService := intelligenceusecase.NewService(intelligenceRepo, logger)
	// intelligenceHandler := intelligencehttp.NewHandler(intelligenceService, logger)

	// Create the HTTP server
	serverConfig := &server.Config{
		Port:           fmt.Sprintf("%d", cfg.Server.Port),
		AllowedOrigins: cfg.GetAllowedOrigins(),
	}
	httpServer := server.NewServer(serverConfig, logger)

	// Create container with all dependencies
	container := &Container{
		// Infrastructure
		DB:      db,
		Queries: queries,
		Logger:  logger,
		Config:  cfg,
		Server:  httpServer,

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
		c.DB.Close()
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
