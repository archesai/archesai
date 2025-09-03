// Package main provides domain generation functionality.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type DomainConfig struct {
	Name        string   // e.g., "billing"
	Package     string   // e.g., "billing"
	Description string   // e.g., "Billing and subscription management"
	Tables      []string // e.g., ["subscription", "invoice", "payment"]
	HasAuth     bool     // Whether domain needs auth middleware
	HasEvents   bool     // Whether domain uses events
}

func main() {
	var config DomainConfig
	var tablesStr string

	flag.StringVar(&config.Name, "name", "", "Domain name (e.g., billing)")
	flag.StringVar(&config.Description, "desc", "", "Domain description")
	flag.StringVar(&tablesStr, "tables", "", "Comma-separated list of database tables")
	flag.BoolVar(&config.HasAuth, "auth", false, "Include auth middleware")
	flag.BoolVar(&config.HasEvents, "events", false, "Include domain events")
	flag.Parse()

	if config.Name == "" {
		log.Fatal("Domain name is required: -name=billing")
	}

	config.Package = strings.ToLower(config.Name)
	if config.Description == "" {
		config.Description = fmt.Sprintf("%s domain functionality", config.Name)
	}

	if tablesStr != "" {
		config.Tables = strings.Split(tablesStr, ",")
		for i := range config.Tables {
			config.Tables[i] = strings.TrimSpace(config.Tables[i])
		}
	}

	if err := generateDomain(config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Domain '%s' generated successfully!\n", config.Name)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Add converter configuration to internal/domains/converters.yaml")
	fmt.Println("2. Wire dependencies in internal/app/deps.go")
	fmt.Println("3. Run 'make generate' to generate converters")
	fmt.Println("4. Implement business logic in service.go")
}

func generateDomain(config DomainConfig) error {
	domainPath := filepath.Join("internal", "domains", config.Package)
	convertersPath := filepath.Join(domainPath, "converters")

	// Create directories
	if err := os.MkdirAll(convertersPath, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files
	files := map[string]string{
		filepath.Join(domainPath, fmt.Sprintf("%s.go", config.Package)): domainTemplate,
		filepath.Join(domainPath, "entities.go"):                        entitiesTemplate,
		filepath.Join(domainPath, "service.go"):                         serviceTemplate,
		filepath.Join(domainPath, "repository.go"):                      repositoryTemplate,
		filepath.Join(domainPath, "handler.go"):                         handlerTemplate,
	}

	if config.HasAuth {
		files[filepath.Join(domainPath, "middleware.go")] = middlewareTemplate
	}

	if config.HasEvents {
		files[filepath.Join(domainPath, "events.go")] = eventsTemplate
	}

	for path, tmplStr := range files {
		if err := generateFile(path, tmplStr, config); err != nil {
			return fmt.Errorf("failed to generate %s: %w", path, err)
		}
	}

	// Create empty .gitkeep in converters directory
	gitkeepPath := filepath.Join(convertersPath, ".gitkeep")
	if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create .gitkeep: %w", err)
	}

	return nil
}

func generateFile(path, tmplStr string, config DomainConfig) error {
	funcMap := template.FuncMap{
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New(filepath.Base(path)).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	return tmpl.Execute(file, config)
}

const domainTemplate = `// Package {{.Package}} provides {{.Description}}.
package {{.Package}}

// ContextKey is a type for context keys specific to this domain.
type ContextKey string

const (
	// Add domain-specific context keys here
	// Example: {{.Name}}ContextKey ContextKey = "{{.Package}}"
)

// Domain-specific constants
const (
	// Add domain constants here
)
`

const entitiesTemplate = `package {{.Package}}

import (
	"errors"
	"time"

	"github.com/archesai/archesai/internal/generated/api"
)

// Domain-specific errors
var (
	ErrNotFound          = errors.New("{{.Package}}: not found")
	ErrInvalidInput      = errors.New("{{.Package}}: invalid input")
	ErrUnauthorized      = errors.New("{{.Package}}: unauthorized")
	ErrAlreadyExists     = errors.New("{{.Package}}: already exists")
)

{{range .Tables}}
// {{title .}} represents a {{.}} in the domain.
// It may extend an API type or be domain-specific.
type {{title .}} struct {
	ID        string    ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	
	// Add domain-specific fields here
	// Consider embedding api.{{title .}}Entity if it exists
}
{{end}}

// Add validation methods, business rules, and domain logic here
`

const serviceTemplate = `package {{.Package}}

import (
	"context"
	"errors"
	"fmt"

	"github.com/archesai/archesai/internal/infrastructure/config"
)

// Repository defines the data access interface for this domain.
// Following the "interface defined by consumer" pattern.
type Repository interface {
	{{range .Tables}}
	// {{title .}} operations
	Get{{title .}}(ctx context.Context, id string) (*{{title .}}, error)
	List{{title .}}s(ctx context.Context, limit, offset int) ([]*{{title .}}, error)
	Create{{title .}}(ctx context.Context, {{lower .}} *{{title .}}) error
	Update{{title .}}(ctx context.Context, {{lower .}} *{{title .}}) error
	Delete{{title .}}(ctx context.Context, id string) error
	{{end}}
}

// Service contains the business logic for the {{.Package}} domain.
type Service struct {
	repo   Repository
	config *config.Config
}

// NewService creates a new {{.Package}} service.
func NewService(repo Repository, config *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

{{range .Tables}}
// Get{{title .}} retrieves a {{.}} by ID.
func (s *Service) Get{{title .}}(ctx context.Context, id string) (*{{title .}}, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	{{lower .}}, err := s.repo.Get{{title .}}(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get {{.}}: %w", err)
	}

	return {{lower .}}, nil
}

// Create{{title .}} creates a new {{.}}.
func (s *Service) Create{{title .}}(ctx context.Context, {{lower .}} *{{title .}}) error {
	// Add validation and business logic here
	
	if err := s.repo.Create{{title .}}(ctx, {{lower .}}); err != nil {
		return fmt.Errorf("failed to create {{.}}: %w", err)
	}

	return nil
}

// Update{{title .}} updates an existing {{.}}.
func (s *Service) Update{{title .}}(ctx context.Context, {{lower .}} *{{title .}}) error {
	// Add validation and business logic here
	
	if err := s.repo.Update{{title .}}(ctx, {{lower .}}); err != nil {
		return fmt.Errorf("failed to update {{.}}: %w", err)
	}

	return nil
}

// Delete{{title .}} deletes a {{.}} by ID.
func (s *Service) Delete{{title .}}(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	if err := s.repo.Delete{{title .}}(ctx, id); err != nil {
		return fmt.Errorf("failed to delete {{.}}: %w", err)
	}

	return nil
}
{{end}}
`

const repositoryTemplate = `package {{.Package}}

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/archesai/archesai/internal/generated/database/postgresql"
)

// Compile-time check that PostgresRepository implements Repository.
var _ Repository = (*PostgresRepository)(nil)

// PostgresRepository is the PostgreSQL implementation of Repository.
type PostgresRepository struct {
	q postgresql.Querier
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(q postgresql.Querier) *PostgresRepository {
	return &PostgresRepository{
		q: q,
	}
}

{{range .Tables}}
// Get{{title .}} retrieves a {{.}} by ID from the database.
func (r *PostgresRepository) Get{{title .}}(ctx context.Context, id string) (*{{title .}}, error) {
	// TODO: Implement using sqlc-generated query
	// Example:
	// db{{title .}}, err := r.q.Get{{title .}}(ctx, id)
	// if err != nil {
	//     if errors.Is(err, sql.ErrNoRows) {
	//         return nil, ErrNotFound
	//     }
	//     return nil, err
	// }
	// return convert{{title .}}FromDB(db{{title .}}), nil
	
	return nil, fmt.Errorf("not implemented")
}

// List{{title .}}s retrieves a list of {{.}}s from the database.
func (r *PostgresRepository) List{{title .}}s(ctx context.Context, limit, offset int) ([]*{{title .}}, error) {
	// TODO: Implement using sqlc-generated query
	return nil, fmt.Errorf("not implemented")
}

// Create{{title .}} creates a new {{.}} in the database.
func (r *PostgresRepository) Create{{title .}}(ctx context.Context, {{lower .}} *{{title .}}) error {
	// TODO: Implement using sqlc-generated query
	return fmt.Errorf("not implemented")
}

// Update{{title .}} updates an existing {{.}} in the database.
func (r *PostgresRepository) Update{{title .}}(ctx context.Context, {{lower .}} *{{title .}}) error {
	// TODO: Implement using sqlc-generated query
	return fmt.Errorf("not implemented")
}

// Delete{{title .}} deletes a {{.}} from the database.
func (r *PostgresRepository) Delete{{title .}}(ctx context.Context, id string) error {
	// TODO: Implement using sqlc-generated query
	return fmt.Errorf("not implemented")
}
{{end}}

// Helper functions for converting between database and domain types
// These will be replaced by generated converters once configured

{{range .Tables}}
func convert{{title .}}FromDB(db *postgresql.{{title .}}) *{{title .}} {
	// TODO: Implement conversion
	return nil
}

func convert{{title .}}ToDB({{lower .}} *{{title .}}) *postgresql.{{title .}} {
	// TODO: Implement conversion
	return nil
}
{{end}}
`

const handlerTemplate = `package {{.Package}}

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/archesai/archesai/internal/generated/api"
)

// Handler handles HTTP requests for the {{.Package}} domain.
type Handler struct {
	service *Service
}

// NewHandler creates a new HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers the HTTP routes for this domain.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	// Register routes here
	// Example:
	{{range .Tables}}
	// e.GET("/api/v1/{{$.Package}}/{{.}}s", h.List{{title .}}s)
	// e.GET("/api/v1/{{$.Package}}/{{.}}s/:id", h.Get{{title .}})
	// e.POST("/api/v1/{{$.Package}}/{{.}}s", h.Create{{title .}})
	// e.PUT("/api/v1/{{$.Package}}/{{.}}s/:id", h.Update{{title .}})
	// e.DELETE("/api/v1/{{$.Package}}/{{.}}s/:id", h.Delete{{title .}})
	{{end}}
}

// Implement OpenAPI ServerInterface methods here
// These methods should match the generated interface from oapi-codegen

{{range .Tables}}
// Get{{title .}} handles GET /{{$.Package}}/{{.}}s/:id
func (h *Handler) Get{{title .}}(c echo.Context) error {
	id := c.Param("id")
	
	{{lower .}}, err := h.service.Get{{title .}}(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get {{.}}")
	}

	return c.JSON(http.StatusOK, {{lower .}})
}

// List{{title .}}s handles GET /{{$.Package}}/{{.}}s
func (h *Handler) List{{title .}}s(c echo.Context) error {
	// TODO: Parse query parameters for pagination
	limit := 20
	offset := 0

	{{.}}s, err := h.service.List{{title .}}s(c.Request().Context(), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list {{.}}s")
	}

	return c.JSON(http.StatusOK, {{.}}s)
}

// Create{{title .}} handles POST /{{$.Package}}/{{.}}s
func (h *Handler) Create{{title .}}(c echo.Context) error {
	var req {{title .}}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.service.Create{{title .}}(c.Request().Context(), &req); err != nil {
		if errors.Is(err, ErrInvalidInput) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create {{.}}")
	}

	return c.JSON(http.StatusCreated, req)
}

// Update{{title .}} handles PUT /{{$.Package}}/{{.}}s/:id
func (h *Handler) Update{{title .}}(c echo.Context) error {
	id := c.Param("id")
	
	var req {{title .}}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	req.ID = id

	if err := h.service.Update{{title .}}(c.Request().Context(), &req); err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if errors.Is(err, ErrInvalidInput) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update {{.}}")
	}

	return c.JSON(http.StatusOK, req)
}

// Delete{{title .}} handles DELETE /{{$.Package}}/{{.}}s/:id
func (h *Handler) Delete{{title .}}(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.Delete{{title .}}(c.Request().Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete {{.}}")
	}

	return c.NoContent(http.StatusNoContent)
}
{{end}}
`

const middlewareTemplate = `package {{.Package}}

import (
	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides authentication middleware for this domain.
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: Implement domain-specific auth logic
			// Example: Check if user has access to this domain's resources
			
			// Get user from context (set by global auth middleware)
			// user := auth.GetUserFromContext(c.Request().Context())
			// if user == nil {
			//     return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			// }
			
			// Check domain-specific permissions
			// if !hasPermission(user, resource) {
			//     return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			// }
			
			return next(c)
		}
	}
}

// Add other domain-specific middleware here
// Examples:
// - RateLimitMiddleware
// - ValidationMiddleware
// - LoggingMiddleware
`

const eventsTemplate = `package {{.Package}}

import (
	"context"
	"time"
)

// EventType represents the type of domain event.
type EventType string

const (
	{{range .Tables}}
	Event{{title .}}Created EventType = "{{$.Package}}.{{lower .}}.created"
	Event{{title .}}Updated EventType = "{{$.Package}}.{{lower .}}.updated"
	Event{{title .}}Deleted EventType = "{{$.Package}}.{{lower .}}.deleted"
	{{end}}
)

// Event represents a domain event.
type Event struct {
	ID         string                 ` + "`json:\"id\"`" + `
	Type       EventType              ` + "`json:\"type\"`" + `
	AggregateID string                ` + "`json:\"aggregate_id\"`" + `
	Timestamp  time.Time              ` + "`json:\"timestamp\"`" + `
	Data       map[string]interface{} ` + "`json:\"data\"`" + `
	UserID     string                 ` + "`json:\"user_id,omitempty\"`" + `
}

// EventHandler handles domain events.
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventPublisher publishes domain events.
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

// Add specific event types with strongly-typed data

{{range .Tables}}
// {{title .}}CreatedEvent is emitted when a {{.}} is created.
type {{title .}}CreatedEvent struct {
	Event
	{{title .}} *{{title .}} ` + "`json:\"{{lower .}}\"`" + `
}

// {{title .}}UpdatedEvent is emitted when a {{.}} is updated.
type {{title .}}UpdatedEvent struct {
	Event
	{{title .}}   *{{title .}} ` + "`json:\"{{lower .}}\"`" + `
	ChangedFields []string       ` + "`json:\"changed_fields\"`" + `
}

// {{title .}}DeletedEvent is emitted when a {{.}} is deleted.
type {{title .}}DeletedEvent struct {
	Event
	{{title .}}ID string ` + "`json:\"{{lower .}}_id\"`" + `
}
{{end}}
`
