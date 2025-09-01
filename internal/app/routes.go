package app

import (
	"github.com/archesai/archesai/gen/api/features/auth/users"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func RegisterRoutes(e *echo.Echo, container *Container) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes (custom endpoints not in OpenAPI)
	container.AuthHandler.RegisterRoutes(v1)

	// Register OpenAPI-generated user routes
	// These implement the users.ServerInterface
	userGroup := v1.Group("/auth")
	users.RegisterHandlers(userGroup, container.AuthHandler)

	// TODO: Register other feature routes as they are implemented
	// Example for intelligence feature:
	// intelligenceGroup := v1.Group("/intelligence")
	// artifacts.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// labels.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// pipelines.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// runs.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
	// tools.RegisterHandlers(intelligenceGroup, container.IntelligenceHandler)
}
