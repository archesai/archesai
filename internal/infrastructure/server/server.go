package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/archesai/archesai/internal/infrastructure/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// Server represents the API server
type Server struct {
	echo   *echo.Echo
	db     *database.DB
	logger *zap.Logger
	config *Config
}

// Config holds server configuration
type Config struct {
	Port           string
	AllowedOrigins []string
	JWTSecret      string
}

// NewServer creates a new API server
func NewServer(db *database.DB, logger *zap.Logger, config *Config) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:   e,
		db:     db,
		logger: logger,
		config: config,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	// Request ID middleware
	s.echo.Use(middleware.RequestID())

	// Logger middleware
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			s.logger.Info("request",
				zap.String("id", v.RequestID),
				zap.String("method", c.Request().Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
				zap.String("remote_ip", c.RealIP()),
				zap.Error(v.Error),
			)
			return nil
		},
	}))

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	// swagger, err := GetSwagger()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
	// 	os.Exit(1)
	// }

	// s.echo.Use(echomiddleware.OapiRequestValidator(swagger))

	// Recover middleware
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			s.logger.Error("panic recovered",
				zap.String("id", c.Response().Header().Get(echo.HeaderXRequestID)),
				zap.Error(err),
				zap.String("stack", string(stack)),
			)
			return nil
		},
	}))

	// CORS middleware
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.config.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Request-ID"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Request-ID"},
		MaxAge:           86400,
	}))

	// Compression middleware
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
	}))

	// Security middleware
	s.echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	// Rate limiting (basic example - consider using a Redis-based solution in production)
	s.echo.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      10,
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			// Use IP address as identifier
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// Body limit middleware
	s.echo.Use(middleware.BodyLimit("10M"))

	// Timeout middleware
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      30 * time.Second,
		ErrorMessage: "Request timeout",
	}))
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Health check
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"database":  "healthy",
		})
	})

	// API version endpoint
	s.echo.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version": "1.0.0",
			"build":   "development",
		})
	})

	// API v1 routes
	// api := s.echo.Group("/api/v1")
	// Register user and organization routes
	// userHandler.RegisterRoutes(api)
	// orgHandler.RegisterRoutes(api)

	// 404 handler
	s.echo.RouteNotFound("/*", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "route not found")
	})
	// 	// 404 handler
	// e.RouteNotFound("/*", func(c echo.Context) error {
	// 	return c.JSON(http.StatusNotFound, map[string]string{
	// 		"error": "Route not found",
	// 	})
	// })
}

// Start starts the server
func (s *Server) Start() error {
	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", s.config.Port)
		s.logger.Info("starting server", zap.String("address", addr))

		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	s.logger.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Error("server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("server shutdown complete")
	return nil
}

// Shutdown shuts down the server gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
