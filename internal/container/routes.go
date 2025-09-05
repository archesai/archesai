package container

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func (c *Container) RegisterRoutes(e *echo.Echo) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes
	c.AuthHandler.RegisterRoutes(v1)

	// Register organizations routes
	orgGroup := v1.Group("/organizations")
	c.OrganizationsHandler.RegisterRoutes(orgGroup)

	// Register workflows routes
	workflowsGroup := v1.Group("/workflows")
	c.WorkflowsHandler.RegisterRoutes(workflowsGroup)

	// Register content routes
	contentGroup := v1.Group("/content")
	c.ContentHandler.RegisterRoutes(contentGroup)
}
