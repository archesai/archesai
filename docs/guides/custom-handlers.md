# Custom Handlers

Generated handlers provide the basic structure. Add your business logic by editing files in the `handlers/` directory.

## Handler Structure

Generated handlers follow this pattern:

```go
// handlers/todo_handler.go
package handlers

type TodoHandler struct {
    repo repositories.TodoRepository
}

func NewTodoHandler(repo repositories.TodoRepository) *TodoHandler {
    return &TodoHandler{repo: repo}
}

func (h *TodoHandler) Create(ctx context.Context, input *models.TodoInput) (*models.Todo, error) {
    // Generated: basic create logic
    todo := &models.Todo{
        ID:        uuid.New(),
        Title:     input.Title,
        Completed: input.Completed,
        CreatedAt: time.Now(),
    }
    return todo, h.repo.Create(ctx, todo)
}
```

## Adding Custom Logic

### Validation

Add validation before saving:

```go
func (h *TodoHandler) Create(ctx context.Context, input *models.TodoInput) (*models.Todo, error) {
    // Add validation
    if len(input.Title) < 3 {
        return nil, errors.New("title must be at least 3 characters")
    }
    if len(input.Title) > 200 {
        return nil, errors.New("title must be less than 200 characters")
    }

    todo := &models.Todo{
        ID:        uuid.New(),
        Title:     input.Title,
        Completed: input.Completed,
        CreatedAt: time.Now(),
    }
    return todo, h.repo.Create(ctx, todo)
}
```

### Business Rules

Add domain-specific logic:

```go
func (h *TodoHandler) Complete(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
    todo, err := h.repo.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // Business rule: can't complete already completed todos
    if todo.Completed {
        return nil, errors.New("todo is already completed")
    }

    todo.Completed = true
    todo.CompletedAt = time.Now()

    if err := h.repo.Update(ctx, todo); err != nil {
        return nil, err
    }

    return todo, nil
}
```

### Adding Dependencies

Inject additional services:

```go
type TodoHandler struct {
    repo     repositories.TodoRepository
    notifier NotificationService  // Add new dependency
}

func NewTodoHandler(
    repo repositories.TodoRepository,
    notifier NotificationService,
) *TodoHandler {
    return &TodoHandler{
        repo:     repo,
        notifier: notifier,
    }
}

func (h *TodoHandler) Create(ctx context.Context, input *models.TodoInput) (*models.Todo, error) {
    todo := &models.Todo{
        ID:        uuid.New(),
        Title:     input.Title,
        CreatedAt: time.Now(),
    }

    if err := h.repo.Create(ctx, todo); err != nil {
        return nil, err
    }

    // Notify after creation
    h.notifier.Send(ctx, "New todo created: " + todo.Title)

    return todo, nil
}
```

## Custom Endpoints

Add new endpoints by creating new handler methods:

```go
// handlers/todo_handler.go

// GetStats returns todo statistics
func (h *TodoHandler) GetStats(ctx context.Context) (*models.TodoStats, error) {
    todos, err := h.repo.List(ctx, nil)
    if err != nil {
        return nil, err
    }

    stats := &models.TodoStats{
        Total:     len(todos),
        Completed: 0,
        Pending:   0,
    }

    for _, t := range todos {
        if t.Completed {
            stats.Completed++
        } else {
            stats.Pending++
        }
    }

    return stats, nil
}
```

Then register the route in your controller:

```go
// controllers/todo_controller.go

func (c *TodoController) RegisterRoutes(e *echo.Echo) {
    g := e.Group("/todos")
    g.GET("", c.List)
    g.POST("", c.Create)
    g.GET("/:id", c.Get)
    g.DELETE("/:id", c.Delete)
    g.GET("/stats", c.GetStats)  // Add new route
}
```

## Error Handling

Use typed errors for better error responses:

```go
var (
    ErrTodoNotFound = errors.New("todo not found")
    ErrTodoAlreadyCompleted = errors.New("todo already completed")
)

func (h *TodoHandler) Get(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
    todo, err := h.repo.Get(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrTodoNotFound
        }
        return nil, err
    }
    return todo, nil
}
```

## Testing Handlers

Use generated mocks for testing:

```go
func TestTodoHandler_Create(t *testing.T) {
    mockRepo := mocks.NewMockTodoRepository(t)

    handler := NewTodoHandler(mockRepo)

    input := &models.TodoInput{
        Title: "Test Todo",
    }

    mockRepo.EXPECT().
        Create(mock.Anything, mock.AnythingOfType("*models.Todo")).
        Return(nil)

    todo, err := handler.Create(context.Background(), input)

    assert.NoError(t, err)
    assert.Equal(t, "Test Todo", todo.Title)
    assert.False(t, todo.Completed)
}
```

## Best Practices

1. **Keep handlers focused** - One responsibility per handler method
2. **Validate early** - Check inputs before processing
3. **Use typed errors** - Makes error handling cleaner
4. **Inject dependencies** - Makes testing easier
5. **Don't modify generated files** - They get overwritten on regeneration
