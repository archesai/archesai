# Quickstart

Generate and run an API in 5 minutes.

## Install

```bash
go install github.com/archesai/archesai/cmd/archesai@latest
```

## Create OpenAPI Spec

Save as `bookstore.yaml`:

```yaml
openapi: 3.1.0
info:
  title: Bookstore API
  version: 1.0.0
paths:
  /books:
    get:
      operationId: listBooks
      tags: [books]
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
          description: Created
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
      operationId: getBook
      tags: [books]
      responses:
        "200":
          description: Book details
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Book"
    delete:
      operationId: deleteBook
      tags: [books]
      responses:
        "204":
          description: Deleted
components:
  schemas:
    Book:
      type: object
      required: [id, title, author, price]
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        author:
          type: string
        price:
          type: number
          format: float
        createdAt:
          type: string
          format: date-time
    BookInput:
      type: object
      required: [title, author, price]
      properties:
        title:
          type: string
        author:
          type: string
        price:
          type: number
          format: float
```

## Generate

```bash
archesai generate --spec bookstore.yaml --output bookstore-app
cd bookstore-app
```

## Run

```bash
archesai dev
# API at http://localhost:3001
```

## Test

```bash
# Create a book
curl -X POST http://localhost:3001/books \
  -H "Content-Type: application/json" \
  -d '{"title": "The Go Programming Language", "author": "Alan Donovan", "price": 39.99}'

# List books
curl http://localhost:3001/books
```

## What's Generated

```text
bookstore-app/
├── main.gen.go          # Entry point
├── spec/                # OpenAPI spec (bundled)
├── models/              # Book, BookInput structs
├── controllers/         # HTTP request handlers
├── application/         # Application layer (use cases)
├── repositories/        # Repository interfaces
├── bootstrap/           # App init, routes, DI container
└── infrastructure/
    ├── postgres/        # PostgreSQL migrations, queries, repos
    └── sqlite/          # SQLite migrations, queries, repos
```

## Next Steps

- Add a database: See [Getting Started](../getting-started.md#database-setup)
- Customize handlers: See [Custom Handlers](custom-handlers.md)
- Add x-codegen annotations: See [Code Generation](code-generation.md)
