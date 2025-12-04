# Getting Started

This guide covers installation and basic usage of Arches.

## Prerequisites

- Go 1.22+
- PostgreSQL 15+ (or Docker)
- Node.js 20+ with pnpm (for platform UI)

## Installation

### Go Install (Recommended)

```bash
go install github.com/archesai/archesai/cmd/archesai@latest
```

### From Source

```bash
git clone https://github.com/archesai/archesai.git
cd archesai
make build
# Binary at ./bin/archesai
```

### Docker

```bash
docker run -p 3000:3000 -p 3001:3001 ghcr.io/archesai/archesai:latest
```

## Create Your First Application

### 1. Write an OpenAPI Spec

Create `api.yaml`:

```yaml
openapi: 3.1.0
info:
  title: My API
  version: 1.0.0
paths:
  /todos:
    get:
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
          description: Created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Todo"
components:
  schemas:
    Todo:
      type: object
      required: [id, title, completed]
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        completed:
          type: boolean
          default: false
        createdAt:
          type: string
          format: date-time
    TodoInput:
      type: object
      required: [title]
      properties:
        title:
          type: string
        completed:
          type: boolean
          default: false
```

### 2. Generate the Application

```bash
archesai generate --spec api.yaml --output ./myapp
```

### 3. Generated Structure

```text
myapp/
├── main.gen.go              # Entry point
├── spec/                    # OpenAPI spec (bundled)
├── models/                  # Go structs from schemas
├── controllers/             # HTTP request handlers
├── application/             # Application layer (use cases)
├── repositories/            # Repository interfaces
├── bootstrap/               # App init, routes, DI container
└── infrastructure/
    ├── postgres/            # PostgreSQL migrations, queries, repos
    └── sqlite/              # SQLite migrations, queries, repos
```

### 4. Run the Application

```bash
cd myapp
archesai dev
```

- API Server: <http://localhost:3001>
- Platform UI: <http://localhost:3000>

### 5. Test Your API

```bash
# Create a todo
curl -X POST http://localhost:3001/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Arches"}'

# List todos
curl http://localhost:3001/todos
```

## Database Setup

For persistence, start PostgreSQL:

```bash
docker run -d \
  --name arches-postgres \
  -e POSTGRES_DB=archesdb \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15
```

## Next Steps

- [Quickstart](guides/quickstart.md) - 5-minute tutorial
- [CLI Reference](cli-reference.md) - All commands
- [Code Generation](guides/code-generation.md) - x-codegen extensions
- [Custom Handlers](guides/custom-handlers.md) - Adding business logic
- [Configuration](configuration.md) - Config options
