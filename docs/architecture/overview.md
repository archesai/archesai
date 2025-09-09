# Architecture Documentation

This section covers the system architecture, design patterns, and structural documentation for
ArchesAI.

## Architecture Overview

- **[System Architecture](system-design.md)** - Complete system design and patterns
- **[Authentication](authentication.md)** - Authentication and authorization architecture
- **[Project Layout](project-layout.md)** - Directory structure and code organization

## Key Architectural Concepts

### Hexagonal Architecture

ArchesAI implements **Hexagonal Architecture** (Ports & Adapters) with **Domain-Driven Design**
principles, ensuring separation of concerns, testability, and business logic independence from
infrastructure.

### Domain-Driven Design

Each bounded context (auth, organizations, workflows, content) operates independently with:

- Own entities and business rules
- Dedicated database tables
- Separate API endpoints
- No cross-domain imports

### Code Generation

The project uses extensive code generation driven by:

- **OpenAPI specs** for types and HTTP handlers
- **SQL queries** for database operations
- **Custom templates** for repositories and adapters

## Architecture Diagrams

Interactive Mermaid diagrams will be added in upcoming iterations.
