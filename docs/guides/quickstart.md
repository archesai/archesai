# Quick Start Guide

Get up and running with Arches in 5 minutes! This guide shows you the fastest path to generating your first application.

## Prerequisites

You'll need:

- Go 1.22+ installed
- A terminal/command prompt
- 5 minutes of your time

## Install Arches

```bash
# Install the Arches CLI
go install github.com/archesai/archesai/cmd/archesai@latest

# Verify installation
archesai version
```

## Generate Your First App

### 1. Create an OpenAPI Spec

Save this as `bookstore.yaml`:

```yaml
openapi: 3.1.0
info:
  title: Bookstore API
  version: 1.0.0
  description: Simple bookstore management API

paths:
  /books:
    get:
      summary: List all books
      operationId: listBooks
      tags: [books]
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 10
      responses:
        "200":
          description: List of books
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Book"

    post:
      summary: Add a new book
      operationId: createBook
      tags: [books]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/BookInput"
      responses:
        "201":
          description: Book created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Book"

  /books/{id}:
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
          format: uuid

    get:
      summary: Get a book by ID
      operationId: getBook
      tags: [books]
      responses:
        "200":
          description: Book details
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Book"
        "404":
          description: Book not found

    delete:
      summary: Delete a book
      operationId: deleteBook
      tags: [books]
      responses:
        "204":
          description: Book deleted
        "404":
          description: Book not found

components:
  schemas:
    Book:
      type: object
      required: [id, title, author, isbn, price]
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
          minLength: 1
          maxLength: 200
        author:
          type: string
          minLength: 1
          maxLength: 100
        isbn:
          type: string
          pattern: "^[0-9]{13}$"
        price:
          type: number
          format: float
          minimum: 0
        publishedAt:
          type: string
          format: date
        createdAt:
          type: string
          format: date-time
        updatedAt:
          type: string
          format: date-time

    BookInput:
      type: object
      required: [title, author, isbn, price]
      properties:
        title:
          type: string
          minLength: 1
          maxLength: 200
        author:
          type: string
          minLength: 1
          maxLength: 100
        isbn:
          type: string
          pattern: "^[0-9]{13}$"
        price:
          type: number
          format: float
          minimum: 0
        publishedAt:
          type: string
          format: date
```

### 2. Generate the Application

```bash
# Generate the complete application
archesai generate --spec bookstore.yaml --output bookstore-app

# Navigate to your new app
cd bookstore-app
```

### 3. What Got Generated?

Your new `bookstore-app` directory contains:

```plaintext
bookstore-app/
├── main.go                 # Application entry point
├── models/                 # Data models (Book, BookInput)
├── controllers/            # HTTP handlers
├── handlers/               # Business logic
├── repositories/           # Database access
├── database/               # Migrations and connection
├── middleware/             # Auth, CORS, logging
├── client/                 # JavaScript/TypeScript SDK
├── docker/                 # Docker configuration
├── kubernetes/             # K8s manifests
└── config.yaml             # Application config
```

### 4. Start Your Application

```bash
# Option 1: If you have the full Arches repo
archesai dev

# Option 2: Direct Go run
go mod init bookstore-app
go mod tidy
go run main.go

# Your API is now running at http://localhost:3001
```

### 5. Test Your API

```bash
# Create a book
curl -X POST http://localhost:3001/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan Donovan",
    "isbn": "9780134190440",
    "price": 39.99
  }'

# List all books
curl http://localhost:3001/books

# Get a specific book (use the ID from the create response)
curl http://localhost:3001/books/{id}

# Delete a book
curl -X DELETE http://localhost:3001/books/{id}
```

## What's Next?

### Add a Database

```bash
# Start PostgreSQL with Docker
docker run -d \
  --name bookstore-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=bookstore \
  -p 5432:5432 \
  postgres:15

# Update your config.yaml with database settings
# Run migrations (check database/ directory)
```

### Customize Your App

1. **Add Business Logic**: Edit files in `handlers/` to add custom logic
2. **Modify Models**: Update your OpenAPI spec and regenerate
3. **Add Authentication**: Add security schemes to your OpenAPI spec
4. **Deploy**: Use the generated Docker/Kubernetes configs

### Learn More

- **[Full Tutorial](../getting-started.md)**: Comprehensive getting started guide
- **[CLI Reference](../cli-reference.md)**: All CLI commands explained
- **[Code Generation](code-generation.md)**: How generation works
- **[Examples](https://github.com/archesai/examples)**: More example applications

## Common Next Steps

### Add Authentication

Add this to your OpenAPI spec:

```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - bearerAuth: []
```

### Add More Endpoints

Extend your spec with more operations:

```yaml
paths:
  /books/{id}/reviews:
    post:
      summary: Add a book review
      # ... rest of the definition
```

### Use the TypeScript Client

```typescript
// In your frontend app
import { BookstoreClient } from "./bookstore-app/client";

const client = new BookstoreClient({
  baseURL: "http://localhost:3001",
});

// Use the type-safe client
const books = await client.listBooks({ limit: 10 });
const newBook = await client.createBook({
  title: "New Book",
  author: "Author Name",
  isbn: "1234567890123",
  price: 29.99,
});
```

### Deploy with Docker

```bash
# Build the image
docker build -t bookstore-api .

# Run the container
docker run -p 3001:3001 bookstore-api
```

## Troubleshooting

### Generation Failed?

```bash
# Validate your OpenAPI spec
npx @apidevtools/swagger-cli validate bookstore.yaml

# Check for common issues:
# - Missing required fields
# - Invalid references
# - Duplicate operation IDs
```

### Port Already in Use?

```bash
# Use a different port
export ARCHES_SERVER_PORT=8080
go run main.go
```

### Need Help?

- Check the [Troubleshooting Guide](../troubleshooting/common-issues.md)
- Open an [issue on GitHub](https://github.com/archesai/archesai/issues)
- Email <support@archesai.com>

---

**Congratulations!** 🎉 You've just generated and run your first Arches application. From here, you can:

- Extend your API with more endpoints
- Add authentication and authorization
- Deploy to production
- Generate clients for multiple languages

Happy building with Arches!
