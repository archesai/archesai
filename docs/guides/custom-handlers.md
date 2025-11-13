# Executor System

A type-safe, generic executor system for running code either locally or in isolated Docker containers with JSON Schema validation.

## Overview

The executor package provides a flexible way to execute code with two execution modes:

- **LocalExecutor**: Run code directly on the host (fast, no isolation)
- **ContainerExecutor**: Run code in Docker containers (secure, isolated, multi-language)

Both executors implement the same `Executor[A, B]` interface, ensuring type safety through generics and validating inputs/outputs using JSON Schema.

## Features

- **Type-safe execution** with Go generics
- **JSON Schema validation** for inputs and outputs
- **Multi-language support** (Python, Node.js, Go) - ContainerExecutor only
- **Security isolation** (network disabled, read-only FS, resource limits) - ContainerExecutor only
- **Resource management** (CPU, memory limits)
- **Timeout control**
- **Local execution** for fast, non-isolated operations

## Choosing an Executor

### Use LocalExecutor when

- Running trusted code (your own functions)
- Performance is critical (nanosecond execution vs seconds for containers)
- Testing and development
- Code doesn't need language/environment isolation
- Docker isn't available or desired

### Use ContainerExecutor when

- Running untrusted or third-party code
- Need multi-language support (Python, Node.js, Go)
- Require security isolation
- Need resource limits enforcement
- Running external tools (e.g., Orval, code generators)

## Usage

### LocalExecutor Example

For fast, local execution of trusted Go code:

```go
import (
    "context"
    "github.com/archesai/archesai/internal/infrastructure/executor"
)

// Define input/output types
type Input struct {
    Values []float64 `json:"values"`
}

type Output struct {
    Count int     `json:"count"`
    Sum   float64 `json:"sum"`
    Mean  float64 `json:"mean"`
}

// Define execution function
executeFunc := func(ctx context.Context, input Input) (Output, error) {
    sum := 0.0
    for _, v := range input.Values {
        sum += v
    }
    mean := sum / float64(len(input.Values))

    return Output{
        Count: len(input.Values),
        Sum:   sum,
        Mean:  mean,
    }, nil
}

// Create local executor
exec, err := executor.NewLocalExecutor[Input, Output](
    executeFunc,
    executor.LocalConfig{
        Timeout: 10 * time.Second,
        // Optional: schema validation
        SchemaIn:  inputSchema,
        SchemaOut: outputSchema,
    },
)
if err != nil {
    return err
}

// Execute (takes ~266 nanoseconds)
ctx := context.Background()
input := Input{Values: []float64{1, 2, 3, 4, 5}}
output, err := exec.Execute(ctx, input)
if err != nil {
    return err
}

log.Infof("Sum: %f, Mean: %f", output.Sum, output.Mean)
```

### ContainerExecutor Example

For isolated, multi-language execution:

```go
import (
    "context"
    "github.com/archesai/archesai/internal/infrastructure/executor"
)

// Define input/output types
type Input struct {
    Values []float64 `json:"values"`
}

type Output struct {
    Count int     `json:"count"`
    Sum   float64 `json:"sum"`
    Mean  float64 `json:"mean"`
}

// Create container executor
config := executor.Config{
    Image:       "archesai/runner-python:latest",
    DisableNet:  true,
    ReadOnlyFS:  true,
    MemoryBytes: 256 * 1024 * 1024, // 256MB
    Timeout:     10 * time.Second,
    SchemaIn:    inputSchema,  // JSON Schema bytes
    SchemaOut:   outputSchema, // JSON Schema bytes
}

exec, err := executor.NewContainerExecutor[Input, Output](config)
if err != nil {
    return err
}

// Execute
ctx := context.Background()
input := Input{Values: []float64{1, 2, 3, 4, 5}}
output, err := exec.Execute(ctx, input)
if err != nil {
    return err
}

log.Infof("Sum: %f, Mean: %f", output.Sum, output.Mean)
```

## Building Containers

Build the runner containers using Make:

```bash
# Build all runner containers
make build-runners

# Build individual runners
make build-runner-python
make build-runner-node
make build-runner-go

# Build generator containers (e.g., orval)
make build-generator-orval
```

## Testing

Run the executor tests:

```bash
# Run all executor tests (includes local and container tests)
go test -v ./internal/infrastructure/executor -timeout 60s

# Run only local executor tests (fast, no Docker required)
go test -v ./internal/infrastructure/executor -run TestLocal

# Run only container tests (requires Docker)
go test -v ./internal/infrastructure/executor -run TestNode

# Skip integration tests (only runs local executor tests)
go test -short ./internal/infrastructure/executor

# Run benchmarks
go test -bench=. ./internal/infrastructure/executor -benchmem
```

Performance comparison:

- **LocalExecutor**: ~266 ns/op (without validation), ~5.2 µs/op (with validation)
- **ContainerExecutor**: ~2-5 seconds/op (includes container startup)

## Container Protocol

Containers communicate via stdin/stdout using JSON:

### Input (stdin)

```json
{
  "schema_in": {
    /* JSON Schema */
  },
  "schema_out": {
    /* JSON Schema */
  },
  "input": {
    /* actual input data */
  }
}
```

### Output (stdout)

Success:

```json
{
  "ok": true,
  "output": {
    /* result data */
  }
}
```

Error:

```json
{
  "ok": false,
  "error": {
    "message": "error description",
    "code": "ERROR_CODE",
    "details": "stack trace or details"
  }
}
```

## Security

- **Network isolation**: `--network none`
- **Read-only filesystem**: `--read-only`
- **Non-root user**: Containers run as UID 1000
- **Resource limits**: Memory and CPU limits enforced
- **Timeouts**: Automatic termination on timeout

## Creating Custom Runners

### Option 1: Mount Custom Execute Module (Node.js)

For quick prototyping or one-off executions, mount a custom `execute.js` to the base node runner:

```javascript
// my-execute.js
export async function executeFunction(input) {
  // Your custom logic here
  return { result: input.value * 2 };
}
```

```bash
docker run -i --rm \
  -v $(pwd)/my-execute.js:/app/execute.js \
  archesai/runner-node:latest < input.json
```

### Option 2: Build Custom Generator Image

For reusable generators, extend the base runner:

```dockerfile
# Dockerfile
FROM archesai/runner-node:latest

# Install additional dependencies
RUN npm install --omit=dev your-package@version

# Copy custom execute module
COPY execute.js ./execute.js

# Entrypoint inherited from base
```

```bash
docker build -t my-custom-generator:latest .
docker run -i --rm my-custom-generator:latest < input.json
```

### Option 3: Other Languages

For Python or Go, implement the container protocol:

1. Read JSON from stdin
2. Parse request and validate schemas
3. Execute your logic
4. Return JSON response to stdout

See the example runners in `deployments/containers/runners/` for reference implementations.
