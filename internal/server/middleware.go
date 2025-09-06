package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware configures all middleware for the server
func (s *Server) SetupMiddleware() {
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

	// TODO: Update to use domain-scoped validation middleware
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	// swagger, err := api.GetSwagger()
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
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(s.config.Cors.Origins, ","),
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Request-ID", "X-Requested-With"},
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
		// FIXME - adjust ContentSecurityPolicy as needed
		//     directives: {
		//       defaultSrc: [`'self'`],
		//       fontSrc: [`'self'`, 'fonts.scalar.com', 'data:'],
		//       imgSrc: [`'self'`, 'data:'],
		//       scriptSrc: [`'self'`, `https: 'unsafe-inline'`, `'unsafe-eval'`],
		//       styleSrc: [`'self'`, `'unsafe-inline'`, 'fonts.scalar.com']
		//     }
		//   }
		ReferrerPolicy: "strict-origin-when-cross-origin",
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
