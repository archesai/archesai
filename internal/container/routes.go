package container

import (
	authhttp "github.com/archesai/archesai/internal/auth/adapters/http"
	contenthttp "github.com/archesai/archesai/internal/content/adapters/http"
	organizationshttp "github.com/archesai/archesai/internal/organizations/adapters/http"
	workflowshttp "github.com/archesai/archesai/internal/workflows/adapters/http"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all application routes with the Echo server
func (c *Container) RegisterRoutes(e *echo.Echo) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Register auth routes using StrictHandler pattern
	strictAuthHandler := authhttp.NewAuthStrictHandlerWithMiddleware(c.AuthHandler)
	authhttp.RegisterHandlers(v1, strictAuthHandler)

	// Register organizations routes using StrictHandler pattern
	strictOrganizationsHandler := organizationshttp.NewOrganizationStrictHandler(c.OrganizationsHandler)
	organizationshttp.RegisterHandlers(v1, strictOrganizationsHandler)

	// Register workflows routes using StrictHandler pattern
	strictWorkflowsHandler := workflowshttp.NewWorkflowStrictHandler(c.WorkflowsHandler)
	workflowshttp.RegisterHandlers(v1, strictWorkflowsHandler)

	// Register content routes using StrictHandler pattern
	strictContentHandler := contenthttp.NewContentStrictHandler(c.ContentHandler)
	contenthttp.RegisterHandlers(v1, strictContentHandler)
}
