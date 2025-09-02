package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server represents the API server
type Server struct {
	echo   *echo.Echo
	config *Config
	logger *slog.Logger
}

// Config holds server configuration
type Config struct {
	Port           string
	AllowedOrigins []string
	DocsEnabled    bool
}

// NewServer creates a new API server
func NewServer(config *Config, logger *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:   e,
		config: config,
		logger: logger,
	}

	server.setupMiddleware()
	server.setupInfrastructureRoutes()

	return server
}

// Echo returns the underlying echo instance for route registration
func (s *Server) Echo() *echo.Echo {
	return s.echo
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
				"id", v.RequestID,
				"method", c.Request().Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"remote_ip", c.RealIP(),
				"error", v.Error,
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
				"id", c.Response().Header().Get(echo.HeaderXRequestID),
				"error", err,
				"stack", string(stack),
			)
			return nil
		},
	}))

	// CORS middleware
	// 	  const allowedOrigins = configService.get('api.cors.origins').split(',')
	//   await app.register(cors, {
	//     allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
	//     credentials: true,
	//     maxAge: 86400,
	//     methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
	//     origin: allowedOrigins
	//   })
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.config.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Request-ID"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Request-ID"},
		MaxAge:           86400,
	}))

	// FIXME
	// // Security Middlewares
	// await httpInstance.register(helmet, {
	//   contentSecurityPolicy: {
	//     directives: {
	//       defaultSrc: [`'self'`],
	//       fontSrc: [`'self'`, 'fonts.scalar.com', 'data:'],
	//       imgSrc: [`'self'`, 'data:'],
	//       scriptSrc: [`'self'`, `https: 'unsafe-inline'`, `'unsafe-eval'`],
	//       styleSrc: [`'self'`, `'unsafe-inline'`, 'fonts.scalar.com']
	//     }
	//   }
	// })

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
			s.logger.Warn("rate limiter error", "error", err, "ip", c.RealIP())
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			s.logger.Info("rate limit exceeded", "identifier", identifier, "path", c.Request().URL.Path, "error", err)
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

// setupInfrastructureRoutes configures infrastructure routes only
func (s *Server) setupInfrastructureRoutes() {
	// Health check - simple liveness probe
	s.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// API version endpoint
	s.echo.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version": "1.0.0",
			"build":   "development",
		})
	})

	// 404 handler - must be registered last (will be overridden when container registers routes)
	s.echo.RouteNotFound("/*", func(_ echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "route not found")
	})
}

// SetReadinessCheck allows the container to provide a readiness check function
func (s *Server) SetReadinessCheck(checkFunc func(echo.Context) error) {
	s.echo.GET("/ready", checkFunc)
}

// Start starts the server
func (s *Server) Start() error {
	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", s.config.Port)
		s.logger.Info("starting server", "address", addr)

		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			s.logger.Error("failed to start server", "error", err)
			os.Exit(1)
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
		s.logger.Error("server forced to shutdown", "error", err)
		return err
	}

	s.logger.Info("server shutdown complete")
	return nil
}

// Shutdown shuts down the server gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
