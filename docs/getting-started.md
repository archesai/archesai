# Getting Started with Arches

Welcome to Arches! This guide will help you get up and running with the Arches code generation platform.

## What is Arches?

Arches is a powerful code generation platform that transforms OpenAPI specifications
into complete, production-ready applications. It generates backend APIs, database layers,
authentication systems, and deployment configurations - all from your OpenAPI schema.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.22+** - [Install Go](https://go.dev/doc/install)
- **Node.js 20+** and pnpm - [Install Node.js](https://nodejs.org/)
- **Docker** (optional) - For containerized development
- **PostgreSQL** - For database (or use Docker)
- **Redis** (optional) - For caching and sessions

## Installation

### Option 1: Install from Source

```bash
# Clone the repository
git clone https://github.com/archesai/archesai.git
cd archesai

# Install dependencies
make deps

# Build the CLI
make build

# The binary is now available at ./bin/archesai
```

### Option 2: Go Install

```bash
# Install directly with go
go install github.com/archesai/archesai/cmd/archesai@latest
```

### Option 3: Docker

```bash
# Run with Docker
docker run -p 3000:3000 -p 3001:3001 ghcr.io/archesai/archesai:latest
```

## Your First Application

Let's create a simple TODO application to demonstrate Arches' capabilities.

### Step 1: Create an OpenAPI Specification

Create a file named `todo-api.yaml`:

```yaml
openapi: 3.1.0
info:
  title: Todo API
  version: 1.0.0
  description: A simple TODO list API

paths:
  /todos:
    get:
      summary: List all todos
      operationId: listTodos
      tags: [todos]
      responses:
        "200":
          description: List of todos
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Todo"

    post:
      summary: Create a new todo
      operationId: createTodo
      tags: [todos]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/TodoInput"
      responses:
        "201":
          description: Created todo
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Todo"

  /todos/{id}:
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
          format: uuid

    get:
      summary: Get a todo by ID
      operationId: getTodo
      tags: [todos]
      responses:
        "200":
          description: Todo details
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Todo"
        "404":
          description: Todo not found

    put:
      summary: Update a todo
      operationId: updateTodo
      tags: [todos]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/TodoInput"
      responses:
        "200":
          description: Updated todo
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Todo"

    delete:
      summary: Delete a todo
      operationId: deleteTodo
      tags: [todos]
      responses:
        "204":
          description: Todo deleted
        "404":
          description: Todo not found

components:
  schemas:
    Todo:
      type: object
      required: [id, title, completed, createdAt]
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
          minLength: 1
          maxLength: 200
        description:
          type: string
          maxLength: 1000
        completed:
          type: boolean
          default: false
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time

    TodoInput:
      type: object
      required: [title]
      properties:
        title:
          type: string
          minLength: 1
          maxLength: 200
        description:
          type: string
          maxLength: 1000
        completed:
          type: boolean
          default: false
```

### Step 2: Generate the Application

```bash
# Generate the complete application
archesai generate openapi todo-api.yaml --output ./todo-app

# Navigate to the generated app
cd todo-app
```

### Step 3: Review Generated Code

Arches has generated:

```plaintext
todo-app/
├── models/           # Data models from schemas
├── repositories/     # Database access layer
├── controllers/      # HTTP request handlers
├── handlers/         # Business logic
├── database/         # Migrations and connection
├── client/           # TypeScript/JavaScript SDK
├── docker/           # Docker configuration
├── kubernetes/       # K8s manifests
└── main.go           # Application entry point
```

### Step 4: Start Development Server

```bash
# Start with hot reload
archesai dev

# Or use the Makefile
make dev-all
```

The application is now running:

- **API Server**: <http://localhost:3001>
- **Platform UI**: <http://localhost:3000>

### Step 5: Test Your API

```bash
# Create a todo
curl -X POST http://localhost:3001/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Arches", "description": "Build amazing APIs"}'

# List todos
curl http://localhost:3001/todos

# Get a specific todo
curl http://localhost:3001/todos/{id}
```

## Development Workflow

### Making Changes

1. **Update OpenAPI Spec**: Modify your `todo-api.yaml`
2. **Regenerate Code**: Run `archesai generate openapi todo-api.yaml --output ./todo-app`
3. **Hot Reload**: Changes are automatically picked up if using `archesai dev`

### Adding Custom Logic

Generated handlers are extensible. Add your business logic in the `handlers/` directory:

```go
// handlers/todo_custom.go
package handlers

func (h *TodoHandler) ValidateTodo(todo *models.Todo) error {
    // Add custom validation logic
    if todo.Title == "" {
        return errors.New("title cannot be empty")
    }
    return nil
}
```

### Database Setup

```bash
# Using Docker
docker run -d \
  --name arches-postgres \
  -e POSTGRES_DB=archesdb \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15

# Run migrations (auto-generated)
make db-migrate
```

## Configuration

Create a `.archesai.yaml` configuration file:

```yaml
app:
  name: todo-app
  version: 1.0.0

server:
  host: localhost
  port: 3001
  cors:
    enabled: true
    origins: ["http://localhost:3000"]

database:
  driver: postgres
  host: localhost
  port: 5432
  name: tododb
  user: postgres
  password: postgres

auth:
  enabled: true
  jwt_secret: your-secret-key-change-in-production
  token_expiry: 24h
```

## Documentation Structure

### 🏗️ [Architecture](architecture/system-design.md)

Learn about the system design, patterns, and overall architecture of Arches.

### 🚀 [Development](guides/development.md)

Everything you need to know about setting up your development environment and contributing to the project.

### 📚 [CLI Reference](cli-reference.md)

Complete reference for all CLI commands.

### 🎯 [Features](features/overview.md)

Detailed guides for platform features.

### 🐳 [Deployment](deployment/overview.md)

Production deployment guides including Docker, Kubernetes, and infrastructure setup.

### 🔧 [Troubleshooting](troubleshooting/common-issues.md)

Common issues, debugging guides, and solutions.

## Next Steps

Now that you have a working application:

1. **Explore the CLI**: See all available commands with `archesai --help`
2. **Read the Docs**:
   - [CLI Reference](cli-reference.md) - Complete command documentation
   - [Code Generation Guide](guides/code-generation.md) - How generation works
   - [Development Guide](guides/development.md) - Advanced development workflows
3. **Deploy Your App**: Check out the [Deployment Guide](deployment/overview.md)
4. **Join the Community**: Report issues and contribute on [GitHub](https://github.com/archesai/archesai)

## Common Issues

### Port Already in Use

If ports 3000 or 3001 are already in use:

```bash
# Change ports via environment variables
export ARCHES_SERVER_PORT=8001
export ARCHES_PLATFORM_PORT=8000
archesai dev
```

### Database Connection Failed

Ensure PostgreSQL is running:

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Or start it
docker-compose up -d postgres
```

### Missing Dependencies

```bash
# Install all dependencies
make deps

# Or separately
make deps-go   # Go dependencies
make deps-ts   # Node.js dependencies
```

## Getting Help

- **Documentation**: [Full documentation](https://archesai.com/docs)
- **GitHub Issues**: [Report bugs or request features](https://github.com/archesai/archesai/issues)
- **Email Support**: <support@archesai.com>
- **Contributing**: See our [Contributing Guide](contributing.md)

## License

See [LICENSE](https://github.com/archesai/archesai/blob/main/LICENSE) file for details.
