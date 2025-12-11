# Coding Standards

This document defines the coding standards, naming conventions, and best practices for the Arches codebase.

## Table of Contents

- [Go Conventions](#go-conventions)
  - [Package Organization](#package-organization)
  - [Struct Design](#struct-design)
  - [Naming Conventions](#naming-conventions)
  - [Error Handling](#error-handling)
  - [Interfaces](#interfaces)
  - [Generics](#generics)
- [File Organization](#file-organization)
- [Code Style](#code-style)
- [Generated Code](#generated-code)
- [Prohibited Patterns](#prohibited-patterns)

---

## Go Conventions

### Package Organization

```text
pkg/           # Public packages that can be imported
internal/      # Private packages for internal use only
cmd/           # Command-line entry points
```

**Package naming:**

- Use short, lowercase, single-word names
- Avoid stuttering (e.g., `pkg/auth/auth.go` should export `Service`, not `AuthService`)
- Group related functionality in the same package

### Struct Design

#### Exported vs Unexported Fields

**Exported fields (uppercase)** for:

- JSON/YAML serialization
- Fields that need external access
- Configuration structs
- API request/response types
- Database models

```go
// Good - exported fields with JSON tags
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
}
```

**Unexported fields (lowercase)** for:

- Internal implementation details
- Dependencies injected via constructor
- State that should not be accessed directly

```go
// Good - unexported fields for internal state
type executorService[A any, B any] struct {
    repo       executorRepository
    builder    *Builder
    baseImages map[models.ExecutorLanguage]string
}
```

#### Constructors

Always use constructor functions (not factory methods or direct initialization):

```go
// Good - constructor function with dependency injection
func NewExecutorService[A any, B any](
    repo executorRepository,
    builder *Builder,
) ExecutorService[A, B] {
    return &executorService[A, B]{
        repo:    repo,
        builder: builder,
        baseImages: map[models.ExecutorLanguage]string{
            models.ExecutorLanguageNodejs: "archesai/runner-node:latest",
        },
    }
}

// Good - simple constructor
func NewParser(doc *OpenAPIDocument) *Parser {
    return &Parser{doc: doc}
}
```

#### No Getters/Setters

Avoid Java-style getters and setters. Access exported fields directly or use methods with meaningful names:

```go
// Bad
func (e *Errors) GetErrors() []FieldError { return e.errors }
func (e *Errors) SetErrors(errs []FieldError) { e.errors = errs }

// Good - meaningful method names
func (e Errors) HasErrors() bool { return len(e) > 0 }
func (e *Errors) Add(field, message string) { *e = append(*e, FieldError{...}) }
func (e *Errors) Merge(other Errors) { *e = append(*e, other...) }
```

#### Value vs Pointer Receivers

- **Value receivers** for small, immutable types and when the method doesn't modify state
- **Pointer receivers** for methods that modify state or for large structs

```go
// Value receiver - doesn't modify, returns computed value
func (e Errors) Error() string { ... }
func (e Errors) HasErrors() bool { return len(e) > 0 }

// Pointer receiver - modifies the struct
func (e *Errors) Add(field, message string) { ... }
func (e *Errors) Merge(other Errors) { ... }
```

### Naming Conventions

#### Variables and Functions

| Type               | Convention                                      | Example                          |
| ------------------ | ----------------------------------------------- | -------------------------------- |
| Local variables    | camelCase                                       | `userName`, `httpClient`         |
| Package-level vars | camelCase (unexported) or PascalCase (exported) | `defaultTimeout`, `MaxRetries`   |
| Functions          | PascalCase (exported) or camelCase (unexported) | `NewService`, `parseConfig`      |
| Acronyms           | Preserve case in middle, uppercase at start     | `userID`, `HTTPClient`, `apiURL` |

```go
// Good
var defaultTimeout = 30 * time.Second
func NewHTTPClient() *http.Client { ... }
func (s *service) getUserByID(id string) (*User, error) { ... }

// Bad
var DefaultTimeOut = 30 * time.Second  // inconsistent case
func NewHttpClient() *http.Client { ... }  // Http should be HTTP
```

#### Types and Interfaces

| Type         | Convention                         | Example                                     |
| ------------ | ---------------------------------- | ------------------------------------------- |
| Structs      | PascalCase, noun                   | `User`, `ExecutorService`, `Parser`         |
| Interfaces   | PascalCase, verb+er or descriptive | `Validator`, `HTTPError`, `ExecutorService` |
| Type aliases | PascalCase                         | `Errors`, `Config`                          |

```go
// Good - interface describes behavior
type Validator interface {
    Validate() Errors
}

// Good - interface with behavior suffix
type HTTPError interface {
    error
    StatusCode() int
    ProblemDetails(instance string) ProblemDetails
}
```

#### Constants and Enums

```go
// Good - grouped constants with type
type ExecutorLanguage string

const (
    ExecutorLanguageNodejs ExecutorLanguage = "nodejs"
    ExecutorLanguagePython ExecutorLanguage = "python"
    ExecutorLanguageGo     ExecutorLanguage = "go"
)

// Good - iota for sequential values
type SchemaType int

const (
    SchemaTypeEntity SchemaType = iota
    SchemaTypeRequest
    SchemaTypeResponse
)
```

#### Files

| Type            | Convention    | Example                                 |
| --------------- | ------------- | --------------------------------------- |
| Go source       | snake_case.go | `executor_service.go`, `http_errors.go` |
| Test files      | \*\_test.go   | `executor_service_test.go`              |
| Generated files | \*.gen.go     | `models.gen.go`, `handlers.gen.go`      |

### Error Handling

#### Custom Error Types

Create typed errors that implement the `error` interface:

```go
// Good - typed errors with behavior
type BadRequestError struct {
    Detail string
}

func (e BadRequestError) Error() string {
    return e.Detail
}

func (e BadRequestError) StatusCode() int {
    return http.StatusBadRequest
}
```

#### Error Wrapping

**Always wrap errors with context using `fmt.Errorf`:**

```go
// Good - wrap with context
exec, err := s.GetExecutor(ctx, executorID)
if err != nil {
    return zero, fmt.Errorf("get executor: %w", err)
}

// Good - chain of context
if err := os.WriteFile(executePath, []byte(exec.ExecuteCode), 0644); err != nil {
    return zero, fmt.Errorf("write execute code: %w", err)
}

// Bad - no context (loses call site information)
if err != nil {
    return err
}
```

**Current violations to fix:**

- `cmd/archesai/generate.go` - multiple `return err` without wrapping
- `pkg/database/migrate.go` - unwrapped errors
- `pkg/redis/queue.go` - unwrapped errors
- `pkg/server/server.go` - unwrapped errors

#### Error Checking Pattern

```go
// Good - handle errors immediately
result, err := doSomething()
if err != nil {
    return nil, fmt.Errorf("do something: %w", err)
}

// Good - multiple error checks with context
if err := step1(); err != nil {
    return fmt.Errorf("step1: %w", err)
}
if err := step2(); err != nil {
    return fmt.Errorf("step2: %w", err)
}
```

### Interfaces

#### Interface Location

Define interfaces where they are used, not where they are implemented:

```go
// Good - interface defined in consumer package
// pkg/executor/service.go
type executorRepository interface {
    Get(ctx context.Context, id uuid.UUID) (*models.Executor, error)
}

type executorService struct {
    repo executorRepository  // accepts interface
}
```

#### Small Interfaces

Prefer small, focused interfaces:

```go
// Good - small, focused interface
type Validator interface {
    Validate() Errors
}

// Good - composed interfaces when needed
type HTTPError interface {
    error
    StatusCode() int
    ProblemDetails(instance string) ProblemDetails
}
```

### Generics

Use generics for type-safe containers and services:

```go
// Good - generic service
type ExecutorService[A any, B any] interface {
    Execute(ctx context.Context, executorID string, input A) (ExecuteResult[B], error)
    ValidateInput(ctx context.Context, executorID string, input A) error
    ValidateOutput(ctx context.Context, executorID string, output B) error
}

// Good - generic result type
type ExecuteResult[B any] struct {
    Output          B
    ExecutionTimeMs int64
    Logs            string
}

// Good - generic helper function
func ValidateStruct[T any](v *T) Errors {
    if v == nil {
        return nil
    }
    if validator, ok := any(v).(Validator); ok {
        return validator.Validate()
    }
    return nil
}
```

**Note:** Using `any` as a generic type constraint is acceptable. Using `any` or `interface{}` as a concrete type is not.

---

## File Organization

### Standard Package Layout

```text
pkg/auth/
    auth.go          # Main types and interfaces
    service.go       # Service implementation
    api/             # OpenAPI spec and generated API types
    handlers/        # HTTP handlers (generated)
    implement/       # Custom handler implementations
    models/          # Data models (generated)
    routes/          # Route handlers (generated)
```

### File Naming

- One primary type per file when the type is substantial
- Group related small types in a single file
- Use descriptive names: `http_errors.go`, `schema_loader.go`

---

## Code Style

### Import Organization

Group imports in this order with blank lines between groups:

```go
import (
    // Standard library
    "context"
    "fmt"
    "net/http"

    // Third-party packages
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"

    // Internal packages
    "github.com/archesai/archesai/pkg/httputil"
    "github.com/archesai/archesai/internal/schema"
)
```

### Comments

#### Package Comments

```go
// Package validation provides validation utilities for struct validation.
package validation
```

#### Type Comments

```go
// FieldError represents a single field validation error.
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// HTTPError is an interface for errors that carry HTTP status information.
// Handlers can return these errors and the route layer will respond appropriately.
type HTTPError interface {
    error
    StatusCode() int
}
```

#### Function Comments

```go
// ValidateStruct validates a struct if it implements the Validator interface.
// Returns nil if the struct does not implement Validator or if validation passes.
func ValidateStruct[T any](v *T) Errors { ... }
```

### Line Length

- Aim for 100 characters per line maximum
- Break long function signatures across multiple lines

```go
// Good - broken across lines
func NewExecutorService[A any, B any](
    repo executorRepository,
    builder *Builder,
) ExecutorService[A, B] {
    ...
}
```

---

## Generated Code

### Conventions

- Generated files end with `.gen.go`
- Never manually modify generated files
- Use `//go:generate` directives

```go
//go:generate mockery --name=ExecutorRepository
```

### Using Generated Types

- Always use generated types for API requests/responses
- Don't create manual type definitions that duplicate generated ones

### Generated Code Exceptions

Generated code may contain patterns that would normally be prohibited:

- TODO comments as scaffolding markers (in handler stubs)
- `any` return types for event interfaces
- `interface{}` in database driver interfaces

These are acceptable in generated code but should not be replicated in hand-written code.

---

## Prohibited Patterns

### Never Use

| Pattern                        | Reason         | Alternative                                 |
| ------------------------------ | -------------- | ------------------------------------------- |
| `interface{}`                  | No type safety | Use generics or specific types              |
| `any` (as concrete type)       | No type safety | Use generics or specific types              |
| `map[string]any`               | No type safety | Define a proper struct                      |
| Hard-coded values in templates | Inflexible     | Use configuration or parameters             |
| Manual mocks                   | Inconsistent   | Use mockery to generate from interfaces     |
| TODO comments in code          | Gets forgotten | Create issues in project management         |
| Deprecated/legacy code         | Technical debt | Remove and use latest patterns              |
| Unwrapped errors               | Loses context  | Always use `fmt.Errorf("context: %w", err)` |

### Code Examples

```go
// Bad - using interface{}
func Process(data interface{}) interface{} { ... }

// Good - using generics
func Process[T any](data T) T { ... }

// Bad - using any as concrete type
var result any = getValue()
output := map[string]any{"key": value}

// Good - using specific types
var result string = getValue()
output := OutputData{Key: value}

// Bad - returning error without context
if err != nil {
    return err
}

// Good - wrapping error with context
if err != nil {
    return fmt.Errorf("failed to process data: %w", err)
}
```

### Current Violations to Address

**`interface{}` usage (non-generated):**

- `internal/codegen/handler_parser.go:249` - Returns `"interface{}"` as fallback

**TODO comments in production code:**

- `pkg/auth/service.go` - Contains multiple TODO comments for incomplete features
- `pkg/executor/service.go:215` - TODO for capturing stderr logs

---

## Tooling

### Required Tools

- `gofmt` - Code formatting
- `golangci-lint` - Linting
- `mockery` - Mock generation

### Makefile Commands

```bash
make format    # Format all Go code
make lint      # Run linter
make generate  # Generate code (mocks, etc.)
make test      # Run tests
```

### Pre-commit Checks

All code must pass:

1. `make format` - No formatting changes
2. `make lint` - No linting errors
3. `make test` - All tests pass
