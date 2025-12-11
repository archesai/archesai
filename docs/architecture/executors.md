# Executors

Executors are the core abstraction for business logic in Arches. Every operation (CreatePipeline,
GetUser, etc.) is an executor with typed input and output.

## The Interface

```go
// pkg/executor/executor.go
type Executor[I, O any] interface {
    Execute(ctx context.Context, input *I) (*O, error)
}
```

That's it. One method, strongly typed, context-aware.

## Implementing an Executor

Given this generated code:

```go
// pipelines/types.gen.go

type CreatePipelineInput struct {
    Name        string
    Description *string
}

type CreatePipelineOutput struct {
    Data Pipeline
}

// Type alias - CreatePipeline IS an Executor
type CreatePipeline = executor.Executor[CreatePipelineInput, CreatePipelineOutput]
```

You implement it:

```go
// pipelines/create_pipeline.go

type createPipeline struct {
    repo PipelineRepository
}

func NewCreatePipeline(repo PipelineRepository) *createPipeline {
    return &createPipeline{repo: repo}
}

func (c *createPipeline) Execute(ctx context.Context, input *CreatePipelineInput) (*CreatePipelineOutput, error) {
    pipeline := &Pipeline{
        ID:          uuid.New(),
        Name:        input.Name,
        Description: input.Description,
        CreatedAt:   time.Now(),
    }

    if err := c.repo.Save(ctx, pipeline); err != nil {
        return nil, err
    }

    return &CreatePipelineOutput{Data: *pipeline}, nil
}
```

## Wiring

```go
// main.go or bootstrap

func main() {
    db := postgres.NewDB(connString)
    repo := postgres.NewPipelineRepository(db)

    // Create executor - returns concrete type
    createPipeline := pipelines.NewCreatePipeline(repo)

    // HTTP handler accepts the interface
    httpHandler := pipelines.NewCreatePipelineHTTPHandler(createPipeline)

    mux := http.NewServeMux()
    pipelines.RegisterCreatePipelineRoute(mux, httpHandler)
}
```

## Patterns

### Unexported struct, exported constructor

```go
// Unexported - implementation detail
type createPipeline struct { ... }

// Exported constructor - returns concrete type
func NewCreatePipeline(repo PipelineRepository) *createPipeline {
    return &createPipeline{repo: repo}
}

// Callers assign to the interface type
var uc pipelines.CreatePipeline = pipelines.NewCreatePipeline(repo)
```

This pattern:

- Hides implementation details (unexported struct)
- Returns concrete type (Go idiom: "accept interfaces, return structs")
- Callers use the interface (can't even reference the concrete type outside the package)

### Returning errors

Use typed HTTP errors for proper status codes:

```go
func (c *createPipeline) Execute(ctx context.Context, input *CreatePipelineInput) (*CreatePipelineOutput, error) {
    existing, _ := c.repo.FindByName(ctx, input.Name)
    if existing != nil {
        return nil, httputil.ConflictError{Detail: "pipeline already exists"}
    }

    pipeline, err := c.repo.FindByID(ctx, input.ID)
    if err != nil {
        return nil, err // 500 Internal Server Error
    }
    if pipeline == nil {
        return nil, httputil.NotFoundError{Detail: "pipeline not found"}
    }

    // ...
}
```

Available error types:

- `httputil.BadRequestError` → 400
- `httputil.UnauthorizedError` → 401
- `httputil.ForbiddenError` → 403
- `httputil.NotFoundError` → 404
- `httputil.ConflictError` → 409

### Multiple dependencies

```go
type createPipelineStep struct {
    pipelines PipelineRepository
    steps     PipelineStepRepository
    tools     ToolRepository
}

func NewCreatePipelineStep(
    pipelines PipelineRepository,
    steps PipelineStepRepository,
    tools ToolRepository,
) *createPipelineStep {
    return &createPipelineStep{
        pipelines: pipelines,
        steps:     steps,
        tools:     tools,
    }
}
```

### Composing executors

Executors can call other executors:

```go
type createPipelineWithSteps struct {
    createPipeline     pipelines.CreatePipeline
    createPipelineStep pipelines.CreatePipelineStep
}

func (c *createPipelineWithSteps) Execute(
    ctx context.Context,
    input *CreatePipelineWithStepsInput,
) (*CreatePipelineWithStepsOutput, error) {
    // Create pipeline
    pipelineResult, err := c.createPipeline.Execute(ctx, &pipelines.CreatePipelineInput{
        Name: input.Name,
    })
    if err != nil {
        return nil, err
    }

    // Create steps
    for _, step := range input.Steps {
        _, err := c.createPipelineStep.Execute(ctx, &pipelines.CreatePipelineStepInput{
            PipelineID: pipelineResult.Data.ID,
            Name:       step.Name,
        })
        if err != nil {
            return nil, err
        }
    }

    return &CreatePipelineWithStepsOutput{Data: pipelineResult.Data}, nil
}
```

## Special Cases

### Container Execution

For sandboxed/isolated execution (e.g., running untrusted code):

```go
containerExec, err := executor.NewContainerExecutor[MyInput, MyOutput](executor.ContainerConfig{
    Image:       "my-executor:latest",
    Timeout:     30 * time.Second,
    MemoryBytes: 256 * 1024 * 1024,
    DisableNet:  true,
})

output, err := containerExec.Execute(ctx, &input)
```

### Local Execution with Schema Validation

For JSON Schema validation on input/output:

```go
localExec, err := executor.NewLocalExecutor[MyInput, MyOutput](
    executor.LocalConfig{
        Timeout:   10 * time.Second,
        SchemaIn:  inputSchemaBytes,
        SchemaOut: outputSchemaBytes,
    },
    func(ctx context.Context, input MyInput) (MyOutput, error) {
        // your logic
    },
)

output, err := localExec.Execute(ctx, input)
```

### Testing

Executors are easy to test - just call Execute:

```go
func TestCreatePipeline(t *testing.T) {
    repo := &mockPipelineRepository{}
    uc := NewCreatePipeline(repo)

    output, err := uc.Execute(context.Background(), &CreatePipelineInput{
        Name: "test-pipeline",
    })

    assert.NoError(t, err)
    assert.Equal(t, "test-pipeline", output.Data.Name)
}
```

Or mock the executor itself:

```go
type mockCreatePipeline struct {
    result *CreatePipelineOutput
    err    error
}

func (m *mockCreatePipeline) Execute(ctx context.Context, input *CreatePipelineInput) (*CreatePipelineOutput, error) {
    return m.result, m.err
}
```

## Generated Structure

```text
app/
  pipelines/
    types.gen.go          # Entities + Executor type aliases + DTOs
    http.gen.go           # HTTP handlers
    create_pipeline.go    # Your implementation
    get_pipeline.go       # Your implementation

  adapters/
    postgres/
      migrations/
      pipelines.gen.go    # Repository implementations
```

## Summary

1. **Executor\[I, O\]** is the only interface you need
2. **Implement directly** - just a struct with an Execute method
3. **Unexported struct** - hide implementation details
4. **Return concrete types** - let callers assign to interface
5. **Compose freely** - executors can call other executors
6. **Test easily** - call Execute or mock the interface
