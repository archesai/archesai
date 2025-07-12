# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

Arches AI is a comprehensive data processing platform built with a modern monorepo architecture:

- **Backend**: NestJS API with modular architecture using workspace packages
- **Frontend**: React with TanStack Router and TanStack Query for state management
- **Database**: PostgreSQL with Drizzle ORM and vector extensions for embeddings
- **Infrastructure**: Kubernetes-based deployment with Skaffold for development
- **Build System**: Turborepo with pnpm workspaces

### Key Applications

- `apps/api`: NestJS backend API server
- `apps/platform`: React frontend application using TanStack Start
- `e2e/api-e2e`: Jest-based API end-to-end tests
- `e2e/platform-e2e`: Playwright-based frontend end-to-end tests

### Core Packages

- `packages/auth`: Authentication & authorization (JWT, OAuth, sessions)
- `packages/billing`: Stripe integration for subscriptions
- `packages/core`: Shared utilities, config, logging, database
- `packages/database`: Drizzle ORM schema and migrations
- `packages/schemas`: Domain entities and business logic
- `packages/intelligence`: AI services (LLM, speech, audio, scraping)
- `packages/orchestration`: Pipeline management and data processing
- `packages/storage`: File handling and cloud storage
- `packages/ui`: Shared React components using shadcn/ui
- `packages/client`: Auto-generated API client from OpenAPI spec

## Development Commands

### Root Level Commands

- `pnpm dev`: Start all development servers with hot reload
- `pnpm dev:api`: Start only the API server
- `pnpm dev:platform`: Start only the platform frontend
- `pnpm build`: Build all packages and applications
- `pnpm lint`: Run ESLint across all packages
- `pnpm lint:fix`: Auto-fix linting issues
- `pnpm typecheck`: Type check all TypeScript code
- `pnpm format`: Check code formatting with Prettier
- `pnpm format:fix`: Auto-fix formatting issues
- `pnpm clean`: Clean all build artifacts and node_modules
- `pnpm clean:workspaces`: Clean individual workspace build artifacts

### Testing Commands

- API E2E tests: `cd e2e/api-e2e && pnpm test`
- Platform E2E tests: `cd e2e/platform-e2e && pnpm test`
- Unit tests: Run `pnpm test` in individual package directories

### API Development

- The API uses path imports with `#app/*` and `#utils/*` prefixes
- Configuration is loaded from Kubernetes config files in `deploy/kubernetes/base/`
- Main entry point: `apps/api/src/main.ts`
- Module structure follows NestJS conventions with dependency injection

### Platform Development

- React application using TanStack Start for SSR capabilities
- Uses path imports: `#app/*`, `#components/*`, `#lib/*`, `#router`
- Routing is file-based with TanStack Router
- State management via TanStack Query with React Query

## Key Architecture Patterns

### Modular NestJS Structure

The API follows a modular architecture where each domain (auth, billing, orchestration, etc.) is a separate package with its own:

- Controllers for HTTP endpoints
- Services for business logic
- Repositories for data access
- DTOs for request/response validation
- Modules for dependency injection setup

### Database Schema

- Uses Drizzle ORM with PostgreSQL
- Vector extensions for embedding storage
- Migration files in `packages/database/migrations/`
- Schema definitions in `packages/database/src/schema/`

### AI/ML Pipeline System

- Tools and pipelines for data processing
- Transformers for various data types (text-to-speech, text-to-image, etc.)
- Job queue system for background processing
- Artifact storage and management

### Authentication Flow

- JWT-based authentication with refresh tokens
- OAuth integration for third-party providers
- Session management with Redis
- API key authentication for external access
- Role-based access control with organization membership

## Infrastructure

### Kubernetes Deployment

- Base configuration in `deploy/kubernetes/base/`
- Environment-specific overlays in `deploy/kubernetes/overlays/`
- Skaffold configuration for development deployment
- Supports local development with Minikube

### Docker Configuration

- Multi-stage Dockerfiles for API and platform
- Optimized for production deployment
- Located in `deploy/dockerfiles/`

## Development Workflow

1. Install dependencies: `pnpm install`
2. Start development servers: `pnpm dev`
3. API runs on configured port (check config)
4. Platform runs on Vite dev server
5. Use `pnpm lint` and `pnpm typecheck` before commits
6. Run relevant tests before pushing changes

## Package Management

- Uses pnpm workspaces for monorepo management
- Shared dependencies managed through workspace catalog
- Each package has independent build and test configurations
- Turborepo handles build caching and task orchestration
